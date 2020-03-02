package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func Request(method, endpoint string, token string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", os.Getenv("BACKEND_ADDRESS"), endpoint)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=%s", "JWT", token))
	return http.DefaultClient.Do(req)
}

func AdminRequest(method, endpoint string, token string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", os.Getenv("BACKEND_ADDRESS"), endpoint)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)
	return http.DefaultClient.Do(req)
}

func GetAuth(cookie string) (uint64, error) {
	res, err := Request(http.MethodGet, "/auth", cookie, nil)
	if err != nil {
		return 0, err
	}

	if res.StatusCode != http.StatusOK {
		return 0, UnableToAuthErr{}
	}

	resBody := new(UserResponse)
	err = json.NewDecoder(res.Body).Decode(resBody)
	if err != nil {
		return 0, err
	}
	return resBody.User.ID, nil
}
