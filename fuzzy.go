package fuzzy

import (
	"fmt"
	"slices"
)

// WordType
const (
	TrapezeLeft = iota + 1
	Triangle
	TrapezeRight
)

func __sum(data []float64) (sum float64) {
	for _, el := range data {
		sum += el
	}
	return
}

type Word struct {
	Min, Max, Middle  float64
	KLeft, KRight, cM float64
	T                 int // WordType
}

func (w *Word) Mu(value float64) (r float64, err error) {
	if value < w.Min && w.T != TrapezeLeft ||
		value > w.Max && w.T != TrapezeRight {
		return
	}

	switch w.T {
	case TrapezeLeft:
		if value < w.Middle {
			r = 1.0
		} else {
			r = 1.0 + w.KRight*(value-w.Middle)
		}
	case TrapezeRight:
		if value >= w.Middle {
			r = 1.0
		} else {
			r = w.KLeft * (value - w.Min)
		}
	case Triangle:
		if value < w.Middle {
			r = w.KLeft * (value - w.Min)
		} else if value > w.Middle {
			r = 1.0 + w.KRight*(value-w.Middle)
		} else {
			r = 1.0
		}
	default:
		err = fmt.Errorf("calculation of Mu: %v with %v", w, value)
	}

	return
}

type Parameter struct {
	Words []Word
}

func NewParameter(data []float64, wordsCount int) Parameter {
	var (
		min       = slices.Min(data)
		max       = slices.Max(data)
		vPerWord  = (max - min) / float64(wordsCount)
		vPerWord2 = vPerWord / 2.0
		words     = make([]Word, wordsCount)
	)
	var _min, _middle float64

	for wordInd := 0; wordInd < wordsCount; wordInd++ {
		_min = min + float64(wordInd)*vPerWord
		_middle = _min + vPerWord2
		words[wordInd] = Word{
			Min:    _min,
			Max:    _min + vPerWord,
			KLeft:  0.0,
			KRight: 0.0,
			Middle: _middle,
			cM:     _middle,
			T:      Triangle,
		}
	}

	words[0].T = TrapezeLeft
	words[wordsCount-1].T = TrapezeRight

	for i := 0; i < wordsCount; i++ {
		if i > 0 {
			words[i].Min = words[i-1].Middle
			words[i].KLeft = 1.0 / (words[i].Middle - words[i].Min)
		}
		if i < wordsCount-1 {
			words[i].Max = words[i+1].Middle
			words[i].KRight = -1.0 / (words[i].Max - words[i].Middle)
		}
	}

	a := words[0].Middle - words[0].Min
	b := words[0].Max - words[0].Min
	words[0].cM = (b*words[0].Min + a*words[0].Max) / (a + b)

	lastWordInd := wordsCount - 1
	a = words[lastWordInd].Max - words[lastWordInd].Middle
	b = words[lastWordInd].Max - words[lastWordInd].Min
	words[lastWordInd].cM = (b*words[lastWordInd].Max + a*words[lastWordInd].Min) / (a + b)

	return Parameter{
		Words: words,
	}
}

func (p *Parameter) Value(data []float64) (r float64, err error) {
	for i, v := range data {
		r += p.Words[i].cM * v
	}

	sum := __sum(data)
	if sum == 0.0 {
		err = fmt.Errorf("Parameter value() error")
	} else {
		r /= sum
	}

	return
}
