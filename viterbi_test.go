package viterbi

import (
	"testing"
	"reflect"
)

func TestEncoder(t *testing.T) {
	t.Log("Test Encoder")
	enc := NewEncoder([]uint{0133, 0165, 0171}, 7, true)
	input := []byte{1, 0, 1, 1, 0, 1, 1, 0, 1, 1}
	output, err := enc.Encode(input)
	expectOutput := []byte{1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0}
	if err == nil {
		if !reflect.DeepEqual(output, expectOutput) {
			t.Errorf("Expect %v but got %v instead", expectOutput, output)
		}
	} else {
		t.Errorf("Something wrong in Encoder")
	}
}

func TestDecoder(t *testing.T) {
	output := []byte{1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0}
	chanout := make([]float64, len(output))
	for i := 0; i < len(output); i++ {
		chanout[i] = 1.0 - 2.0 * float64(output[i])
	}
	dec := NewDecoder([]uint{0133, 0165, 0171}, 7, true)
	decoded, err := dec.Decode(chanout)
	expectDecoded := []byte{1, 0, 1, 1, 0, 1, 1, 0, 1, 1}
	if err == nil {
		if !reflect.DeepEqual(decoded, expectDecoded) {
			t.Errorf("Expect %v but got %v instead", expectDecoded, decoded)
		}
	} else {
		t.Errorf("Something wrong in Decoder")
	}
}

