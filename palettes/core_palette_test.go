package palettes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCorePaletteFromInt(t *testing.T) {
	blue := NewCorePaletteFromInt(0xff0000FF)

	assert.Equal(t, blue.A1.Tone(0), 0xff000000)
	assert.Equal(t, blue.A1.Tone(10), 0xff00006e)
	assert.Equal(t, blue.A1.Tone(20), 0xff0001ac)
	assert.Equal(t, blue.A1.Tone(30), 0xff0000ef)
	assert.Equal(t, blue.A1.Tone(40), 0xff343dff)
	assert.Equal(t, blue.A1.Tone(50), 0xff5a64ff)
	assert.Equal(t, blue.A1.Tone(60), 0xff7c84ff)
	assert.Equal(t, blue.A1.Tone(70), 0xff9da3ff)
	assert.Equal(t, blue.A1.Tone(80), 0xffbec2ff)
	assert.Equal(t, blue.A1.Tone(90), 0xffe0e0ff)
	assert.Equal(t, blue.A1.Tone(95), 0xfff1efff)
	assert.Equal(t, blue.A1.Tone(100), 0xffffffff)

	assert.Equal(t, blue.A2.Tone(0), 0xff000000)
	assert.Equal(t, blue.A2.Tone(10), 0xff191a2c)
	assert.Equal(t, blue.A2.Tone(20), 0xff2e2f42)
	assert.Equal(t, blue.A2.Tone(30), 0xff444559)
	assert.Equal(t, blue.A2.Tone(40), 0xff5c5d72)
	assert.Equal(t, blue.A2.Tone(50), 0xff75758b)
	assert.Equal(t, blue.A2.Tone(60), 0xff8f8fa6)
	assert.Equal(t, blue.A2.Tone(70), 0xffa9a9c1)
	assert.Equal(t, blue.A2.Tone(80), 0xffc5c4dd)
	assert.Equal(t, blue.A2.Tone(90), 0xffe1e0f9)
	assert.Equal(t, blue.A2.Tone(95), 0xfff1efff)
	assert.Equal(t, blue.A2.Tone(100), 0xffffffff)

	blueContent := NewContentCorePaletteFromInt(0xff0000FF)

	assert.Equal(t, blueContent.A1.Tone(0), 0xff000000)
	assert.Equal(t, blueContent.A1.Tone(10), 0xff00006e)
	assert.Equal(t, blueContent.A1.Tone(20), 0xff0001ac)
	assert.Equal(t, blueContent.A1.Tone(30), 0xff0000ef)
	assert.Equal(t, blueContent.A1.Tone(40), 0xff343dff)
	assert.Equal(t, blueContent.A1.Tone(50), 0xff5a64ff)
	assert.Equal(t, blueContent.A1.Tone(60), 0xff7c84ff)
	assert.Equal(t, blueContent.A1.Tone(70), 0xff9da3ff)
	assert.Equal(t, blueContent.A1.Tone(80), 0xffbec2ff)
	assert.Equal(t, blueContent.A1.Tone(90), 0xffe0e0ff)
	assert.Equal(t, blueContent.A1.Tone(95), 0xfff1efff)
	assert.Equal(t, blueContent.A1.Tone(100), 0xffffffff)

	assert.Equal(t, blueContent.A2.Tone(0), 0xff000000)
	assert.Equal(t, blueContent.A2.Tone(10), 0xff14173f)
	assert.Equal(t, blueContent.A2.Tone(20), 0xff2a2d55)
	assert.Equal(t, blueContent.A2.Tone(30), 0xff40436d)
	assert.Equal(t, blueContent.A2.Tone(40), 0xff585b86)
	assert.Equal(t, blueContent.A2.Tone(50), 0xff7173a0)
	assert.Equal(t, blueContent.A2.Tone(60), 0xff8b8dbb)
	assert.Equal(t, blueContent.A2.Tone(70), 0xffa5a7d7)
	assert.Equal(t, blueContent.A2.Tone(80), 0xffc1c3f4)
	assert.Equal(t, blueContent.A2.Tone(90), 0xffe0e0ff)
	assert.Equal(t, blueContent.A2.Tone(95), 0xfff1efff)
	assert.Equal(t, blueContent.A2.Tone(100), 0xffffffff)
}
