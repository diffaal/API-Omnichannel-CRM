package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func PostRequest(url url.URL, requestBody interface{}, authToken string) (*http.Response, error) {
	bodyBytes, err := json.Marshal(&requestBody)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(bodyBytes)

	request, err := http.NewRequest(http.MethodPost, url.String(), reader)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	if authToken != "" {
		bearerToken := fmt.Sprintf("Bearer %s", authToken)
		request.Header.Set("Authorization", bearerToken)
	}

	httpClient := &http.Client{}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetRequest(url url.URL, authToken string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	if authToken != "" {
		bearerToken := fmt.Sprintf("Bearer %s", authToken)
		request.Header.Set("Authorization", bearerToken)
	}

	httpClient := &http.Client{}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func CloseBody(response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}
