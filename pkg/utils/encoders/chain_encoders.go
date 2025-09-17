package encoders

type IChainEncoder interface {
	Encode(data []byte) ([]byte, error)
	Decode(data []byte) ([]byte, error)
}

type ChainEncoder struct {
	encoders []Encoder
}

func NewChainEncoder(encoders []Encoder) *ChainEncoder {
	return &ChainEncoder{encoders: encoders}
}

func (ce *ChainEncoder) Encode(data []byte) ([]byte, error) {
	result := data
	var err error
	for _, encoder := range ce.encoders {
		result, err = encoder.Encode(result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil

}

func (ce *ChainEncoder) Decode(data []byte) ([]byte, error) {
	result := data
	var err error
	for i := len(ce.encoders) - 1; i >= 0; i-- {
		result, err = ce.encoders[i].Decode(result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
