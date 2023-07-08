package client

import (
	"bufio"
	"fmt"
	"os"
)

func LoadPhoto() []byte {
	fileToBeUploaded := "/home/ubuntu/starlink/utils/client/airport.jpg"
	file, err := os.Open(fileToBeUploaded)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	//将文件读成字节
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return bytes
}
