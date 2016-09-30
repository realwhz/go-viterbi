package viterbi

import "errors"

type Encoder struct {
	GenPolys []uint
	ConsLen int
	NumOut int
	Tailbiting bool
	State uint
}

func NewEncoder(GenPolys []uint, ConsLen int, Tailbiting bool) *Encoder {
	enc := new(Encoder)
	enc.GenPolys = GenPolys
	enc.ConsLen = ConsLen
	enc.NumOut = len(GenPolys)
	enc.Tailbiting = Tailbiting
	return enc
}

func (enc *Encoder) Encode(input []byte) (output []byte, err error) {
	if enc.Tailbiting {
		tmpbuf := make([]byte, enc.ConsLen-1)
		if len(input) < enc.ConsLen-1 {
			err = errors.New("input too short for tailbiting encoder")
			for i := range(input) {
				tmpbuf[len(tmpbuf)-1-i] = input[len(input)-1-i]
			}
		}

		for i := range(tmpbuf) {
			tmpbuf[len(tmpbuf)-1-i] = input[len(input)-1-i]
		}

		enc.State, _ = BinToUint(tmpbuf)
	}

	output = make([]byte, len(input)*enc.NumOut)
	for i := range(input) {
		for j := 0; j < enc.NumOut; j++ {
			coeffBin, _ := UintToBin(enc.GenPolys[j], uint(enc.ConsLen))
			enc.State |= uint(input[i] << uint(enc.ConsLen-1))
			stateBin, _ := UintToBin(enc.State, uint(enc.ConsLen))
			output[i*enc.NumOut+j], _ = Correlation(stateBin, coeffBin)
		}
		enc.State >>= 1
	}

	return output, err
}
