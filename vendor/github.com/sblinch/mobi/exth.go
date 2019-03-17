package mobi

// Type of EXTH record. If it's Binary/Numberic then read/write
// it using BigEndian, String is read/write using LittleEndian
type ExthType uint32

const (
	EXTH_TYPE_NUMERIC ExthType = 0
	EXTH_TYPE_STRING  ExthType = 1
	EXTH_TYPE_BINARY  ExthType = 2
)

// EXTH record IDs
const (
	EXTH_DRMSERVER       uint32 = 1
	EXTH_DRMCOMMERCE            = 2
	EXTH_DRMEBOOKBASE           = 3
	EXTH_TITLE                  = 99  /**< <dc:title> */
	EXTH_AUTHOR                 = 100 /**< <dc:creator> */
	EXTH_PUBLISHER              = 101 /**< <dc:publisher> */
	EXTH_IMPRINT                = 102 /**< <imprint> */
	EXTH_DESCRIPTION            = 103 /**< <dc:description> */
	EXTH_ISBN                   = 104 /**< <dc:identifier opf:scheme="ISBN"> */
	EXTH_SUBJECT                = 105 // Could appear multiple times /**< <dc:subject> */
	EXTH_PUBLISHINGDATE         = 106 /**< <dc:date> */
	EXTH_REVIEW                 = 107 /**< <review> */
	EXTH_CONTRIBUTOR            = 108 /**< <dc:contributor> */
	EXTH_RIGHTS                 = 109 /**< <dc:rights> */
	EXTH_SUBJECTCODE            = 110 /**< <dc:subject BASICCode="subjectcode"> */
	EXTH_TYPE                   = 111 /**< <dc:type> */
	EXTH_SOURCE                 = 112 /**< <dc:source> */
	EXTH_ASIN                   = 113 // Kindle Paperwhite labels books with "Personal" if they don't have this record.
	EXTH_VERSION                = 114
	EXTH_SAMPLE                 = 115 // 0x0001 if the book content is only a sample of the full book
	EXTH_STARTREADING           = 116 // Position (4-byte offset) in file at which to open when first opened /**< Start reading */
	EXTH_ADULT                  = 117 // Mobipocket Creator adds this if Adult only is checked on its GUI; contents: "yes" /**< <adult> */
	EXTH_PRICE                  = 118 // As text, e.g. "4.99" /**< <srp> */
	EXTH_CURRENCY               = 119 // As text, e.g. "USD" /**< <srp currency="currency"> */
	EXTH_KF8BOUNDARY            = 121
	EXTH_FIXEDLAYOUT            = 122 /**< <fixed-layout> */
	EXTH_BOOKTYPE               = 123 /**< <book-type> */
	EXTH_ORIENTATIONLOCK        = 124 /**< <orientation-lock> */
	EXTH_COUNTRESOURCES         = 125
	EXTH_ORIGRESOLUTION         = 126 /**< <original-resolution> */
	EXTH_ZEROGUTTER             = 127 /**< <zero-gutter> */
	EXTH_ZEROMARGIN             = 128 /**< <zero-margin> */
	EXTH_KF8COVERURI            = 129
	EXTH_RESCOFFSET             = 131
	EXTH_REGIONMAGNI            = 132 /**< <region-mag> */

	EXTH_DICTNAME     = 200 // As text /**< <DictionaryVeryShortName> */
	EXTH_COVEROFFSET  = 201 // Add to first image field in Mobi Header to find PDB record containing the cover image/**< <EmbeddedCover> */
	EXTH_THUMBOFFSET  = 202 // Add to first image field in Mobi Header to find PDB record containing the thumbnail cover image
	EXTH_HASFAKECOVER = 203
	EXTH_CREATORSOFT  = 204 //Known Values: 1=mobigen, 2=Mobipocket Creator, 200=kindlegen (Windows), 201=kindlegen (Linux), 202=kindlegen (Mac).
	EXTH_CREATORMAJOR = 205
	EXTH_CREATORMINOR = 206
	EXTH_CREATORBUILD = 207
	EXTH_WATERMARK    = 208
	EXTH_TAMPERKEYS   = 209

	EXTH_FONTSIGNATURE = 300

	EXTH_CLIPPINGLIMIT  = 401 // Integer percentage of the text allowed to be clipped. Usually 10.
	EXTH_PUBLISHERLIMIT = 402
	EXTH_UNK403         = 403
	EXTH_TTSDISABLE     = 404 // 1 - Text to Speech disabled; 0 - Text to Speech enabled
	EXTH_UNK405         = 405 // 1 in this field seems to indicate a rental book
	EXTH_RENTAL         = 406 // If this field is removed from a rental, the book says it expired in 1969
	EXTH_UNK407         = 407
	EXTH_UNK450         = 450
	EXTH_UNK451         = 451
	EXTH_UNK452         = 452
	EXTH_UNK453         = 453

	EXTH_DOCTYPE         = 501 // PDOC - Personal Doc; EBOK - ebook; EBSP - ebook sample;
	EXTH_LASTUPDATE      = 502
	EXTH_UPDATEDTITLE    = 503
	EXTH_ASIN504         = 504 // ?? ASIN in this record.
	EXTH_TITLEFILEAS     = 508
	EXTH_CREATORFILEAS   = 517
	EXTH_PUBLISHERFILEAS = 522
	EXTH_LANGUAGE        = 524 /**< <dc:language> */
	EXTH_ALIGNMENT       = 525 // ?? horizontal-lr in this record /**< <primary-writing-mode> */
	EXTH_PAGEDIR         = 527
	EXTH_OVERRIDEFONTS   = 528 /**< <override-kindle-fonts> */
	EXTH_SORCEDESC       = 529
	EXTH_DICTLANGIN      = 531
	EXTH_DICTLANGOUT     = 532
	EXTH_UNK534          = 534
	EXTH_CREATORBUILDREV = 535
)

