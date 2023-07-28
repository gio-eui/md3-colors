package colorUtils

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
	mathUtils "github.com/gio-eui/md3-colors/utils/math"
	"math"
)

var srgbToXyz = [][]float64{
	{0.41233895, 0.35762064, 0.18051042},
	{0.2126, 0.7152, 0.0722},
	{0.01932141, 0.11916382, 0.95034478},
}

var xyzToSrgb = [][]float64{
	{3.2413774792388685, -1.5376652402851851, -0.49885366846268053},
	{-0.9691452513005321, 1.8758853451067872, 0.04156585616912061},
	{0.05562093689691305, -0.20395524564742123, 1.0571799111220335},
}

var whitePointD65 = []float64{95.047, 100.0, 108.883}

// ArgbFromRgb converts a color from RGB components to ARGB format.
func ArgbFromRgb(red, green, blue int) int {
	return (255 << 24) | ((red & 255) << 16) | ((green & 255) << 8) | (blue & 255)
}

// ArgbFromLinrgb converts a color from linear RGB components to ARGB format
func ArgbFromLinrgb(linrgb []float64) int {
	r := Delinearized(linrgb[0])
	g := Delinearized(linrgb[1])
	b := Delinearized(linrgb[2])
	return ArgbFromRgb(r, g, b)
}

// AlphaFromArgb returns the alpha component of a color in ARGB format
func AlphaFromArgb(argb int) int {
	return (argb >> 24) & 255
}

// RedFromArgb returns the red component of a color in ARGB format
func RedFromArgb(argb int) int {
	return (argb >> 16) & 255
}

// GreenFromArgb returns the green component of a color in ARGB format
func GreenFromArgb(argb int) int {
	return (argb >> 8) & 255
}

// BlueFromArgb returns the blue component of a color in ARGB format
func BlueFromArgb(argb int) int {
	return argb & 255
}

// IsOpaque returns whether a color in ARGB format is opaque
func IsOpaque(argb int) bool {
	return AlphaFromArgb(argb) >= 255
}

// ArgbFromXyz converts a color from XYZ components to ARGB format
func ArgbFromXyz(x, y, z float64) int {
	matrix := xyzToSrgb
	linearR := matrix[0][0]*x + matrix[0][1]*y + matrix[0][2]*z
	linearG := matrix[1][0]*x + matrix[1][1]*y + matrix[1][2]*z
	linearB := matrix[2][0]*x + matrix[2][1]*y + matrix[2][2]*z
	r := Delinearized(linearR)
	g := Delinearized(linearG)
	b := Delinearized(linearB)
	return ArgbFromRgb(r, g, b)
}

// XyzFromArgb converts a color from ARGB format to XYZ components
func XyzFromArgb(argb int) []float64 {
	r := Linearized(RedFromArgb(argb))
	g := Linearized(GreenFromArgb(argb))
	b := Linearized(BlueFromArgb(argb))
	row := []float64{r, g, b}
	return mathUtils.MatrixMultiply(row, srgbToXyz)
}

// ArgbFromLab converts a color represented in Lab color space into an ARGB integer
func ArgbFromLab(l, a, b float64) int {
	whitePoint := whitePointD65
	fy := (l + 16.0) / 116.0
	fx := a/500.0 + fy
	fz := fy - b/200.0
	xNormalized := labInvf(fx)
	yNormalized := labInvf(fy)
	zNormalized := labInvf(fz)
	x := xNormalized * whitePoint[0]
	y := yNormalized * whitePoint[1]
	z := zNormalized * whitePoint[2]
	return ArgbFromXyz(x, y, z)
}

// LabFromArgb converts a color from ARGB representation to L*a*b*  representation.
//
// [argb] the ARGB representation of a color
// Returns a Lab object representing the color
func LabFromArgb(argb int) []float64 {
	linearR := Linearized(RedFromArgb(argb))
	linearG := Linearized(GreenFromArgb(argb))
	linearB := Linearized(BlueFromArgb(argb))
	matrix := srgbToXyz
	x := matrix[0][0]*linearR + matrix[0][1]*linearG + matrix[0][2]*linearB
	y := matrix[1][0]*linearR + matrix[1][1]*linearG + matrix[1][2]*linearB
	z := matrix[2][0]*linearR + matrix[2][1]*linearG + matrix[2][2]*linearB
	whitePoint := whitePointD65
	xNormalized := x / whitePoint[0]
	yNormalized := y / whitePoint[1]
	zNormalized := z / whitePoint[2]
	fx := labF(xNormalized)
	fy := labF(yNormalized)
	fz := labF(zNormalized)
	l := 116.0*fy - 16
	a := 500.0 * (fx - fy)
	b := 200.0 * (fy - fz)
	return []float64{l, a, b}
}

