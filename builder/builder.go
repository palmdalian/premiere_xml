package builder

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/google/uuid"

	premiere "github.com/palmdalian/premiere_xml"
)

type Timing struct {
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	Rate      int64  `json:"rate"`
	StartTick string `json:"startTick"`
	EndTick   string `json:"endTick"`
	Path      string `json:"path"`
}

type PremiereBuilder struct {
	XML         *premiere.PremiereXML
	Masterclips map[string]*MasterClip
	CurrentClip int64
	FrameRate   int64
}

type MasterClip struct {
	id              string
	File            *premiere.File
	alreadyInserted bool
}

func newAudioClip() *premiere.ClipItem {
	blob := "CQkJCQk8Y2xpcGl0ZW0gaWQ9ImNsaXBpdGVtLTciIHByZW1pZXJlQ2hhbm5lbFR5cGU9Im1vbm8iPgoJCQkJCQk8bWFzdGVyY2xpcGlkPm1hc3RlcmNsaXAtNjwvbWFzdGVyY2xpcGlkPgoJCQkJCQk8bmFtZT5CLmFpZmY8L25hbWU+CgkJCQkJCTxlbmFibGVkPlRSVUU8L2VuYWJsZWQ+CgkJCQkJCTxkdXJhdGlvbj44NTU8L2R1cmF0aW9uPgoJCQkJCQk8cmF0ZT4KCQkJCQkJCTx0aW1lYmFzZT42MDwvdGltZWJhc2U+CgkJCQkJCQk8bnRzYz5GQUxTRTwvbnRzYz4KCQkJCQkJPC9yYXRlPgoJCQkJCQk8c3RhcnQ+MDwvc3RhcnQ+CgkJCQkJCTxlbmQ+ODU1PC9lbmQ+CgkJCQkJCTxpbj4wPC9pbj4KCQkJCQkJPG91dD44NTU8L291dD4KCQkJCQkJPGZpbGUgaWQ9ImZpbGUtNiI+CgkJCQkJCQk8bmFtZT5CLmFpZmY8L25hbWU+CgkJCQkJCQk8cGF0aHVybD5maWxlOi8vbG9jYWxob3N0L0IuYWlmZjwvcGF0aHVybD4KCQkJCQkJCTxyYXRlPgoJCQkJCQkJCTx0aW1lYmFzZT4zMDwvdGltZWJhc2U+CgkJCQkJCQkJPG50c2M+VFJVRTwvbnRzYz4KCQkJCQkJCTwvcmF0ZT4KCQkJCQkJCTxkdXJhdGlvbj40Mjc8L2R1cmF0aW9uPgoJCQkJCQkJPHRpbWVjb2RlPgoJCQkJCQkJCTxyYXRlPgoJCQkJCQkJCQk8dGltZWJhc2U+MzA8L3RpbWViYXNlPgoJCQkJCQkJCQk8bnRzYz5UUlVFPC9udHNjPgoJCQkJCQkJCTwvcmF0ZT4KCQkJCQkJCQk8c3RyaW5nPjAwOzAwOzAwOzAwPC9zdHJpbmc+CgkJCQkJCQkJPGZyYW1lPjA8L2ZyYW1lPgoJCQkJCQkJCTxkaXNwbGF5Zm9ybWF0PkRGPC9kaXNwbGF5Zm9ybWF0PgoJCQkJCQkJCTxyZWVsPgoJCQkJCQkJCQk8bmFtZT48L25hbWU+CgkJCQkJCQkJPC9yZWVsPgoJCQkJCQkJPC90aW1lY29kZT4KCQkJCQkJCTxtZWRpYT4KCQkJCQkJCQk8YXVkaW8+CgkJCQkJCQkJCTxzYW1wbGVjaGFyYWN0ZXJpc3RpY3M+CgkJCQkJCQkJCQk8ZGVwdGg+MTY8L2RlcHRoPgoJCQkJCQkJCQkJPHNhbXBsZXJhdGU+NDgwMDA8L3NhbXBsZXJhdGU+CgkJCQkJCQkJCTwvc2FtcGxlY2hhcmFjdGVyaXN0aWNzPgoJCQkJCQkJCQk8Y2hhbm5lbGNvdW50PjE8L2NoYW5uZWxjb3VudD4KCQkJCQkJCQkJPGF1ZGlvY2hhbm5lbD4KCQkJCQkJCQkJCTxzb3VyY2VjaGFubmVsPjE8L3NvdXJjZWNoYW5uZWw+CgkJCQkJCQkJCTwvYXVkaW9jaGFubmVsPgoJCQkJCQkJCTwvYXVkaW8+CgkJCQkJCQk8L21lZGlhPgoJCQkJCQk8L2ZpbGU+CgkJCQkJCTxzb3VyY2V0cmFjaz4KCQkJCQkJCTxtZWRpYXR5cGU+YXVkaW88L21lZGlhdHlwZT4KCQkJCQkJCTx0cmFja2luZGV4PjE8L3RyYWNraW5kZXg+CgkJCQkJCTwvc291cmNldHJhY2s+CgkJCQkJCTxsb2dnaW5naW5mbz4KCQkJCQkJCTxkZXNjcmlwdGlvbj48L2Rlc2NyaXB0aW9uPgoJCQkJCQkJPHNjZW5lPjwvc2NlbmU+CgkJCQkJCQk8c2hvdHRha2U+PC9zaG90dGFrZT4KCQkJCQkJCTxsb2dub3RlPjwvbG9nbm90ZT4KCQkJCQkJPC9sb2dnaW5naW5mbz4KCQkJCQkJPGxhYmVscz4KCQkJCQkJCTxsYWJlbDI+Q2FyaWJiZWFuPC9sYWJlbDI+CgkJCQkJCTwvbGFiZWxzPgoJCQkJCTwvY2xpcGl0ZW0+"
	clip := &premiere.ClipItem{}
	data, _ := base64.StdEncoding.DecodeString(blob)
	xml.Unmarshal(data, clip)
	return clip
}

