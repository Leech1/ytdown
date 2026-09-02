package main

import (
	"fmt"
	"os"
	"ytdown/extractor"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ytdown <video-link>")
		os.Exit(1)
	}

	videoURL := os.Args[1]

	pr, err := extractor.GetPlayerResponse(videoURL)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Printf("title: %s\n", pr.VideoDetails.Title)
	fmt.Printf("progressive formats found: %d\n", len(pr.StreamingData.Formats))

	format, err := extractor.Find720pProgressive(pr.StreamingData.Formats)
	if err != nil {
		fmt.Println("no 720p format found:", err)
		os.Exit(1)
	}

	fmt.Printf("720p format: itag=%d mimeType=%s hasURL=%v\n",
		format.Itag, format.MimeType, format.URL != "")
}
