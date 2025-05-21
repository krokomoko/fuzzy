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
	t                 int // WordType
}

func (w *Word) Mu(value float64) (r float64, err error) {
	if value < w.min && w.t != TrapezeLeft ||
		value > w.max && w.t != TrapezeRight {
		return
	}

	switch w.t {
	case TrapezeLeft:
		if value < w.middle {
			r = 1.0
		} else {
			r = 1.0 + w.kRight*(value-w.middle)
		}
	case TrapezeRight:
		if value >= w.middle {
			r = 1.0
		} else {
			r = w.kLeft * (value - w.min)
		}
	case Triangle:
		if value < w.middle {
			r = w.kLeft * (value - w.min)
		} else if value > w.middle {
			r = 1.0 + w.kRight*(value-w.middle)
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
			min:    _min,
			max:    _min + vPerWord,
			kLeft:  0.0,
			kRight: 0.0,
			middle: _middle,
			cM:     _middle,
			t:      Triangle,
		}
	}

	words[0].t = TrapezeLeft
	words[wordsCount-1].t = TrapezeRight

	for i := 0; i < wordsCount; i++ {
		if i > 0 {
			words[i].min = words[i-1].middle
			words[i].kLeft = 1.0 / (words[i].middle - words[i].min)
		}
		if i < wordsCount-1 {
			words[i].max = words[i+1].middle
			words[i].kRight = -1.0 / (words[i].max - words[i].middle)
		}
	}

	a := words[0].middle - words[0].min
	b := words[0].max - words[0].min
	words[0].cM = (b*words[0].min + a*words[0].max) / (a + b)

	lastWordInd := wordsCount - 1
	a = words[lastWordInd].max - words[lastWordInd].middle
	b = words[lastWordInd].max - words[lastWordInd].min
	words[lastWordInd].cM = (b*words[lastWordInd].max + a*words[lastWordInd].min) / (a + b)

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
