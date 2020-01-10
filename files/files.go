package files

import (
	"bufio"
	"errors"
	"io"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// WriteToGridFile 은 GridFile로 파일 쓰기를 수행합니다.
func WriteToGridFile(file multipart.File, uploadStream *gridfs.UploadStream) error {
	reader := bufio.NewReader(file)
	defer func() { file.Close() }()
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
