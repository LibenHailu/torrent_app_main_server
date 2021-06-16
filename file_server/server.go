package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	filesource "github.com/LibenHailu/grpc_file_stream/file_stream/file_source"
	"github.com/LibenHailu/grpc_file_stream/file_stream/filepb"
	"github.com/LibenHailu/grpc_file_stream/file_stream/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// const (
// 	maxFileSize = 1 << 20
// )

type server struct {
	fileStore service.FileStore
	filepb.UnimplementedFileServiceServer
}

func NewServer(fileStore service.FileStore) *server {
	return &server{
		fileStore: fileStore,
	}
}
func (s *server) RegisterPeers(ctx context.Context, req *filepb.RegisterPeersRequest) (*filepb.RegisterPeersResponse, error) {

	fmt.Println("reciving address")
	ip := req.Ip
	port := req.Port
	fileNames := req.FileNames

	filesource.AddToSources(fmt.Sprintf("%s:%d", ip, port), fileNames)

	res := &filepb.RegisterPeersResponse{
		ServerAddress: fmt.Sprintf("%s:%d", ip, port),
	}

	return res, nil
}

func (s *server) DownloadFile(req *filepb.ServeFileRequest, res filepb.FileService_DownloadFileServer) error {
	bufferSize := 64 * 1024 //64KiB, tweak this as desired
	file, err := os.Open("../file/" + req.GetFileName())
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	buff := make([]byte, bufferSize)
	for {
		bytesRead, err := file.Read(buff)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		resp := &filepb.ServeFileResponse{
			ChunkData: buff[:bytesRead],
		}
		err = res.Send(resp)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}

	}
	return nil
}

func (s *server) UploadFile(stream filepb.FileService_UploadFileServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Printf("couldn't recive file info ", err)
		return status.Errorf(codes.Unknown, "couldn't recive file info")
	}

	fileID := req.GetInfo().GetFileId()
	fileType := req.GetInfo().GetFileType()
	fileName := req.GetInfo().GetFileName()

	log.Println("recived an upload file with id %s with type %s", fileID, fileType)

	fileData := bytes.Buffer{}
	fileSize := 0

	for {
		log.Println("waiting to recive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("no more data")
			break
		}

		if err != nil {
			return logError(status.Errorf(codes.Unknown, "couldn't recive chunk data: %v", err))
		}
		chunk := req.GetChunkData()
		size := len(chunk)

		fileSize += size

		// if fileSize > maxFileSize {
		// 	return logError(status.Errorf(codes.InvalidArgument, "file is too large: %d > %d", fileSize, maxFileSize))
		// }

		_, err = fileData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "couldn't write chunk data: %v", err))
		}
	}

	result, err := s.fileStore.Save(fileID, fileType, fileData, fileName)

	if err != nil {
		return logError(status.Errorf(codes.Internal, "couldn't save file: %v", err))
	}

	res := &filepb.UploadFileResponse{
		Id:   result,
		Size: uint32(fileSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "couldn't send response: %v", err))
	}

	log.Printf("saved file with id: %s, size: %d", result, fileSize)

	return nil
}

func logError(err error) error {
	if err != nil {
		log.Println(err)
	}
	return err
}

func main() {

	fmt.Println("file stream server")

	fileStore := service.NewDiskFileStore("../file")

	fileServer := NewServer(fileStore)
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v ", err)
	}

	s := grpc.NewServer()
	filepb.RegisterFileServiceServer(s, fileServer)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve:", err)
	}
}
