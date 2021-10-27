package genWords

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetResp(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		fmt.Println("error status:", resp.StatusCode)
		return nil, fmt.Errorf("error status: %d", resp.StatusCode)
	}
	return body, nil
}

func Bytes2Map(b []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func GetWords() (string, error) {
	body, err := GetResp("https://soul-soup.fe.workers.dev/:id")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// bytes to map
	m, err := Bytes2Map(body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return m["title"].(string), nil
}
