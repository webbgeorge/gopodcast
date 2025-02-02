package gopodcast_test

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/webbgeorge/gopodcast"
)

func TestParseFeed_RequiredFieldsOnly(t *testing.T) {
	parser := gopodcast.NewParser()

	f, err := os.Open("testdata/test-feed-minimum.xml")
	if err != nil {
		t.Fatal(err)
	}

	podcast, err := parser.ParseFeed(f)
	if err != nil {
		t.Fatal(err)
	}

	// channel fields
	assertNotNil(t, podcast)
	assertStr(t, "http://www.example.com/feed", podcast.AtomLink.Href)
	assertStr(t, "self", podcast.AtomLink.Rel)
	assertStr(t, "application/rss+xml", podcast.AtomLink.Type)
	assertStr(t, "Test podcast 1", podcast.Title)
	assertStr(t, "http://www.example.com/podcast-site", podcast.Link)
	assertStr(t, "en", podcast.Language)
	assertStr(t, "Test podcast description goes here", podcast.Description.Text)
	assertBool(t, true, podcast.ITunesExplicit)
	assertStr(t, "http://www.example.com/image.jpg", podcast.ITunesImage.Href)
	assertInt(t, 1, len(podcast.ITunesCategory))
	assertStr(t, "Comedy", podcast.ITunesCategory[0].Text)

	// non-required channel fields should be zero values
	assertNil(t, podcast.ITunesCategory[0].SubCategory)
	assertStr(t, "", podcast.PodcastLocked)
	assertStr(t, "", podcast.PodcastGUID)
	assertStr(t, "", podcast.ITunesAuthor)
	assertStr(t, "", podcast.Copyright)
	assertNil(t, podcast.PodcastText)
	assertNil(t, podcast.PodcastFunding)
	assertStr(t, "", podcast.ITunesType)
	assertStr(t, "", podcast.ITunesComplete)

	// item fields
	assertInt(t, 2, len(podcast.Items))
	item := podcast.Items[0]
	item2 := podcast.Items[1]
	assertStr(t, "Test episode 1", item.Title)
	assertStr(t, "Test episode 2", item2.Title)
	assertStr(t, "http://www.example.com/episode-1.mp3", item.Enclosure.URL)
	assertStr(t, "audio/mpeg", item.Enclosure.Type)
	assertInt(t, 1001, int(item.Enclosure.Length))
	assertStr(t, "12345-67890-abcdef", item.GUID.Text)

	// non-required item fields should be zero values
	assertStr(t, "", item.Link)
	assertStr(t, "", item.PubDate)
	assertNil(t, item.Description)
	assertStr(t, "", item.ITunesDuration)
	assertNil(t, item.ITunesImage)
	assertNil(t, item.ITunesExplicit)
	assertInt(t, 0, len(item.PodcastTranscript))
	assertStr(t, "", item.ITunesEpisode)
	assertStr(t, "", item.ITunesSeason)
	assertStr(t, "", item.ITunesBlock)
}