func newVideoClip() *premiere.ClipItem {
	blob := "CQkJCQk8Y2xpcGl0ZW0gaWQ9ImNsaXBpdGVtLTYiPgoJCQkJCQk8bWFzdGVyY2xpcGlkPm1hc3RlcmNsaXAtNTwvbWFzdGVyY2xpcGlkPgoJCQkJCQk8bmFtZT5BLm1wNDwvbmFtZT4KCQkJCQkJPGVuYWJsZWQ+VFJVRTwvZW5hYmxlZD4KCQkJCQkJPGR1cmF0aW9uPjQyNzwvZHVyYXRpb24+CgkJCQkJCTxyYXRlPgoJCQkJCQkJPHRpbWViYXNlPjYwPC90aW1lYmFzZT4KCQkJCQkJCTxudHNjPkZBTFNFPC9udHNjPgoJCQkJCQk8L3JhdGU+CgkJCQkJCTxzdGFydD4wPC9zdGFydD4KCQkJCQkJPGVuZD40Mjc8L2VuZD4KCQkJCQkJPGluPjA8L2luPgoJCQkJCQk8b3V0PjQyNzwvb3V0PgoJCQkJCQk8YWxwaGF0eXBlPm5vbmU8L2FscGhhdHlwZT4KCQkJCQkJPHBpeGVsYXNwZWN0cmF0aW8+c3F1YXJlPC9waXhlbGFzcGVjdHJhdGlvPgoJCQkJCQk8YW5hbW9ycGhpYz5GQUxTRTwvYW5hbW9ycGhpYz4KCQkJCQkJPGZpbGUgaWQ9ImZpbGUtNSI+CgkJCQkJCQk8bmFtZT5BLm1wNDwvbmFtZT4KCQkJCQkJCTxwYXRodXJsPmZpbGU6Ly9sb2NhbGhvc3QvQS5tcDQ8L3BhdGh1cmw+CgkJCQkJCQk8cmF0ZT4KCQkJCQkJCQk8dGltZWJhc2U+NjA8L3RpbWViYXNlPgoJCQkJCQkJCTxudHNjPkZBTFNFPC9udHNjPgoJCQkJCQkJPC9yYXRlPgoJCQkJCQkJPGR1cmF0aW9uPjQyNzwvZHVyYXRpb24+CgkJCQkJCQk8dGltZWNvZGU+CgkJCQkJCQkJPHJhdGU+CgkJCQkJCQkJCTx0aW1lYmFzZT42MDwvdGltZWJhc2U+CgkJCQkJCQkJCTxudHNjPkZBTFNFPC9udHNjPgoJCQkJCQkJCTwvcmF0ZT4KCQkJCQkJCQk8c3RyaW5nPjAwOjAwOjAwOjAwPC9zdHJpbmc+CgkJCQkJCQkJPGZyYW1lPjA8L2ZyYW1lPgoJCQkJCQkJCTxkaXNwbGF5Zm9ybWF0Pk5ERjwvZGlzcGxheWZvcm1hdD4KCQkJCQkJCQk8cmVlbD4KCQkJCQkJCQkJPG5hbWU+PC9uYW1lPgoJCQkJCQkJCTwvcmVlbD4KCQkJCQkJCTwvdGltZWNvZGU+CgkJCQkJCQk8bWVkaWE+CgkJCQkJCQkJPHZpZGVvPgoJCQkJCQkJCQk8c2FtcGxlY2hhcmFjdGVyaXN0aWNzPgoJCQkJCQkJCQkJPHJhdGU+CgkJCQkJCQkJCQkJPHRpbWViYXNlPjYwPC90aW1lYmFzZT4KCQkJCQkJCQkJCQk8bnRzYz5GQUxTRTwvbnRzYz4KCQkJCQkJCQkJCTwvcmF0ZT4KCQkJCQkJCQkJCTx3aWR0aD4xOTIwPC93aWR0aD4KCQkJCQkJCQkJCTxoZWlnaHQ+MTA4MDwvaGVpZ2h0PgoJCQkJCQkJCQkJPGFuYW1vcnBoaWM+RkFMU0U8L2FuYW1vcnBoaWM+CgkJCQkJCQkJCQk8cGl4ZWxhc3BlY3RyYXRpbz5zcXVhcmU8L3BpeGVsYXNwZWN0cmF0aW8+CgkJCQkJCQkJCQk8ZmllbGRkb21pbmFuY2U+bm9uZTwvZmllbGRkb21pbmFuY2U+CgkJCQkJCQkJCTwvc2FtcGxlY2hhcmFjdGVyaXN0aWNzPgoJCQkJCQkJCTwvdmlkZW8+CgkJCQkJCQk8L21lZGlhPgoJCQkJCQk8L2ZpbGU+CgkJCQkJCTxsb2dnaW5naW5mbz4KCQkJCQkJCTxkZXNjcmlwdGlvbj48L2Rlc2NyaXB0aW9uPgoJCQkJCQkJPHNjZW5lPjwvc2NlbmU+CgkJCQkJCQk8c2hvdHRha2U+PC9zaG90dGFrZT4KCQkJCQkJCTxsb2dub3RlPjwvbG9nbm90ZT4KCQkJCQkJPC9sb2dnaW5naW5mbz4KCQkJCQkJPGxhYmVscz4KCQkJCQkJCTxsYWJlbDI+VmlvbGV0PC9sYWJlbDI+CgkJCQkJCTwvbGFiZWxzPgoJCQkJCTwvY2xpcGl0ZW0+"
	clip := &premiere.ClipItem{}
	data, _ := base64.StdEncoding.DecodeString(blob)
	xml.Unmarshal(data, clip)
	return clip
}

