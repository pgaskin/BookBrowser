package mobi

import (
	"bytes"
	"encoding/binary"
)

type mobiFlis struct { //  RECORD -2
	Identifier uint32 //ID <comment="FLIS">;
	Fixed0     uint32 //UINT   fixed1  <comment="fixed value 8">;
	Fixed1     uint16 //USHORT fixed2  <comment="fixed value 65">;
	Fixed2     uint16 //USHORT fixed3  <comment="fixed value 0">;
	Fixed3     uint32 //UINT   fixed4  <comment="fixed value 0">;
	Fixed4     uint32 //UINT   fixed5  <comment="fixed value -1">;
	Fixed5     uint16 //USHORT fixed6  <comment="fixed value 1">;
	Fixed6     uint16 //USHORT fixed7  <comment="fixed value 3">;
	Fixed7     uint32 //UINT   fixed8  <comment="fixed value 3">;
	Fixed8     uint32 //UINT   fixed9  <comment="fixed value 1">;
	Fixed9     uint32 //UINT   fixed10 <comment="fixed value -1">;
} //FLISRECORD;

func (w *MobiWriter) generateFlis() []byte {
	c := mobiFlis{}
	c.Identifier = 1179404627 //StringToBytes("FLIS", &c.Identifier)
	c.Fixed0 = 8
	c.Fixed1 = 65
	//c.Fixed2
	//c.Fixed3
	c.Fixed4 = 4294967295
	c.Fixed5 = 1
	c.Fixed6 = 3
	c.Fixed7 = 3
	c.Fixed8 = 1
	c.Fixed9 = 4294967295

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, c)
	return buf.Bytes()
}
