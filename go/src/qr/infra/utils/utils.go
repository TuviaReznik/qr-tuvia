package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"runtime"

	"os"
	"strconv"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
	"github.com/tuvirz/qr/go/src/qr/infra/types"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

const DummyInfo = " " + "abcdeABCDE12345xyzXYZ678"

var singleton = false

func EncodeFileToBytes(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}

func SaveQrCodeToImageFile(body []byte, fileName string) error {

	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to decode qr code: %w", err)
	}

	out, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create temporary qr code file: %w", err)
	}
	defer out.Close()

	opt := &jpeg.Options{}
	opt.Quality = 100

	err = jpeg.Encode(out, img, opt)
	if err != nil {
		return fmt.Errorf("failed to write qr code to file: %w", err)
	}

	return nil
}

func SaveTextAsQRCode(text, writeFile string) error {
	enc := qrcode.NewQRCodeWriter()
	img, err := enc.Encode(text, gozxing.BarcodeFormat_QR_CODE,
		types.FrameSize, types.FrameSize, nil)
	if err != nil {
		return fmt.Errorf("failed to encode text: %w", err)
	}

	// os.Remove(writeFile)
	file, err := os.Create(writeFile)
	if err != nil {
		return fmt.Errorf("failed to create temporary qr code file: %w", err)
	}
	defer file.Close()

	opt := &jpeg.Options{}
	opt.Quality = 100

	return jpeg.Encode(file, img, opt)
}

func QrCodeToTextWithRetry(fileName string, times int) (string, error) {
	var res string
	var err error
	for i := 0; i < times; i++ {
		res, err = QrCodeToText(fileName)
		if err != nil {
			time.Sleep(time.Millisecond * types.WaitInterval)
			continue
		}
		break
	}
	return res, err
}

func QrCodeToText(fileName string) (string, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to open qr code file: %w", err)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return "", fmt.Errorf("failed to decode image file: %w", err)
	}

	bbm, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("failed to prepare bitmap: %w", err)
	}

	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bbm, nil)
	fmt.Println("--- res:", result)
	if err != nil {
		return "", fmt.Errorf("failed to decode QRCode: %w", err)
	}

	return result.String(), nil
}

func GetSerialAndDataFromText(text string) (int, string, error) {

	spaceIndex := strings.Index(text, " ")
	if spaceIndex == -1 {
		return 0, "", fmt.Errorf("broken data")
	}

	serial := text[:spaceIndex]
	serialNum, err := strconv.Atoi(serial)
	if err != nil {
		return 0, "", fmt.Errorf("missing serial number")
	}

	// no data - EOF
	if serialNum == -1 {
		return serialNum, "", nil
	}

	// normal data
	if len(text) == spaceIndex+1 && serialNum != -1 {
		return 0, "", fmt.Errorf("missing content")
	}
	data := text[spaceIndex+1:]

	return serialNum, data, nil
}

func BuildSerialAndDataPack(serial int, data string) string {
	return fmt.Sprintf("%d %s", serial, data)
}

func UpdateImageDisplay(fileName string) error {
	if !singleton {
		err := DisplayImage(fileName)
		if err != nil {
			return fmt.Errorf("failed to dispaly image on screen: %w", err)
		}
		singleton = true
	}

	time.Sleep(time.Second * 1)
	return nil
}

func DisplayImage(fileName string) error {
	if runtime.GOOS == "linux" {
		return open.Run(fileName)
	}
	if runtime.GOOS == "darwin" {
		// file display isn't updated on regular apps
		err := open.RunWith(fileName, "Visual Studio Code")
		if err != nil {
			return fmt.Errorf("error: %w. did you install vscode?", err)
		}
	}
	return nil
}

func CapturePictureToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create qr code file: %w", err)
	}

	err = CapturePicture(filename)
	if err != nil {
		return fmt.Errorf("failed to capture picture: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close qr code file: %w", err)
	}
	return nil
}
