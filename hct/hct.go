package hct

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
	colorUtils "github.com/gio-eui/md3-colors/utils/color"
)

// A color system built using CAM16 hue and chroma, and L* from L*a*b*.
//
// Using L* creates a link between the color system, contrast, and thus accessibility. Contrast
// ratio depends on relative luminance, or Y in the XYZ color space. L*, or perceptual luminance can
// be calculated from Y.
//
// Unlike Y, L* is linear to human perception, allowing trivial creation of accurate color tones.
//
// Unlike contrast ratio, measuring contrast in L* is linear, and simple to calculate. A
// difference of 40 in HCT tone guarantees a contrast ratio >= 3.0, and a difference of 50
// guarantees a contrast ratio >= 4.5.
//

// Hct hue, chroma, and tone. A color system that provides a perceptually accurate color
// measurement system that can also accurately render what colors will appear as in different
// lighting environments.
type Hct struct {
	hue    float64
	chroma float64
	tone   float64
	argb   int
}

// NewHct creates an HCT color from hue, chroma, and tone.
//
// 0 <= [hue] < 360; invalid values are corrected.
// 0 <= [chroma] <= ?; Informally, colorfulness. The color returned may be lower than
// the requested chroma. Chroma has a different maximum for any given hue and tone.
// 0 <= [tone] <= 100; informally, lightness. Invalid values are corrected.
func NewHct(hue, chroma, tone float64) *Hct {
	argb := solveToInt(hue, chroma, tone)
	return &Hct{
		hue:    hue,
		chroma: chroma,
		tone:   tone,
		argb:   argb,
	}
}

// NewHctFromInt creates an HCT color from an ARGB color representation.
//
// [argb] ARGB representation of a color.
func NewHctFromInt(argb int) *Hct {
	hct := &Hct{}
	hct.setInternalState(argb)
	return hct
}

func (h *Hct) setInternalState(argb int) {
	cam := Cam16FromInt(argb)
	h.hue = cam.GetHue()
	h.chroma = cam.GetChroma()
	h.tone = colorUtils.LstarFromArgb(argb)
	h.argb = argb
}

// GetHue returns the hue component of the HCT color.
func (h *Hct) GetHue() float64 {
	return h.hue
}

// GetChroma returns the chroma component of the HCT color.
func (h *Hct) GetChroma() float64 {
	return h.chroma
}

// GetTone returns the tone component of the HCT color.
func (h *Hct) GetTone() float64 {
	return h.tone
}

// ToInt returns the ARGB representation of the HCT color.
func (h *Hct) ToInt() int {
	return h.argb
}

// SetHue sets the hue of the HCT color.
// Chroma may decrease because chroma has a different maximum for any given hue and tone.
//
// newHue 0 <= newHue < 360; invalid values are corrected.
func (h *Hct) SetHue(newHue float64) {
	h.setInternalState(solveToInt(newHue, h.chroma, h.tone))
}

// SetChroma sets the chroma of the HCT color.
// Chroma may decrease because chroma has a different maximum for any given hue and tone.
//
// newChroma 0 <= newChroma < ?; Informally, colorfulness.
func (h *Hct) SetChroma(newChroma float64) {
	h.setInternalState(solveToInt(h.hue, newChroma, h.tone))
}

// SetTone sets the tone of the HCT color.
// Chroma may decrease because chroma has a different maximum for any given hue and tone.
//
// newTone 0 <= newTone <= 100; invalid valids are corrected.
func (h *Hct) SetTone(newTone float64) {
	h.setInternalState(solveToInt(h.hue, h.chroma, newTone))
}

// InViewingConditions translates the color into different viewing conditions.
//
// Colors change appearance. They look different with lights on versus off, the same color, as
// in hex code, on white looks different when on black. This is called color relativity, most
// famously explicated by Josef Albers in Interaction of Color.
//
// In color science, color appearance models can account for this and calculate the appearance
// of a color in different settings. HCT is based on CAM16, a color appearance model, and uses it
// to make these calculations.
//
// See MakeViewingConditions for parameters affecting color appearance.
func (h *Hct) InViewingConditions(vc ViewingConditions) *Hct {
	// 1. Use CAM16 to find XYZ coordinates of color in specified VC.
	c16 := Cam16FromInt(h.ToInt())
	viewedInVc := c16.XyzInViewingConditions(vc, nil)

	// 2. Create CAM16 of those XYZ coordinates in default VC.
	recastInVc := Cam16FromXyzInViewingConditions(viewedInVc[0], viewedInVc[1], viewedInVc[2], DefaultViewingConditions)

	// 3. Create HCT from:
	// - CAM16 using default VC with XYZ coordinates in specified VC.
	// - L* converted from Y in XYZ coordinates in specified VC.
	return NewHct(recastInVc.GetHue(), recastInVc.GetChroma(), colorUtils.LstarFromY(viewedInVc[1]))
}
