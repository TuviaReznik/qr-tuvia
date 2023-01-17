package receiver

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/tuvirz/qr/go/src/qr/infra/types"
	"github.com/tuvirz/qr/go/src/qr/infra/utils"
)

const (
	TmpQrFileRead  = "../qr_test/data/___TMP_QR_CODE_FILE___.jpeg"
	TmpQrFileWrite = "../qr_test/data/___TMP_ACK_FILE___.jpeg"
	Terminator     = -1
)

func RunReceiver() (string, error) {

	serialNum := 0
	fileName, err := getPackage(serialNum)
	if err != nil {
		return "", err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer file.Close()

	for {
		serialNum++
		data, err := getPackage(serialNum)
		if err != nil {
			return "", err
		}
		if data == "" {
			break
		}

		_, err = file.WriteString(data)
		if err != nil {
			return "", fmt.Errorf("failed to write to destination file: %w", err)
		}
	}

	_, err = getPackage(Terminator)
	if err != nil {
		return "", err
	}

	// cleanup
	return fileName, nil
}

func getPackage(serialNum int) (string, error) {

	data, err := getTextFromQrCode(serialNum)
	if err != nil {
		if err.Error() == "EOF" {
			serialNum = Terminator
		} else {
			return "", err
		}
	}

	err = sendAck(serialNum)
	if err != nil {
		return "", err
	}

	return data, nil
}

func getTextFromQrCode(expSerialNum int) (string, error) {
	retries := 0
	for {
		time.Sleep(time.Millisecond * types.WaitInterval)

		err := utils.CapturePicture(TmpQrFileRead)
		if err != nil {
			return "", fmt.Errorf("failed to capture picture: %w", err)
		}

		text, err := utils.QrCodeToTextWithRetry(TmpQrFileRead, 10)
		if err != nil {
			if expSerialNum == 0 {
				continue
			}
			if retries < 10 {
				retries++
				continue
			}
			return "", fmt.Errorf("failed to convert qr code to text: %w", err)
		}
		retries = 0

		serial, data, err := utils.GetSerialAndDataFromText(text)
		if err != nil {
			return "", fmt.Errorf("failed to parse text: %w", err)
		}

		fmt.Println("--- receive:", serial)
		if serial == Terminator {
			return "", fmt.Errorf("EOF")
		}

		if serial != expSerialNum {
			// fmt.Println("--- serial:", serial)
			continue
		}
		return data, nil
	}
}

func sendAck(serialNum int) error {
	fmt.Println("--- send:", serialNum)

	err := utils.SaveTextAsQRCode(addDummyInfoToAck(serialNum), TmpQrFileWrite)
	if err != nil {
		return fmt.Errorf("failed to send ack as qr code: %w", err)
	}

	utils.DisplayImage(TmpQrFileWrite)
	return nil
}

func addDummyInfoToAck(serialNum int) string {
	return strconv.Itoa(serialNum) + " " + "abcdeABCDE12345xyzXYZ678s"
}
