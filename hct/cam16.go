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

// Cam16 a color appearance model. Colors are not just defined by their hex
// code, but rather, a hex code and viewing conditions.
//
// CAM16 instances also have coordinates in the CAM16-UCS space, called J*, a*,
// b*, or jstar, astar, bstar in code. CAM16-UCS is included in the CAM16
// specification, and should be used when measuring distances between colors.
//
// In traditional color spaces, a color can be identified solely by the
// observer's measurement of the color. Color appearance models such as CAM16
// also use information about the environment where the color was
// observed, known as the viewing conditions.
//
// For example, white under the traditional assumption of a midday sun white
// point is accurately measured as a slightly chromatic blue by CAM16.
// (roughly, hue 203, chroma 3, lightness 100)
// CAM16, a color appearance model. Colors are not just defined by their hex
// code, but rather, a hex code and viewing conditions.
//
// CAM16 instances also have coordinates in the CAM16-UCS space, called J*, a*,
// b*, or jstar, astar, bstar in code. CAM16-UCS is included in the CAM16
// specification, and should be used when measuring distances between colors.
//
// In traditional color spaces, a color can be identified solely by the
// observer's measurement of the color. Color appearance models such as CAM16
// also use information about the environment where the color was
// observed, known as the viewing conditions.
//
// For example, white under the traditional assumption of a midday sun white
// point is accurately measured as a slightly chromatic blue by CAM16.
// (roughly, hue 203, chroma 3, lightness 100)
type Cam16 struct {
	// CAM16 color dimensions
	hue    float64
	chroma float64
	j      float64
	q      float64
	m      float64
	s      float64

	// Coordinates in UCS space
	jstar float64
	astar float64
	bstar float64
}

// XYZToCam16RGB transforms XYZ color space coordinates to 'cone'/'RGB' responses in CAM16.
var XYZToCam16RGB = [][]float64{
	{0.401288, 0.650173, -0.051461},
	{-0.250268, 1.204414, 0.045854},
	{-0.002079, 0.048952, 0.953127},
}

// CAM16RGBToXYZ transforms 'cone'/'RGB' responses in CAM16 to XYZ color space coordinates.
var CAM16RGBToXYZ = [][]float64{
	{1.8620678, -1.0112547, 0.14918678},
	{0.38752654, 0.62144744, -0.00897398},
	{-0.01584150, -0.03412294, 1.0499644},
}

// distance calculates the color distance between two CAM16 instances.
func (c *Cam16) distance(other *Cam16) float64 {
	dJ := c.GetJstar() - other.GetJstar()
	dA := c.GetAstar() - other.GetAstar()
	dB := c.GetBstar() - other.GetBstar()
	dEPrime := math.Sqrt(dJ*dJ + dA*dA + dB*dB)
	dE := 1.41 * math.Pow(dEPrime, 0.63)
	return dE
}

// GetHue returns the hue in CAM16.
func (c *Cam16) GetHue() float64 {
	return c.hue
}

// GetChroma returns the chroma in CAM16.
func (c *Cam16) GetChroma() float64 {
	return c.chroma
}

// GetJ returns the lightness in CAM16.
func (c *Cam16) GetJ() float64 {
	return c.j
}

// GetQ returns the brightness in CAM16.
func (c *Cam16) GetQ() float64 {
	return c.q
}

// GetM returns the colorfulness in CAM16.
func (c *Cam16) GetM() float64 {
	return c.m
}

// GetS returns the saturation in CAM16.
func (c *Cam16) GetS() float64 {
	return c.s
}

// GetJstar returns the lightness coordinate in CAM16-UCS.
func (c *Cam16) GetJstar() float64 {
	return c.jstar
}

// GetAstar returns the a* coordinate in CAM16-UCS.
func (c *Cam16) GetAstar() float64 {
	return c.astar
}

// GetBstar returns the b* coordinate in CAM16-UCS.
func (c *Cam16) GetBstar() float64 {
	return c.bstar
}

