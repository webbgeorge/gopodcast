//go:generate go run generate/main.go

package gopodcast

import (
	"encoding/xml"
	"io"
)

// TODO better types on XML structs, plus custom marshallers

type feed struct {
	XMLName      xml.Name `xml:"rss"`
	Version      string   `xml:"version,attr"`
	XMLNSContent string   `xml:"xmlns:content,attr"`
	XMLNSPodcast string   `xml:"xmlns:podcast,attr"`
	XMLNSAtom    string   `xml:"xmlns:atom,attr"`
	XMLNSITunes  string   `xml:"xmlns:itunes,attr"`
	Channel      *Podcast `xml:"channel"`
}

var emptyFeed = &feed{
	Version:      "2.0",
	XMLNSContent: "http://purl.org/rss/1.0/modules/content/",
	XMLNSPodcast: "https://podcastindex.org/namespace/1.0",
	XMLNSAtom:    "http://www.w3.org/2005/Atom",
	XMLNSITunes:  "http://www.itunes.com/dtds/podcast-1.0.dtd",
}

type Podcast struct {
	// PSP Required
	AtomLink       AtomLink         `xml:"atom:link"`
	Title          string           `xml:"title"`
	Description    Description      `xml:"description"`
	Link           string           `xml:"link"`
	Language       string           `xml:"language"`
	ITunesCategory []ITunesCategory `xml:"itunes:category"`
	ITunesExplicit FlexBool         `xml:"itunes:explicit"`
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

func (p *Podcast) WriteFeedXML(w io.Writer) error {
	feed := emptyFeed
	feed.Channel = p
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	return xml.NewEncoder(w).Encode(feed)
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
	Text        string          `xml:"text,attr"`
	SubCategory *ITunesCategory `xml:"itunes:category,omitempty"`
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
	ITunesDuration    string              `xml:"itunes:duration,omitempty"` // TODO custom marshaller
	ITunesImage       *ITunesImage        `xml:"itunes:image,omitempty"`
	ITunesExplicit    *FlexBool           `xml:"itunes:explicit,omitempty"`
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
	IsPermaLink *FlexBool `xml:"isPermaLink,attr,omitempty"`
	Text        string    `xml:",chardata"`
}

type PodcastTranscript struct {
	URL      string `xml:"url,attr"`
	Type     string `xml:"type,attr"`
	Rel      string `xml:"rel,attr,omitempty"`
	Language string `xml:"language,attr,omitempty"`
}
