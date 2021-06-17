package filesource

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// in memory ip->file names store

var (
	sources = map[string][]string{}
	mutex   sync.RWMutex
)

type Source struct {
	Address   string   `json:"address"`
	FileNames []string `json:file_names`
}

// add address of connected clients to the server
func AddToSources(serverAddress string, fileNames []string) {
	mutex.Lock()
	defer mutex.Unlock()
	sources, err := os.Open("C:/Users/Liben/go/src/github.com/LibenHailu/grpc_file_stream/file_stream/file_source/sources.json")
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer sources.Close()

	content, err := ioutil.ReadFile("C:/Users/Liben/go/src/github.com/LibenHailu/grpc_file_stream/file_stream/file_source/sources.json")

	if err != nil {
		log.Fatal(err)
	}

	var tmpSources []Source

	err = json.Unmarshal(content, &tmpSources)

	if err != nil {
		log.Fatal(err)
	}

	var newSource Source
	newSource.Address = serverAddress
	newSource.FileNames = fileNames

	tmpSources = append(tmpSources, newSource)

	jsonData, err := json.Marshal(tmpSources)

	err = ioutil.WriteFile("C:/Users/Liben/go/src/github.com/LibenHailu/grpc_file_stream/file_stream/file_source/sources.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// sources.Write(jsonData)

	// sources[serverAddress] = fileNames
	// for _, fileNames := range sources {
	// 	for _, fName := range fileNames {
	// 		fmt.Println(fName)
	// 	}
	// }
}

// searching files to grant access
func SearchAddressForThefile(fileName string) (*string, *string) {
	var serverAddr string
	// client := "client"
	server := ""

	content, err := ioutil.ReadFile("C:/Users/Liben/go/src/github.com/LibenHailu/grpc_file_stream/file_stream/file_source/sources.json")

	if err != nil {
		log.Fatal(err)
	}

	var curSources []Source

	err = json.Unmarshal(content, &curSources)

	for _, source := range curSources {
		for _, fName := range source.FileNames {
			if fName == fileName {
				serverAddr = source.Address
				server = "client"
				return &serverAddr, &server
			}
		}
	}
	// for address, fileNames := range  {
	// 	fmt.Println(address)
	// 	for _, fName := range fileNames {
	// 		if fileName == fName {
	// 			serverAddr = address
	// 			server = "client"
	// 			return &serverAddr, &server
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
			server = "server"
			return &serverAddr, &server
		}
	}
	return nil, nil
}
