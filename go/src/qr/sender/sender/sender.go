package sender

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tuvirz/qr/go/src/qr/infra/types"
	"github.com/tuvirz/qr/go/src/qr/infra/utils"
)

const (
	TmpQrFileWrite = "../qr_test/data/___TMP_QR_CODE_FILE___.jpeg"
	TmpQrFileRead  = "../qr_test/data/___TMP_ACK_FILE___.jpeg"
	Terminator     = -1
)

func RunSender(srcFileName, dstFileName string) error {

	os.Remove(TmpQrFileRead)
	os.Remove(TmpQrFileWrite)
	_, err := os.Create(TmpQrFileWrite)
	if err != nil {
		return fmt.Errorf("failed to create qr code file: %w", err)
	}
	defer os.Remove(TmpQrFileWrite)

	file, err := os.Open(srcFileName)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer file.Close()

	// max ascii size to qr code is 2953 chars. we will take no more than 2048.
	serialNum := 0
	err = sendPackage(serialNum, dstFileName)
	if err != nil {
		return err
	}

	for {
		serialNum++
		buf := make([]byte, types.MaxSize)
		charsNum, err := file.Read(buf)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("failed to read source file: %w", err)
			}

			break
		}
		fileContent := string(buf[:charsNum])

		err = sendPackage(serialNum, fileContent)
		if err != nil {
			return err
		}
	}

	err = sendPackage(Terminator, utils.DummyInfo)
	if err != nil {
		return err
	}

	return nil
}

func sendPackage(serialNum int, fileContent string) error {
	text := utils.BuildSerialAndDataPack(serialNum, fileContent)
	fmt.Println("--- send:", serialNum)
	err := utils.SaveTextAsQRCode(text, TmpQrFileWrite)
	if err != nil {
		return fmt.Errorf("failed to send text as qr code: %w", err)
	}

	err = utils.UpdateImageDisplay(TmpQrFileWrite)
	if err != nil {
		return err
	}

	err = waitForAck(serialNum)
	if err != nil {
		return err
	}
	return nil
}

func waitForAck(expSerialNum int) error {
	retries := 0
	for {
		ack, err := getAck()
		fmt.Println("--- expect:", expSerialNum)
		fmt.Println("--- ack:", ack)
		if err != nil {
			fmt.Println("--- error:", err.Error())

			if expSerialNum == 0 {
				continue
			}
			// fmt.Println("failed to convert qr code to text:", err)
			if retries < 10 {
				retries++
				continue
			}
			return fmt.Errorf("failed to get ack from receiver: %w", err)
		}
		retries = 0

		ack = strings.Fields(ack)[0]
		ackNum, err := strconv.Atoi(ack)
		if err != nil {
			return fmt.Errorf("failed to read ack number: %w", err)
		}
		if expSerialNum == Terminator {
			// retry once
			time.Sleep(time.Millisecond * types.WaitInterval * 10)
			ack, err = getAck()
			if err != nil {
				return nil
			}
		}
		if ackNum == expSerialNum {
			break
		}
	}

	return nil
}

func getAck() (string, error) {

	err := utils.CapturePictureToFile(TmpQrFileRead)
	if err != nil {
		return "", err
	}

	time.Sleep(time.Millisecond * 500)

	return utils.QrCodeToText(TmpQrFileRead)
}
