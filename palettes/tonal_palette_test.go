package palettes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTonalPaletteFromInt(t *testing.T) {

	blue := NewTonalPaletteFromInt(0xFF0000FF)

	assert.Equal(t, blue.Tone(0), 0xFF000000)
	assert.Equal(t, blue.Tone(10), 0xff00006e)
	assert.Equal(t, blue.Tone(20), 0xff0001ac)
	assert.Equal(t, blue.Tone(30), 0xff0000ef)
	assert.Equal(t, blue.Tone(40), 0xff343dff)
	assert.Equal(t, blue.Tone(50), 0xff5a64ff)
	assert.Equal(t, blue.Tone(60), 0xff7c84ff)
	assert.Equal(t, blue.Tone(70), 0xff9da3ff)
	assert.Equal(t, blue.Tone(80), 0xffbec2ff)
	assert.Equal(t, blue.Tone(90), 0xffe0e0ff)
	assert.Equal(t, blue.Tone(95), 0xfff1efff)
	assert.Equal(t, blue.Tone(100), 0xffffffff)
}
