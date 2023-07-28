package palettes

// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"github.com/gio-eui/md3-colors/hct"
	"math"
)

type TonalPalette struct {
	cache    map[int]int
	keyColor *hct.Hct
	hue      float64
	chroma   float64
}

// NewTonalPaletteFromInt creates a TonalPalette from an ARGB color.
// for example, NewTonalPaletteFromInt(0xFF000000) will return a TonalPalette with black tones.
// NewTonalPaletteFromInt(0xFFFF0000) will return a TonalPalette with red tones.
func NewTonalPaletteFromInt(argb int) *TonalPalette {
	return NewTonalPaletteFromHct(hct.NewHctFromInt(argb))
}

// NewTonalPaletteFromHct creates a TonalPalette from an Hct.
func NewTonalPaletteFromHct(hct *hct.Hct) *TonalPalette {
	return NewTonalPaletteFromHueChroma(hct.GetHue(), hct.GetChroma())
}

// NewTonalPaletteFromHueChroma creates a TonalPalette from a hue and chroma.
func NewTonalPaletteFromHueChroma(hue, chroma float64) *TonalPalette {
	return &TonalPalette{
		cache:    make(map[int]int),
		keyColor: createKeyColor(hue, chroma),
		hue:      hue,
		chroma:   chroma,
	}
}

// createKeyColor creates the key color of the TonalPalette.
func createKeyColor(hue, chroma float64) *hct.Hct {
	startTone := 50.0
	smallestDeltaHct := hct.NewHct(hue, chroma, startTone)
	smallestDelta := math.Abs(smallestDeltaHct.GetChroma() - chroma)

	for delta := 1.0; delta < 50.0; delta += 1.0 {
		if math.Round(chroma) == math.Round(smallestDeltaHct.GetChroma()) {
			return smallestDeltaHct
		}

		hctAdd := hct.NewHct(hue, chroma, startTone+delta)
		hctAddDelta := math.Abs(hctAdd.GetChroma() - chroma)
		if hctAddDelta < smallestDelta {
			smallestDelta = hctAddDelta
			smallestDeltaHct = hctAdd
		}

		hctSubtract := hct.NewHct(hue, chroma, startTone-delta)
		hctSubtractDelta := math.Abs(hctSubtract.GetChroma() - chroma)
		if hctSubtractDelta < smallestDelta {
			smallestDelta = hctSubtractDelta
			smallestDeltaHct = hctSubtract
		}
	}

	return smallestDeltaHct
}

// Tone returns an ARGB color with the HCT hue and chroma of the TonalPalette and the provided tone.
func (tp *TonalPalette) Tone(tone int) int {
	color, ok := tp.cache[tone]
	if !ok {
		color = hct.NewHct(tp.hue, tp.chroma, float64(tone)).ToInt()
		tp.cache[tone] = color
	}
	return color
}

// GetHct returns the HCT color with the specified tone.
func (tp *TonalPalette) GetHct(tone float64) *hct.Hct {
	return hct.NewHct(tp.hue, tp.chroma, tone)
}

// GetChroma returns the chroma of the TonalPalette.
func (tp *TonalPalette) GetChroma() float64 {
	return tp.chroma
}

// GetHue returns the hue of the TonalPalette.
func (tp *TonalPalette) GetHue() float64 {
	return tp.hue
}

// GetKeyColor returns the key color of the TonalPalette.
func (tp *TonalPalette) GetKeyColor() *hct.Hct {
	return tp.keyColor
}