// EXTH Tag ID - Name - Type relationship
var ExthMeta = []mobiExthMeta{
	{0, 0, ""},
	{EXTH_SAMPLE, EXTH_TYPE_NUMERIC, "Sample"},
	{EXTH_STARTREADING, EXTH_TYPE_NUMERIC, "Start offset"},
	{EXTH_KF8BOUNDARY, EXTH_TYPE_NUMERIC, "K8 Boundary Offset"},
	{EXTH_COUNTRESOURCES, EXTH_TYPE_NUMERIC, "K8 Resources Count"}, // of , fonts, images
	{EXTH_RESCOFFSET, EXTH_TYPE_NUMERIC, "RESC Offset"},
	{EXTH_COVEROFFSET, EXTH_TYPE_NUMERIC, "Cover Offset"},
	{EXTH_THUMBOFFSET, EXTH_TYPE_NUMERIC, "Thumbnail Offset"},
	{EXTH_HASFAKECOVER, EXTH_TYPE_NUMERIC, "Has Fake Cover"},
	{EXTH_CREATORSOFT, EXTH_TYPE_NUMERIC, "Creator Software"},
	{EXTH_CREATORMAJOR, EXTH_TYPE_NUMERIC, "Creator Major Version"},
	{EXTH_CREATORMINOR, EXTH_TYPE_NUMERIC, "Creator Minor Version"},
	{EXTH_CREATORBUILD, EXTH_TYPE_NUMERIC, "Creator Build Number"},
	{EXTH_CLIPPINGLIMIT, EXTH_TYPE_NUMERIC, "Clipping Limit"},
	{EXTH_PUBLISHERLIMIT, EXTH_TYPE_NUMERIC, "Publisher Limit"},
	{EXTH_TTSDISABLE, EXTH_TYPE_NUMERIC, "Text-to-Speech Disabled"},
	{EXTH_RENTAL, EXTH_TYPE_NUMERIC, "Rental Indicator"},
	{EXTH_DRMSERVER, EXTH_TYPE_STRING, "DRM Server ID"},
	{EXTH_DRMCOMMERCE, EXTH_TYPE_STRING, "DRM Commerce ID"},
	{EXTH_DRMEBOOKBASE, EXTH_TYPE_STRING, "DRM Ebookbase Book ID"},
	{EXTH_TITLE, EXTH_TYPE_STRING, "Title"},
	{EXTH_AUTHOR, EXTH_TYPE_STRING, "Creator"},
	{EXTH_PUBLISHER, EXTH_TYPE_STRING, "Publisher"},
	{EXTH_IMPRINT, EXTH_TYPE_STRING, "Imprint"},
	{EXTH_DESCRIPTION, EXTH_TYPE_STRING, "Description"},
	{EXTH_ISBN, EXTH_TYPE_STRING, "ISBN"},
	{EXTH_SUBJECT, EXTH_TYPE_STRING, "Subject"},
	{EXTH_PUBLISHINGDATE, EXTH_TYPE_STRING, "Published"},
	{EXTH_REVIEW, EXTH_TYPE_STRING, "Review"},
	{EXTH_CONTRIBUTOR, EXTH_TYPE_STRING, "Contributor"},
	{EXTH_RIGHTS, EXTH_TYPE_STRING, "Rights"},
	{EXTH_SUBJECTCODE, EXTH_TYPE_STRING, "Subject Code"},
	{EXTH_TYPE, EXTH_TYPE_STRING, "Type"},
	{EXTH_SOURCE, EXTH_TYPE_STRING, "Source"},
	{EXTH_ASIN, EXTH_TYPE_STRING, "ASIN"},
	{EXTH_VERSION, EXTH_TYPE_STRING, "Version Number"},
	{EXTH_ADULT, EXTH_TYPE_STRING, "Adult"},
	{EXTH_PRICE, EXTH_TYPE_STRING, "Price"},
	{EXTH_CURRENCY, EXTH_TYPE_STRING, "Currency"},
	{EXTH_FIXEDLAYOUT, EXTH_TYPE_STRING, "Fixed Layout"},
	{EXTH_BOOKTYPE, EXTH_TYPE_STRING, "Book Type"},
	{EXTH_ORIENTATIONLOCK, EXTH_TYPE_STRING, "Orientation Lock"},
	{EXTH_ORIGRESOLUTION, EXTH_TYPE_STRING, "Original Resolution"},
	{EXTH_ZEROGUTTER, EXTH_TYPE_STRING, "Zero Gutter"},
	{EXTH_ZEROMARGIN, EXTH_TYPE_STRING, "Zero margin"},
	{EXTH_KF8COVERURI, EXTH_TYPE_STRING, "K8 Masthead/Cover Image"},
	{EXTH_REGIONMAGNI, EXTH_TYPE_STRING, "Region Magnification"},
	{EXTH_DICTNAME, EXTH_TYPE_STRING, "Dictionary Short Name"},
	{EXTH_WATERMARK, EXTH_TYPE_STRING, "Watermark"},
	{EXTH_DOCTYPE, EXTH_TYPE_STRING, "Document Type"},
	{EXTH_LASTUPDATE, EXTH_TYPE_STRING, "Last Update Time"},
	{EXTH_UPDATEDTITLE, EXTH_TYPE_STRING, "Updated Title"},
	{EXTH_ASIN504, EXTH_TYPE_STRING, "ASIN (504)"},
	{EXTH_TITLEFILEAS, EXTH_TYPE_STRING, "Title File As"},
	{EXTH_CREATORFILEAS, EXTH_TYPE_STRING, "Creator File As"},
	{EXTH_PUBLISHERFILEAS, EXTH_TYPE_STRING, "Publisher File As"},
	{EXTH_LANGUAGE, EXTH_TYPE_STRING, "Language"},
	{EXTH_ALIGNMENT, EXTH_TYPE_STRING, "Primary Writing Mode"},
	{EXTH_PAGEDIR, EXTH_TYPE_STRING, "Page Progression Direction"},
	{EXTH_OVERRIDEFONTS, EXTH_TYPE_STRING, "Override Kindle Fonts"},
	{EXTH_SORCEDESC, EXTH_TYPE_STRING, "Original Source description"},
	{EXTH_DICTLANGIN, EXTH_TYPE_STRING, "Dictionary Input Language"},
	{EXTH_DICTLANGOUT, EXTH_TYPE_STRING, "Dictionary output Language"},
	{EXTH_UNK534, EXTH_TYPE_STRING, "Unknown (534)"},
	{EXTH_CREATORBUILDREV, EXTH_TYPE_STRING, "Kindlegen BuildRev Number"},
	{EXTH_TAMPERKEYS, EXTH_TYPE_BINARY, "Tamper Proof Keys"},
	{EXTH_FONTSIGNATURE, EXTH_TYPE_BINARY, "Font Signature"},
	{EXTH_UNK403, EXTH_TYPE_BINARY, "Unknown (403)"},
	{EXTH_UNK405, EXTH_TYPE_BINARY, "Unknown (405)"},
	{EXTH_UNK407, EXTH_TYPE_BINARY, "Unknown (407)"},
	{EXTH_UNK450, EXTH_TYPE_BINARY, "Unknown (450)"},
	{EXTH_UNK451, EXTH_TYPE_BINARY, "Unknown (451)"},
	{EXTH_UNK452, EXTH_TYPE_BINARY, "Unknown (452)"},
	{EXTH_UNK453, EXTH_TYPE_BINARY, "Unknown (453)"}}

