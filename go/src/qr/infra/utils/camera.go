package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"github.com/blackjack/webcam"
)

const (
	topQuality = 100
)

func getCameraDevice() string {
	return "/dev/video0"
	// return "disp0"
	// return "disp0:dcpav-video-interface-epic:0"
	// return "dcpav-video-interface-epic"
	// return "disp0:dcpav-video-interface-epi"
	// return "dcpav-video-interface-epi"
}

func getVideoStream(fileName string, frame []byte) error {

	if len(frame) == 0 {
		return fmt.Errorf("empty frame")
	}

	err := saveContentToImageFile(fileName, frame)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
}

func saveContentToImageFile(fileName string, body []byte) error {
	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	out, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create image file: %w", err)
	}
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = topQuality

	err = jpeg.Encode(out, img, &opts)
	if err != nil {
		return fmt.Errorf("failed to write jpeg image: %w", err)
	}
	return nil
}

func capturePicture(targetFileName string) error {
	cam, err := webcam.Open(getCameraDevice())
	if err != nil {
		return fmt.Errorf("failed to open camera: %w", err)
	}
	defer cam.Close()

	err = cam.StartStreaming()
	if err != nil {
		return fmt.Errorf("failed to stream with camera: %w", err)
	}

	err = cam.WaitForFrame(5)
	switch err.(type) {
	case nil:
	case *webcam.Timeout:
		fmt.Println(err.Error())
		return fmt.Errorf("failed to take a picture with camera: %w", err)
	default:
		return fmt.Errorf("failed to take a picture with camera: %w", err)
	}

	frame, err := cam.ReadFrame()
	if err != nil {
		return fmt.Errorf("failed to read from camera: %w", err)
	}
	return getVideoStream(targetFileName, frame)
}

func CapturePicture(targetFileName string) error {
	return capturePicture(targetFileName)
}
