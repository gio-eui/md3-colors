package stringsUtils

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
	"fmt"
	colorUtils "github.com/gio-eui/md3-colors/utils/color"
)

func HexFromArgb(argb int) string {
	red := colorUtils.RedFromArgb(argb)
	blue := colorUtils.BlueFromArgb(argb)
	green := colorUtils.GreenFromArgb(argb)
	return fmt.Sprintf("#%02x%02x%02x", red, green, blue)
}