func NewPremiereBuilder() (*PremiereBuilder, error) {
	blob := "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPCFET0NUWVBFIHhtZW1sPgo8eG1lbWwgdmVyc2lvbj0iNCI+Cgk8c2VxdWVuY2UgaWQ9InNlcXVlbmNlLTMiIFRMLlNRQXVkaW9WaXNpYmxlQmFzZT0iMCIgVEwuU1FWaWRlb1Zpc2libGVCYXNlPSIwIiBUTC5TUVZpc2libGVCYXNlVGltZT0iMCIgVEwuU1FBVkRpdmlkZXJQb3NpdGlvbj0iMC41IiBUTC5TUUhpZGVTaHlUcmFja3M9IjAiIFRMLlNRSGVhZGVyV2lkdGg9IjIzNiIgTW9uaXRvci5Qcm9ncmFtWm9vbU91dD0iMzYxOTcyODAwMDAwMCIgTW9uaXRvci5Qcm9ncmFtWm9vbUluPSIwIiBUTC5TUVRpbWVQZXJQaXhlbD0iMC40MzM1MjYwMTE1NjA2OTM2NSIgTVouRWRpdExpbmU9IjU1MDM2ODAwMDAwMCIgTVouU2VxdWVuY2UuUHJldmlld0ZyYW1lU2l6ZUhlaWdodD0iMTA4MCIgTVouU2VxdWVuY2UuUHJldmlld0ZyYW1lU2l6ZVdpZHRoPSIxOTIwIiBNWi5TZXF1ZW5jZS5BdWRpb1RpbWVEaXNwbGF5Rm9ybWF0PSIyMDAiIE1aLlNlcXVlbmNlLlByZXZpZXdSZW5kZXJpbmdDbGFzc0lEPSIxMjk3MTA2NzYxIiBNWi5TZXF1ZW5jZS5QcmV2aWV3UmVuZGVyaW5nUHJlc2V0Q29kZWM9IjEyOTcxMDcyNzgiIE1aLlNlcXVlbmNlLlByZXZpZXdSZW5kZXJpbmdQcmVzZXRQYXRoPSJFbmNvZGVyUHJlc2V0cy9TZXF1ZW5jZVByZXZpZXcvNzk1NDU0ZDktZDNjMi00MjlkLTk0NzQtOTIzYWIxM2I3MDE4L0ktRnJhbWUgT25seSBNUEVHLmVwciIgTVouU2VxdWVuY2UuUHJldmlld1VzZU1heFJlbmRlclF1YWxpdHk9ImZhbHNlIiBNWi5TZXF1ZW5jZS5QcmV2aWV3VXNlTWF4Qml0RGVwdGg9ImZhbHNlIiBNWi5TZXF1ZW5jZS5FZGl0aW5nTW9kZUdVSUQ9Ijc5NTQ1NGQ5LWQzYzItNDI5ZC05NDc0LTkyM2FiMTNiNzAxOCIgTVouU2VxdWVuY2UuVmlkZW9UaW1lRGlzcGxheUZvcm1hdD0iMTA4IiBNWi5Xb3JrT3V0UG9pbnQ9IjM2MTk3MjgwMDAwMDAiIE1aLldvcmtJblBvaW50PSIwIiBNWi5aZXJvUG9pbnQ9IjAiIGV4cGxvZGVkVHJhY2tzPSJ0cnVlIj4KCQk8dXVpZD5mZjE5ZmY1Zi0zYjFiLTQ5ZDctODE5MC1iZDY2Y2EwN2MyZDk8L3V1aWQ+CgkJPGR1cmF0aW9uPjg1NTwvZHVyYXRpb24+CgkJPHJhdGU+CgkJCTx0aW1lYmFzZT42MDwvdGltZWJhc2U+CgkJCTxudHNjPkZBTFNFPC9udHNjPgoJCTwvcmF0ZT4KCQk8bmFtZT5TZXF1ZW5jZTwvbmFtZT4KCQk8bWVkaWE+CgkJCTx2aWRlbz4KCQkJCTxmb3JtYXQ+CgkJCQkJPHNhbXBsZWNoYXJhY3RlcmlzdGljcz4KCQkJCQkJPHJhdGU+CgkJCQkJCQk8dGltZWJhc2U+NjA8L3RpbWViYXNlPgoJCQkJCQkJPG50c2M+RkFMU0U8L250c2M+CgkJCQkJCTwvcmF0ZT4KCQkJCQkJPGNvZGVjPgoJCQkJCQkJPG5hbWU+QXBwbGUgUHJvUmVzIDQyMjwvbmFtZT4KCQkJCQkJCTxhcHBzcGVjaWZpY2RhdGE+CgkJCQkJCQkJPGFwcG5hbWU+RmluYWwgQ3V0IFBybzwvYXBwbmFtZT4KCQkJCQkJCQk8YXBwbWFudWZhY3R1cmVyPkFwcGxlIEluYy48L2FwcG1hbnVmYWN0dXJlcj4KCQkJCQkJCQk8YXBwdmVyc2lvbj43LjA8L2FwcHZlcnNpb24+CgkJCQkJCQkJPGRhdGE+CgkJCQkJCQkJCTxxdGNvZGVjPgoJCQkJCQkJCQkJPGNvZGVjbmFtZT5BcHBsZSBQcm9SZXMgNDIyPC9jb2RlY25hbWU+CgkJCQkJCQkJCQk8Y29kZWN0eXBlbmFtZT5BcHBsZSBQcm9SZXMgNDIyPC9jb2RlY3R5cGVuYW1lPgoJCQkJCQkJCQkJPGNvZGVjdHlwZWNvZGU+YXBjbjwvY29kZWN0eXBlY29kZT4KCQkJCQkJCQkJCTxjb2RlY3ZlbmRvcmNvZGU+YXBwbDwvY29kZWN2ZW5kb3Jjb2RlPgoJCQkJCQkJCQkJPHNwYXRpYWxxdWFsaXR5PjEwMjQ8L3NwYXRpYWxxdWFsaXR5PgoJCQkJCQkJCQkJPHRlbXBvcmFscXVhbGl0eT4wPC90ZW1wb3JhbHF1YWxpdHk+CgkJCQkJCQkJCQk8a2V5ZnJhbWVyYXRlPjA8L2tleWZyYW1lcmF0ZT4KCQkJCQkJCQkJCTxkYXRhcmF0ZT4wPC9kYXRhcmF0ZT4KCQkJCQkJCQkJPC9xdGNvZGVjPgoJCQkJCQkJCTwvZGF0YT4KCQkJCQkJCTwvYXBwc3BlY2lmaWNkYXRhPgoJCQkJCQk8L2NvZGVjPgoJCQkJCQk8d2lkdGg+MTkyMDwvd2lkdGg+CgkJCQkJCTxoZWlnaHQ+MTA4MDwvaGVpZ2h0PgoJCQkJCQk8YW5hbW9ycGhpYz5GQUxTRTwvYW5hbW9ycGhpYz4KCQkJCQkJPHBpeGVsYXNwZWN0cmF0aW8+c3F1YXJlPC9waXhlbGFzcGVjdHJhdGlvPgoJCQkJCQk8ZmllbGRkb21pbmFuY2U+bm9uZTwvZmllbGRkb21pbmFuY2U+CgkJCQkJCTxjb2xvcmRlcHRoPjI0PC9jb2xvcmRlcHRoPgoJCQkJCTwvc2FtcGxlY2hhcmFjdGVyaXN0aWNzPgoJCQkJPC9mb3JtYXQ+CgkJCQk8dHJhY2sgVEwuU1FUcmFja1NoeT0iMCIgVEwuU1FUcmFja0V4cGFuZGVkSGVpZ2h0PSIyNSIgVEwuU1FUcmFja0V4cGFuZGVkPSIwIiBNWi5UcmFja1RhcmdldGVkPSIxIj4KCQkJCQk8ZW5hYmxlZD5UUlVFPC9lbmFibGVkPgoJCQkJCTxsb2NrZWQ+RkFMU0U8L2xvY2tlZD4KCQkJCTwvdHJhY2s+CgkJCTwvdmlkZW8+CgkJCTxhdWRpbz4KCQkJCTxudW1PdXRwdXRDaGFubmVscz4yPC9udW1PdXRwdXRDaGFubmVscz4KCQkJCTxmb3JtYXQ+CgkJCQkJPHNhbXBsZWNoYXJhY3RlcmlzdGljcz4KCQkJCQkJPGRlcHRoPjE2PC9kZXB0aD4KCQkJCQkJPHNhbXBsZXJhdGU+NDgwMDA8L3NhbXBsZXJhdGU+CgkJCQkJPC9zYW1wbGVjaGFyYWN0ZXJpc3RpY3M+CgkJCQk8L2Zvcm1hdD4KCQkJCTxvdXRwdXRzPgoJCQkJCTxncm91cD4KCQkJCQkJPGluZGV4PjE8L2luZGV4PgoJCQkJCQk8bnVtY2hhbm5lbHM+MTwvbnVtY2hhbm5lbHM+CgkJCQkJCTxkb3dubWl4PjA8L2Rvd25taXg+CgkJCQkJCTxjaGFubmVsPgoJCQkJCQkJPGluZGV4PjE8L2luZGV4PgoJCQkJCQk8L2NoYW5uZWw+CgkJCQkJPC9ncm91cD4KCQkJCQk8Z3JvdXA+CgkJCQkJCTxpbmRleD4yPC9pbmRleD4KCQkJCQkJPG51bWNoYW5uZWxzPjE8L251bWNoYW5uZWxzPgoJCQkJCQk8ZG93bm1peD4wPC9kb3dubWl4PgoJCQkJCQk8Y2hhbm5lbD4KCQkJCQkJCTxpbmRleD4yPC9pbmRleD4KCQkJCQkJPC9jaGFubmVsPgoJCQkJCTwvZ3JvdXA+CgkJCQk8L291dHB1dHM+CgkJCQk8dHJhY2sgVEwuU1FUcmFja0F1ZGlvS2V5ZnJhbWVTdHlsZT0iMCIgVEwuU1FUcmFja1NoeT0iMCIgVEwuU1FUcmFja0V4cGFuZGVkSGVpZ2h0PSIyNSIgVEwuU1FUcmFja0V4cGFuZGVkPSIwIiBNWi5UcmFja1RhcmdldGVkPSIxIiBQYW5uZXJDdXJyZW50VmFsdWU9IjAuNSIgUGFubmVySXNJbnZlcnRlZD0idHJ1ZSIgUGFubmVyU3RhcnRLZXlmcmFtZT0iLTkxNDQ1NzYwMDAwMDAwMDAwLDAuNSwwLDAsMCwwLDAsMCIgUGFubmVyTmFtZT0iQmFsYW5jZSIgY3VycmVudEV4cGxvZGVkVHJhY2tJbmRleD0iMCIgdG90YWxFeHBsb2RlZFRyYWNrQ291bnQ9IjEiIHByZW1pZXJlVHJhY2tUeXBlPSJTdGVyZW8iPgoJCQkJCTxlbmFibGVkPlRSVUU8L2VuYWJsZWQ+CgkJCQkJPGxvY2tlZD5GQUxTRTwvbG9ja2VkPgoJCQkJCTxvdXRwdXRjaGFubmVsaW5kZXg+MTwvb3V0cHV0Y2hhbm5lbGluZGV4PgoJCQkJPC90cmFjaz4KCQkJPC9hdWRpbz4KCQk8L21lZGlhPgoJCTx0aW1lY29kZT4KCQkJPHJhdGU+CgkJCQk8dGltZWJhc2U+NjA8L3RpbWViYXNlPgoJCQkJPG50c2M+RkFMU0U8L250c2M+CgkJCTwvcmF0ZT4KCQkJPHN0cmluZz4wMDowMDowMDowMDwvc3RyaW5nPgoJCQk8ZnJhbWU+MDwvZnJhbWU+CgkJCTxkaXNwbGF5Zm9ybWF0Pk5ERjwvZGlzcGxheWZvcm1hdD4KCQk8L3RpbWVjb2RlPgoJCTxsYWJlbHM+CgkJCTxsYWJlbDI+Rm9yZXN0PC9sYWJlbDI+CgkJPC9sYWJlbHM+Cgk8L3NlcXVlbmNlPgo8L3htZW1sPgo="
	pxml := &premiere.PremiereXML{}
	data, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return nil, err
	}
	if err := xml.Unmarshal(data, pxml); err != nil {
		return nil, err
	}
	newID := uuid.New()
	pxml.Sequence.UUID = newID.String()
	builder := &PremiereBuilder{
		XML:         pxml,
		Masterclips: make(map[string]*MasterClip),
		FrameRate:   pxml.Sequence.Rate.Timebase,
	}
	return builder, nil
}

