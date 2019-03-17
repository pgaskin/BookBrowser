package mobi

type mobiIdxt struct {
	Identifier [4]byte  `format:"string"`
	Offset     []uint16 /* mobiIndx.HeaderLenght + len(mobiTagx.HeaderLenght) */
	//Unk1       uint16   // Pad with zeros to make it multiples of 4?
}
