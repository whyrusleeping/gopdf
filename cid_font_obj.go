package gopdf

import (
	"bytes"
	"fmt"
	"sort"
)

type CIDFontObj struct {
	buffer                    bytes.Buffer
	PtrToSubsetFontObj        *SubsetFontObj
	indexObjSubfontDescriptor int
}

func (ci *CIDFontObj) init(funcGetRoot func() *GoPdf) {
}

func (ci *CIDFontObj) build(objID int) error {

	ci.buffer.WriteString("<<\n")
	ci.buffer.WriteString(fmt.Sprintf("/BaseFont /%s\n", CreateEmbeddedFontSubsetName(ci.PtrToSubsetFontObj.GetFamily())))
	ci.buffer.WriteString("/CIDSystemInfo\n")
	ci.buffer.WriteString("<<\n")
	ci.buffer.WriteString("  /Ordering (Identity)\n")
	ci.buffer.WriteString("  /Registry (Adobe)\n")
	ci.buffer.WriteString("  /Supplement 0\n")
	ci.buffer.WriteString(">>\n")
	ci.buffer.WriteString(fmt.Sprintf("/FontDescriptor %d 0 R\n", ci.indexObjSubfontDescriptor+1)) //TODO fix
	ci.buffer.WriteString("/Subtype /CIDFontType2\n")
	ci.buffer.WriteString("/Type /Font\n")
	characterToGlyphIndex := ci.PtrToSubsetFontObj.CharacterToGlyphIndex
	ci.buffer.WriteString("/W [")

	sortedForEachRI(characterToGlyphIndex, func(k rune, v uint) {
		width := ci.PtrToSubsetFontObj.GlyphIndexToPdfWidth(v)
		ci.buffer.WriteString(fmt.Sprintf("%d[%d]", v, width))
	})
	ci.buffer.WriteString("]\n")
	ci.buffer.WriteString(">>\n")
	return nil
}

func sortedForEachRI(m map[rune]uint, f func(k rune, v uint)) {
	var keys []rune
	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		f(k, m[k])
	}
}

func sortedForEachIR(m map[int]rune, f func(k int, v rune)) {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		f(k, m[k])
	}
}

//SetIndexObjSubfontDescriptor set  indexObjSubfontDescriptor
func (ci *CIDFontObj) SetIndexObjSubfontDescriptor(index int) {
	ci.indexObjSubfontDescriptor = index
}

func (ci *CIDFontObj) getType() string {
	return "CIDFont"
}

func (ci *CIDFontObj) getObjBuff() *bytes.Buffer {
	//fmt.Printf("%s\n", me.buffer.String())
	return &ci.buffer
}

//SetPtrToSubsetFontObj set PtrToSubsetFontObj
func (ci *CIDFontObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	ci.PtrToSubsetFontObj = ptr
}

//GetObjBuff get buffer
func (ci *CIDFontObj) GetObjBuff() *bytes.Buffer {
	return ci.getObjBuff()
}

//Build build buffer
func (ci *CIDFontObj) Build(objID int) error {
	return ci.build(objID)
}