// Cam16FromInt convert [argb] to CAM16, assuming the color was viewed in default viewing conditions.
func Cam16FromInt(argb int) Cam16 {
	return Cam16FromIntInViewingConditions(argb, DefaultViewingConditions)
}

// Cam16FromIntInViewingConditions Create a CAM16 color from a color in defined viewing conditions.
//
// [argb] ARGB representation of a color.
// [viewingConditions] Information about the environment where the color was observed.
//
// The RGB => XYZ conversion matrix elements are derived scientific constants. While the values
// may differ at runtime due to floating point imprecision, keeping the values the same, and
// accurate, across implementations takes precedence.
func Cam16FromIntInViewingConditions(argb int, viewingConditions ViewingConditions) Cam16 {
	red := (argb & 0x00ff0000) >> 16
	green := (argb & 0x0000ff00) >> 8
	blue := argb & 0x000000ff
	redL := colorUtils.Linearized(red)
	greenL := colorUtils.Linearized(green)
	blueL := colorUtils.Linearized(blue)
	x := 0.41233895*redL + 0.35762064*greenL + 0.18051042*blueL
	y := 0.2126*redL + 0.7152*greenL + 0.0722*blueL
	z := 0.01932141*redL + 0.11916382*greenL + 0.95034478*blueL

	return Cam16FromXyzInViewingConditions(x, y, z, viewingConditions)
}

// Cam16FromXyzInViewingConditions Create a CAM16 color from a color in defined viewing conditions.
func Cam16FromXyzInViewingConditions(x, y, z float64, viewingConditions ViewingConditions) Cam16 {
	// Transform XYZ to 'cone'/'rgb' responses
	matrix := XYZToCam16RGB
	rT := (x * matrix[0][0]) + (y * matrix[0][1]) + (z * matrix[0][2])
	gT := (x * matrix[1][0]) + (y * matrix[1][1]) + (z * matrix[1][2])
	bT := (x * matrix[2][0]) + (y * matrix[2][1]) + (z * matrix[2][2])

	// Discount illuminant
	rD := viewingConditions.GetRgbD()[0] * rT
	gD := viewingConditions.GetRgbD()[1] * gT
	bD := viewingConditions.GetRgbD()[2] * bT

	// Chromatic adaptation
	rAF := math.Pow(viewingConditions.GetFl()*math.Abs(rD)/100.0, 0.42)
	gAF := math.Pow(viewingConditions.GetFl()*math.Abs(gD)/100.0, 0.42)
	bAF := math.Pow(viewingConditions.GetFl()*math.Abs(bD)/100.0, 0.42)
	rA := mathUtils.Signum(rD) * 400.0 * rAF / (rAF + 27.13)
	gA := mathUtils.Signum(gD) * 400.0 * gAF / (gAF + 27.13)
	bA := mathUtils.Signum(bD) * 400.0 * bAF / (bAF + 27.13)

	// redness-greenness
	a := (11.0*rA - 12.0*gA + bA) / 11.0
	// yellowness-blueness
	b := (rA + gA - 2.0*bA) / 9.0

	// auxiliary components
	u := (20.0*rA + 20.0*gA + 21.0*bA) / 20.0
	p2 := (40.0*rA + 20.0*gA + bA) / 20.0

	// hue
	atan2 := math.Atan2(b, a)
	atanDegrees := mathUtils.ToDegrees(atan2)
	hue := atanDegrees
	if atanDegrees < 0 {
		hue += 360.0
	} else if atanDegrees >= 360 {
		hue -= 360.0
	}
	hueRadians := mathUtils.ToRadians(hue)

	// achromatic response to color
	ac := p2 * viewingConditions.GetNbb()

	// CAM16 lightness and brightness
	j := 100.0 * math.Pow(ac/viewingConditions.GetAw(), viewingConditions.GetC()*viewingConditions.GetZ())
	q := 4.0 / viewingConditions.GetC() * math.Sqrt(j/100.0) * (viewingConditions.GetAw() + 4.0) * viewingConditions.GetFlRoot()

	// CAM16 chroma, colorfulness, and saturation.
	huePrime := hue
	if hue < 20.14 {
		huePrime += 360
	}
	eHue := 0.25 * (math.Cos(mathUtils.ToRadians(huePrime)+2.0) + 3.8)
	p1 := 50000.0 / 13.0 * eHue * viewingConditions.GetNc() * viewingConditions.GetNcb()
	t := p1 * math.Hypot(a, b) / (u + 0.305)
	alpha := math.Pow(1.64-math.Pow(0.29, viewingConditions.GetN()), 0.73) * math.Pow(t, 0.9)
	// CAM16 chroma, colorfulness, saturation
	c := alpha * math.Sqrt(j/100.0)
	m := c * viewingConditions.GetFlRoot()
	s := 50.0 * math.Sqrt((alpha*viewingConditions.GetC())/(viewingConditions.GetAw()+4.0))

	// CAM16-UCS components
	jstar := (1.0 + 100.0*0.007) * j / (1.0 + 0.007*j)
	mstar := 1.0 / 0.0228 * math.Log1p(0.0228*m)
	astar := mstar * math.Cos(hueRadians)
	bstar := mstar * math.Sin(hueRadians)

	return Cam16{
		hue:    hue,
		chroma: c,
		j:      j,
		q:      q,
		m:      m,
		s:      s,
		jstar:  jstar,
		astar:  astar,
		bstar:  bstar,
	}
}

