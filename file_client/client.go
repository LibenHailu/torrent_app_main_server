package fileclient

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	myFilePath "path/filepath"

	filesource "github.com/LibenHailu/grpc_file_stream/file_stream/file_source"
	"github.com/LibenHailu/grpc_file_stream/file_stream/filepb"
	"github.com/LibenHailu/grpc_file_stream/file_stream/service/client"
	"google.golang.org/grpc"
)

var (
	c filepb.FileServiceClient
)

func Connect() *grpc.ClientConn {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not found connect %v ", err)
	}

	defer cc.Close()

	return cc
}

func main() {

	fmt.Println("Hello i am file stream client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not found connect %v ", err)
	}

	defer cc.Close()

	c = filepb.NewFileServiceClient(cc)

	// UploadFile(c, "bini", "C:/Users/Liben/Desktop/Liben.jpg")
	// DownloadFile(c, "Liben.jpg")
	fileNames := []string{"sdf", "asdf"}
	RegisterPeers(c, "as", 253, fileNames)
}

func UploadFile(c filepb.FileServiceClient, fileID string, filepath string) {
	file, err := os.Open(filepath)

	if err != nil {
		log.Fatalf("couldn't open file: %v", err)
	}

	defer file.Close()

	stream, err := c.UploadFile(context.Background())

	if err != nil {
		log.Fatalf("couldn't upload file %v", err)
	}

	fileStat, err := os.Stat(filepath)

	if err != nil {
		log.Fatal(err)
	}

	req := &filepb.UploadFileRequest{
		Data: &filepb.UploadFileRequest_Info{
			Info: &filepb.FileInfo{
				FileId:   fileID,
				FileType: myFilePath.Ext(filepath),
				FileName: fileStat.Name(),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatalf("couldn't send file %v", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {

		n, err := reader.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("couldn't read chunk to buffer %v", err)
		}

		req := &filepb.UploadFileRequest{
			Data: &filepb.UploadFileRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		err = stream.Send(req)
		if err != nil {
			log.Fatalf("couldn't send chunk to server %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("couldn't recive response %v", err)
	}

	log.Printf("file uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
}

func DownloadFile(c filepb.FileServiceClient, fileName string) {
	req := &filepb.ServeFileRequest{
		FileName: fileName,
	}
	resStream, err := c.DownloadFile(context.Background(), req)
	if err != nil {
		log.Fatalf("error downloading file: %v", err)
	}
	fileData := bytes.Buffer{}
	fileSize := 0
	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			// we've reached the end of stream
			log.Println("recived all chunks")
			break
		}
		if err != nil {
			log.Fatalf("error while reciving chunk %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v ", msg.ChunkData)
		chunk := msg.GetChunkData()
		size := len(chunk)

		fileSize += size

		// if fileSize > maxFileSize {
		// 	return logError(status.Errorf(codes.InvalidArgument, "file is too large: %d > %d", fileSize, maxFileSize))
		// }

		_, err = fileData.Write(chunk)
		if err != nil {
			log.Fatal("couldn't write chunk data: %v", err)
		}

	}

	clientSave := client.NewDiskFileStore("C:/Users/Liben/Desktop/dsp")
	clientSave.Save(fileData, fileName)

}

func RegisterPeers(c filepb.FileServiceClient, ip string, port int32, fileNames []string) {
	filesource.SearchAddressForThefile("Liben")
	// myPort, _ := strconv.Atoi(port)
	req := &filepb.RegisterPeersRequest{
		Ip:        ip,
		Port:      port,
		FileNames: fileNames,
	}

	res, err := c.RegisterPeers(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Register Peers RPC: %v", err)
	}

	log.Printf("Response form Sum: %v", res.ServerAddress)
}
