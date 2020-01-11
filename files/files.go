package files

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// WriteToGridFileFile 은 GridFile로 파일 쓰기를 수행합니다.
func WriteToGridFileFile(file multipart.File, uploadStream *gridfs.UploadStream) error {
	reader := bufio.NewReader(file)
	defer func() { file.Close() }()
	return WriteToGridFile(reader, uploadStream)
}

// WriteToGridFileString 은 String을 GridFile로 저장하는 작업을 수행합니다
func WriteToGridFileString(s string, uploadStream *gridfs.UploadStream) error {
	byteData := []byte(s)
	reader := bytes.NewReader(byteData)
	return WriteToGridFile(reader, uploadStream)
}

type readerType interface {
	Read(p []byte) (n int, err error)
}

// WriteToGridFile 은 GridFile로 쓰기 작업을 수행합니다
func WriteToGridFile(reader readerType, uploadStream *gridfs.UploadStream) error {
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return errors.New("Could not read the input file")
		}
		if n == 0 {
			break
		}
		if _, err := uploadStream.Write(buf[:n]); err != nil {
			return errors.New("Could not write to GridFs")
		}
	}
	uploadStream.Close()
	return nil
}
