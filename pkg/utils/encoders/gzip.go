package encoders

import (
	"bytes"
	"compress/gzip"
	"io"
)

type GzipEncoder struct{}

func (e GzipEncoder) Encode(data []byte) ([]byte, error) {
	return GzipBuf(data)
}

func (e GzipEncoder) Decode(data []byte) ([]byte, error) {
	result := GunzipBuf(data)
	return result, nil
}

func GzipBuf(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	if _, err := writer.Write(data); err != nil {
		writer.Close()
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func GunzipBuf(data []byte) []byte {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return data
	}
	defer reader.Close()

	result, err := io.ReadAll(reader)
	if err != nil {
		return data
	}

	return result
}
