package osmpbf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringTableEncode(t *testing.T) {
	sTable := NewStringTable()

	assert.Equal(t, 1, len(sTable.getStrTable()))

	idx1 := sTable.getStrIndex("hello")
	idx2 := sTable.getStrIndex("osm")
	idx3 := sTable.getStrIndex("world")
	idx4 := sTable.getStrIndex("hello")
	idx5 := sTable.getStrIndex("hello")
	idx6 := sTable.getStrIndex("world")
	idx7 := sTable.getStrIndex("help")
	idx8 := sTable.getStrIndex("help")

	assert.Equal(t, 1, idx1)
	assert.Equal(t, 2, idx2)
	assert.Equal(t, 3, idx3)
	assert.Equal(t, 1, idx4)
	assert.Equal(t, 1, idx5)
	assert.Equal(t, 3, idx6)
	assert.Equal(t, 4, idx7)
	assert.Equal(t, 4, idx8)

	assert.Equal(t, 5, len(sTable.getStrTable()))
	assert.Equal(t, "", sTable.getStrTable()[0])
	assert.Equal(t, "hello", sTable.getStrTable()[1])
	assert.Equal(t, "osm", sTable.getStrTable()[2])
	assert.Equal(t, "world", sTable.getStrTable()[3])
	assert.Equal(t, "help", sTable.getStrTable()[4])

	assert.Equal(t, 17, sTable.getStrTableSize())
}
