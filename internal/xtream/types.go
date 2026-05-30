package xtream

import (
	"encoding/json"
	"fmt"
)

// Category is shared across live/vod/series — only `type` differs.
type Category struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	ParentID     int    `json:"parent_id"`
}

// FlexibleString unmarshals both JSON strings and numbers into a string.
// This handles providers that inconsistently return rating as "8.5" vs 8.5.
type FlexibleString string

func (fs *FlexibleString) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch val := v.(type) {
	case string:
		*fs = FlexibleString(val)
	case float64:
		*fs = FlexibleString(fmt.Sprintf("%v", val))
	case nil:
		*fs = ""
	default:
		return fmt.Errorf("cannot unmarshal %T into FlexibleString", v)
	}
	return nil
}

// FlexibleStringArray unmarshals both JSON strings and arrays into []string.
// This handles providers that inconsistently return backdrop_path as a string or array.
type FlexibleStringArray []string

func (fsa *FlexibleStringArray) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch val := v.(type) {
	case string:
		if val != "" {
			*fsa = FlexibleStringArray{val}
		} else {
			*fsa = FlexibleStringArray{}
		}
	case []interface{}:
		arr := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				arr = append(arr, s)
			}
		}
		*fsa = FlexibleStringArray(arr)
	case nil:
		*fsa = FlexibleStringArray{}
	default:
		return fmt.Errorf("cannot unmarshal %T into FlexibleStringArray", v)
	}
	return nil
}

// VOD represents a movie entry from get_vod_streams.
type VOD struct {
	StreamID           int            `json:"stream_id"`
	Name               string         `json:"name"`
	StreamIcon         string         `json:"stream_icon"`
	Rating             FlexibleString `json:"rating"`
	Year               string         `json:"year"`
	Added              string         `json:"added"`
	CategoryID         string         `json:"category_id"`
	ContainerExtension string         `json:"container_extension"`
	Plot               string         `json:"plot"`
}

// SeriesListing is what get_series returns: the show, not its seasons.
type SeriesListing struct {
	SeriesID    int                 `json:"series_id"`
	Name        string              `json:"name"`
	Cover       string              `json:"cover"`
	Plot        string              `json:"plot"`
	ReleaseDate string              `json:"releaseDate"`
	CategoryID  string              `json:"category_id"`
	Backdrop    FlexibleStringArray `json:"backdrop_path"`
}

// SeriesInfo is the get_series_info response: seasons + episodes.
type SeriesInfo struct {
	Info struct {
		Name     string              `json:"name"`
		Cover    string              `json:"cover"`
		Plot     string              `json:"plot"`
		Backdrop FlexibleStringArray `json:"backdrop_path"`
	} `json:"info"`
	Seasons []struct {
		SeasonNumber int    `json:"season_number"`
		Name         string `json:"name"`
		Overview     string `json:"overview"`
		Cover        string `json:"cover"`
	} `json:"seasons"`
	// Episodes is keyed by season number as a string.
	Episodes map[string][]Episode `json:"episodes"`
}

// Episode is one entry in SeriesInfo.Episodes[season].
type Episode struct {
	ID                 string `json:"id"`
	EpisodeNum         int    `json:"episode_num"`
	Title              string `json:"title"`
	ContainerExtension string `json:"container_extension"`
	Info               struct {
		Plot     string `json:"plot"`
		Duration string `json:"duration"`
	} `json:"info"`
	Season int `json:"season"`
}