type mobiExth struct {
	Identifier   [4]uint8 `format:"string"`
	HeaderLenght uint32   // The length of the EXTH header, including the previous 4 bytes - but not including the final padding.
	RecordCount  uint32   // The number of records in the EXTH header. the rest of the EXTH header consists of repeated EXTH records to the end of the EXTH length.

	Records []mobiExthRecord // Lenght of RecordCount

	// []uint8 - lenght of X. Where X is the amount of bytes needed to reach multiples of 4 for the whole EXTH record

	// According to Wiki padding null bytes are not included into header lenght calculation, but from what
	// I see in mobi files, those bytes are included in total calculation.
}

type mobiExthRecord struct {
	RecordType   uint32 // Exth Record type. Just a number identifying what's stored in the record
	RecordLength uint32 // Length of EXTH record = L , including the 8 bytes in the type and length fields
	Value        []uint8
}

// Copy from https://github.com/bfabiszewski/libmobi/blob/f4f75982f0c00b592c418bfcf3f9920600e81573/src/util.c
type mobiExthMeta struct {
	ID   uint32
	Type ExthType
	Name string
}

func (w *mobiExth) GetHeaderLenght() int {
	elen := 12

	for _, k := range w.Records {
		elen += int(k.RecordLength)
	}

	Padding := elen % 4
	elen += Padding

	return elen
}

