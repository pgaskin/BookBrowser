package mobi

type mobiHeader struct {
	Identifier          [4]uint8 `format:"string"` // Must be characters MOBI
	HeaderLength        uint32   // The length of the MOBI header, including the previous 4 bytes
	MobiType            uint32   // Mobi type enum
	TextEncoding        uint32   // 1252 = CP1252 (WinLatin1); 65001 = UTF-8
	UniqueID            uint32   // Some kind of unique ID number (random?)
	FileVersion         uint32   // Version of the Mobipocket format used in this file. //If FileVersion == 8. Then it's KF8
	OrthographicIndex   uint32   // Section number of orthographic meta index. 0xFFFFFFFF if index is not available.
	InflectionIndex     uint32   // Section number of inflection meta index. 0xFFFFFFFF if index is not available.
	IndexNames          uint32   // 0xFFFFFFFF if index is not available.
	IndexKeys           uint32   // 0xFFFFFFFF if index is not available.
	ExtraIndex0         uint32   // Section number of extra 0 meta index. 0xFFFFFFFF if index is not available.
	ExtraIndex1         uint32   // Section number of extra 1 meta index. 0xFFFFFFFF if index is not available.
	ExtraIndex2         uint32   // Section number of extra 2 meta index. 0xFFFFFFFF if index is not available.
	ExtraIndex3         uint32   // Section number of extra 3 meta index. 0xFFFFFFFF if index is not available.
	ExtraIndex4         uint32   // Section number of extra 4 meta index. 0xFFFFFFFF if index is not available.
	ExtraIndex5         uint32   // Section number of extra 5 meta index. 0xFFFFFFFF if index is not available.
	FirstNonBookIndex   uint32   // First record number (starting with 0) that's not the book's text
	FullNameOffset      uint32   // Offset in record 0 (not from start of file) of the full name of the book
	FullNameLength      uint32   // Length in bytes of the full name of the book
	Locale              uint32   // Book locale code. Low byte is main language 09=English, next byte is dialect, 08=British, 04=US. Thus US English is 1033, UK English is 2057.
	InputLanguage       uint32   //Input language for a dictionary
	OutputLanguage      uint32   //Output language for a dictionary
	MinVersion          uint32   //Minimum mobipocket version support needed to read this file.
	FirstImageIndex     uint32   //First record number (starting with 0) that contains an image. Image records should be sequential.
	HuffmanRecordOffset uint32   //The record number of the first huffman compression record.
	HuffmanRecordCount  uint32   //The number of huffman compression records.
	HuffmanTableOffset  uint32
	HuffmanTableLength  uint32
	ExthFlags           uint32   //Bitfield. If bit 6 (0x40) is set, then there's an EXTH record
	Unknown1            [32]byte //Unknown values
	DrmOffset           uint32   //Offset to DRM key info in DRMed files. 0xFFFFFFFF if no DRM
	DrmCount            uint32   //Number of entries in DRM info. 0xFFFFFFFF if no DRM
	DrmSize             uint32   //Number of bytes in DRM info.
	DrmFlags            uint32   //Some flags concerning the DRM info.
	Unknown0            [12]byte //Unknown values

	// If it's KF8
	// 		FdstRecordIndex uint32
	// else
	FirstContentRecordNumber uint16 //Number of first text record. Normally 1.
	LastContentRecordNumber  uint16 //Number of last image record or number of last text record if it contains no images. Includes Image, DATP, HUFF, DRM.
	//End else

	Unknown6        uint32 //FdstRecordCount? //Use 0x00000001.
	FcisRecordIndex uint32
	FcisRecordCount uint32 //Use 0x00000001. // Always 1
	FlisRecordIndex uint32
	FlisRecordCount uint32 //Use 0x00000001. // Always 1
	Unknown7        uint32
	Unknown8        uint32
	SrcsRecordIndex uint32
	SrcsRecordCount uint32
	Unknown9        uint32
	Unknown10       uint32

	// A set of binary flags, some of which indicate extra data at the end of each text block. This only
	// seems to be valid for Mobipocket format version 5 and 6 (and higher?), when the header length is 228 (0xE4) or 232 (0xE8).
	// 		bit 1 (0x1): <extra multibyte bytes><size>
	// 		bit 2 (0x2): <TBS indexing description of this HTML record><size>
	// 		bit 3 (0x4): <uncrossable breaks><size>
	// Setting bit 2 (0x2) disables <guide><reference type="start"> functionality.
	ExtraRecordDataFlags uint32 `format:"bits"`
	IndxRecodOffset      uint32 //(If not 0xFFFFFFFF) The record number of the first INDX record created from an ncx file.

	//If header lenght is 248 then there's 16 extra bytes.

	/*
			If KF8
				FragmentIndex uint32
				SkeletonIndex uint32
			Else
				unknown14 uint32
				unknown15 uint32

			DatpIndex uint32

			If KF8
				GuideIndex uint32
			Else
				unknown16 uint32

			unknown17 uint32
		    unknown18 uint32
		    unknown19 uint32 ?
		    unknown20 uint32 ?
	*/
}
