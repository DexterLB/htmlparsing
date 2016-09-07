package htmlparsing

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/jbowtie/gokogiri"
	htmlParser "github.com/jbowtie/gokogiri/html"
	"github.com/jbowtie/gokogiri/xml"
	"github.com/sethgrid/pester"
)

type Client struct {
	pester.Client

	maxServerErrorRetries    int
	serverErrorRetryInterval time.Duration
}

func NewClient(settings *Settings) *Client {
	client := &Client{
		Client:                   *pester.New(),
		maxServerErrorRetries:    settings.MaxServerErrorRetries,
		serverErrorRetryInterval: settings.ServerErrorRetryInterval,
	}

	client.Transport = settings.Transport
	client.Timeout = settings.Timeout
	client.MaxRetries = settings.MaxHttpRetries
	client.Backoff = func(retry int) time.Duration {
		return settings.HttpRetryInterval
	}
	return client
}

func NewCookiedClient(settings *Settings) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cannot initialise cookie jar: %s", err)
	}

	client := NewClient(settings)
	client.Jar = jar

	return client, nil
}

func (client *Client) ParsePage(
	url string, formData url.Values,
) (*htmlParser.HtmlDocument, error) {
	data, err := client.OpenPage(url, formData)
	if err != nil {
		return nil, err
	}

	page, err := gokogiri.ParseHtml(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing html: %s", err)
	}
	return page, nil
}

func (client *Client) OpenPage(
	url string, formData url.Values,
) ([]byte, error) {
	var errors []string

	for i := 0; i < client.maxServerErrorRetries+1; i++ {
		if i > 0 {
			time.Sleep(client.serverErrorRetryInterval)
		}

		data, err := client.request(url, formData)
		if err == nil {
			return data, nil
		}
		errors = append(errors, err.Error())
	}

	return nil, fmt.Errorf(
		"unable to open %s: %s", url,
		strings.Join(errors, ", "),
	)
}

func (client *Client) request(
	url string, formData url.Values,
) ([]byte, error) {
	var (
		resp *http.Response
		err  error
	)

	if formData == nil {
		resp, err = client.Get(url)
	} else {
		resp, err = client.PostForm(url, formData)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to download page: %s", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read page: %s", err)
	}

	if strings.Contains(string(data), "An internal system error has occured.") {
		return nil, fmt.Errorf("internal server error")
	}

	return data, nil
}

// First returns the first child node of node which matches expression
func First(node xml.Node, expression string) (xml.Node, error) {
	nodes, err := node.Search(expression)
	if err != nil {
		return nil, err
	}
	if len(nodes) < 1 {
		return nil, fmt.Errorf("node not present")
	}
	return nodes[0], nil
}
