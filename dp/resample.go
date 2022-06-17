package dp

import "C"

type AudioResampler struct {
 	inRate   int
	outRate  int
}

func (rs *AudioResampler) Init(channels int, inRate int, outRate int, quality int) error {

	return nil
}

func (rs *AudioResampler) ProcessInt(chanIndex int, bufIn []byte) (bufOut []byte, err error) {

	return bufIn,nil
}

func (rs *AudioResampler) Destroy() error {

	return nil
}
