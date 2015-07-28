package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func IPFSUpload(url string, state []byte, w io.Writer) (string, error) {

	input := bytes.NewReader(state)
	request, err := http.NewRequest("POST", url, input)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	hash, ok := response.Header["Ipfs-Hash"]
	if !ok || hash[0] == "" {
		return "", fmt.Errorf("No hash returned")
	}
	return hash[0], nil
}

func IPFSCat(url, fileHash string, api bool, w io.Writer) ([]byte, error) {
	var request *http.Request
	var err error
	if api {
		reqUrl := url + "cat?arg=" + fileHash
		request, err = http.NewRequest("POST", reqUrl, nil)
		if err != nil {
			return []byte(""), err
		}
	} else {
		reqUrl := url + fileHash
		request, err = http.NewRequest("GET", reqUrl, nil)
		if err != nil {
			return []byte(""), err
		}
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return []byte(""), err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte(""), err
	}

	var errs struct {
		Message string
		Code    int
	}
	if response.StatusCode >= http.StatusBadRequest {
		//TODO better err handling; this is a (very) slimed version of how IPFS does it.
		if err = json.Unmarshal(body, &errs); err != nil {
			return []byte(""), fmt.Errorf("error json unmarshaling body1 %v", err)
		}
		return []byte(errs.Message), nil

		if response.StatusCode == http.StatusNotFound {
			if err = json.Unmarshal(body, &errs); err != nil {
				return []byte(""), fmt.Errorf("error json unmarshaling body2 %v", err)
			}
			return []byte(errs.Message), nil
		}
	}
	return body, nil
}

//-------------------------------------------------------------------------------
//helpers

func composeIPFSUrl(host string, api bool) string {
	var url string
	var urlExt string
	if api {
		urlExt = ":5001/api/v0/"
	} else {
		urlExt = ":8080/ipfs/"
	}
	if host == "" {
		url = "http://0.0.0.0" + urlExt
	} else {
		url = "http://" + host + urlExt
	}
	return url
}
