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

	videoUrl := os.Args[1]

	raw, err := extractor.GetPlayerResponseJSON(videoUrl)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// temporarily
	fmt.Printf("got player response JSON: %d bytes\n", len(raw))
}