func (e *mobiExth) Add(recType uint32, Value interface{}) *mobiExth {
	e.RecordCount++

	var MetaType = getExthMetaByTag(recType)
	var ExthRec mobiExthRecord = mobiExthRecord{RecordType: recType}

	switch MetaType.Type {
	case EXTH_TYPE_BINARY:
		ExthRec.Value = Value.([]uint8)
	case EXTH_TYPE_NUMERIC:
		var castValue uint32
		switch Value.(type) {
		case int:
			castValue = uint32(Value.(int))
		case uint16:
			castValue = uint32(Value.(uint16))
		case uint32:
			castValue = uint32(Value.(uint32))
		case uint64:
			castValue = uint32(Value.(uint64))
		case int16:
			castValue = uint32(Value.(int16))
		case int32:
			castValue = uint32(Value.(int32))
		case int64:
			castValue = uint32(Value.(int64))
		default:
			panic("EXTH_TYPE_NUMERIC type is unsupported")
		}
		ExthRec.Value = int32ToBytes(castValue)
	case EXTH_TYPE_STRING:
		switch Value.(type) {
		case []uint8:
			ExthRec.Value = Value.([]uint8)
		case string:
			ExthRec.Value = []uint8(Value.(string))
		}
	default:
		panic("Unknown EXTH meta type")
	}

	ExthRec.RecordLength = uint32(8 + len(ExthRec.Value))
	e.Records = append(e.Records, ExthRec)
	return e
}
