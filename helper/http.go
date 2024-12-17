package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"sync"

	"project/domain"
)

func Upload(wg *sync.WaitGroup, files []*multipart.FileHeader) ([]domain.CdnResponse, error) {
	var results []domain.CdnResponse
	var err error
	for _, file := range files {
		wg.Add(1)

		go func() {
			defer wg.Done()
			var f multipart.File
			f, err = file.Open()
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			var part io.Writer
			part, err = writer.CreateFormFile("image", filepath.Base(""))
			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(part, f)

			err = writer.Close()
			if err != nil {
				log.Fatal(err)
			}

			var request *http.Request
			request, err = http.NewRequest("POST", "https://cdn-lumoshive-academy.vercel.app/api/v1/upload", body)
			if err != nil {
				log.Fatal(err)
			}
			request.Header.Add("Content-Type", writer.FormDataContentType())
			client := &http.Client{}
			var response *http.Response
			response, err = client.Do(request)
			if err != nil {
				log.Fatal(err)
			}

			defer response.Body.Close()

			var res []byte
			res, err = io.ReadAll(response.Body)
			if err != nil {
				log.Fatal("Error reading response:", err)
			}

			var result domain.CdnResponse
			json.Unmarshal(res, &result)
			results = append(results, result)
		}()
	}
	wg.Wait()
	return results, err
}
