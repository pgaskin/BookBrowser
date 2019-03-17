package mobi

import (
	"bytes"
	"encoding/binary"
)

func (w *MobiWriter) chapterIsDeep() bool {
	for _, node := range w.chapters {
		if node.SubChapterCount() > 0 {
			return true
		}
	}
	return false
}

func (w *MobiWriter) writeINDX_1() {
	buf := new(bytes.Buffer)
	// Tagx
	tagx := mobiTagx{}
	if w.chapterIsDeep() {
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_Pos])
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_Len])
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_NameOffset])
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_DepthLvl])
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_Parent])
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_Child1])
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_ChildN])
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_END])
	} else {
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_Pos])
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_Len])
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_NameOffset])
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_DepthLvl])
		tagx.Tags = append(tagx.Tags, mobiTagxMap[TagEntry_END])
	}

	/*************************************/

	/*************************************/
	magicTagx.WriteTo(&tagx.Identifier)
	tagx.ControlByteCount = 1
	tagx.HeaderLenght = uint32(tagx.TagCount()*4) + 12

	TagX := new(bytes.Buffer)
	binary.Write(TagX, binary.BigEndian, tagx.Identifier)
	binary.Write(TagX, binary.BigEndian, tagx.HeaderLenght)
	binary.Write(TagX, binary.BigEndian, tagx.ControlByteCount)
	binary.Write(TagX, binary.BigEndian, tagx.Tags)

	// Indx
	//	IndxBin := new(bytes.Buffer)
	indx := mobiIndx{}
	magicIndx.WriteTo(&indx.Identifier)
	indx.HeaderLen = MOBI_INDX_HEADER_LEN
	indx.Indx_Type = INDX_TYPE_INFLECTION
	indx.Idxt_Count = 1
	indx.Idxt_Encoding = MOBI_ENC_UTF8
	indx.SetUnk2 = 4294967295
	indx.Cncx_Records_Count = 1
	indx.Idxt_Entry_Count = uint32(w.chapterCount)
	indx.Tagx_Offset = MOBI_INDX_HEADER_LEN

	//binary.Write(IndxBin, binary.BigEndian, indx)
	// Idxt

	/************/

	IdxtLast := len(w.Idxt.Offset)
	Offset := w.Idxt.Offset[IdxtLast-1]
	Rec := w.cncxBuffer.Bytes()[Offset-MOBI_INDX_HEADER_LEN:]

	Rec = Rec[0 : Rec[0]+1]
	RLen := len(Rec)

	//w.File.Write(Rec)

	Padding := (RLen + 2) % 4

	//IDXT_OFFSET, := w.File.Seek(0, 1)

	indx.Idxt_Offset = MOBI_INDX_HEADER_LEN + uint32(TagX.Len()) + uint32(RLen+2+Padding) // Offset to Idxt Record
	//w.Idxt1.Offset = []uint16{uint16(offset)}
	/************/

	binary.Write(buf, binary.BigEndian, indx)
	buf.Write(TagX.Bytes())
	buf.Write(Rec)
	binary.Write(buf, binary.BigEndian, uint16(IdxtLast))

	for Padding != 0 {
		buf.Write([]byte{0})
		Padding--
	}

	buf.WriteString(magicIdxt.String())

	binary.Write(buf, binary.BigEndian, uint16(MOBI_INDX_HEADER_LEN+uint32(TagX.Len())))

	//ioutil.WriteFile("TAGX_TEST", TagX.Bytes(), 0644)
	//ioutil.WriteFile("INDX_TEST", IndxBin.Bytes(), 0644)
	buf.Write([]uint8{0, 0})
	w.Header.IndxRecodOffset = w.AddRecord(buf.Bytes()).UInt32()
}

func (w *MobiWriter) writeINDX_2() {
	buf := new(bytes.Buffer)
	indx := mobiIndx{}
	magicIndx.WriteTo(&indx.Identifier)
	indx.HeaderLen = MOBI_INDX_HEADER_LEN
	indx.Indx_Type = INDX_TYPE_NORMAL
	indx.Unk1 = uint32(1)
	indx.Idxt_Encoding = 4294967295
	indx.SetUnk2 = 4294967295
	indx.Idxt_Offset = uint32(MOBI_INDX_HEADER_LEN + w.cncxBuffer.Len())
	indx.Idxt_Count = uint32(len(w.Idxt.Offset))

	binary.Write(buf, binary.BigEndian, indx)
	buf.Write(w.cncxBuffer.Bytes())

	buf.WriteString(magicIdxt.String())
	for _, offset := range w.Idxt.Offset {
		//Those offsets are not relative INDX record.
		//So we need to adjust that.
		binary.Write(buf, binary.BigEndian, offset) //+MOBI_INDX_HEADER_LEN)

	}

	Padding := (len(w.Idxt.Offset) + 4) % 4
	for Padding != 0 {
		buf.Write([]byte{0})
		Padding--
	}

	w.AddRecord(buf.Bytes())
	w.AddRecord(w.cncxLabelBuffer.Bytes())
}
