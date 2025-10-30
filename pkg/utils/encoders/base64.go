package encoders

import "encoding/base64"

type Base64 struct{}

var base64Alphabet = "a0b2c5def6hijklmnopqr_st-uvwxyzA1B3C4DEFGHIJKLM7NO9PQR8ST+UVWXYZ"
var Base64Encoding = base64.NewEncoding(base64Alphabet).WithPadding(base64.NoPadding)

func (e Base64) Encode(data []byte) ([]byte, error) {
	return []byte(Base64Encoding.EncodeToString(data)), nil
}

func (e Base64) Decode(data []byte) ([]byte, error) {
	return Base64Encoding.DecodeString(string(data))
}
