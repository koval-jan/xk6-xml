package xml

import (
	"html"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/clbanning/mxj"
	"go.k6.io/k6/js/modules"
)

// XML is the k6 xml extension.
type XML struct {
	vu modules.VU
}

// Parse parses xml
func (*XML) Parse(body string) (map[string]interface{}, error) {
	return mxj.NewMapXml([]byte(body))
}

// Query xml with xpath
func (*XML) FindOne(path string, body string) (string, error) {
	doc, err := xmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return "", err
	}
	result, err := xmlquery.Query(doc, path)
	if err != nil {
		return "", err
	}
	if result == nil {
		return "", nil
	}
	return result.OutputXMLWithOptions(xmlquery.WithoutComments(), xmlquery.WithPreserveSpace()), nil
}

func (*XML) FindAll(path string, body string) ([]string, error) {
	doc, err := xmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return make([]string, 0), err
	}

	nodes, err := xmlquery.QueryAll(doc, path)
	if err != nil {
		return make([]string, 0), err
	}
	var result []string
	for _, node := range nodes {
		result = append(result, node.OutputXMLWithOptions(xmlquery.WithoutComments(), xmlquery.WithOutputSelf()))
	}
	return result, nil
}

func (*XML) FindOneNS(ns map[string]string, path string, body string) (string, error) {
	doc, err := xmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return "", err
	}

	exp, err := xpath.CompileWithNS(path, ns)
	if err != nil {
		return "", err
	}

	result := xmlquery.QuerySelector(doc, exp)
	if result == nil {
		return "", nil
	}
	return result.OutputXMLWithOptions(xmlquery.WithoutComments()), nil
}

func (*XML) FindAllNS(ns map[string]string, path string, body string) ([]string, error) {
	doc, err := xmlquery.Parse(strings.NewReader(body))
	if err != nil {
		return make([]string, 0), err
	}

	exp, err := xpath.CompileWithNS(path, ns)
	if err != nil {
		return make([]string, 0), err
	}

	nodes := xmlquery.QuerySelectorAll(doc, exp)
	if nodes == nil {
		return make([]string, 0), nil
	}
	var result []string
	for _, node := range nodes {
		result = append(result, node.OutputXMLWithOptions(xmlquery.WithoutComments(), xmlquery.WithOutputSelf()))
	}
	return result, nil
}

// Encode xml
func (*XML) EncodeXml(input string) (string, error) {
	return html.EscapeString(input), nil
}

// Decode xml
func (*XML) DecodeXml(input string) (string, error) {
	return html.UnescapeString(input), nil
}

func newXml(vu modules.VU) *XML {
	return &XML{vu: vu}
}
