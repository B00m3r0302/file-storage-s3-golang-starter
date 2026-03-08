package main

import (
	"bytes"
	"encoding/json"
	"math"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	type ffprobeOutput struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	data := out.Bytes()
	result := ffprobeOutput{}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return "", err
	}

	landscape := 16.0 / 9.0
	portrait := 9.0 / 16.0

	ratio := float64(result.Streams[0].Width) / float64(result.Streams[0].Height)
	tolerance := 0.05

	if math.Abs(ratio-landscape) < tolerance {
		return "landscape", nil
	}

	if math.Abs(ratio-portrait) < tolerance {
		return "portrait", nil
	}

	return "other", nil
}
