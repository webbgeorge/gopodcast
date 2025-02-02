package gopodcast_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
	"time"

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
	assertBool(t, true, bool(podcast.ITunesExplicit))
	assertStr(t, "http://www.example.com/image.jpg", podcast.ITunesImage.Href)
	assertInt(t, 1, len(podcast.ITunesCategory))
	assertStr(t, "Comedy", podcast.ITunesCategory[0].Text)

	// non-required channel fields should be zero values
	assertNil(t, podcast.ITunesCategory[0].SubCategory)
	assertNil(t, podcast.PodcastLocked)
	assertStr(t, "", podcast.PodcastGUID)
	assertStr(t, "", podcast.ITunesAuthor)
	assertStr(t, "", podcast.Copyright)
	assertNil(t, podcast.PodcastText)
	assertNil(t, podcast.PodcastFunding)
	assertStr(t, "", podcast.ITunesType)
	assertNil(t, podcast.ITunesComplete)

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
	assertNil(t, item.PubDate)
	assertNil(t, item.Description)
	assertStr(t, "", item.ITunesDuration)
	assertNil(t, item.ITunesImage)
	assertNil(t, item.ITunesExplicit)
	assertInt(t, 0, len(item.PodcastTranscript))
	assertStr(t, "", item.ITunesEpisode)
	assertStr(t, "", item.ITunesSeason)
	assertNil(t, item.ITunesBlock)
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
	assertBool(t, true, bool(podcast.ITunesExplicit))
	assertStr(t, "http://www.example.com/image.jpg", podcast.ITunesImage.Href)
	assertInt(t, 2, len(podcast.ITunesCategory))
	assertStr(t, "Comedy", podcast.ITunesCategory[0].Text)
	assertStr(t, "Drama", podcast.ITunesCategory[1].Text)
	assertStr(t, "Thriller", podcast.ITunesCategory[1].SubCategory.Text)
	assertBool(t, true, bool(*podcast.PodcastLocked))
	assertStr(t, "podcast-123456", podcast.PodcastGUID)
	assertStr(t, "Dr Tester", podcast.ITunesAuthor)
	assertStr(t, "Tester Inc.", podcast.Copyright)
	assertStr(t, "abcdef", podcast.PodcastText.Text)
	assertStr(t, "validation", podcast.PodcastText.Purpose)
	assertStr(t, "Money please", podcast.PodcastFunding.Text)
	assertStr(t, "http://www.example.com/money", podcast.PodcastFunding.URL)
	assertStr(t, "Serialised", podcast.ITunesType)
	assertBool(t, true, bool(*podcast.ITunesComplete))

	// item fields
	assertInt(t, 1, len(podcast.Items))
	item := podcast.Items[0]
	assertStr(t, "Test episode 1", item.Title)
	assertStr(t, "http://www.example.com/episode-1.mp3", item.Enclosure.URL)
	assertStr(t, "audio/mpeg", item.Enclosure.Type)
	assertInt(t, 1001, int(item.Enclosure.Length))
	assertStr(t, "12345-67890-abcdef", item.GUID.Text)
	assertStr(t, "http://www.example.com/ep-link", item.Link)
	assertStr(t, "2024-12-26T11:12:13Z", time.Time(*item.PubDate).Format(time.RFC3339))
	assertStr(t, "Episode test description", item.Description.Text)
	assertStr(t, "1234", item.ITunesDuration)
	assertStr(t, "http://www.example.com/ep-image.png", item.ITunesImage.Href)
	assertBool(t, false, bool(*item.ITunesExplicit))
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
	assertBool(t, false, bool(*item.ITunesBlock))
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
		PodcastLocked: yesNoPtr(true),
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
		ITunesComplete: yesNoPtr(true),
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
				PubDate: timeFromStr("2024-12-25T10:11:12"),
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
				ITunesBlock:       yesNoPtr(false),
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

func TestParseFeedFromURL(t *testing.T) {
	testFeedURL := "https://feeds.captivate.fm/elis-james-and-john-robins/"

	parser := gopodcast.NewParser()
	feed, err := parser.ParseFeedFromURL(context.Background(), testFeedURL)
	if err != nil {
		t.Fatal(err)
	}

	checkRequiredFeedValuesPresent(t, feed)
}

func TestParseFeedFromURL_Non200(t *testing.T) {
	testFeedURL := "http://www.example.com/feed"

	parser := gopodcast.NewParser()
	parser.HTTPClient = newTestClient(500, "error")

	feed, err := parser.ParseFeedFromURL(context.Background(), testFeedURL)

	assertNil(t, feed)
	assertStr(t, "non-200 http response '500'", err.Error())
}

