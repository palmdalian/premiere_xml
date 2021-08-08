package builder

import (
	"encoding/json"
	"math"
	"os/exec"
	"strconv"
	"strings"

	premiere "github.com/palmdalian/premiere_xml"
)

func videoFromStreams(streams []*ffprobeStream) *premiere.Video {
	var stream *ffprobeStream
	for _, s := range streams {
		if s.CodecType == "video" {
			stream = s
			break
		}
	}

	if stream == nil {
		return nil
	}

	rate := int64(0)
	ntsc := true

	split := strings.Split(stream.RFrameRate, "/")
	if len(split) != 2 {
		return nil
	}
	num, err := strconv.Atoi(split[0])
	if err != nil {
		return nil
	}
	den, err := strconv.Atoi(split[1])
	if err != nil {
		return nil
	}
	if num == 0 || den == 0 {
		return nil
	}
	// TODO test if this is correct.
	if num%den == 0 {
		ntsc = false
	}
	rate = int64(math.Round(float64(num) / float64(den)))

	video := &premiere.Video{
		SampleCharacteristics: &premiere.SampleCharacteristics{
			Rate: &premiere.Rate{
				Timebase: rate,
				NTSC:     ntsc,
			},
			Width:            stream.Width,
			Height:           stream.Height,
			Anamorphic:       false,
			PixelAspectRatio: "square",
			FieldDominance:   "none",
		},
		Tracks: []*premiere.Track{},
	}
	return video
}

func audioFromStreams(streams []*ffprobeStream) *premiere.Audio {
	var stream *ffprobeStream
	for _, s := range streams {
		if s.CodecType == "audio" {
			stream = s
			break
		}
	}

	if stream == nil {
		return nil
	}

	audio := &premiere.Audio{
		SampleCharacteristics: &premiere.SampleCharacteristics{
			SampleRate: stream.SampleRate,
			Depth:      16, // TODO figure out how to extract this
		},
		Tracks:       []*premiere.Track{},
		ChannelCount: stream.Channels,
	}
	return audio
}

type ffprobeOutput struct {
	Streams []*ffprobeStream `json:"streams"`
	Format  *ffprobeFormat   `json:"format"`
}

