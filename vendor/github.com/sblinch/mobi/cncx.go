package mobi

type mobiCncx struct {
	Len       uint8   `init:"Id"`       //Lenght of Cncx ID
	Id        []uint8 `format:"string"` //String ID,
	NCX_Count uint16  // Number of IndxEntries
	// Pad with zeros to reach a multiple of 4
	/*
		0 - 2: IDLen 	Lenght of ID
		2 - *: ID

	*/
}
