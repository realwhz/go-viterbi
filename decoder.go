package viterbi

import "errors"

type Decoder struct {
	GenPolys []uint
	ConsLen int
	NumOut int
	Tailbiting bool
}

func NewDecoder(GenPolys []uint, ConsLen int, Tailbiting bool) *Decoder {
	dec := new(Decoder)
	dec.GenPolys = GenPolys
	dec.ConsLen = ConsLen
	dec.NumOut = len(GenPolys)
	dec.Tailbiting = Tailbiting
	return dec
}

func (dec *Decoder) Decode(input []float64) (output []byte, err error) {
	numStates := 1 << uint(dec.ConsLen-1)
	initMetrics := make([]float64, numStates)
	if dec.Tailbiting {
		return dec.DecodeWAVA(input, initMetrics)
	}
	for i := 1; i < len(initMetrics); i++ {
		initMetrics[i] = 100.0
	}
	return dec.DecodeWAVA(input, initMetrics)
}

func (dec *Decoder) DecodeWAVA(input []float64, initMetrics []float64) (output []byte, err error) {
	resMetrics, survPaths := dec.ViterbiTrial(input, initMetrics)
	minSM := resMetrics[0]
	minIndex := uint(0)
	for i := 1; i < len(resMetrics); i++ {
		if resMetrics[i] < minSM {
			minSM = resMetrics[i]
			minIndex = uint(i)
		}
	}
	var bstate uint
	output, bstate = dec.TraceBack(survPaths, minIndex)
	if minIndex == bstate {
		return
	}

	minIndex = bstate
	output, bstate = dec.TraceBack(survPaths, minIndex)
	if minIndex == bstate {
		err = errors.New("didn't find tailbiting path")
	}
	return
}

func (dec *Decoder) ViterbiTrial(input, initMetrics []float64) (resMetrics []float64, survPaths [][]uint) {
	numStates := 1 << uint(dec.ConsLen-1)
	outLen := len(input) / dec.NumOut
	survPaths = make([][]uint, numStates)
	for i := range(survPaths) {
		survPaths[i] = make([]uint, outLen)
	}

	pathMetrics := make([]float64, numStates)
	tmpPathMetrics := make([]float64, numStates)
	copy(pathMetrics, initMetrics)

	for i := 0; i < outLen; i++ {
		seg := input[i*dec.NumOut:(i+1)*dec.NumOut]
		for s := uint(0); s < uint(numStates); s++ {
			prev0 := (s << 1) & ^(1 << uint(dec.ConsLen-1))
			prev1 := prev0 + 1
			y := make([]byte, dec.NumOut)
			for j := 0; j < dec.NumOut; j++ {
				coeffBin, _ := UintToBin(dec.GenPolys[j], uint(dec.ConsLen))
				stateBin, _ := UintToBin(s << 1, uint(dec.ConsLen))
				y[j], _ = Correlation(stateBin, coeffBin)
			}
			bm0 := 0.0
			for j := 0; j < dec.NumOut; j++ {
				bm0 += -(1.0-2.0*float64(y[j]))*seg[j]
			}
			bm1 := -bm0
			if pathMetrics[prev0]+bm0 < pathMetrics[prev1]+bm1 {
				tmpPathMetrics[s] = pathMetrics[prev0]+bm0
				survPaths[s][i] = prev0
			} else {
				tmpPathMetrics[s] = pathMetrics[prev1]+bm1
				survPaths[s][i] = prev1
			}
		}
		copy(pathMetrics, tmpPathMetrics)
	}
	resMetrics = make([]float64, numStates)
	copy(resMetrics, pathMetrics)
	return
}

func (dec *Decoder) TraceBack(survPaths [][]uint, survState uint) (output []byte, bstate uint) {
	outLen := len(survPaths[0])
	output = make([]byte, outLen)
	var prev uint
	for i := 0; i < outLen; i++ {
		output[outLen-1-i] = byte((survState & (1 << uint(dec.ConsLen-2))) >> uint(dec.ConsLen-2))
		prev = survPaths[survState][outLen-1-i]
		survState = prev
	}
	bstate = survState
	return output, bstate
}
