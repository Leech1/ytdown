package main

import (
	"fmt"
	"os"
	"ytdown/extractor"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ytdl <video-url>")
		os.Exit(1)
	}

	videoURL := os.Args[1]

	videoID, err := extractor.ExtractVideoID(videoURL)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	pr, err := extractor.FetchPlayerResponseViaAPI(videoID)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("title: %s\n", pr.VideoDetails.Title)
	fmt.Printf("playability status: %s\n", pr.PlayabilityStatus.Status)
	fmt.Printf("playability reason: %s\n", pr.PlayabilityStatus.Reason)

	// Look for itag 136 (720p video-only mp4) in adaptive formats.
	var target *extractor.Format
	for i, f := range pr.StreamingData.AdaptiveFormats {
		if f.Itag == 136 {
			target = &pr.StreamingData.AdaptiveFormats[i]
			break
		}
	}
	if target == nil {
		fmt.Println("itag 136 not found")
		os.Exit(1)
	}

	playerJS, err := extractor.GetPlayerJS(videoURL)
	if err != nil {
		fmt.Println("error fetching player JS:", err)
		os.Exit(1)
	}

	finalURL, err := extractor.ResolveFormatURL(*target, playerJS)
	if err != nil {
		fmt.Println("error resolving format URL:", err)
		os.Exit(1)
	}

	fmt.Println("resolved URL:", finalURL)
}
