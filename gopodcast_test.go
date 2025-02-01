package gopodcast_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/webbgeorge/gopodcast"
)

func TestParseFeed(t *testing.T) {
	parser := gopodcast.NewParser()

	f, err := os.Open("testdata/test-feed-minimum.xml")
	if err != nil {
		t.Fatal(err)
	}

	feed, err := parser.ParseFeed(f)
	if err != nil {
		t.Fatal(err)
	}

	assertNotNil(t, feed.Channel)
	assertStr(t, "Test podcast 1", feed.Channel.Title)
	assertStr(t, "http://www.example.com/podcast-site", feed.Channel.Link)
	assertStr(t, "en", feed.Channel.Language)
	assertStr(t, "Test podcast description goes here", feed.Channel.Description.Text)
	// TODO fix XML parsing of fields with namespace
	// assertBool(t, true, feed.Channel.ITunesExplicit)
	// assertStr(t, "http://www.example.com/image.jpg", feed.Channel.ITunesImage.Href)
	// assertInt(t, 1, len(feed.Channel.ITunesCategory))
	// assertStr(t, "Comedy", feed.Channel.ITunesCategory[0].Text)

	assertInt(t, 2, len(feed.Channel.Items))
	assertStr(t, "Test episode 1", feed.Channel.Items[0].Title)
	assertStr(t, "http://www.example.com/episode-1.mp3", feed.Channel.Items[0].Enclosure.URL)
	assertStr(t, "audio/mpeg", feed.Channel.Items[0].Enclosure.Type)
	assertInt(t, 1001, int(feed.Channel.Items[0].Enclosure.Length))
	assertStr(t, "12345-67890-abcdef", feed.Channel.Items[0].GUID.Text)
}

func TestWriteFeed_RequiredFieldsOnly(t *testing.T) {
	feed := gopodcast.Feed{
		Channel: &gopodcast.Channel{
			AtomLink: gopodcast.AtomLink{
				Href: "http://www.example.com/feed",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Title:          "Test title",
			Description:    gopodcast.Description{Text: "Test description"},
			Link:           "http://www.example.com/podcast-site",
			Language:       "fr",
			ITunesExplicit: true,
			ITunesImage: gopodcast.ITunesImage{
				Href: "http://www.example.com/image.png",
			},
			ITunesCategory: []gopodcast.ITunesCategory{{Text: "Drama"}},
			Items: []*gopodcast.Item{
				{
					Title: "A podcast 1",
					Enclosure: gopodcast.Enclosure{
						URL:    "http://www.example.com/pod1.mp3",
						Type:   "audio/mpeg",
						Length: 2001,
					},
					GUID: gopodcast.ItemGUID{
						Text: "abcdef-123456",
					},
				},
				{
					Title: "A podcast 2",
					Enclosure: gopodcast.Enclosure{
						URL:    "http://www.example.com/pod2.mp3",
						Type:   "audio/mpeg",
						Length: 2002,
					},
					GUID: gopodcast.ItemGUID{
						Text: "abcdef-223456",
					},
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	feed.WriteFeedXML(buf)

	exp, err := os.ReadFile("testdata/test-feed-write-minimum.xml")
	if err != nil {
		t.Fatal(err)
	}

	assertStr(
		t,
		strings.TrimSpace(string(exp)),
		strings.TrimSpace(buf.String()),
	)
}

func TestWriteFeed_AllFields(t *testing.T) {
	feed := gopodcast.Feed{
		Channel: &gopodcast.Channel{
			AtomLink: gopodcast.AtomLink{
				Href: "http://www.example.com/feed",
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Title: "Test title",
			Description: gopodcast.Description{
				Text: "Test description",
			},
			Link:           "http://www.example.com/podcast-site",
			Language:       "fr",
			ITunesExplicit: true,
			ITunesImage: gopodcast.ITunesImage{
				Href: "http://www.example.com/image.png",
			},
			ITunesCategory: []gopodcast.ITunesCategory{
				{Text: "Drama", SubCategoryText: "Thriller"},
				{Text: "Comedy"},
			},
			PodcastLocked: "yes",
			PodcastGUID:   "podcast-123-abc",
			ITunesAuthor:  "Mr Author",
			Copyright:     "Mr Author's Boss",
			PodcastText: &gopodcast.PodcastText{
				Purpose: "validation",
				Text:    "text test",
			},
			PodcastFunding: &gopodcast.PodcastFunding{
				URL:  "http://www.example.com/funding",
				Text: "Money please",
			},
			ITunesType:     "episodic",
			ITunesComplete: "yes",
			Items: []*gopodcast.Item{
				{
					Title: "A podcast 1",
					Enclosure: gopodcast.Enclosure{
						URL:    "http://www.example.com/pod1.mp3",
						Type:   "audio/mpeg",
						Length: 2001,
					},
					GUID: gopodcast.ItemGUID{
						IsPermaLink: boolPtr(false),
						Text:        "abcdef-123456",
					},
					Link:    "http://www.example.com/ep-link",
					PubDate: "test date",
					Description: &gopodcast.Description{
						Text: "Test episode description",
					},
					ITunesDuration: "12345",
					ITunesImage: &gopodcast.ITunesImage{
						Href: "http://www.example.com/ep-image.jpg",
					},
					ITunesExplicit: boolPtr(true),
					PodcastTranscript: []gopodcast.PodcastTranscript{
						{
							URL:      "http://www.example.com/ep/trans.fr.txt",
							Type:     "text/plain",
							Rel:      "something",
							Language: "fr",
						},
						{
							URL:      "http://www.example.com/ep/trans.en.txt",
							Type:     "text/plain",
							Rel:      "something",
							Language: "en",
						},
					},
					ITunesEpisode:     "1",
					ITunesSeason:      "2",
					ITunesEpisodeType: "long",
					ITunesBlock:       "no",
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	feed.WriteFeedXML(buf)

	exp, err := os.ReadFile("testdata/test-feed-write-all.xml")
	if err != nil {
		t.Fatal(err)
	}

	assertStr(
		t,
		strings.TrimSpace(string(exp)),
		strings.TrimSpace(buf.String()),
	)
}

func boolPtr(b bool) *bool {
	return &b
}

// aim is for this library to have no dependencies, hence the assert funcs here
func assertNotNil(t *testing.T, act any) {
	if act == nil {
		t.Fatalf("expected '%+v' to not be nil", act)
	}
}

func assertStr(t *testing.T, exp string, act string) {
	if exp != act {
		t.Fatalf("expected '%s', got '%s'", exp, act)
	}
}

func assertBool(t *testing.T, exp bool, act bool) {
	if exp != act {
		t.Fatalf("expected '%t', got '%t'", exp, act)
	}
}

func assertInt(t *testing.T, exp int, act int) {
	if act != exp {
		t.Fatalf("expected '%d', got '%d'", exp, act)
	}
}
