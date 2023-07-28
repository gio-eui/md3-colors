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
	mathUtils "github.com/gio-eui/md3-colors/utils/math"
	"math"
)

// In traditional color spaces, a color can be identified solely by the observer's measurement of
// the color. Color appearance models such as CAM16 also use information about the environment where
// the color was observed, known as the viewing conditions.
//
// For example, white under the traditional assumption of a midday sun white point is accurately
// measured as a slightly chromatic blue by CAM16. (roughly, hue 203, chroma 3, lightness 100)
//
// This class caches intermediate values of the CAM16 conversion process that depend only on
// viewing conditions, enabling speed ups.

type ViewingConditions struct {
	Aw     float64
	Nbb    float64
	Ncb    float64
	C      float64
	Nc     float64
	N      float64
	RgbD   []float64
	Fl     float64
	FlRoot float64
	Z      float64
}

// DefaultViewingConditions represents the default sRGB-like viewing conditions.
var DefaultViewingConditions = DefaultViewingConditionsWithBackgroundLstar(50.0)

// MakeViewingConditions create a ViewingConditions from a simple, physically relevant, set of parameters.
//
// Parameters affecting color appearance include:
// [whitePoint]: coordinates of white in XYZ color space.
// [adaptingLuminance]: light strength, in lux.
// [backgroundLstar]: average luminance of 10 degrees around color.
// [surround]: brightness of the entire environment.
// [discountingIlluminant]: whether eyes have adjusted to lighting.
func MakeViewingConditions(whitePoint []float64, adaptingLuminance, backgroundLstar, surround float64, discountingIlluminant bool) ViewingConditions {
	backgroundLstar = math.Max(0.1, backgroundLstar)
	matrix := XYZToCam16RGB
	xyz := whitePoint
	rW := (xyz[0] * matrix[0][0]) + (xyz[1] * matrix[0][1]) + (xyz[2] * matrix[0][2])
	gW := (xyz[0] * matrix[1][0]) + (xyz[1] * matrix[1][1]) + (xyz[2] * matrix[1][2])
	bW := (xyz[0] * matrix[2][0]) + (xyz[1] * matrix[2][1]) + (xyz[2] * matrix[2][2])
	f := 0.8 + (surround / 10.0)
	c := 0.0
	d := 0.0

	if f >= 0.9 {
		c = mathUtils.Lerp(0.59, 0.69, (f-0.9)*10.0)
	} else {
		c = mathUtils.Lerp(0.525, 0.59, (f-0.8)*10.0)
	}

	if discountingIlluminant {
		d = 1.0
	} else {
		d = f * (1.0 - ((1.0 / 3.6) * math.Exp((-adaptingLuminance-42.0)/92.0)))
		d = math.Min(math.Max(0.0, d), 1.0)
	}

	nc := f
	rgbD := []float64{
		d*(100.0/rW) + 1.0 - d,
		d*(100.0/gW) + 1.0 - d,
		d*(100.0/bW) + 1.0 - d,
	}

	k := 1.0 / (5.0*adaptingLuminance + 1.0)
	k4 := k * k * k * k
	k4F := 1.0 - k4
	fl := (k4 * adaptingLuminance) + (0.1 * k4F * k4F * math.Cbrt(5.0*adaptingLuminance))
	n := (colorUtils.YFromLstar(backgroundLstar) / whitePoint[1])
	z := 1.48 + math.Sqrt(n)
	nbb := 0.725 / math.Pow(n, 0.2)
	ncb := nbb
	rgbAFactors := []float64{
		math.Pow(fl*rgbD[0]*rW/100.0, 0.42),
		math.Pow(fl*rgbD[1]*gW/100.0, 0.42),
		math.Pow(fl*rgbD[2]*bW/100.0, 0.42),
	}
	rgbA := []float64{
		(400.0 * rgbAFactors[0]) / (rgbAFactors[0] + 27.13),
		(400.0 * rgbAFactors[1]) / (rgbAFactors[1] + 27.13),
		(400.0 * rgbAFactors[2]) / (rgbAFactors[2] + 27.13),
	}
	aw := ((2.0 * rgbA[0]) + rgbA[1] + (0.05 * rgbA[2])) * nbb

	return ViewingConditions{
		Aw:     aw,
		Nbb:    nbb,
		Ncb:    ncb,
		C:      c,
		Nc:     nc,
		N:      n,
		RgbD:   rgbD,
		Fl:     fl,
		FlRoot: math.Pow(fl, 0.25),
		Z:      z,
	}
}

// DefaultViewingConditionsWithBackgroundLstar sRGB-like viewing conditions with a custom background lstar.
//
// Default viewing conditions have a lstar of 50, midgray.
func DefaultViewingConditionsWithBackgroundLstar(backgroundLstar float64) ViewingConditions {
	return MakeViewingConditions(colorUtils.WhitePointD65(), (200.0/math.Pi*colorUtils.YFromLstar(50.0))/100.0, backgroundLstar, 2.0, false)
}

// GetAw returns the value of aw.
func (vc ViewingConditions) GetAw() float64 {
	return vc.Aw
}

// GetN returns the value of n.
func (vc ViewingConditions) GetN() float64 {
	return vc.N
}

// GetNbb returns the value of nbb.
func (vc ViewingConditions) GetNbb() float64 {
	return vc.Nbb
}

// GetNcb returns the value of ncb.
func (vc ViewingConditions) GetNcb() float64 {
	return vc.Ncb
}

// GetC returns the value of c.
func (vc ViewingConditions) GetC() float64 {
	return vc.C
}

// GetNc returns the value of nc.
func (vc ViewingConditions) GetNc() float64 {
	return vc.Nc
}

// GetRgbD returns the value of rgbD.
func (vc ViewingConditions) GetRgbD() []float64 {
	return vc.RgbD
}

// GetFl returns the value of fl.
func (vc ViewingConditions) GetFl() float64 {
	return vc.Fl
}

// GetFlRoot returns the value of flRoot.
func (vc ViewingConditions) GetFlRoot() float64 {
	return vc.FlRoot
}

// GetZ returns the value of z.
func (vc ViewingConditions) GetZ() float64 {
	return vc.Z
}
