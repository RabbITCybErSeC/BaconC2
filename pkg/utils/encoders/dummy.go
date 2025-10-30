package encoders

type DummyEncoder struct{}

func (e DummyEncoder) Encode(data []byte) ([]byte, error) {
	return data, nil
}

func (e DummyEncoder) Decode(data []byte) ([]byte, error) {
	return data, nil
}
