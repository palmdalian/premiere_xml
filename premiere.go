package premiere

import (
	"encoding/xml"
	"net/url"
	"strings"
)

const (
	PProTicksConstant = 254016000000
)

type PremiereXML struct {
	XMLName  xml.Name  `xml:"xmeml"`
	Sequence *Sequence `xml:"sequence"`
}

type Sequence struct {
	UUID     string    `xml:"uuid"`
	Duration int64     `xml:"duration"`
	Rate     *Rate     `xml:"rate"`
	Name     string    `xml:"name"`
	Media    *Media    `xml:"media"`
	Timecode *Timecode `xml:"timecode"`
	Markers  []*Marker `xml:"marker"`
	Labels   []*Label  `xml:"labels"`
}

type Rate struct {
	Timebase int64 `xml:"timebase"`
	NTSC     bool  `xml:"ntsc"`
}

type Timecode struct {
	Rate   *Rate  `xml:"rate"`
	String string `xml:"string"`
	Frame  int64  `xml:"frame"`
}

type Marker struct {
	Comment string `xml:"comment"`
	Name    string `xml:"name"`
	In      int64  `xml:"in"`
	Out     int64  `xml:"out"`
}

type Label struct {
	Color string `xml:"label2"`
}

type Media struct {
	Video *Video `xml:"video"`
	Audio *Audio `xml:"audio"`
}

type Video struct {
	Format *Format  `xml:"format"`
	Tracks []*Track `xml:"track"`
}

type Audio struct {
	Format           *Format  `xml:"format"`
	Tracks           []*Track `xml:"track"`
	NumOutputChannes int      `xml:"numOutputChannels"`
}

type Format struct {
	SampleCharacteristics *SampleCharacteristics `xml:"samplecharacteristics"`
}

type SampleCharacteristics struct {
	Depth      int64 `xml:"depth"`
	SampleRate int64 `xml:"samplerate"`
}

type Track struct {
	ClipItems          []*ClipItem       `xml:"clipitem"`
	TransitionItems    []*TransitionItem `xml:"transitionitem"`
	Enabled            bool              `xml:"enabled"`
	Locked             bool              `xml:"locked"`
	OutputChannelIndex int               `xml:"outputchannelindex"`
}

type ClipItem struct {
	Id             string       `xml:"id,attr"`
	MasterClipId   string       `xml:"masterclipid"`
	Name           string       `xml:"name"`
	Enabled        bool         `xml:"enabled"`
	Duration       int64        `xml:"duration"`
	Rate           *Rate        `xml:"rate"`
	Start          int64        `xml:"start"`        // In point within the sequence
	End            int64        `xml:"end"`          // Out point within the sequence
	In             int64        `xml:"in"`           // In point within the media file
	Out            int64        `xml:"out"`          // Out point within the media file
	PProTicksIn    int64        `xml:"pproTicksIn"`  // Premiere specific in point (use const to get seconds)
	PProTicksInOut int64        `xml:"pproTicksOut"` // Premiere specific out point (use const to get seconds)
	Label          *Label       `xml:"labels"`
	File           *File        `xml:"file"`
	LoggingInfo    *LoggingInfo `xml:"logginginfo"`
	Links          []*Link      `xml:"link"`
	SourceTrack    *SourceTrack `xml:"sourcetrack"`
}

type File struct {
	Id       string    `xml:"id,attr"`
	Name     string    `xml:"name"`
	PathUrl  string    `xml:"pathurl"`
	Rate     *Rate     `xml:"rate"`
	Duration int64     `xml:"duration"`
	Timecode *Timecode `xml:"timecode"`
	Media    *Media    `xml:"media"`
}

type TransitionItem struct {
	Start         int64  `xml:"start"`
	End           int64  `xml:"end"`
	Alignment     string `xml:"alignment"`
	CutPointTicks int64  `xml:"cutPointTicks"`
	Rate          *Rate  `xml:"rate"`
}

type LoggingInfo struct {
	Description           string `xml:"description"`
	Scene                 string `xml:"scene"`
	ShotTake              string `xml:"shottake"`
	LogNote               string `xml:"lognote"`
	Good                  string `xml:"good"`
	OriginalVideoFilename string `xml:"originalvideofilename"`
	OriginalAudioFilename string `xml:"originalaudiofilename"`
}

type SourceTrack struct {
	MediaType  string `xml:"sourcetrack"`
	TrackIndex int    `xml:"trackindex"`
}

type Link struct {
	LinkClipRef string `xml:"linkclipref"`
	MediaType   string `xml:"mediatype"`
	TrackIndex  int    `xml:"trackindex"`
	ClipIndex   int    `xml:"clipindex"`
	GroupIndex  int    `xml:"groupindex"`
}

func (seq *Sequence) AllFilePaths() (fileMap map[string]string) {
	fileMap = make(map[string]string)
	for _, track := range seq.Media.Video.Tracks {
		for _, clip := range track.ClipItems {
			if clip.File.PathUrl != "" {
				path, _ := url.PathUnescape(clip.File.PathUrl)
				trimmedPath := strings.Trim(path, "file://localhost")
				fileMap[clip.File.Id] = trimmedPath
			}
		}
	}
	for _, track := range seq.Media.Audio.Tracks {
		for _, clip := range track.ClipItems {
			if clip.File.PathUrl != "" {
				path, _ := url.PathUnescape(clip.File.PathUrl)
				trimmedPath := strings.Trim(path, "file://localhost")
				fileMap[clip.File.Id] = trimmedPath
			}
		}
	}
	return
}