func (builder *PremiereBuilder) AddNewMasterclip(filePath string, masterFile *premiere.File) *MasterClip {
	masterClip := &MasterClip{id: fmt.Sprintf("masterclip-%d", len(builder.Masterclips))}
	masterClip.File = masterFile
	builder.Masterclips[filePath] = masterClip
	return masterClip
}

// masterClipFromFilePath return masterclip if already exists
// otherwise use ffprobe to figure out file attributes
func (builder *PremiereBuilder) masterClipFromFilePath(filePath string) *MasterClip {
	masterClip, ok := builder.Masterclips[filePath]
	if ok {
		return masterClip
	}
	masterFile := masterFileFromFilePath(filePath, len(builder.Masterclips))
	masterClip = builder.AddNewMasterclip(filePath, masterFile)
	return masterClip
}

func (builder *PremiereBuilder) AddNewClipItem(clipType, name, filePath string, start, end, insert int64) {
	var tempClip *premiere.ClipItem
	var referenceTrack *premiere.Track
	masterClip := builder.masterClipFromFilePath(filePath)
	if clipType == "audio" {
		tempClip = newAudioClip()
		if masterClip.File.Media.Audio.NumOutputChannels > 1 {
			tempClip.PremiereChannelType = "stereo"
		}
		referenceTrack = builder.XML.Sequence.Media.Audio.Tracks[0]
	} else {
		tempClip = newVideoClip()
		builder.XML.Sequence.Media.Video.Format.SampleCharacteristics.Rate = masterClip.File.Rate
		builder.XML.Sequence.Media.Video.Format.SampleCharacteristics.Width = masterClip.File.Media.Video.SampleCharacteristics.Width
		builder.XML.Sequence.Media.Video.Format.SampleCharacteristics.Height = masterClip.File.Media.Video.SampleCharacteristics.Height
		referenceTrack = builder.XML.Sequence.Media.Video.Tracks[0]
	}

	if name == "" {
		name = path.Base(filePath)
	}
	tempClip.Name = name
	tempClip.Id = fmt.Sprintf("clipitem-%v", builder.CurrentClip)

	tempClip.MasterClipId = masterClip.id
	tempClip.File = masterClip.File
	if masterClip.alreadyInserted {
		tempClip.File = &premiere.File{Id: masterClip.File.Id}
	}

	tempClip.Rate = masterClip.File.Rate
	tempClip.In = start
	tempClip.Out = end
	tempClip.PProTicksIn = start * premiere.PProTicksConstant
	tempClip.PProTicksInOut = end * premiere.PProTicksConstant
	tempClip.Start = insert
	tempClip.End = (insert + end - start)
	tempClip.Duration = masterClip.File.Duration

	masterClip.alreadyInserted = true
	referenceTrack.ClipItems = append(referenceTrack.ClipItems, tempClip)
	builder.CurrentClip += 1
}

