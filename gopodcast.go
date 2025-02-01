package gopodcast

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// TODO better types on XML structs, plus custom marshallers

type Feed struct {
	XMLName      xml.Name `xml:"rss"`
	Version      string   `xml:"version,attr"`
	XMLNSContent string   `xml:"xmlns:content,attr"`
	XMLNSPodcast string   `xml:"xmlns:podcast,attr"`
	XMLNSAtom    string   `xml:"xmlns:atom,attr"`
	XMLNSITunes  string   `xml:"xmlns:itunes,attr"`
	Channel      *Channel `xml:"channel"`
}

func (f *Feed) WriteFeedXML(w io.Writer) error {
	// add standard feed attrs
	f.Version = "2.0"
	f.XMLNSContent = "http://purl.org/rss/1.0/modules/content/"
	f.XMLNSPodcast = "https://podcastindex.org/namespace/1.0"
	f.XMLNSAtom = "http://www.w3.org/2005/Atom"
	f.XMLNSITunes = "http://www.itunes.com/dtds/podcast-1.0.dtd"

	w.Write([]byte(xml.Header))
	return xml.NewEncoder(w).Encode(f)
}

type Channel struct {
	// PSP Required
	AtomLink       AtomLink         `xml:"atom:link"`
	Title          string           `xml:"title"`
	Description    Description      `xml:"description"`
	Link           string           `xml:"link"`
	Language       string           `xml:"language"`
	ITunesCategory []ITunesCategory `xml:"itunes:category"`
	ITunesExplicit bool             `xml:"itunes:explicit"`
	ITunesImage    ITunesImage      `xml:"itunes:image"`

	// PSP Recommended
	PodcastLocked string `xml:"podcast:locked,omitempty"` // TODO custom marshaller (yes/no)
	PodcastGUID   string `xml:"podcast:guid,omitempty"`
	ITunesAuthor  string `xml:"itunes:author,omitempty"`

	// PSP Optional
	Copyright      string          `xml:"copyright,omitempty"`
	PodcastText    *PodcastText    `xml:"podcast:txt,omitempty"`
	PodcastFunding *PodcastFunding `xml:"podcast:funding,omitempty"`
	ITunesType     string          `xml:"itunes:type,omitempty"`
	ITunesComplete string          `xml:"itunes:complete,omitempty"` // TODO custom marshaller

	// Other fields
	// TODO other podcast index namespace fields
	// TODO other itunes fields

	Items []*Item `xml:"item"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Description struct {
	Text string `xml:",cdata"`
}

type ITunesCategory struct {
	Text            string `xml:"text,attr"`
	SubCategoryText string `xml:"itunes:category,omitempty>text,attr,omitempty"`
}

type ITunesImage struct {
	Href string `xml:"href,attr"`
}

type PodcastText struct {
	Purpose string `xml:"purpose,attr,omitempty"`
	Text    string `xml:",chardata"`
}

type PodcastFunding struct {
	URL  string `xml:"url,attr"`
	Text string `xml:",chardata"`
}

type Item struct {
	// PSP required
	Title     string    `xml:"title"`
	Enclosure Enclosure `xml:"enclosure"`
	GUID      ItemGUID  `xml:"guid"`

	// PSP Recommended
	Link              string              `xml:"link,omitempty"`
	PubDate           string              `xml:"pubDate,omitempty"` // TODO time.Time with custom marshaller
	Description       *Description        `xml:"description,omitempty"`
	ITunesDuration    string              `xml:"itunes:duration,omitempty"`
	ITunesImage       *ITunesImage        `xml:"itunes:image,omitempty"`
	ITunesExplicit    *bool               `xml:"itunes:explicit,omitempty"`
	PodcastTranscript []PodcastTranscript `xml:"podcast:transcript,omitempty"`

	// PSP Optional
	ITunesEpisode     string `xml:"itunes:episode,omitempty"`
	ITunesSeason      string `xml:"itunes:season,omitempty"`
	ITunesEpisodeType string `xml:"itunes:episodeType,omitempty"`
	ITunesBlock       string `xml:"itunes:block,omitempty"`

	// Other Fields
	// TODO itunes, podcast index namespace
}

type Enclosure struct {
	Length int64  `xml:"length,attr"`
	Type   string `xml:"type,attr"`
	URL    string `xml:"url,attr"`
}

type ItemGUID struct {
	IsPermaLink *bool  `xml:"isPermaLink,attr,omitempty"`
	Text        string `xml:",chardata"`
}

type PodcastTranscript struct {
	URL      string `xml:"url,attr"`
	Type     string `xml:"type,attr"`
	Rel      string `xml:"rel,attr,omitempty"`
	Language string `xml:"language,attr,omitempty"`
}

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

func (p *Parser) ParseFeedFromURL(ctx context.Context, url string) (*Feed, error) {
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
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("non-200 http response '%d'", res.StatusCode)
	}

	return p.ParseFeed(res.Body)
}

func (p *Parser) ParseFeed(r io.Reader) (*Feed, error) {
	var feed Feed
	err := xml.NewDecoder(r).Decode(&feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}
