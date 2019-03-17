package mobi

import (
	"os"
	"reflect"
)

type Mobi struct {
	file     *os.File
	fileStat os.FileInfo

	Pdf     mobiPDF            // Palm Database Format: http://wiki.mobileread.com/wiki/PDB#Palm_Database_Format
	Offsets []mobiRecordOffset // Offsets for all the records. Starting from beginning of a file.
	Pdh     mobiPDH

	Header mobiHeader
	Exth   mobiExth

	//Index
	Indx  []mobiIndx
	Idxt  mobiIdxt
	Cncx  mobiCncx
	Tagx  mobiTagx
	PTagx []mobiPTagx
}

const (
	MOBI_MAX_RECORD_SIZE    = 4096
	MOBI_PALMDB_HEADER_LEN  = 78
	MOBI_INDX_HEADER_LEN    = 192
	MOBI_PALMDOC_HEADER_LEN = 16
	MOBI_MOBIHEADER_LEN     = 232
)

type mobiRecordOffset struct {
	Offset     uint32 //The offset of record {N} from the start of the PDB of this record
	Attributes uint8  //Bit Field. The least significant four bits are used to represent the category values.
	Skip       uint8  //UniqueID is supposed to take 3 bytes, but for our inteded purposes uint16(UniqueID) should work. Let me know if there's any mobi files with more than 32767 records
	UniqueID   uint16 //The unique ID for this record. Often just a sequential count from 0
}

const (
	magicMobi     mobiMagicType = "MOBI"
	magicExth     mobiMagicType = "EXTH"
	magicHuff     mobiMagicType = "HUFF"
	magicCdic     mobiMagicType = "CDIC"
	magicFdst     mobiMagicType = "FDST"
	magicIdxt     mobiMagicType = "IDXT"
	magicIndx     mobiMagicType = "INDX"
	magicLigt     mobiMagicType = "LIGT"
	magicOrdt     mobiMagicType = "ORDT"
	magicTagx     mobiMagicType = "TAGX"
	magicFont     mobiMagicType = "FONT"
	magicAudi     mobiMagicType = "AUDI"
	magicVide     mobiMagicType = "VIDE"
	magicResc     mobiMagicType = "RESC"
	magicBoundary mobiMagicType = "BOUNDARY"
)

type mobiMagicType string

func (m mobiMagicType) String() string {
	return string(m)
}

func (m mobiMagicType) WriteTo(output interface{}) {
	out := reflect.ValueOf(output).Elem()

	if out.Type().Len() != len(m) {
		panic("Magic lenght is larger than target size")
	}

	for i := 0; i < out.Type().Len(); i++ {
		if i > len(m)-1 {
			break
		}
		out.Index(i).Set(reflect.ValueOf(byte(m[i])))
	}
}

const (
	MOBI_ENC_CP1252 = 1252  /**< cp-1252 encoding */
	MOBI_ENC_UTF8   = 65001 /**< utf-8 encoding */
	MOBI_ENC_UTF16  = 65002 /**< utf-16 encoding */
)