package gopdf

import (
	"bytes"
	"fmt"

	"github.com/signintech/gopdf/fontmaker/core"
)

type SubfontDescriptorObj struct {
	buffer                bytes.Buffer
	PtrToSubsetFontObj    *SubsetFontObj
	indexObjPdfDictionary int
}

func (s *SubfontDescriptorObj) Init(func() *GoPdf) {}

func (s *SubfontDescriptorObj) GetType() string {
	return "SubFontDescriptor"
}
func (s *SubfontDescriptorObj) GetObjBuff() *bytes.Buffer {
	return &s.buffer
}

func (s *SubfontDescriptorObj) Build() error {
	ttfp := s.PtrToSubsetFontObj.GetTTFParser()
	s.buffer.WriteString("<<\n")
	s.buffer.WriteString("/Type /FontDescriptor\n")
	s.buffer.WriteString(fmt.Sprintf("/Ascent %d\n", DesignUnitsToPdf(ttfp.Ascender(), ttfp.UnitsPerEm())))
	s.buffer.WriteString(fmt.Sprintf("/CapHeight %d\n", DesignUnitsToPdf(ttfp.CapHeight(), ttfp.UnitsPerEm())))
	s.buffer.WriteString(fmt.Sprintf("/Descent %d\n", DesignUnitsToPdf(ttfp.Descender(), ttfp.UnitsPerEm())))
	s.buffer.WriteString(fmt.Sprintf("/Flags %d\n", ttfp.Flag()))
	s.buffer.WriteString(fmt.Sprintf("/FontBBox [%d %d %d %d]\n",
		DesignUnitsToPdf(ttfp.XMin(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.YMin(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.XMax(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.YMax(), ttfp.UnitsPerEm()),
	))
	s.buffer.WriteString(fmt.Sprintf("/FontFile2 %d 0 R\n", s.indexObjPdfDictionary+1))
	s.buffer.WriteString(fmt.Sprintf("/FontName /%s\n", CreateEmbeddedFontSubsetName(s.PtrToSubsetFontObj.GetFamily())))
	s.buffer.WriteString(fmt.Sprintf("/ItalicAngle %d\n", ttfp.ItalicAngle()))
	s.buffer.WriteString("/StemV 0\n")
	s.buffer.WriteString(fmt.Sprintf("/XHeight %d\n", DesignUnitsToPdf(ttfp.XHeight(), ttfp.UnitsPerEm())))
	s.buffer.WriteString(">>\n")
	return nil
}

func (s *SubfontDescriptorObj) SetIndexObjPdfDictionary(index int) {
	s.indexObjPdfDictionary = index
}

func (s *SubfontDescriptorObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	s.PtrToSubsetFontObj = ptr
}

func DesignUnitsToPdf(val int64, unitsPerEm uint64) int64 {
	return core.Round(float64(float64(val) * 1000.00 / float64(unitsPerEm)))
}