func (builder *PremiereBuilder) ProcessVideoTimings(timings []*Timing) {
	for _, timing := range timings {
		builder.AddNewClipItem("video", "", timing.Path, timing.Start, timing.End, timing.Start)
	}
	if len(timings) > 0 {
		builder.XML.Sequence.Rate.Timebase = timings[0].Rate
		builder.XML.Sequence.Timecode.Rate.Timebase = timings[0].Rate
	}
}

func (builder *PremiereBuilder) ProcessAudioTimings(timings []*Timing) {
	for _, timing := range timings {
		builder.AddNewClipItem("audio", "", timing.Path, timing.Start, timing.End, timing.Start)
	}
}

func (builder *PremiereBuilder) SaveToPath(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	xmlWriter := io.Writer(file)

	header := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE xmeml>\n"
	xmlWriter.Write([]byte((header)))
	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("", "    ")
	if err := enc.Encode(builder.XML); err != nil {
		return err
	}
	return nil
}

func masterFileFromFilePath(filePath string, fileNumber int) *premiere.File {
	ffOut, err := getFFProbeOutput(filePath)
	if err != nil {
		fmt.Printf("Error getting clip attributes with ffprobe. File will be nil %v\n", err)
		return nil
	}

	if len(ffOut.Streams) == 0 {
		fmt.Printf("No streams found. File will be nil %v\n", err)
		return nil
	}

	media := &premiere.Media{
		Video: videoFromStreams(ffOut.Streams),
		Audio: audioFromStreams(ffOut.Streams),
	}

	rate := int64(0)
	ntsc := false
	if media.Video != nil {
		rate = media.Video.SampleCharacteristics.Rate.Timebase
		ntsc = media.Video.SampleCharacteristics.Rate.NTSC
	} else if media.Audio != nil {
		rate = media.Audio.SampleCharacteristics.SampleRate
	}

	name := path.Base(filePath)
	masterFile := &premiere.File{
		Id:       fmt.Sprintf("file-%d", fileNumber),
		Name:     name,
		Duration: int64(ffOut.Format.Duration * float64(rate)),
		Rate: &premiere.Rate{
			Timebase: rate,
			NTSC:     ntsc,
		},
		Timecode: &premiere.Timecode{
			Rate: &premiere.Rate{
				Timebase: rate,
				NTSC:     ntsc,
			},
			String:        "00:00:00:00",
			DisplayFormat: "NDF",
		},
		Media:   media,
		PathUrl: fmt.Sprintf("file://localhost%s", filePath),
	}

	return masterFile
}
