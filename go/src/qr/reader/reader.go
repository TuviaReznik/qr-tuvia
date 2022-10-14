package reader

import (

	// support gif, jpeg, png
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	// support bmp, tiff, webp
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode/decoder"
	"github.com/makiuchi-d/gozxing/qrcode/detector"
)

type QRCodeReader struct {
	decoder *decoder.Decoder

	QRCodeReaderI
}

func NewQRCodeReader() *QRCodeReader {
	return &QRCodeReader{
		decoder: decoder.NewDecoder(),
	}
}

type QRCodeReaderI interface {
	Decode(image *gozxing.BinaryBitmap) (*gozxing.Result, error)
}

func (r *QRCodeReader) Decode(image *gozxing.BinaryBitmap) (*gozxing.Result, error) {

	matrix, err := image.GetBlackMatrix()
	if err != nil {
		return nil, err
	}

	detectorResult, err := detector.NewDetector(matrix).Detect(nil)
	if err != nil {
		return nil, err
	}

	decoderResult, err := r.decoder.Decode(detectorResult.GetBits(), nil)
	if err != nil {
		return nil, err
	}

	points := detectorResult.GetPoints()
	// If the code was mirrored: swap the bottom-left and the top-right points.
	if metadata, ok := decoderResult.GetOther().(*decoder.QRCodeDecoderMetaData); ok {
		metadata.ApplyMirroredCorrection(points)
	}

	result := gozxing.NewResult(decoderResult.GetText(), decoderResult.GetRawBytes(), points,
		gozxing.BarcodeFormat_QR_CODE)

	byteSegments := decoderResult.GetByteSegments()
	if byteSegments != nil {
		result.PutMetadata(gozxing.ResultMetadataType_BYTE_SEGMENTS, byteSegments)
	}

	ecLevel := decoderResult.GetECLevel()
	if ecLevel != "" {
		result.PutMetadata(gozxing.ResultMetadataType_ERROR_CORRECTION_LEVEL, ecLevel)
	}

	if decoderResult.HasStructuredAppend() {
		result.PutMetadata(gozxing.ResultMetadataType_STRUCTURED_APPEND_SEQUENCE,
			decoderResult.GetStructuredAppendSequenceNumber())
		result.PutMetadata(gozxing.ResultMetadataType_STRUCTURED_APPEND_PARITY,
			decoderResult.GetStructuredAppendParity())
	}

	return result, nil
}
