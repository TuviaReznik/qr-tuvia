package handler

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"

	"github.com/makiuchi-d/gozxing"
	"github.com/tuvirz/qr/go/src/qr/reader"
)

func HandleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	b := new(bytes.Buffer)
	if _, err := io.Copy(b, r.Body); err != nil {
		msg := fmt.Sprintf("Failed to read request body: %v", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	res, err := scan(b.Bytes())
	if err != nil {
		msg := fmt.Sprintf("Internal server error: %v", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Write([]byte(res))
}

func scan(b []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("failed to read image: %v", err)
	}

	source := gozxing.NewLuminanceSourceFromImage(img)
	bin := gozxing.NewHybridBinarizer(source)
	bbm, err := gozxing.NewBinaryBitmap(bin)

	if err != nil {
		return "", fmt.Errorf("error during processing: %w", err)
	}

	qrReader := reader.NewQRCodeReader()
	result, err := qrReader.Decode(bbm)
	if err != nil {
		return "", fmt.Errorf("unable to decode QRCode: %w", err)
	}

	res := result.String()
	return res, nil
}
