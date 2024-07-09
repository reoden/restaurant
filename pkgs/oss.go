package pkgs

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func getBucket(isGlobal bool) (*oss.Bucket, error) {
	if isGlobal {
		return getGlobalBucket()
	} else {
		return getLocalBucket()
	}
}

func getLocalBucket() (*oss.Bucket, error) {
	client, err := oss.New("https://file.shiflow.com", "LTAI5tMH1obRmTtMboBpJqvU", "bEiRZV34a0CZO6s6MGx3yXW232wTpR", oss.UseCname(true))
	if err != nil {
		// HandleError(err)
		return nil, err
	}

	bucket, err := client.Bucket("shiflow-file")
	if err != nil {
		// HandleError(err)
		return nil, err
	}
	return bucket, nil
}

func getGlobalBucket() (*oss.Bucket, error) {
	client, err := oss.New("global.shiflow.com", "LTAI5tMH1obRmTtMboBpJqvU", "bEiRZV34a0CZO6s6MGx3yXW232wTpR", oss.UseCname(true))
	if err != nil {
		// HandleError(err)
		return nil, err
	}

	bucket, err := client.Bucket("shiflow-global")
	if err != nil {
		// HandleError(err)
		return nil, err
	}
	return bucket, nil
}

// 上传文件到oss
func UploadOSS(file []byte, filename string, isGlobal bool) error {
	bucket, err := getBucket(isGlobal)
	if err != nil {
		return err
	}
	return bucket.PutObject(filename, bytes.NewReader(file))
}

func UploadOSSByFilePath(filePath string, filename string, isGlobal bool) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("UploadOSSByFilePath: %v", err)
		return err
	}

	return UploadOSS(data, filename, isGlobal)
}

func DeleteOSS(filename string, isGlobal bool) error {
	bucket, err := getBucket(isGlobal)
	if err != nil {
		return err
	}
	return bucket.DeleteObject(filename)
}

// 获取加密链接
func SignedUrl(filename string, isGlobal bool) (string, error) {
	bucket, err := getBucket(isGlobal)
	if err != nil {
		return "", err
	}

	return bucket.SignURL(filename, oss.HTTPGet, 3600)
}

func UploadAndFilename(ctx context.Context, bs []byte, fileType string, isGlobal bool) (string, error) {
	filename := fmt.Sprintf("%d.%s", time.Now().UnixMilli(), fileType)
	log.Println("UploadAndFilename_filename", fmt.Sprintf("filename:%s", filename))
	for i := 0; i < 3; i++ {
		err := UploadOSS(bs, filename, isGlobal)
		if err != nil {
			log.Printf("UploadAndFilename_UploadOSS_err: %+v", err)
			// SendFeishuMsg(ctx, "error", fmt.Sprintf("Error UploadAndFilename-UploadOSS-try:%d", i), fmt.Sprintf("filename:\n%s\n\n, err:\n %+v", filename, err))
			continue
		}
	}
	return filename, nil
}
