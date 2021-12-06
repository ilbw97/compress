package compress

import (
	"bytes"
	"compress/gzip"

	"github.com/google/brotli/go/cbrotli"
	"github.com/sirupsen/logrus"
)

func CompressBr(data []byte) []byte {
	var b bytes.Buffer

	br := cbrotli.NewWriter(&b, cbrotli.WriterOptions{Quality: 5, LGWin: 18})
	if n, err := br.Write(data); err != nil || n < 1 {
		logrus.Errorf("Compress by brotli err : %v, Write %v", err, n)
		return nil
	}

	if err := br.Close(); err != nil {
		logrus.Errorf("Close to compress by brotli err : %v", err)
	}

	if !checkCompressedBrData(b.Bytes(), data) {
		return nil
	}

	logrus.Info("### Compress body by Brotli SUCCESSFULY!!")
	return b.Bytes()
}

func checkCompressedBrData(compressedData, wantOriginalData []byte) bool {
	uncompressed, err := cbrotli.Decode(compressedData)
	if err != nil {
		logrus.Errorf("brotli decompress failed: %v", err)
	}
	if !bytes.Equal(uncompressed, wantOriginalData) {
		if len(wantOriginalData) != len(uncompressed) {
			logrus.Errorf(""+
				"Data doesn't uncompress to the original value.\n"+
				"Length of original: %v\n"+
				"Length of uncompressed: %v",
				len(wantOriginalData), len(uncompressed))
			return false
		}
		for i := range wantOriginalData {
			if wantOriginalData[i] != uncompressed[i] {
				logrus.Errorf(""+
					"Data doesn't uncompress to the original value.\n"+
					"Original at %v is %v\n"+
					"Uncompressed at %v is %v",
					i, wantOriginalData[i], i, uncompressed[i])
				return false
			}
		}
	}
	return true
}

func CompressGzip(data []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	if n, err := gz.Write(data); err != nil || n < 1 {
		logrus.Errorf("Compress by gzip err : %v, Write %v", err, n)
		return nil
	}

	if err := gz.Close(); err != nil {
		logrus.Errorf("Close to compress by gzip err : %v", err)
	}

	// logrus.Info("### Compress body by gzip SUCCESSFULY!!")
	return b.Bytes()
}
