package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func MakePostRequest(url string, data interface{}, headers map[string]string) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataBytes))
	if err != nil {
		fmt.Printf("sendWechatyMsg NewRequest error: %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("sendWechatyMsg client do error: %v\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func MakePostFormRequest(url string, data map[string]string, headers map[string]string) ([]byte, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	for key, value := range data {
		_ = writer.WriteField(key, value)
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Printf("sendWechatyMsg NewRequest error: %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("sendWechatyMsg client do error: %v\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func MakeGetRequest(url string, headers map[string]string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("makeGetRequest NewRequest error: %v\n", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("makeGetRequest client do error: %v\n", err)
		return nil
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body
}

func MakePostFileRequest(url string, file *multipart.FileHeader, form map[string]string, headers map[string]string) ([]byte, error) {
	uploadFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer uploadFile.Close()
	// 创建一个用于存储文件内容的字节缓冲区
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件字段和值
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, uploadFile)
	if err != nil {
		return nil, err
	}
	for k, v := range form {
		err := writer.WriteField(k, v)
		if err != nil {
			return nil, err
		}
	}
	// 完成multipart数据的写入并获取Content-Type
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

func MakeDeleteRequest(url string, data interface{}, headers map[string]string) ([]byte, error) {
	dataBytes := []byte{}
	if data != nil {
		var err error
		dataBytes, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(dataBytes))
	if err != nil {
		fmt.Printf("sendWechatyMsg NewRequest error: %v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("sendWechatyMsg client do error: %v\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func GetImageBuffer(imageUrl string) (*bytes.Buffer, error) {
	// 创建一个Buffer来存储图片内容
	imageBuffer := new(bytes.Buffer)
	imageResponse, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println("无法下载图片:", err)
		return imageBuffer, err
	}
	defer imageResponse.Body.Close()

	// 将图片内容复制到Buffer中
	_, err = io.Copy(imageBuffer, imageResponse.Body)
	if err != nil {
		fmt.Println("无法复制图片内容:", err)
		return imageBuffer, err
	}
	return imageBuffer, err
}
