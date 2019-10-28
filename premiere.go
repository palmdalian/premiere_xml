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
	UUID                                  string    `xml:"uuid"`
	Duration                              int64     `xml:"duration"`
	Rate                                  *Rate     `xml:"rate"`
	Name                                  string    `xml:"name"`
	Media                                 *Media    `xml:"media"`
	Timecode                              *Timecode `xml:"timecode"`
	Markers                               []*Marker `xml:"marker"`
	Labels                                []*Label  `xml:"labels"`
	TLSQAudioVisibleBase                  string    `xml:"TL.SQAudioVisibleBase,attr"`
	TLSQVideoVisibleBase                  string    `xml:"TL.SQVideoVisibleBase,attr"`
	TLSQVisibleBaseTime                   string    `xml:"TL.SQVisibleBaseTime,attr"`
	TLSQAVDividerPosition                 string    `xml:"TL.SQAVDividerPosition,attr"`
	TLSQHideShyTracks                     string    `xml:"TL.SQHideShyTracks,attr"`
	TLSQHeaderWidth                       string    `xml:"TL.SQHeaderWidth,attr"`
	MonitorProgramZoomOut                 string    `xml:"Monitor.ProgramZoomOut,attr"`
	MonitorProgramZoomIn                  string    `xml:"Monitor.ProgramZoomIn,attr"`
	TLSQTimePerPixel                      string    `xml:"TL.SQTimePerPixel,attr"`
	MZEditLine                            string    `xml:"MZ.EditLine,attr"`
	MZSequencePreviewFrameSizeHeight      string    `xml:"MZ.Sequence.PreviewFrameSizeHeight,attr"`
	MZSequencePreviewFrameSizeWidth       string    `xml:"MZ.Sequence.PreviewFrameSizeWidth,attr"`
	MZSequenceAudioTimeDisplayFormat      string    `xml:"MZ.Sequence.AudioTimeDisplayFormat,attr"`
	MZSequencePreviewRenderingClassID     string    `xml:"MZ.Sequence.PreviewRenderingClassID,attr"`
	MZSequencePreviewRenderingPresetCodec string    `xml:"MZ.Sequence.PreviewRenderingPresetCodec,attr"`
	MZSequencePreviewRenderingPresetPath  string    `xml:"MZ.Sequence.PreviewRenderingPresetPath,attr"`
	MZSequencePreviewUseMaxRenderQuality  string    `xml:"MZ.Sequence.PreviewUseMaxRenderQuality,attr"`
	MZSequencePreviewUseMaxBitDepth       string    `xml:"MZ.Sequence.PreviewUseMaxBitDepth,attr"`
	MZSequenceEditingModeGUID             string    `xml:"MZ.Sequence.EditingModeGUID,attr"`
	MZSequenceVideoTimeDisplayFormat      string    `xml:"MZ.Sequence.VideoTimeDisplayFormat,attr"`
	MZWorkOutPoint                        string    `xml:"MZ.WorkOutPoint,attr"`
	MZWorkInPoint                         string    `xml:"MZ.WorkInPoint,attr"`
	MZZeroPoint                           string    `xml:"MZ.ZeroPoint,attr"`
	ExplodedTracks                        string    `xml:"explodedTracks,attr"`
}

type Rate struct {
	Timebase int64 `xml:"timebase"`
	NTSC     bool  `xml:"ntsc"`
}

type Timecode struct {
	Rate          *Rate  `xml:"rate"`
	String        string `xml:"string"`
	Frame         int64  `xml:"frame"`
	DisplayFormat string `xml:"displayformat"`
	Reel          *Reel  `xml:"reel"`
}
type Reel struct {
	Name string `xml:"name"`
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
	Format            *Format  `xml:"format"`
	Tracks            []*Track `xml:"track"`
	NumOutputChannels int      `xml:"numOutputChannels"`
	Outputs           *Outputs `xml:"outputs"`
}

type Outputs struct {
	Groups []*Group `xml:"group"`
}

type Group struct {
	Index       int      `xml:"index"`
	NumChannels int      `xml:"numchannels"`
	Downmix     int      `xml:"downmix"`
	Channel     *Channel `xml:"channel"`
}
type Channel struct {
	Index int `xml:"index"`
}

type Format struct {
	SampleCharacteristics *SampleCharacteristics `xml:"samplecharacteristics"`
}

type SampleCharacteristics struct {
	Depth            int64  `xml:"depth"`
	SampleRate       int64  `xml:"samplerate"`
	Codec            *Codec `xml:"codec"`
	Width            int    `xml:"width"`
	Height           int    `xml:"height"`
	Anamorphic       bool   `xml:"anamorphic"`
	PixelAspectRatio string `xml:"pixelaspectratio"`
	FieldDominance   string `xml:"fielddominance"`
	ColorDepth       int    `xml:"colordepth"`
}

type Codec struct {
	Name            string           `xml:"name"`
	AppSpecificData *AppSpecificData `xml:"appspecificdata"`
}

type AppSpecificData struct {
	AppName         string   `xml:"appname"`
	AppmMnufacturer string   `xml:"appmanufacturer"`
	AppVersion      string   `xml:"appversion"`
	Data            *AppData `xml:"data"`
}

type AppData struct {
	QTCodec *QTCodec `xml:"qtcodec"`
}

type QTCodec struct {
	CodecName       string `xml:"codecname"`
	CodecTypeName   string `xml:"codectypename"`
	CodecTypeCode   string `xml:"codectypecode"`
	CodecVendorCode string `xml:"codecvendorcode"`
	SpacialQuality  int    `xml:"spatialquality"`
	TemporalQuality int    `xml:"temporalquality"`
	KeyframeRate    int    `xml:"keyframerate"`
	DataRate        int    `xml:"datarate"`
}

type Track struct {
	ClipItems          []*ClipItem       `xml:"clipitem"`
	TransitionItems    []*TransitionItem `xml:"transitionitem"`
	Enabled            bool              `xml:"enabled"`
	Locked             bool              `xml:"locked"`
	OutputChannelIndex int               `xml:"outputchannelindex"`
}

type ClipItem struct {
	Id               string       `xml:"id,attr"`
	MasterClipId     string       `xml:"masterclipid"`
	Name             string       `xml:"name"`
	Enabled          bool         `xml:"enabled"`
	Duration         int64        `xml:"duration"`
	Rate             *Rate        `xml:"rate"`
	AlphaType        string       `xml:"alphatype"`
	Anamorphic       bool         `xml:"anamorphic"`
	PixelAspectRatio string       `xml:"pixelaspectratio"`
	Start            int64        `xml:"start"`        // In point within the sequence
	End              int64        `xml:"end"`          // Out point within the sequence
	In               int64        `xml:"in"`           // In point within the media file
	Out              int64        `xml:"out"`          // Out point within the media file
	PProTicksIn      int64        `xml:"pproTicksIn"`  // Premiere specific in point (use const to get seconds)
	PProTicksInOut   int64        `xml:"pproTicksOut"` // Premiere specific out point (use const to get seconds)
	Label            *Label       `xml:"labels"`
	File             *File        `xml:"file"`
	LoggingInfo      *LoggingInfo `xml:"logginginfo"`
	Links            []*Link      `xml:"link"`
	SourceTrack      *SourceTrack `xml:"sourcetrack"`
	Sequence         *Sequence    `xml:"sequence"`
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
