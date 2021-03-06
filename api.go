package politico

import (
	// "html"
	"net/url"

	"strings"
	"time"
)

// TODO: URL, tags

type politicoTime struct {
	time.Time
}

func (t *politicoTime) UnmarshalJSON(buf []byte) error {
	if string(buf) == "\"\"" {
		return nil
	}
	buf = buf[:len(buf)-4]
	/* Get rid of the "ET". This isn't a real timezone, and golang doesn't
	 * know how to turn that into EST. Frickn'. Anyway. I'm going to go against
	 * my own advice in https://notes.pault.ag/its-all-relative/ and
	 * ignore timezones. I guess. SOMEONE PLEASE FIX ME */
	tt, err := time.Parse("01/02/2006 03:04:05 PM", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

type Response struct {
	Stories Stories `json:"stories"`
}

type Stories struct {
	Header        string  `json:"header"`
	SectionLayout string  `json:"section_layout"`
	Stories       []Story `json:"story"`
}

type escapedString string

func (t *escapedString) UnmarshalJSON(buf []byte) error {
	s, err := url.QueryUnescape(string(buf[1 : len(buf)-1]))
	if err != nil {
		return err
	}
	*t = escapedString(s)
	return nil
}

type csvString []string

func (t *csvString) UnmarshalJSON(buf []byte) error {
	*t = strings.Split(strings.Trim(string(buf), `"`), ",")
	return nil
}

type Story struct {
	GUID      escapedString `json:"guid"`
	Permalink escapedString `json:"permalink"`
	StoryType escapedString `json:"story_type"`
	Tags      csvString     `json:"tags"`

	AuthorBioHTML string `json:"author_bio_html"`
	By            escapedString
	BylineHTML    string `json:"byline_html"`

	Date    politicoTime `json:"date"`
	Updated politicoTime `json:"updated"`

	Title escapedString `json:"title"`
	Dek   escapedString `json:"dek"`

	HTML string `json:"html"`

	Links []struct {
		Length   int    `json:"len"`
		Position int    `json:"posn"`
		URL      string `json:"url"`
	}

	Media struct {
		Caption    escapedString `json:"caption"`
		Credit     escapedString `json:"credit"`
		URL        escapedString `json:"url"`
		URLSmall   escapedString `json:"url_small"`
		Renditions []struct {
			DisplayName   escapedString `json:"displayName"`
			EncodingRate  escapedString `json:"encodingRate"`
			FrameHeight   escapedString `json:"frameHeight"`
			FrameWidth    escapedString `json:"frameWidth"`
			URL           escapedString `json:"url"`
			VideoDuration escapedString `json:"videoDuration"`
		} `json:"renditions"`
	}

	Related []struct {
		Text            escapedString `json:"text"`
		Type            escapedString `json:"type"`
		URL             escapedString `json:"url"`
		SectionOfOrigin escapedString `json:"section_of_origin"`
	}
}

func (s Story) DekOrTitle() escapedString {
	if s.Dek != "" {
		return s.Dek
	}
	return s.Title
}
