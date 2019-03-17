package mobi

type Mint int

func (i Mint) UInt16() uint16 {
	return uint16(i)
}

func (i Mint) UInt32() uint32 {
	return uint32(i)
}

func (i Mint) Int() int {
	return int(i)
}
