package util

import (
	"bytes"
	"encoding/json"
	"github.com/josexy/gochatroom/logx"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

const uploadImageUrl = "https://freeimage.host/json"

type uploadResp struct {
	Image struct {
		DisplayUrl string `json:"display_url"`
	} `json:"image"`
}

func UploadFile(uploadFile string) string {
	body := bytes.NewBuffer(nil)

	writer := multipart.NewWriter(body)
	_ = writer.WriteField("action", "upload")
	_ = writer.WriteField("type", "file")
	fileWriter, err := writer.CreateFormFile("source", path.Base(uploadFile))

	file, err := os.Open(uploadFile)
	_, _ = io.Copy(fileWriter, file)
	_ = file.Close()

	_ = writer.Close()
	req, err := http.NewRequest(http.MethodPost, uploadImageUrl, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logx.ErrorBy(err)
		return ""
	}
	respBody := bytes.NewBuffer(nil)

	_, err = respBody.ReadFrom(resp.Body)
	defer resp.Body.Close()

	logx.Debug("%s", respBody.String())

	if err != nil {
		logx.ErrorBy(err)
		return ""
	}

	var data uploadResp
	if err = json.Unmarshal(respBody.Bytes(), &data); err != nil {
		return ""
	}
	return data.Image.DisplayUrl
}
