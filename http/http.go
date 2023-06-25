package http

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func CreateFormBody(formData url.Values) io.Reader {
	encodedForm := formData.Encode()
	return strings.NewReader(encodedForm)
}

func PostRequest(url string, body io.Reader) (*html.Node, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return html.Parse(strings.NewReader(string(respBody)))
}