func getFFProbeOutput(filePath string) (*ffprobeOutput, error) {
	args := []string{"-v", "quiet", "-i", filePath,
		"-print_format", "json", "-show_format", "-show_streams"}

	outBytes, err := exec.Command("ffprobe", args...).Output()
	if err != nil {
		return nil, err
	}
	output := &ffprobeOutput{}
	err = json.Unmarshal(outBytes, output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

/*
   "streams": [
       {
           "index": 0,
           "codec_name": "h264",
           "codec_long_name": "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
           "profile": "High",
           "codec_type": "video",
           "codec_time_base": "31008770/1858667763",
           "codec_tag_string": "avc1",
           "codec_tag": "0x31637661",
           "width": 1080,
           "height": 1920,
           "coded_width": 1088,
           "coded_height": 1920,
           "closed_captions": 0,
           "has_b_frames": 2,
           "pix_fmt": "yuv420p",
           "level": 40,
           "color_range": "tv",
           "color_space": "bt709",
           "color_transfer": "bt709",
           "color_primaries": "bt709",
           "chroma_location": "left",
           "refs": 1,
           "is_avc": "true",
           "nal_length_size": "4",
           "r_frame_rate": "30000/1001",
           "avg_frame_rate": "1344455921/44860007",
           "time_base": "1/90000",
           "start_pts": 5400,
           "start_time": "0.060000",
           "duration_ts": 217687443,
           "duration": "2418.749367",
           "bit_rate": "5002442",
           "bits_per_raw_sample": "8",
           "nb_frames": "72490",
           "tags": {
               "language": "und",
               "handler_name": "VideoHandler"
           }
       },
       {
           "index": 1,
           "codec_name": "aac",
           "codec_long_name": "AAC (Advanced Audio Coding)",
           "profile": "LC",
           "codec_type": "audio",
           "codec_time_base": "1/48000",
           "codec_tag_string": "mp4a",
           "codec_tag": "0x6134706d",
           "sample_fmt": "fltp",
           "sample_rate": "48000",
           "channels": 2,
           "channel_layout": "stereo",
           "bits_per_sample": 0,
           "r_frame_rate": "0/0",
           "avg_frame_rate": "0/0",
           "time_base": "1/48000",
           "start_pts": 0,
           "start_time": "0.000000",
           "duration_ts": 116147200,
           "duration": "2419.733333",
           "bit_rate": "191974",
           "max_bit_rate": "191974",
           "nb_frames": "113425",
           "tags": {
               "language": "und",
               "handler_name": "SoundHandler"
           }
       }
   ],
*/

type ffprobeStream struct {
	Index            int          `json:"index"`
	CodecName        string       `json:"codec_name"`
	CodecLongName    string       `json:"codec_long_name"`
	Profile          string       `json:"profile"`
	CodecType        string       `json:"codec_type"`
	CodecTimeBase    string       `json:"codec_time_base"`
	CodecTagString   string       `json:"codec_tag_string"`
	CodecTag         string       `json:"codec_tag"`
	SampleFormat     string       `json:"sample_fmt"`
	SampleRate       int64        `json:"sample_rate,string"`
	Channels         int          `json:"channels"`
	ChannelLayout    string       `json:"channel_layout"`
	Width            int          `json:"width"`
	Height           int          `json:"height"`
	CodedWidth       int          `json:"coded_width"`
	CodedHeight      int          `json:"coded_height"`
	ClosedCaption    int          `json:"closed_captions"`
	HasBFrames       int          `json:"has_b_frames"`
	PixFmt           string       `json:"pix_fmt"`
	Level            int          `json:"level"`
	ColorRange       string       `json:"color_range"`
	ColorSpace       string       `json:"color_space"`
	ColorTransfer    string       `json:"color_transfer"`
	ColorPrimaries   string       `json:"color_primaries"`
	ChromaLocation   string       `json:"chroma_location"`
	Refs             int          `json:"refs"`
	IsAvc            bool         `json:"is_avc,string"`
	NalLengthSize    string       `json:"nal_length_size"`
	RFrameRate       string       `json:"r_frame_rate"`
	AvgFrameRate     string       `json:"avg_frame_rate"`
	TimeBase         string       `json:"time_base"`
	StartPTS         int64        `json:"start_pts"`
	StartTime        float64      `json:"start_time,string"`
	DurationTs       int64        `json:"duration_ts"`
	Duration         float64      `json:"duration,string"`
	BitRate          int64        `json:"bit_rate,string"`
	MaxBitRate       int64        `json:"max_bit_rate,string"`
	BitsPerRawSample int          `json:"bits_per_raw_sample,string"`
	NBFrames         int64        `json:"nb_frames,string"`
	Tags             *ffprobeTags `json:"tags"`
}

/*
	"filename": "/Users/mhand/Desktop/AmyWorkout/AmyWorkout.mp4",
	"nb_streams": 2,
	"nb_programs": 0,
	"format_name": "mov,mp4,m4a,3gp,3g2,mj2",
	"format_long_name": "QuickTime / MOV",
	"start_time": "0.000000",
	"duration": "2419.734000",
	"size": "1573598525",
	"bit_rate": "5202550",
	"probe_score": 100,
	"tags": {
		"major_brand": "isom",
		"minor_version": "512",
		"compatible_brands": "isomiso2avc1mp41",
		"encoder": "Lavf58.45.100"
	}
*/
type ffprobeFormat struct {
	Filename       string       `json:"filename"`
	NBStreams      int          `json:"nb_streams"`
	NBPrograms     int          `json:"nb_programs"`
	FormatName     string       `json:"format_name"`
	FormatLongName string       `json:"format_long_name"`
	StartTime      float64      `json:"start_time,string"`
	Duration       float64      `json:"duration,string"`
	Size           int64        `json:"size,string"`
	BitRate        int64        `json:"bit_rate,string"`
	ProbeScore     int64        `json:"probe_score"`
	Tags           *ffprobeTags `json:"tags"`
}

type ffprobeTags struct {
	MajorBrand       string `json:"major_brand"`
	MinorVersion     int64  `json:"minor_version,string"`
	CompatibleBrands string `json:"compatible_brands"`
	Encoder          string `json:"encoder"`
	Language         string `json:"language"`
	HandlerName      string `json:"handler_name"`
}
