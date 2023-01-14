package qr_test

// installations::
//  ubuntu - sudo apt-get install golang gcc libgl1-mesa-dev xorg-dev

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tuvirz/qr/go/src/qr/infra/utils"
)

const (
	QrFile  = "./data/qrcode_hello_world.bmp"
	TmpFile = "___qr_code_tmp___.bmp"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Second * 1)
}

func TestEncodeDecode(t *testing.T) {
	text := GenerateRandomAlphaNumericString(20)
	fmt.Println("- origin:", text)

	err := utils.SendTextAsQRCode(text, TmpFile)
	require.NoError(t, err)
	defer os.Remove(TmpFile)

	res, err := utils.QrCodeToText(TmpFile)
	fmt.Println("- result:", res)
	require.NoError(t, err)

	require.Equal(t, text, res)
}

func TestCaptureAndDecode(t *testing.T) {
	// preparations:
	//  put your phone in front of the computer camera,
	//  make sure there is QRCode open on your phone screen.

	tmpCaptureFile := "___tmp_test_qr_code___.jpeg"
	err := utils.CapturePicture(tmpCaptureFile)
	require.NoError(t, err)

	res, err := utils.QrCodeToText(tmpCaptureFile)
	fmt.Println("- result:", res)
	require.NoError(t, err)
}

func TestQrCodeToText(t *testing.T) {
	res, err := utils.QrCodeToText(QrFile)
	fmt.Println("--- result:", res)
	require.NoError(t, err)
}

func TestTextToQrCode(t *testing.T) {
	text := GenerateRandomAlphaNumericString(20)
	err := utils.SendTextAsQRCode(text, TmpFile)
	require.NoError(t, err)
}
