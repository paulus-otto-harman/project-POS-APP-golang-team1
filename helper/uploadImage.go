package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func UploadFileThirdPartyAPI(file io.Reader, filename string) (string, error) {
	// Siapkan body multipart/form-data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Tambahkan file ke form-data
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return "", fmt.Errorf("could not create form file: %w", err)
	}

	// Salin konten file ke bagian multipart
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("could not copy file content: %w", err)
	}

	err = writer.WriteField("folder", "example_folder")
	if err != nil {
		return "", fmt.Errorf("could not add folder field: %w", err)
	}

	// Menutup writer untuk mengirimkan boundary dan header
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("could not close writer: %w", err)
	}

	// Membuat request ke API pihak ketiga
	thirdPartyURL := "https://cdn-lumoshive-academy.vercel.app/api/v1/upload" // Ganti dengan URL API pihak ketiga
	req, err := http.NewRequest("POST", thirdPartyURL, body)
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}

	// Set Content-Type dari multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Kirim request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	respThirdParty, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}
	var result map[string]interface{}
	err = json.Unmarshal(respThirdParty, &result)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to upload to third party, status: %s", resp.Status)
	}

	data, _ := result["data"].(map[string]interface{})
	imageUrl, _ := data["url"].(string)

	return imageUrl, nil
}
