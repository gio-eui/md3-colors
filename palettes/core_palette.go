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
	hct2 "github.com/gio-eui/md3-colors/hct"
	"math"
)

// CorePalette represents a collection of TonalPalettes.
type CorePalette struct {
	A1    *TonalPalette
	A2    *TonalPalette
	A3    *TonalPalette
	N1    *TonalPalette
	N2    *TonalPalette
	Error *TonalPalette
}

// NewCorePaletteFromInt creates key tones from an ARGB color.
// for example, NewCorePaletteFromInt(0xFF000000) will return a core palette with black tones.
// NewCorePaletteFromInt(0xFFFF0000) will return a core palette with red tones.
func NewCorePaletteFromInt(argb int) *CorePalette {
	return newCorePalette(argb, false)
}

// NewContentCorePaletteFromInt creates content key tones from an ARGB color.
// for example, NewContentCorePaletteFromInt(0xFF000000) will return a content core palette with black tones.
// NewContentCorePaletteFromInt(0xFFFF0000) will return a content core palette with red tones.
func NewContentCorePaletteFromInt(argb int) *CorePalette {
	return newCorePalette(argb, true)
}

// newCorePalette creates a new CorePalette.
func newCorePalette(argb int, isContent bool) *CorePalette {
	hct := hct2.Cam16FromInt(argb)
	hue := hct.GetHue()
	chroma := hct.GetChroma()

	corePalette := &CorePalette{
		A1:    NewTonalPaletteFromHueChroma(hue, chroma),
		A2:    NewTonalPaletteFromHueChroma(hue, chroma/3.0),
		A3:    NewTonalPaletteFromHueChroma(hue+60.0, chroma/2.0),
		N1:    NewTonalPaletteFromHueChroma(hue, math.Min(chroma/12.0, 4.0)),
		N2:    NewTonalPaletteFromHueChroma(hue, math.Min(chroma/6.0, 8.0)),
		Error: NewTonalPaletteFromHueChroma(25.0, 84.0),
	}

	if !isContent {
		corePalette.A1 = NewTonalPaletteFromHueChroma(hue, math.Max(48.0, chroma))
		corePalette.A2 = NewTonalPaletteFromHueChroma(hue, 16.0)
		corePalette.A3 = NewTonalPaletteFromHueChroma(hue+60.0, 24.0)
		corePalette.N1 = NewTonalPaletteFromHueChroma(hue, 4.0)
		corePalette.N2 = NewTonalPaletteFromHueChroma(hue, 8.0)
	}

	return corePalette
}
