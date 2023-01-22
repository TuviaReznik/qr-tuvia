package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/blackjack/webcam"
)

const (
	topQuality        = 100
	linuxCameraDevice = "/dev/video0"
)

func saveFrame(fileName string, frame []byte) error {

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

	time.Sleep(time.Second)

	var opts jpeg.Options
	opts.Quality = topQuality

	err = jpeg.Encode(out, img, &opts)
	if err != nil {
		return fmt.Errorf("failed to write jpeg image: %w", err)
	}
	return nil
}

func capturePictureLinux(targetFileName string) error {
	cam, err := webcam.Open(linuxCameraDevice)
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
		return fmt.Errorf("failed to take a picture with camera: timeout: %w", err)
	default:
		return fmt.Errorf("failed to take a picture with camera: %w", err)
	}

	frame, err := cam.ReadFrame()
	if err != nil {
		return fmt.Errorf("failed to read from camera: %w", err)
	}
	return saveFrame(targetFileName, frame)
}

func capturePictureMac(targetFileName string) error {
	cmd := exec.Command("imagesnap", "-w", "1.00", targetFileName)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to run command : imagesnap")
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("failed to wait for command : imagesnap")
	}
	return nil
}

func CapturePicture(targetFileName string) error {
	if runtime.GOOS == "linux" {
		return capturePictureLinux(targetFileName)
	}
	if runtime.GOOS == "darwin" {
		return capturePictureMac(targetFileName)
	}
	return fmt.Errorf("unsupported operating system")
}