func TestParseFeed_AllFields(t *testing.T) {
	parser := gopodcast.NewParser()

	f, err := os.Open("testdata/test-feed-all.xml")
	if err != nil {
		t.Fatal(err)
	}

	podcast, err := parser.ParseFeed(f)
	if err != nil {
		t.Fatal(err)
	}

	// channel fields
	assertNotNil(t, podcast)
	assertStr(t, "http://www.example.com/feed", podcast.AtomLink.Href)
	assertStr(t, "self", podcast.AtomLink.Rel)
	assertStr(t, "application/rss+xml", podcast.AtomLink.Type)
	assertStr(t, "Test podcast 1", podcast.Title)
	assertStr(t, "http://www.example.com/podcast-site", podcast.Link)
	assertStr(t, "en", podcast.Language)
	assertStr(t, "Test podcast description goes here", podcast.Description.Text)
	assertBool(t, true, podcast.ITunesExplicit)
	assertStr(t, "http://www.example.com/image.jpg", podcast.ITunesImage.Href)
	assertInt(t, 2, len(podcast.ITunesCategory))
	assertStr(t, "Comedy", podcast.ITunesCategory[0].Text)
	assertStr(t, "Drama", podcast.ITunesCategory[1].Text)
	assertStr(t, "Thriller", podcast.ITunesCategory[1].SubCategory.Text)
	assertStr(t, "yes", podcast.PodcastLocked)
	assertStr(t, "podcast-123456", podcast.PodcastGUID)
	assertStr(t, "Dr Tester", podcast.ITunesAuthor)
	assertStr(t, "Tester Inc.", podcast.Copyright)
	assertStr(t, "abcdef", podcast.PodcastText.Text)
	assertStr(t, "validation", podcast.PodcastText.Purpose)
	assertStr(t, "Money please", podcast.PodcastFunding.Text)
	assertStr(t, "http://www.example.com/money", podcast.PodcastFunding.URL)
	assertStr(t, "Serialised", podcast.ITunesType)
	assertStr(t, "yes", podcast.ITunesComplete)

	// item fields
	assertInt(t, 1, len(podcast.Items))
	item := podcast.Items[0]
	assertStr(t, "Test episode 1", item.Title)
	assertStr(t, "http://www.example.com/episode-1.mp3", item.Enclosure.URL)
	assertStr(t, "audio/mpeg", item.Enclosure.Type)
	assertInt(t, 1001, int(item.Enclosure.Length))
	assertStr(t, "12345-67890-abcdef", item.GUID.Text)
	assertStr(t, "http://www.example.com/ep-link", item.Link)
	assertStr(t, "someDate", item.PubDate)
	assertStr(t, "Episode test description", item.Description.Text)
	assertStr(t, "1234", item.ITunesDuration)
	assertStr(t, "http://www.example.com/ep-image.png", item.ITunesImage.Href)
	assertBool(t, false, *item.ITunesExplicit)
	assertInt(t, 2, len(item.PodcastTranscript))
	assertStr(t, "http://www.example.com/transcript-1-en.txt", item.PodcastTranscript[0].URL)
	assertStr(t, "text/plain", item.PodcastTranscript[0].Type)
	assertStr(t, "self", item.PodcastTranscript[0].Rel)
	assertStr(t, "en", item.PodcastTranscript[0].Language)
	assertStr(t, "http://www.example.com/transcript-1-fr.txt", item.PodcastTranscript[1].URL)
	assertStr(t, "text/plain", item.PodcastTranscript[1].Type)
	assertStr(t, "self", item.PodcastTranscript[1].Rel)
	assertStr(t, "fr", item.PodcastTranscript[1].Language)
	assertStr(t, "1", item.ITunesEpisode)
	assertStr(t, "2", item.ITunesSeason)
	assertStr(t, "no", item.ITunesBlock)
}

func TestWriteFeed_RequiredFieldsOnly(t *testing.T) {
	podcast := &gopodcast.Podcast{
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
	}

	buf := &bytes.Buffer{}
	podcast.WriteFeedXML(buf)

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
	podcast := &gopodcast.Podcast{
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
			{
				Text:        "Drama",
				SubCategory: &gopodcast.ITunesCategory{Text: "Thriller"},
			},
			{
				Text: "Comedy",
			},
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
	}

	buf := &bytes.Buffer{}
	podcast.WriteFeedXML(buf)

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
	t.Helper()
	if reflect.ValueOf(act).IsNil() {
		t.Fatalf("expected '%+v' to not be nil", act)
	}
}

func assertNil(t *testing.T, act any) {
	t.Helper()
	if !reflect.ValueOf(act).IsNil() {
		t.Fatalf("expected '%+v' to be nil", act)
	}
}

func assertStr(t *testing.T, exp string, act string) {
	t.Helper()
	if exp != act {
		t.Fatalf("expected '%s', got '%s'", exp, act)
	}
}

func assertBool(t *testing.T, exp bool, act bool) {
	t.Helper()
	if exp != act {
		t.Fatalf("expected '%t', got '%t'", exp, act)
	}
}

func assertInt(t *testing.T, exp int, act int) {
	t.Helper()
	if act != exp {
		t.Fatalf("expected '%d', got '%d'", exp, act)
	}
}
