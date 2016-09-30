package viterbi

import "errors"

func UintToBin(x uint, width uint) (bin []byte, err error) {
	bin = make([]byte, width)
	bin[0] = byte(x & 1)
	for i := uint(1); i < width; i++ {
		x >>= 1
		bin[i] = byte(x & 1)
	}

	if x != 0 {
		err = errors.New("width too small")
	}

	return bin, err
}

func BinToUint(bin []byte) (x uint, err error) {
	if len(bin) > 31 {
		err = errors.New("overflowed")
		return x, err
	}

	for i := 0; i < len(bin); i++ {
		x += uint(bin[i]) << uint(i)
	}

	return x, err
}

func Correlation(x, y []byte) (z byte, err error) {
	if len(x) != len(y) {
		err = errors.New("unequal input byte arrays")
	}

	for i := 0; i < len(x); i++ {
		z ^= x[i] & y[i]
	}

	return z, err
}
