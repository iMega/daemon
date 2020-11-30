package helper

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// Request getting data from testing http server
func Request(uri string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("not ok")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := res.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil
}
