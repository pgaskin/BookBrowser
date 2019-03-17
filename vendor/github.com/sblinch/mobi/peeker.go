package mobi

type Peeker []uint8

func (p Peeker) Magic() mobiMagicType {
	return mobiMagicType(p)
}

func (p Peeker) String() string {
	return string(p)
}

func (p Peeker) Bytes() []uint8 {
	return p
}

func (p Peeker) Len() int {
	return len(p)
}
