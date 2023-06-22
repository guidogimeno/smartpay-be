package services

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/guidogimeno/smartpay-be/types"
	"github.com/shopspring/decimal"
	"golang.org/x/net/html"
)

const (
	bcraUrl = "https://www.bcra.gob.ar/PublicacionesEstadisticas/Principales_variables_datos.asp"

	trTag = "tr"
	tdTag = "td"

	comma = ","
	dot   = "."

	inflationSerie = "7931"
)

func ScrapInflation(startDate string, finishDate string) ([]*types.Inflation, error) {
	formData := url.Values{
		"fecha_desde": {startDate},
		"fecha_hasta": {finishDate},
		"primeravez":  {"1"},
		"serie":       {inflationSerie},
	}

	body := createFormBody(formData)
	document, err := postRequest(body)
	if err != nil {
		return nil, err
	}

	return parseInflationIndexes(document)
}

func createFormBody(formData url.Values) io.Reader {
	encodedForm := formData.Encode()
	return strings.NewReader(encodedForm)
}

func postRequest(body io.Reader) (*html.Node, error) {
	req, err := http.NewRequest(http.MethodPost, bcraUrl, body)
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

func parseInflationIndexes(document *html.Node) ([]*types.Inflation, error) {
	inflationIndexes := []*types.Inflation{}
	htmlIterator(document, func(n *html.Node) {
		if !isTableRow(n) {
			return
		}

		inflation, err := buildInflation(n)
		if err != nil {
			return
		}

		inflationIndexes = append(inflationIndexes, inflation)
	})

	return inflationIndexes, nil
}

func htmlIterator(n *html.Node, callback func(*html.Node)) {
	callback(n)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		htmlIterator(c, callback)
	}
}

func isTableRow(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == trTag && n.FirstChild.NextSibling.Data == tdTag
}

func buildInflation(n *html.Node) (*types.Inflation, error) {
	firstNode := n.FirstChild.NextSibling
	secondNode := firstNode.NextSibling.NextSibling

	rawDate := firstNode.FirstChild.Data
	rawIndex := secondNode.FirstChild.Data

	trimmedIndex := strings.TrimSpace(rawIndex)
	formattedIndex := strings.ReplaceAll(trimmedIndex, comma, dot)

	index, err := decimal.NewFromString(formattedIndex)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	inflation := &types.Inflation{
		Date:  rawDate,
		Index: index,
	}

	return inflation, nil
}
