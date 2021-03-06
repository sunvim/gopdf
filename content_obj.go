package gopdf

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
)

type ContentObj struct { //impl IObj
	buffer bytes.Buffer
	stream bytes.Buffer

	//text bytes.Buffer
	getRoot func() *GoPdf
}

func (c *ContentObj) Init(funcGetRoot func() *GoPdf) {
	c.getRoot = funcGetRoot
}

func (c *ContentObj) Build() error {
	streamlen := c.stream.Len()
	c.buffer.WriteString("<<\n")
	c.buffer.WriteString("/Length " + strconv.Itoa(streamlen) + "\n")
	c.buffer.WriteString(">>\n")
	c.buffer.WriteString("stream\n")
	c.buffer.Write(c.stream.Bytes())
	c.buffer.WriteString("endstream\n")
	return nil
}

func (c *ContentObj) GetType() string {
	return "Content"
}

func (c *ContentObj) GetObjBuff() *bytes.Buffer {
	return &(c.buffer)
}

func (c *ContentObj) AppendStreamSubsetFont(rectangle *Rect, text string) {

	sumWidth := uint64(0)
	var buff bytes.Buffer
	for _, r := range text {
		index, err := c.getRoot().Curr.Font_ISubset.CharIndex(r)
		if err != nil {
			log.Fatalf("err:%s", err.Error())
		}
		buff.WriteString(fmt.Sprintf("%04X", index))
		width, err := c.getRoot().Curr.Font_ISubset.CharWidth(r)
		if err != nil {
			log.Fatalf("err:%s", err.Error())
		}
		sumWidth += width
	}

	fontSize := c.getRoot().Curr.Font_Size
	x := fmt.Sprintf("%0.2f", c.getRoot().Curr.X)
	y := fmt.Sprintf("%0.2f", c.getRoot().config.PageSize.H-c.getRoot().Curr.Y-(float64(fontSize)*0.7))

	c.stream.WriteString("BT\n")
	c.stream.WriteString(x + " " + y + " TD\n")
	c.stream.WriteString("/F" + strconv.Itoa(c.getRoot().Curr.Font_FontCount+1) + " " + strconv.Itoa(fontSize) + " Tf\n")
	c.stream.WriteString("<" + buff.String() + "> Tj\n")
	c.stream.WriteString("ET\n")
	if rectangle == nil {
		fontSize := c.getRoot().Curr.Font_Size
		c.getRoot().Curr.X += float64(sumWidth) * (float64(fontSize) / 1000.0)
	} else {
		c.getRoot().Curr.X += rectangle.W
	}
}

func (c *ContentObj) AppendStream(rectangle *Rect, text string) {

	fontSize := c.getRoot().Curr.Font_Size

	x := fmt.Sprintf("%0.2f", c.getRoot().Curr.X)
	y := fmt.Sprintf("%0.2f", c.getRoot().config.PageSize.H-c.getRoot().Curr.Y-(float64(fontSize)*0.7))

	c.stream.WriteString("BT\n")
	c.stream.WriteString(x + " " + y + " TD\n")
	c.stream.WriteString("/F" + strconv.Itoa(c.getRoot().Curr.Font_FontCount+1) + " " + strconv.Itoa(fontSize) + " Tf\n")
	c.stream.WriteString("(" + text + ") Tj\n")
	c.stream.WriteString("ET\n")
	if rectangle == nil {
		c.getRoot().Curr.X += StrHelperGetStringWidth(text, fontSize, c.getRoot().Curr.Font_IFont)
	} else {
		c.getRoot().Curr.X += rectangle.W
	}

}

func (c *ContentObj) AppendStreamLine(x1 float64, y1 float64, x2 float64, y2 float64) {

	h := c.getRoot().config.PageSize.H
	c.stream.WriteString(fmt.Sprintf("%0.2f %0.2f m %0.2f %0.2f l s\n", x1, h-y1, x2, h-y2))
}

func (c *ContentObj) AppendUnderline(startX float64, y float64, endX float64, endY float64, text string) {

	h := c.getRoot().config.PageSize.H
	ut := int(0)
	if c.getRoot().Curr.Font_IFont != nil {
		ut = c.getRoot().Curr.Font_IFont.GetUt()
	} else if c.getRoot().Curr.Font_ISubset != nil {
		ut = int(c.getRoot().Curr.Font_ISubset.GetUt())
	} else {
		log.Fatal("error AppendUnderline not found font")
	}

	textH := ContentObj_CalTextHeight(c.getRoot().Curr.Font_Size)
	arg3 := float64(h) - float64(y) - textH - textH*0.07
	arg4 := (float64(ut) / 1000.00) * float64(c.getRoot().Curr.Font_Size)
	c.stream.WriteString(fmt.Sprintf("%0.2f %0.2f %0.2f -%0.2f re f\n", startX, arg3, endX-startX, arg4))
}

func (c *ContentObj) AppendStreamSetLineWidth(w float64) {

	c.stream.WriteString(fmt.Sprintf("%.2f w\n", w))

}

//  Set the grayscale fills
func (c *ContentObj) AppendStreamSetGrayFill(w float64) {
	w = fixRange10(w)
	c.stream.WriteString(fmt.Sprintf("%.2f g\n", w))
}

//  Set the grayscale stroke
func (c *ContentObj) AppendStreamSetGrayStroke(w float64) {
	w = fixRange10(w)
	c.stream.WriteString(fmt.Sprintf("%.2f G\n", w))
}

func (c *ContentObj) AppendStreamImage(index int, x float64, y float64, rect *Rect) {
	//fmt.Printf("index = %d",index)
	h := c.getRoot().config.PageSize.H
	c.stream.WriteString(fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", rect.W, rect.H, x, h-(y+rect.H), index+1))
}

//cal text height
func ContentObj_CalTextHeight(fontsize int) float64 {
	return (float64(fontsize) * 0.7)
}

// When setting colour and grayscales the value has to be between 0.00 and 1.00
// This function takes a float64 and returns 0.0 if it is less than 0.0 and 1.0 if it
// is more than 1.0
func fixRange10(val float64) float64 {
	if val < 0.0 {
		return 0.0
	}
	if val > 1.0 {
		return 1.0
	}
	return val
}
