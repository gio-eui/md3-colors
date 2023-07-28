package mathUtils

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

import "math"

// Signum returns the sign of the given number.
//
// Returns -1 if the number is negative, 0 if the number is 0, and 1 if the number is positive.
func Signum(num float64) float64 {
	if num < 0 {
		return -1.0
	} else if num == 0.0 {
		return 0.0
	} else {
		return 1.0
	}
}

// Lerp performs linear interpolation between two values.
// It returns the interpolated value based on the given amount.
// If amount = 0, it returns start; if amount = 1, it returns stop.
func Lerp(start, stop, amount float64) float64 {
	return (1.0-amount)*start + amount*stop
}

// ClampInt clamps an integer between two integers.
// It returns the input value when min <= input <= max,
// and either min or max otherwise.
func ClampInt(min, max, input int) int {
	if input < min {
		return min
	} else if input > max {
		return max
	}

	return input
}

// ClampDouble clamps a double value between two floating-point numbers.
// It returns the input value when min <= input <= max,
// and either min or max otherwise.
func ClampDouble(min, max, input float64) float64 {
	if input < min {
		return min
	} else if input > max {
		return max
	}

	return input
}

// ToDegrees converts a radian measure to a degree measure.
func ToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// SanitizeDegreesInt sanitizes a degree measure as an integer.
// It returns a degree measure between 0 (inclusive) and 360 (exclusive).
func SanitizeDegreesInt(degrees int) int {
	degrees = degrees % 360
	if degrees < 0 {
		degrees = degrees + 360
	}
	return degrees
}

// SanitizeDegreesDouble sanitizes a degree measure as a floating-point number.
// It returns a degree measure between 0.0 (inclusive) and 360.0 (exclusive).
func SanitizeDegreesDouble(degrees float64) float64 {
	degrees = math.Mod(degrees, 360.0)
	if degrees < 0 {
		degrees = degrees + 360.0
	}
	return degrees
}

// RotationDirection returns the sign of the difference between two
// angles, in degrees.
//
// For angles that are 180 degrees apart from each other, both
// directions have the same travel distance, so either direction is
// shortest. The value 1.0 is returned in this case.
//
// [from] The angle travel starts from, in degrees.
// [to] The angle travel ends at, in degrees.
// Returns -1 if decreasing from leads to the shortest travel
// distance, 1 if increasing from leads to the shortest travel
// distance.
func RotationDirection(from, to float64) float64 {
	increasingDifference := SanitizeDegreesDouble(to - from)
	if increasingDifference <= 180.0 {
		return 1.0
	}
	return -1.0
}

// DifferenceDegrees calculates the distance between two points on a circle, represented using degrees.
func DifferenceDegrees(a, b float64) float64 {
	return 180.0 - math.Abs(math.Abs(a-b)-180.0)
}

// ToRadians converts a degree measure to a radian measure.
func ToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// MatrixMultiply multiplies a 1x3 row vector with a 3x3 matrix.
func MatrixMultiply(row []float64, matrix [][]float64) []float64 {
	a := row[0]*matrix[0][0] + row[1]*matrix[0][1] + row[2]*matrix[0][2]
	b := row[0]*matrix[1][0] + row[1]*matrix[1][1] + row[2]*matrix[1][2]
	c := row[0]*matrix[2][0] + row[1]*matrix[2][1] + row[2]*matrix[2][2]
	return []float64{a, b, c}
}
