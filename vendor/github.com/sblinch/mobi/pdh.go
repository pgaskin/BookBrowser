package mobi

type mobiPDHCompression uint16

// Compression Enum
const (
	// CompressionNone uint16(1). Text is stored without any compression
	CompressionNone mobiPDHCompression = 1
	// CompressionPalmDoc uint16(2). Text is compressed using simple LZ77 algorithm
	CompressionPalmDoc mobiPDHCompression = 2
	// CompressionHuffCdic uint16(17480). Text is compressed using HuffCdic
	CompressionHuffCdic mobiPDHCompression = 17480
)

//PalmDoc Header
type mobiPDH struct {
	Compression mobiPDHCompression //0  // 1 == no compression, 2 = PalmDOC compression, 17480 = HUFF/CDIC compression
	Unk1        uint16             //2  // Always zero
	TextLength  uint32             //4  // Uncompressed length of the entire text of the book
	RecordCount uint16             //8  // Number of PDB records used for the text of the book.
	RecordSize  uint16             //10 // Maximum size of each record containing text, always 4096
	Encryption  uint16             //12 // 0 == no encryption, 1 = Old Mobipocket Encryption, 2 = Mobipocket Encryption
	Unk2        uint16             //12 // Usually zero
}
