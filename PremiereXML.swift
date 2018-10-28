//
//  PremiereXML.swift
//  AudioSplitter
//
//  Created by Michael Hand on 8/25/18.
//  Copyright © 2018 Michael Hand. All rights reserved.
//
// WIP
// Need https://github.com/ShawnMoore/XMLParsing for this to work.

import Foundation

//const PProTicksConstant = 254016000000

struct PremiereXML: Codable {
    var XMLBase: XMLBase

    enum CodingKeys: String, CodingKey {
        case XMLBase = "xmeml"
    }
}

struct XMLBase: Codable {
    var sequence: [Sequence]
}

struct Sequence: Codable {
    var uuid: String
    var duration: Int64
    var rate: Rate
    var name: String
    var media: Media
    var timecode: Timecode
    var markers: [Marker]
    var labels: [Label]
}

struct Rate: Codable {
    var timebase: Int64
    var ntsc: Bool
}

struct Timecode: Codable {
    var rate: Rate
    var string: String
    var frame: Int64
}

struct Marker: Codable {
    var comment: String
    var name: String
    var inPoint: Int64
    var outPoint: Int64

    enum CodingKeys: String, CodingKey {
        case comment, name
        case inPoint = "in"
        case outPoint = "out"
    }
}

struct Label: Codable {
    var color: String

    enum CodingKeys: String, CodingKey {
        case color = "label2"
    }
}

struct Media: Codable {
    var video: Video
    var audio: Audio
}

struct Video: Codable {
    var format: Format
    var tracks: [Track]

    enum CodingKeys: String, CodingKey {
        case format
        case tracks = "track"
    }
}

struct Audio: Codable {
    var format: Format
    var tracks: [Track]
    var numOutputChannels: Int

    enum CodingKeys: String, CodingKey {
        case format, numOutputChannels
        case tracks = "track"
    }
}

struct Format: Codable {
    var sampleCharacteristics: SampleCharacteristics

    enum CodingKeys: String, CodingKey {
        case sampleCharacteristics = "samplecharacteristics"
    }
}

struct SampleCharacteristics: Codable {
    var depth: Int64
    var sampleRate: Int64

    enum CodingKeys: String, CodingKey {
        case depth
        case sampleRate = "samplerate"
    }
}

struct Track: Codable {
    var clipItems: [ClipItem]
    var transitionItems: [TransitionItem]
    var enabled: Bool
    var locked: Bool
    var outputChannelIndex: Int

    enum CodingKeys: String, CodingKey {
        case enabled, locked
        case clipItems = "clipitem"
        case outputChannelIndex = "outputchannelindex"
        case transitionItems = "transitionitem"
    }
}

struct ClipItem: Codable {
    var id: String
    var masterClipId: String
    var name: String
    var enabled: Bool
    var duration: Int64
    var rate: Rate
    var start: Int64 // In point within the sequence
    var end: Int64 // Out point within the sequence
    var inPoint: Int64 // In point within the media file
    var outPoint: Int64 // Out point within the media file
    var pproTicksIn: Int64 // Premiere specific in point (use const to get seconds)
    var pproTicksOut: Int64 // Premiere specific out point (use const to get seconds)
    var label: Label
    var file: File
    var loggingInfo: LoggingInfo
    var links: [Link]
    var sourceTrack: SourceTrack

    enum CodingKeys: String, CodingKey {
        case id, name, enabled, duration, rate, start, end, pproTicksIn, pproTicksOut, label, file
        case sourceTrack = "sourcetrack"
        case links = "link"
        case masterClipId = "masterclipid"
        case loggingInfo = "logginginfo"
        case inPoint = "in"
        case outPoint = "out"
    }
}

struct File: Codable {
    var id: String
    var name: String
    var pathUrl: String
    var rate: Rate
    var duration: Int64
    var timecode: Timecode
    var media: Media

    enum CodingKeys: String, CodingKey {
        case id, name, duration, rate, timecode, media
        case pathUrl = "pathurl"
    }
}

struct LoggingInfo: Codable {
    var description: String
    var scene: String
    var shotTake: String
    var logNote: String
    var good: String
    var originalVideoFilename: String
    var originalAudioFilename: String

    enum CodingKeys: String, CodingKey {
        case description, scene, good
        case shotTake = "shottake"
        case logNote = "lognote"
        case originalVideoFilename = "originalvideofilename"
        case originalAudioFilename = "originalaudiofilename"
    }
}

struct SourceTrack: Codable {
    var mediaType: String
    var trackIndex: Int

    enum CodingKeys: String, CodingKey {
        case mediaType = "sourcetrack"
        case trackIndex = "trackindex"
    }
}

struct TransitionItem: Codable {
    var start: Int64
    var end: Int64
    var alignment: String
    var cutPointTicks: Int64
    var rate: Rate
}

struct Link: Codable {
    var linkClipRef: String
    var mediaType: String
    var trackIndex: Int
    var clipIndex: Int
    var groupIndex: Int

    enum CodingKeys: String, CodingKey {
        case linkClipRef = "linkclipref"
        case mediaType = "mediatype"
        case trackIndex = "trackindex"
        case clipIndex = "clipindex"
        case groupIndex = "groupindex"
    }
}