// Cam16FromJch constructs a CAM16 color from the given CAM16 lightness, chroma, and hue.
//
// [j] CAM16 lightness
// [c] CAM16 chroma
// [h] CAM16 hue
func Cam16FromJch(j, c, h float64) Cam16 {
	return cam16FromJchInViewingConditions(j, c, h, DefaultViewingConditions)
}

// cam16FromJchInViewingConditions constructs a CAM16 color from the given CAM16 lightness [j], chroma [c], and hue [h],
// in the given viewing conditions.
func cam16FromJchInViewingConditions(j, c, h float64, viewingConditions ViewingConditions) Cam16 {
	q := (4.0 / viewingConditions.GetC()) * math.Sqrt(j/100.0) * (viewingConditions.GetAw() + 4.0) * math.Sqrt(viewingConditions.GetFl())
	m := c * math.Sqrt(viewingConditions.GetFl())
	alpha := c / math.Sqrt(j/100.0)
	s := 50.0 * math.Sqrt((alpha*viewingConditions.GetC())/(viewingConditions.GetAw()+4.0))

	hueRadians := math.Pi * h / 180.0
	jstar := (1.0 + 100.0*0.007) * j / (1.0 + 0.007*j)
	mstar := 1.0 / 0.0228 * math.Log1p(0.0228*m)
	astar := mstar * math.Cos(hueRadians)
	bstar := mstar * math.Sin(hueRadians)

	return Cam16{
		hue:    h,
		chroma: c,
		j:      j,
		q:      q,
		m:      m,
		s:      s,
		jstar:  jstar,
		astar:  astar,
		bstar:  bstar,
	}
}

// Cam16FromUcs constructs a CAM16 color from the given CAM16-UCS lightness, a dimension, and b dimension.
//
// [jstar] CAM16-UCS lightness.
// [astar] CAM16-UCS a dimension. Like a* in L*a*b*, it is a Cartesian coordinate on the Y axis.
// [bstar] CAM16-UCS b dimension. Like a* in L*a*b*, it is a Cartesian coordinate on the X axis.
func Cam16FromUcs(jstar, astar, bstar float64) Cam16 {
	return cam16FromUcsInViewingConditions(jstar, astar, bstar, DefaultViewingConditions)
}