func TestParseFeedFromURL_InvalidFeed(t *testing.T) {
	testFeedURL := "http://www.example.com/feed"

	parser := gopodcast.NewParser()
	parser.HTTPClient = newTestClient(200, "NOT a valid feed")

	feed, err := parser.ParseFeedFromURL(context.Background(), testFeedURL)

	assertNil(t, feed)
	assertNotNil(t, err)
}

func TestParseFeedFromURL_SendsAuthCreds(t *testing.T) {
	interceptTransport := &interceptAuthTransport{
		transport: http.DefaultTransport,
	}

	interceptClient := &http.Client{
		Transport: interceptTransport,
	}

	testFeedURL := "http://www.example.com/feed"

	parser := gopodcast.NewParser()
	parser.HTTPClient = interceptClient
	parser.AuthCredentials = &gopodcast.AuthCredentials{
		Username: "user1",
		Password: "password1",
	}

	_, _ = parser.ParseFeedFromURL(context.Background(), testFeedURL)

	// basic auth: base64(user:pass)
	assertStr(t, "Basic dXNlcjE6cGFzc3dvcmQx", interceptTransport.authHeader)
}

// TestParseFeed_TopPodcasts tests the parser against many different real podcasts,
// taken from the Apple charts.
func TestParseFeed_TopPodcasts(t *testing.T) {
	files, err := os.ReadDir("testdata/top-podcasts")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if !file.Type().IsRegular() {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			parser := gopodcast.NewParser()
			f, err := os.Open(path.Join("testdata/top-podcasts", file.Name()))
			if err != nil {
				t.Fatal(err)
			}
			podcast, err := parser.ParseFeed(f)
			if err != nil {
				t.Fatal(err)
			}
			checkRequiredFeedValuesPresent(t, podcast)
		})
	}
}

// checkRequiredFeedValuesPresent does some simple checks to make sure key
// fields are present in a podcast feed. This is used for running the parser
// tests across a large number of real podcast feeds.
//
// note that some podcast feeds don't include channel->link, enclosure->length,
// or channel->atom:link despite these being required by the PSP spec, so we
// don't test for them here.
func checkRequiredFeedValuesPresent(t *testing.T, podcast *gopodcast.Podcast) {
	// channel fields
	assertNotNil(t, podcast)
	assertStrNotEmpty(t, podcast.Title)
	assertStrNotEmpty(t, podcast.Language)
	assertStrNotEmpty(t, podcast.Description.Text)
	assertStrNotEmpty(t, podcast.ITunesImage.Href)
	assertTrue(t, len(podcast.ITunesCategory) > 0)
	assertStrNotEmpty(t, podcast.ITunesCategory[0].Text)

	assertTrue(t, len(podcast.Items) > 0)
	// use the first episode, as some feeds publish the latest episode ahead
	// of release time, without the audio file.
	item := podcast.Items[len(podcast.Items)-1]

	// item fields
	assertStrNotEmpty(t, item.Title)
	assertStrNotEmpty(t, item.Enclosure.URL)
	assertStrNotEmpty(t, item.Enclosure.Type)
	assertStrNotEmpty(t, item.GUID.Text)
}

func boolPtr(b bool) *gopodcast.Bool {
	bb := gopodcast.Bool(b)
	return &bb
}

func yesNoPtr(b bool) *gopodcast.YesNo {
	bb := gopodcast.YesNo(b)
	return &bb
}

func timeFromStr(str string) *gopodcast.Time {
	t, err := time.Parse("2006-01-02T15:04:05", str)
	if err != nil {
		panic(err)
	}
	tt := gopodcast.Time(t)
	return &tt
}

// aim is for this library to have no dependencies, hence the assert funcs here
func assertTrue(t *testing.T, act bool) {
	t.Helper()
	if !act {
		t.Fatal("expected to be true")
	}
}

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

func assertStrNotEmpty(t *testing.T, act string) {
	t.Helper()
	if act == "" {
		t.Fatal("expected string to not be empty")
	}
}

func newTestClient(httpStatus int, content string) *http.Client {
	return &http.Client{
		Transport: &testTransport{
			httpStatus: httpStatus,
			content:    content,
		},
	}
}

// testTransport returns the given http status code and content to enable testing of http clients
type testTransport struct {
	httpStatus int
	content    string
}

func (t *testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	buf := bytes.NewBufferString(t.content)
	body := io.NopCloser(buf)
	return &http.Response{
		StatusCode: t.httpStatus,
		Body:       body,
	}, nil
}

// interceptAuthTransport captures the value of the Authorization header to be used in tests
type interceptAuthTransport struct {
	transport  http.RoundTripper
	authHeader string
}

func (t *interceptAuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.authHeader = r.Header.Get("Authorization")

	return t.transport.RoundTrip(r)
}
