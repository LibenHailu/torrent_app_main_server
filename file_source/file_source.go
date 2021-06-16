package filesource

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// in memory ip->file names store

var (
	sources = map[string][]string{}
	mutex   sync.RWMutex
)

// add address of connected clients to the server
func AddToSources(serverAddress string, fileNames []string) {
	mutex.Lock()
	defer mutex.Unlock()
	sources[serverAddress] = fileNames
	fmt.Println("hehehe")
	for _, fileNames := range sources {
		for _, fName := range fileNames {
			fmt.Println(fName)
		}
	}
}

// searching files to grant access
func SearchAddressForThefile(fileName string) (*string, *string) {
	var serverAddr string
	// client := "client"
	server := "server"
	// for address, fileNames := range sources {
	// 	for _, fName := range fileNames {
	// 		if fileName == fName {
	// 			return address
	// 		}
	// 	}
	// }

	file, err := os.Open("C:/Users/Liben/go/src/github.com/LibenHailu/grpc_file_stream/file_stream/file")
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	list, _ := file.Readdirnames(0) // 0 to read all files and folders
	for _, name := range list {
		if name == fileName {
			serverAddr = "127.0.0.1:50051"
			return &serverAddr, &server
		}
	}
	return nil, nil
}