func cam16FromUcsInViewingConditions(jstar, astar, bstar float64, viewingConditions ViewingConditions) Cam16 {
	m := math.Hypot(astar, bstar)
	m2 := (math.Exp(m*0.0228) - 1.0) / 0.0228
	c := m2 / math.Sqrt(viewingConditions.GetFl())
	h := math.Atan2(bstar, astar) * (180.0 / math.Pi)
	if h < 0.0 {
		h += 360.0
	}
	j := jstar / (1.0 - (jstar-100.0)*0.007)
	return cam16FromJchInViewingConditions(j, c, h, viewingConditions)
}

func (c *Cam16) ToInt() int {
	return c.Viewed(DefaultViewingConditions)
}

func (c *Cam16) Viewed(viewingConditions ViewingConditions) int {
	xyz := c.XyzInViewingConditions(viewingConditions, nil)
	return colorUtils.ArgbFromXyz(xyz[0], xyz[1], xyz[2])
}

func (c *Cam16) XyzInViewingConditions(viewingConditions ViewingConditions, returnArray []float64) []float64 {
	alpha := 0.0
	if c.GetChroma() != 0.0 && c.GetJ() != 0.0 {
		alpha = c.GetChroma() / math.Sqrt(c.GetJ()/100.0)
	}

	t := math.Pow(alpha/math.Pow(1.64-math.Pow(0.29, viewingConditions.GetN()), 0.73), 1.0/0.9)
	hRad := math.Pi * c.GetHue() / 180.0

	eHue := 0.25 * (math.Cos(hRad+2.0) + 3.8)
	ac := viewingConditions.GetAw() * math.Pow(c.GetJ()/100.0, 1.0/(viewingConditions.GetC()*viewingConditions.GetZ()))
	p1 := eHue * (50000.0 / 13.0) * viewingConditions.GetNc() * viewingConditions.GetNcb()
	p2 := ac / viewingConditions.GetNbb()

	hSin := math.Sin(hRad)
	hCos := math.Cos(hRad)

	gamma := 23.0 * (p2 + 0.305) * t / (23.0*p1 + 11.0*t*hCos + 108.0*t*hSin)
	a := gamma * hCos
	b := gamma * hSin
	rA := (460.0*p2 + 451.0*a + 288.0*b) / 1403.0
	gA := (460.0*p2 - 891.0*a - 261.0*b) / 1403.0
	bA := (460.0*p2 - 220.0*a - 6300.0*b) / 1403.0

	rCBase := math.Max(0, (27.13*math.Abs(rA))/(400.0-math.Abs(rA)))
	rC := mathUtils.Signum(rA) * (100.0 / viewingConditions.GetFl()) * math.Pow(rCBase, 1.0/0.42)
	gCBase := math.Max(0, (27.13*math.Abs(gA))/(400.0-math.Abs(gA)))
	gC := mathUtils.Signum(gA) * (100.0 / viewingConditions.GetFl()) * math.Pow(gCBase, 1.0/0.42)
	bCBase := math.Max(0, (27.13*math.Abs(bA))/(400.0-math.Abs(bA)))
	bC := mathUtils.Signum(bA) * (100.0 / viewingConditions.GetFl()) * math.Pow(bCBase, 1.0/0.42)
	rF := rC / viewingConditions.GetRgbD()[0]
	gF := gC / viewingConditions.GetRgbD()[1]
	bF := bC / viewingConditions.GetRgbD()[2]

	matrix := CAM16RGBToXYZ
	x := rF*matrix[0][0] + gF*matrix[0][1] + bF*matrix[0][2]
	y := rF*matrix[1][0] + gF*matrix[1][1] + bF*matrix[1][2]
	z := rF*matrix[2][0] + gF*matrix[2][1] + bF*matrix[2][2]

	if returnArray != nil {
		returnArray[0] = x
		returnArray[1] = y
		returnArray[2] = z
		return returnArray
	} else {
		return []float64{x, y, z}
	}
}
