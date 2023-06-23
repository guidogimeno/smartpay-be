package scrapper

import (
	"net/url"
	"strings"

	"github.com/guidogimeno/smartpay-be/http"
	"github.com/guidogimeno/smartpay-be/parser"
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
	tnaSerie       = "7935"
)

func ScrapInflation(startDate string, finishDate string) ([]*types.FinancialData, error) {
	return scrapBCRAStadistics(inflationSerie, startDate, finishDate)
}

func ScrapTNA(startDate string, finishDate string) ([]*types.FinancialData, error) {
	return scrapBCRAStadistics(tnaSerie, startDate, finishDate)
}

func scrapBCRAStadistics(serie string,
	startDate string,
	finishDate string,
) ([]*types.FinancialData, error) {
	formData := url.Values{
		"fecha_desde": {startDate},
		"fecha_hasta": {finishDate},
		"primeravez":  {"1"},
		"serie":       {serie},
	}

	body := http.CreateFormBody(formData)
	document, err := http.PostRequest(bcraUrl, body)
	if err != nil {
		return nil, err
	}

	return parseFinancialData(document)
}

func parseFinancialData(document *html.Node) ([]*types.FinancialData, error) {
	financialDataOverTime := []*types.FinancialData{}
	parser.HtmlIterator(document, func(n *html.Node) {
		if !isTableRow(n) {
			return
		}

		financialData, err := buildFinancialData(n)
		if err != nil {
			return
		}

		financialDataOverTime = append(financialDataOverTime, financialData)
	})

	return financialDataOverTime, nil
}

func isTableRow(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == trTag && n.FirstChild.NextSibling.Data == tdTag
}

func buildFinancialData(n *html.Node) (*types.FinancialData, error) {
	firstNode := n.FirstChild.NextSibling
	secondNode := firstNode.NextSibling.NextSibling

	rawDate := firstNode.FirstChild.Data
	rawIndex := secondNode.FirstChild.Data

	trimmedIndex := strings.TrimSpace(rawIndex)
	formattedIndex := strings.ReplaceAll(trimmedIndex, comma, dot)

	index, err := decimal.NewFromString(formattedIndex)
	if err != nil {
		return nil, err
	}

	financialData := &types.FinancialData{
		Date:  rawDate,
		Index: index.Div(decimal.NewFromInt(100)),
	}

	return financialData, nil
}
