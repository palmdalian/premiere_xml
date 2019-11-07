package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/palmdalian/premiere_xml/builder"
)

func main() {
	var jsonPath string
	var outputPath string

	flag.StringVar(&jsonPath, "i", "/tmp/test.json", "inputPath")
	flag.StringVar(&outputPath, "o", "out.xml", "outputPath")
	flag.Parse()

	fd, err := os.Open(jsonPath)
	defer fd.Close()
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	bytes, _ := ioutil.ReadAll(fd)
	timings := []*builder.Timing{}
	json.Unmarshal(bytes, &timings)

	builder, err := builder.NewPremiereBuilder()
	if err != nil {
		log.Fatal(err)
	}

	builder.ProcessAudioTimings(timings)
	if err := builder.SaveToPath(outputPath); err != nil {
		log.Printf("err: %v", err)
	}

}