// ArgbFromLstar converts an L* value to an ARGB representation.
//
// [lstar] L* in L*a*b*
// Returns ARGB representation of grayscale color with lightness
// matching L*
func ArgbFromLstar(lstar float64) int {
	var x, y, z, fx, fy, fz float64
	fy = (lstar + 16.0) / 116.0
	fz = fy
	fx = fy
	kappa := 24389.0 / 27.0
	epsilon := 216.0 / 24389.0
	lExceedsEpsilonKappa := lstar > 8.0
	if lExceedsEpsilonKappa {
		y = fy * fy * fy
	} else {
		y = lstar / kappa
	}
	cubeExceedEpsilon := fy*fy*fy > epsilon
	if cubeExceedEpsilon {
		x = fx * fx * fx
		z = fz * fz * fz
	} else {
		x = lstar / kappa
		z = lstar / kappa
	}
	whitePoint := whitePointD65
	return ArgbFromXyz(x*whitePoint[0], y*whitePoint[1], z*whitePoint[2])
}

// LstarFromArgb computes the L* value of a color in ARGB representation.
//
// [argb] ARGB representation of a color
// Returns L*, from L*a*b*, coordinate of the color
func LstarFromArgb(argb int) float64 {
	xyz := XyzFromArgb(argb)
	y := xyz[1] / 100.0
	e := 216.0 / 24389.0
	if y <= e {
		return 24389.0 / 27.0 * y
	} else {
		yIntermediate := math.Pow(y, 1.0/3.0)
		return 116.0*yIntermediate - 16.0
	}
}

// YFromLstar converts an L* value to a Y value.
//
// L* in L*a*b* and Y in XYZ measure the same quantity, luminance.
//
// L* measures perceptual luminance, a linear scale. Y in XYZ
// measures relative luminance, a logarithmic scale.
//
// [lstar] L* in L*a*b*
// Returns Y in XYZ
func YFromLstar(lstar float64) float64 {
	ke := 8.0
	if lstar > ke {
		return math.Pow((lstar+16.0)/116.0, 3.0) * 100.0
	} else {
		return lstar / (24389.0 / 27.0) * 100.0
	}
}

// LstarFromY converts a Y value to an L* value.
//
// L* in L*a*b* and Y in XYZ measure the same quantity, luminance.
//
// L* measures perceptual luminance, a linear scale. Y in XYZ
// measures relative luminance, a logarithmic scale.
//
// [y] Y in XYZ
// Returns L* in L*a*b*
func LstarFromY(y float64) float64 {
	return labF(y/100.0)*116.0 - 16.0
}

// Linearized linearizes an RGB component
//
// [rgbComponent] 0 <= rgb_component <= 255, represents R/G/B channel
// Returns 0.0 <= output <= 100.0, color channel converted to linear RGB space
func Linearized(rgbComponent int) float64 {
	normalized := float64(rgbComponent) / 255.0
	if normalized <= 0.040449936 {
		return normalized / 12.92 * 100.0
	} else {
		return math.Pow((normalized+0.055)/1.055, 2.4) * 100.0
	}
}

// Delinearized an RGB component.
//
// [rgbComponent] 0.0 <= rgb_component <= 100.0, represents linear R/G/B channel
// Returns 0 <= output <= 255, color channel converted to regular RGB space
func Delinearized(rgbComponent float64) int {
	normalized := rgbComponent / 100.0
	var delinearized float64
	if normalized <= 0.0031308 {
		delinearized = normalized * 12.92
	} else {
		delinearized = 1.055*math.Pow(normalized, 1.0/2.4) - 0.055
	}
	return mathUtils.ClampInt(0, 255, int(math.Round(delinearized*255.0)))
}

// WhitePointD65 returns the standard white point; white on a sunny day
func WhitePointD65() []float64 {
	return whitePointD65
}

func labF(t float64) float64 {
	e := 216.0 / 24389.0
	kappa := 24389.0 / 27.0
	if t > e {
		return math.Pow(t, 1.0/3.0)
	} else {
		return (kappa*t + 16) / 116
	}
}

func labInvf(ft float64) float64 {
	e := 216.0 / 24389.0
	kappa := 24389.0 / 27.0
	ft3 := ft * ft * ft
	if ft3 > e {
		return ft3
	} else {
		return (116*ft - 16) / kappa
	}
}
