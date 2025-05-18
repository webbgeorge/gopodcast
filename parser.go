package gopodcast

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type Parser struct {
	HTTPClient      *http.Client
	UserAgent       string
	AuthCredentials *AuthCredentials
}

type AuthCredentials struct {
	Username string
	Password string
}

const defaultUserAgent = "gopodcast/1.0"

func NewParser() *Parser {
	return &Parser{
		HTTPClient: http.DefaultClient,
		UserAgent:  defaultUserAgent,
	}
}

func (p *Parser) ParseFeedFromURL(ctx context.Context, url string) (pc *Podcast, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", p.UserAgent)

	if p.AuthCredentials != nil && p.AuthCredentials.Username != "" && p.AuthCredentials.Password != "" {
		req.SetBasicAuth(p.AuthCredentials.Username, p.AuthCredentials.Password)
	}

	res, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		errVal := res.Body.Close()
		if errVal != nil {
			err = errVal
		}
	}()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("non-200 http response '%d'", res.StatusCode)
	}

	return p.ParseFeed(res.Body)
}

func (p *Parser) ParseFeed(r io.Reader) (*Podcast, error) {
	var feed xmlFixfeed
	err := xml.NewDecoder(r).Decode(&feed)
	if err != nil {
		return nil, err
	}
	return feed.Translate().Channel, nil
}
