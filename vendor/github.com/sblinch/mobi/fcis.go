package mobi

import (
	"bytes"
	"encoding/binary"
)

type mobiFcis struct { //  RECORD -1
	Identifier uint32 //UINT   ID <comment="FCIS">;
	Fixed0     uint32 //UINT  // fixed1  <comment="fixed value 20">;
	Fixed1     uint32 //UINT //  fixed2  <comment="fixed value 16">;
	Fixed2     uint32 //UINT   fixed3  <comment="fixed value 1">;
	Fixed3     uint32 //UINT  // fixed4  <comment="fixed value 0">;
	Fixed4     uint32 //UINT  // fixed5  <comment="text length (the same value as \"text length\" in the PalmDoc header)">;
	Fixed5     uint32 //UINT   fixed6  <comment="fixed value 0">;
	Fixed6     uint32 //UINT   fixed7  <comment="fixed value 32">;
	Fixed7     uint32 //UINT   fixed8  <comment="fixed value 8">;
	Fixed8     uint16 //USHORT fixed9  <comment="fixed value 1">;
	Fixed9     uint16 //USHORT fixed10 <comment="fixed value 1">;
	Fixed10    uint32 //UINT   fixed11 <comment="fixed value 0">;
} //FCISRECORD;*/

func (w *MobiWriter) generateFcis() []byte {
	c := mobiFcis{}
	c.Identifier = 1178814803 //StringToBytes("FLIS", &c.Identifier)
	c.Fixed0 = 20
	c.Fixed1 = 16
	c.Fixed2 = 1
	//c.Fixed3
	c.Fixed4 = w.Pdh.TextLength
	//c.Fixed5 = 0
	c.Fixed6 = 32
	c.Fixed7 = 8
	c.Fixed8 = 1
	c.Fixed9 = 1
	//c.Fixed10 = 0

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, c)
	return buf.Bytes()
}
