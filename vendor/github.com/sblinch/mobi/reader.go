package mobi

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"image"
	"io"
)

type MobiReader struct {
	Mobi

	fullName string		// full name of the book (from the main header)

	sample bool
	startReading int64 // Start offset //Position (4-byte offset) in file at which to open when first opened /**< Start reading */
	coverOffset int64
	coverLength int64
	thumbOffset int64
	thumbLength int64
	hasFakeCover bool
	creatorSoft int64 // Creator Software
	creatorMajor int64 // Creator Major Version
	creatorMinor int64 // Creator Minor Version
	creatorBuild int64 // Creator Build Number
	creatorBuildRev string // Kindlegen BuildRev Number
	clippingLimit int64 // Clipping Limit
	publisherLimit int64 // Publisher Limit
	ttsDisable bool // Text-to-Speech Disabled
	rental bool // Rental Indicator
	drmServer string // DRM Server ID
	drmCommerce string // DRM Commerce ID
	drmEbookbase string // DRM Ebookbase Book ID
	title string // Title
	authors []string // Creator
	publisher string // Publisher
	imprint string // Imprint
	description string // Description
	isbn string // ISBN
	subjects []string // Subject
	publishingDate string // Published
	review string // Review
	contributor string // Contributor
	rights string // Rights
	subjectCode string // Subject Code
	typeStr string // Type
	source string // Source
	asin string // ASIN Kindle Paperwhite labels books with "Personal" if they don't have this record.
	version string // Version Number
	adult bool // Adult Mobipocket Creator adds this if Adult only is checked on its GUI; contents: "yes" /**< <adult> */
	price string // Price // As text, e.g. "4.99" /**< <srp> */
	currency string // Currency // As text, e.g. "USD" /**< <srp currency="currency"> */
	fixedLayout string // Fixed Layout
	bookType string // Book Type
	orientationLock string // Orientation Lock
	origResolution string // Original Resolution
	zeroGutter string // Zero Gutter
	zeroMargin string // Zero margin
	kf8CoverUri string // K8 Masthead/Cover Image
	regionMagni string // Region Magnification
	dictName string // Dictionary Short Name
	watermark string // Watermark
	docType string // Document Type
	lastUpdate string // Last Update Time
	updatedTitle string // Updated Title
	asin504 string // ASIN (504)
	titleFileAs string // Title File As
	creatorFileAs string // Creator File As
	publisherFileAs string // Publisher File As
	language string // Language
	alignment string // Primary Writing Mode
	pageDir string // Page Progression Direction
	overrideFonts string // Override Kindle Fonts
	sourceDesc string // Original Source description
	dictLangIn string // Dictionary Input Language
	dictLangOut string // Dictionary output Language
}

func NewReader(filename string) (out *MobiReader, err error) {
	out = &MobiReader{}
	out.file, err = os.Open(filename)
	if err != nil {
		return nil, err
	}

	out.fileStat, err = out.file.Stat()
	if err != nil {
		return nil, err
	}

	return out, out.parse()
}

func (r *MobiReader) Close() error {
	if r.file == nil {
		return os.ErrClosed
	}

	if err := r.file.Close(); err != nil {
		return err
	}
	r.file = nil

	return nil
}

func (r *MobiReader) parse() (err error) {
	if err = r.parsePdf(); err != nil {
		return
	}

	if err = r.parsePdh(); err != nil {
		return
	}

	// Check if INDX offset is set + attempt to parse INDX
	if r.Header.IndxRecodOffset > 0 && r.Header.IndxRecodOffset != 4294967295 {
		err = r.parseIndexRecord(r.Header.IndxRecodOffset)
		if err != nil {
			return
		}
	}

	return
}

// parseHeader reads Palm Database Format header, and record offsets
func (r *MobiReader) parsePdf() error {
	//First we read PDF Header, this will help us parse subsequential data
	//binary.Read will take struct and fill it with data from mobi File
	err := binary.Read(r.file, binary.BigEndian, &r.Pdf)
	if err != nil {
		return err
	}

	if r.Pdf.RecordsNum < 1 {
		return errors.New("Number of records in this file is less than 1.")
	}

	r.Offsets = make([]mobiRecordOffset, r.Pdf.RecordsNum)
	err = binary.Read(r.file, binary.BigEndian, &r.Offsets)
	if err != nil {
		return err
	}

	//After the records offsets there's a 2 byte padding
	r.file.Seek(2, 1)

	return nil
}

// parsePdh processes record 0 that contains PalmDoc Header, Mobi Header and Exth meta data
func (r *MobiReader) parsePdh() error {
	// Palm Doc Header
	// Now we go onto reading record 0 that contains Palm Doc Header, Mobi Header, Exth Header...
	binary.Read(r.file, binary.BigEndian, &r.Pdh)

	// Check and see if there's a record encryption
	if r.Pdh.Encryption != 0 {
		return errors.New("Records are encrypted.")
	}

	// Mobi Header
	// Now it's time to read Mobi Header
	if r.matchMagic(magicMobi) {
		binary.Read(r.file, binary.BigEndian, &r.Header)
	} else {
		return errors.New("Can not find MOBI header. File might be corrupt.")
	}

	/*
	log.Printf("Full name: ")
	FullNameOffset      uint32   // Offset in record 0 (not from start of file) of the full name of the book
	FullNameLength      uint32   // Length in bytes of the full name of the book
	*/


	// Current header struct only reads 232 bytes. So if actual header lenght is greater, then we need to skip to Exth.
	Skip := int64(r.Header.HeaderLength) - int64(reflect.TypeOf(r.Header).Size())
	r.file.Seek(Skip, 1)

	// Exth Record
	// To check whenever there's EXTH record or not, we need to check and see if 6th bit of r.Header.ExthFlags is set.
	if hasBit(int(r.Header.ExthFlags), 6) {
		err := r.exthParse()

		if err != nil {
			return errors.New("Can not read EXTH record")
		}

		r.populateMetadata()
	}

	if r.Header.FullNameLength>0 {
		name := make([]byte, r.Header.FullNameLength)
		_, err := r.file.Seek(int64(r.Offsets[0].Offset+r.Header.FullNameOffset), 0)
		if err == nil {
			_, err = r.file.Read(name)
			if err == nil {
				r.fullName = string(name)
			}
		}
	}


	return nil
}

func (r *MobiReader) parseIndexRecord(n uint32) error {
	_, err := r.OffsetToRecord(n)
	if err != nil {
		return err
	}

	RecPos, _ := r.file.Seek(0, 1)

	if !r.matchMagic(magicIndx) {
		return errors.New("Index record not found at specified at given offset")
	}
	//fmt.Printf("Index %s %v\n", r.Peek(4), RecLen)

	//if len(r.Indx) == 0 {
	r.Indx = append(r.Indx, mobiIndx{})
	//}

	idx := &r.Indx[len(r.Indx)-1]

	err = binary.Read(r.file, binary.BigEndian, idx)
	if err != nil {
		return err
	}

	/* Tagx Record Parsing + Last CNCX */
	if idx.Tagx_Offset != 0 {
		_, err = r.file.Seek(RecPos+int64(idx.Tagx_Offset), 0)
		if err != nil {
			return err
		}

		err = r.parseTagx()
		if err != nil {
			return err
		}

		// Last CNCX record follows TAGX
		if idx.Cncx_Records_Count > 0 {
			r.Cncx = mobiCncx{}
			binary.Read(r.file, binary.BigEndian, &r.Cncx.Len)

			r.Cncx.Id = make([]uint8, r.Cncx.Len)
			binary.Read(r.file, binary.LittleEndian, &r.Cncx.Id)
			r.file.Seek(1, 1) //Skip 0x0 termination

			binary.Read(r.file, binary.BigEndian, &r.Cncx.NCX_Count)

			// PrintStruct(r.Cncx)
		}
	}

	/* Ordt Record Parsing */
	if idx.Idxt_Encoding == MOBI_ENC_UTF16 || idx.Ordt_Entries_Count > 0 {
		// ignore
		//return errors.New("ORDT parser not implemented")
	}

	/* Ligt Record Parsing */
	if idx.Ligt_Entries_Count > 0 {
		// ignore
		//return errors.New("LIGT parser not implemented")
	}

	/* Idxt Record Parsing */
	if idx.Idxt_Count > 0 {
		_, err = r.file.Seek(RecPos+int64(idx.Idxt_Offset), 0)
		if err != nil {
			return err
		}

		err = r.parseIdxt(idx.Idxt_Count)
		if err != nil {
			return err
		}
	}

	//CNCX Data?
	var Count = 0
	if idx.Indx_Type == INDX_TYPE_NORMAL {
		//r.file.Seek(RecPos+int64(idx.HeaderLen), 0)

		var PTagxLen = []uint8{0}
		for i, offset := range r.Idxt.Offset {
			r.file.Seek(RecPos+int64(offset), 0)

			// Read Byte containing the lenght of a label
			r.file.Read(PTagxLen)

			// Read label
			PTagxLabel := make([]uint8, PTagxLen[0])
			r.file.Read(PTagxLabel)

			PTagxLen1 := uint16(idx.Idxt_Offset) - r.Idxt.Offset[i]
			if i+1 < len(r.Idxt.Offset) {
				PTagxLen1 = r.Idxt.Offset[i+1] - r.Idxt.Offset[i]
			}

			PTagxData := make([]uint8, PTagxLen1)
			r.file.Read(PTagxData)
			//fmt.Printf("\n------ %v --------\n", i)
			r.parsePtagx(PTagxData)
			Count++
			//fmt.Printf("Len: %v | Label: %s | %v\n", PTagxLen, PTagxLabel, Count)
		}
	}

	// Check next record
	//r.OffsetToRecord(n + 1)

	//
	// Process remaining INDX records
	if idx.Indx_Type == INDX_TYPE_INFLECTION {
		r.parseIndexRecord(n + 1)
	}
	//fmt.Printf("%s", )
	// Read Tagx
	//		if idx.Tagx_Offset > 0 {
	//			err := r.parseTagx()
	//			if err != nil {
	//				return err
	//			}
	//		}

	return nil
}

// matchMagic matches next N bytes (based on lenght of magic word)
func (r *MobiReader) matchMagic(magic mobiMagicType) bool {
	if r.peek(len(magic)).Magic() == magic {
		return true
	}
	return false
}

// peek returns next N bytes without advancing the reader.
func (r *MobiReader) peek(n int) Peeker {
	buf := make([]uint8, n)
	r.file.Read(buf)
	r.file.Seek(int64(n)*-1, 1)
	return buf
}

// Parse reads/parses Exth meta data records from file
func (r *MobiReader) exthParse() error {
	// If next 4 bytes are not EXTH then we have a problem
	if !r.matchMagic(magicExth) {
		return errors.New("Currect reading position does not contain EXTH record")
	}

	binary.Read(r.file, binary.BigEndian, &r.Exth.Identifier)
	binary.Read(r.file, binary.BigEndian, &r.Exth.HeaderLenght)
	binary.Read(r.file, binary.BigEndian, &r.Exth.RecordCount)

	r.Exth.Records = make([]mobiExthRecord, r.Exth.RecordCount)
	for i := range r.Exth.Records {
		binary.Read(r.file, binary.BigEndian, &r.Exth.Records[i].RecordType)
		binary.Read(r.file, binary.BigEndian, &r.Exth.Records[i].RecordLength)

		r.Exth.Records[i].Value = make([]uint8, r.Exth.Records[i].RecordLength-8)

		Tag := getExthMetaByTag(r.Exth.Records[i].RecordType)
		switch Tag.Type {
		case EXTH_TYPE_BINARY:
			binary.Read(r.file, binary.BigEndian, &r.Exth.Records[i].Value)
			//			fmt.Printf("%v: %v\n", Tag.Name, r.Exth.Records[i].Value)
		case EXTH_TYPE_STRING:
			binary.Read(r.file, binary.LittleEndian, &r.Exth.Records[i].Value)
			//			fmt.Printf("%v: %s\n", Tag.Name, r.Exth.Records[i].Value)
		case EXTH_TYPE_NUMERIC:
			binary.Read(r.file, binary.BigEndian, &r.Exth.Records[i].Value)
			//			fmt.Printf("%v: %d\n", Tag.Name, binary.BigEndian.Uint32(r.Exth.Records[i].Value))
		}
	}

	return nil
}

// OffsetToRecord set s reading position to record N, returns total record lenght
func (r *MobiReader) OffsetToRecord(nu uint32) (uint32, error) {
	n := int(nu)
	if n > int(r.Pdf.RecordsNum)-1 {
		return 0, fmt.Errorf("Record ID requested (%d) is greater than total amount of records (%d)",n,int(r.Pdf.RecordsNum))
	}

	RecLen := uint32(0)
	if n+1 < int(r.Pdf.RecordsNum) {
		RecLen = r.Offsets[n+1].Offset
	} else {
		RecLen = uint32(r.fileStat.Size())
	}

	_, err := r.file.Seek(int64(r.Offsets[n].Offset), 0)

	return RecLen - r.Offsets[n].Offset, err
}

func (r *MobiReader) parseTagx() error {
	if !r.matchMagic(magicTagx) {
		return errors.New("TAGX record not found at given offset.")
	}

	r.Tagx = mobiTagx{}

	binary.Read(r.file, binary.BigEndian, &r.Tagx.Identifier)
	binary.Read(r.file, binary.BigEndian, &r.Tagx.HeaderLenght)
	if r.Tagx.HeaderLenght < 12 {
		return errors.New("TAGX record too short")
	}
	binary.Read(r.file, binary.BigEndian, &r.Tagx.ControlByteCount)

	TagCount := (r.Tagx.HeaderLenght - 12) / 4
	r.Tagx.Tags = make([]mobiTagxTags, TagCount)

	for i := 0; i < int(TagCount); i++ {
		err := binary.Read(r.file, binary.BigEndian, &r.Tagx.Tags[i])
		if err != nil {
			return err
		}
	}

	//fmt.Println("TagX called")
	// PrintStruct(r.Tagx)

	return nil
}

func (r *MobiReader) parseIdxt(IdxtCount uint32) error {
	//fmt.Println("parseIdxt called")
	if !r.matchMagic(magicIdxt) {
		return errors.New("IDXT record not found at given offset.")
	}

	binary.Read(r.file, binary.BigEndian, &r.Idxt.Identifier)

	r.Idxt.Offset = make([]uint16, IdxtCount)

	binary.Read(r.file, binary.BigEndian, &r.Idxt.Offset)
	//for id, _ := range r.Idxt.Offset {
	//	binary.Read(r.Buffer, binary.BigEndian, &r.Idxt.Offset[id])
	//}

	//Skip two bytes? Or skip necessary amount to reach total lenght in multiples of 4?
	r.file.Seek(2, 1)

	// PrintStruct(r.Idxt)
	return nil
}

func (r *MobiReader) parsePtagx(data []byte) {
	//control_byte_count
	//tagx
	control_bytes := data[:r.Tagx.ControlByteCount]
	data = data[r.Tagx.ControlByteCount:]

	var Ptagx []mobiPTagx //= make([]mobiPTagx, r.Tagx.TagCount())

	for _, x := range r.Tagx.Tags {
		if x.Control_Byte == 0x01 {
			control_bytes = control_bytes[1:]
			continue
		}

		value := control_bytes[0] & x.Bitmask
		if value != 0 {
			var value_count uint32
			var value_bytes uint32

			if value == x.Bitmask {
				if setBits[x.Bitmask] > 1 {
					// If all bits of masked value are set and the mask has more
					// than one bit, a variable width value will follow after
					// the control bytes which defines the length of bytes (NOT
					// the value count!) which will contain the corresponding
					// variable width values.
					var consumed uint32
					value_bytes, consumed = vwiDec(data, true)
					//fmt.Printf("\nConsumed %v", consumed)
					data = data[consumed:]
				} else {
					value_count = 1
				}
			} else {
				mask := x.Bitmask
				for {
					if mask&1 != 0 {
						//fmt.Printf("Break")
						break
					}
					mask >>= 1
					value >>= 1
				}
				value_count = uint32(value)
			}

			Ptagx = append(Ptagx, mobiPTagx{x.Tag, x.TagNum, value_count, value_bytes})
			//						ptagx[ptagx_count].tag = tagx->tags[i].tag;
			//       ptagx[ptagx_count].tag_value_count = tagx->tags[i].values_count;
			//       ptagx[ptagx_count].value_count = value_count;
			//       ptagx[ptagx_count].value_bytes = value_bytes;

			//fmt.Printf("TAGX %v %v VC:%v VB:%v\n", x.Tag, x.TagNum, value_count, value_bytes)
		}
	}
	//fmt.Printf("%+v", Ptagx)
	var IndxEntry []mobiIndxEntry
	for _, x := range Ptagx {
		var values []uint32

		if x.Value_Count != 0 {
			// Read value_count * values_per_entry variable width values.
			//fmt.Printf("\nDec: ")
			for i := 0; i < int(x.Value_Count)*int(x.Tag_Value_Count); i++ {
				byts, consumed := vwiDec(data, true)
				data = data[consumed:]

				values = append(values, byts)
				IndxEntry = append(IndxEntry, mobiIndxEntry{x.Tag, byts})
				//fmt.Printf("%v %s: %v ", iz, tagEntryMap[x.Tag], byts)
			}
		} else {
			// Convert value_bytes to variable width values.
			total_consumed := 0
			for {
				if total_consumed < int(x.Value_Bytes) {
					byts, consumed := vwiDec(data, true)
					data = data[consumed:]

					total_consumed += int(consumed)

					values = append(values, byts)
					IndxEntry = append(IndxEntry, mobiIndxEntry{x.Tag, byts})
				} else {
					break
				}
			}
			if total_consumed != int(x.Value_Bytes) {
				panic("Error not enough bytes are consumed. Consumed " + strconv.Itoa(total_consumed) + " out of " + strconv.Itoa(int(x.Value_Bytes)))
			}
		}
	}
	//fmt.Println("---------------------------")
}

func rvToInt(v []uint8) int64 {
	switch(len(v)) {
		case 1:
			return int64(v[0])
		case 2:
			return int64(binary.BigEndian.Uint16([]byte(v)))
		case 4:
			return int64(binary.BigEndian.Uint32([]byte(v)))
		case 8:
			return int64(binary.BigEndian.Uint64([]byte(v)))
		default:
			return 0
	}
}

func rvToBool(v []uint8) bool {
	return rvToInt(v) != 0
}

func rvToStringAppend(v []uint8, existing string) string {
	if len(existing)>0 {
		return existing+"; "+string(v)
	} else {
		return string(v)
	}
}

func (r *MobiReader) rvToImageOffset(v []uint8) (int64, int64) {
	imageOffset := int64(0)
	imageLength := int64(0)

	imagePDBOffset := r.Header.FirstImageIndex + uint32(rvToInt(v))

	n := int(imagePDBOffset)
	if n <= int(r.Pdf.RecordsNum)-1 {
		imageOffset = int64(r.Offsets[n].Offset)
		if n+1 < int(r.Pdf.RecordsNum) {
			imageLength = int64(r.Offsets[n+1].Offset) - imageOffset
		} else {
			imageLength = r.fileStat.Size() - imageOffset
		}
	}

	return imageOffset, imageLength
}

func (r *MobiReader) populateMetadata() {
	for _, rec := range r.Exth.Records {
		switch rec.RecordType {
		case EXTH_COVEROFFSET:
			r.coverOffset, r.coverLength = r.rvToImageOffset(rec.Value)
		case EXTH_THUMBOFFSET:
			r.thumbOffset, r.thumbLength = r.rvToImageOffset(rec.Value)
		case EXTH_SAMPLE:
			r.sample = rvToBool(rec.Value)
		case EXTH_STARTREADING:
			r.startReading = rvToInt(rec.Value)
		case EXTH_HASFAKECOVER:
			r.hasFakeCover = rvToBool(rec.Value)
		case EXTH_CREATORSOFT:
			r.creatorSoft = rvToInt(rec.Value)
		case EXTH_CREATORMAJOR:
			r.creatorMajor = rvToInt(rec.Value)
		case EXTH_CREATORMINOR:
			r.creatorMinor = rvToInt(rec.Value)
		case EXTH_CREATORBUILD:
			r.creatorBuild = rvToInt(rec.Value)
		case EXTH_CREATORBUILDREV:
			r.creatorBuildRev = rvToStringAppend(rec.Value,r.creatorBuildRev)
		case EXTH_CLIPPINGLIMIT:
			r.clippingLimit = rvToInt(rec.Value)
		case EXTH_PUBLISHERLIMIT:
			r.publisherLimit = rvToInt(rec.Value)
		case EXTH_TTSDISABLE:
			r.ttsDisable = rvToBool(rec.Value)
		case EXTH_RENTAL:
			r.rental = rvToBool(rec.Value)
		case EXTH_DRMSERVER:
			r.drmServer = rvToStringAppend(rec.Value,r.drmServer)
		case EXTH_DRMCOMMERCE:
			r.drmCommerce = rvToStringAppend(rec.Value,r.drmCommerce)
		case EXTH_DRMEBOOKBASE:
			r.drmEbookbase = rvToStringAppend(rec.Value,r.drmEbookbase)
		case EXTH_TITLE:
			r.title = rvToStringAppend(rec.Value,r.title)
		case EXTH_AUTHOR:
			r.authors = append(r.authors,string(rec.Value))
		case EXTH_PUBLISHER:
			r.publisher = rvToStringAppend(rec.Value,r.publisher)
		case EXTH_IMPRINT:
			r.imprint = rvToStringAppend(rec.Value,r.imprint)
		case EXTH_DESCRIPTION:
			r.description = rvToStringAppend(rec.Value,r.description)
		case EXTH_ISBN:
			r.isbn = rvToStringAppend(rec.Value,r.isbn)
		case EXTH_SUBJECT:
			r.subjects = append(r.subjects,string(rec.Value))
		case EXTH_PUBLISHINGDATE:
			r.publishingDate = rvToStringAppend(rec.Value,r.publishingDate)
		case EXTH_REVIEW:
			r.review = rvToStringAppend(rec.Value,r.review)
		case EXTH_CONTRIBUTOR:
			r.contributor = rvToStringAppend(rec.Value,r.contributor)
		case EXTH_RIGHTS:
			r.rights = rvToStringAppend(rec.Value,r.rights)
		case EXTH_SUBJECTCODE:
			r.subjectCode = rvToStringAppend(rec.Value,r.subjectCode)
		case EXTH_TYPE:
			r.typeStr = rvToStringAppend(rec.Value,r.typeStr)
		case EXTH_SOURCE:
			r.source = rvToStringAppend(rec.Value,r.source)
		case EXTH_ASIN:
			r.asin = rvToStringAppend(rec.Value,r.asin)
		case EXTH_VERSION:
			r.version = rvToStringAppend(rec.Value,r.version)
		case EXTH_ADULT:
			r.adult = rvToBool(rec.Value)
		case EXTH_PRICE:
			r.price = rvToStringAppend(rec.Value,r.price)
		case EXTH_CURRENCY:
			r.currency = rvToStringAppend(rec.Value,r.currency)
		case EXTH_FIXEDLAYOUT:
			r.fixedLayout = rvToStringAppend(rec.Value,r.fixedLayout)
		case EXTH_BOOKTYPE:
			r.bookType = rvToStringAppend(rec.Value,r.bookType)
		case EXTH_ORIENTATIONLOCK:
			r.orientationLock = rvToStringAppend(rec.Value,r.orientationLock)
		case EXTH_ORIGRESOLUTION:
			r.origResolution = rvToStringAppend(rec.Value,r.origResolution)
		case EXTH_ZEROGUTTER:
			r.zeroGutter = rvToStringAppend(rec.Value,r.zeroGutter)
		case EXTH_ZEROMARGIN:
			r.zeroMargin = rvToStringAppend(rec.Value,r.zeroMargin)
		case EXTH_KF8COVERURI:
			r.kf8CoverUri = rvToStringAppend(rec.Value,r.kf8CoverUri)
		case EXTH_REGIONMAGNI:
			r.regionMagni = rvToStringAppend(rec.Value,r.regionMagni)
		case EXTH_DICTNAME:
			r.dictName = rvToStringAppend(rec.Value,r.dictName)
		case EXTH_WATERMARK:
			r.watermark = rvToStringAppend(rec.Value,r.watermark)
		case EXTH_DOCTYPE:
			r.docType = rvToStringAppend(rec.Value,r.docType)
		case EXTH_LASTUPDATE:
			r.lastUpdate = rvToStringAppend(rec.Value,r.lastUpdate)
		case EXTH_UPDATEDTITLE:
			r.updatedTitle = rvToStringAppend(rec.Value,r.updatedTitle)
		case EXTH_ASIN504:
			r.asin504 = rvToStringAppend(rec.Value,r.asin504)
		case EXTH_TITLEFILEAS:
			r.titleFileAs = rvToStringAppend(rec.Value,r.titleFileAs)
		case EXTH_CREATORFILEAS:
			r.creatorFileAs = rvToStringAppend(rec.Value,r.creatorFileAs)
		case EXTH_PUBLISHERFILEAS:
			r.publisherFileAs = rvToStringAppend(rec.Value,r.publisherFileAs)
		case EXTH_LANGUAGE:
			r.language = rvToStringAppend(rec.Value,r.language)
		case EXTH_ALIGNMENT:
			r.alignment = rvToStringAppend(rec.Value,r.alignment)
		case EXTH_PAGEDIR:
			r.pageDir = rvToStringAppend(rec.Value,r.pageDir)
		case EXTH_OVERRIDEFONTS:
			r.overrideFonts = rvToStringAppend(rec.Value,r.overrideFonts)
		case EXTH_SORCEDESC:
			r.sourceDesc = rvToStringAppend(rec.Value,r.sourceDesc)
		case EXTH_DICTLANGIN:
			r.dictLangIn = rvToStringAppend(rec.Value,r.dictLangIn)
		case EXTH_DICTLANGOUT:
			r.dictLangOut = rvToStringAppend(rec.Value,r.dictLangOut)
		}
	}

}

func (r *MobiReader) BestTitle() string {
	if len(r.updatedTitle)>0 {
		return r.updatedTitle

	} else if len(r.title)>0 {
		return r.title

	} else {
		return r.fullName
	}
}

func (r *MobiReader) FullName() string {
	return r.fullName
}

func (r *MobiReader) Sample() bool {
	return r.sample
}

func (r *MobiReader) StartReading() int64 {
	return r.startReading
}

func (r *MobiReader) CoverOffsetLength() (int64, int64) {
	return r.coverOffset, r.coverLength
}

func (r *MobiReader) image(offset, length int64) (image.Image, error) {
	if r.file == nil {
		return nil, os.ErrClosed
	}

	var i image.Image
	var err error

	if _, err := r.file.Seek(offset, 0); err != nil {
		return nil, err
	}

	ltd := io.LimitReader(r.file,length)
	if i, _, err = image.Decode(ltd); err != nil {
		return nil, err
	}

	return i, nil

}

func (r *MobiReader) HasCover() bool {
	return r.coverOffset > 0
}

func (r *MobiReader) Cover() (image.Image, error) {
	return r.image(r.CoverOffsetLength())
}

func (r *MobiReader) ThumbnailOffsetLength() (int64, int64) {
	return r.thumbOffset, r.thumbLength
}

func (r *MobiReader) HasThumbnail() bool {
	return r.thumbOffset > 0
}

func (r *MobiReader) Thumbnail() (image.Image, error) {
	return r.image(r.ThumbnailOffsetLength())
}


func (r *MobiReader) HasFakeCover() bool {
	return r.hasFakeCover
}

func (r *MobiReader) CreatorSoft() int64 {
	return r.creatorSoft
}

func (r *MobiReader) CreatorMajor() int64 {
	return r.creatorMajor
}

func (r *MobiReader) CreatorMinor() int64 {
	return r.creatorMinor
}

func (r *MobiReader) CreatorBuild() int64 {
	return r.creatorBuild
}

func (r *MobiReader) CreatorBuildRev() string {
	return r.creatorBuildRev
}

func (r *MobiReader) ClippingLimit() int64 {
	return r.clippingLimit
}

func (r *MobiReader) PublisherLimit() int64 {
	return r.publisherLimit
}

func (r *MobiReader) TtsDisable() bool {
	return r.ttsDisable
}

func (r *MobiReader) Rental() bool {
	return r.rental
}

func (r *MobiReader) DrmServer() string {
	return r.drmServer
}

func (r *MobiReader) DrmCommerce() string {
	return r.drmCommerce
}

func (r *MobiReader) DrmEbookbase() string {
	return r.drmEbookbase
}

func (r *MobiReader) Title() string {
	return r.title
}

func (r *MobiReader) Authors() []string {
	return r.authors
}

func (r *MobiReader) Publisher() string {
	return r.publisher
}

func (r *MobiReader) Imprint() string {
	return r.imprint
}

func (r *MobiReader) Description() string {
	return r.description
}

func (r *MobiReader) Isbn() string {
	return r.isbn
}

func (r *MobiReader) Subjects() []string {
	return r.subjects
}

func (r *MobiReader) PublishingDate() string {
	return r.publishingDate
}

func (r *MobiReader) Review() string {
	return r.review
}

func (r *MobiReader) Contributor() string {
	return r.contributor
}

func (r *MobiReader) Rights() string {
	return r.rights
}

func (r *MobiReader) SubjectCode() string {
	return r.subjectCode
}

func (r *MobiReader) TypeStr() string {
	return r.typeStr
}

func (r *MobiReader) Source() string {
	return r.source
}

func (r *MobiReader) Asin() string {
	return r.asin
}

func (r *MobiReader) Version() string {
	return r.version
}

func (r *MobiReader) Adult() bool {
	return r.adult
}

func (r *MobiReader) Price() string {
	return r.price
}

func (r *MobiReader) Currency() string {
	return r.currency
}

func (r *MobiReader) FixedLayout() string {
	return r.fixedLayout
}

func (r *MobiReader) BookType() string {
	return r.bookType
}

func (r *MobiReader) OrientationLock() string {
	return r.orientationLock
}

func (r *MobiReader) OrigResolution() string {
	return r.origResolution
}

func (r *MobiReader) ZeroGutter() string {
	return r.zeroGutter
}

func (r *MobiReader) ZeroMargin() string {
	return r.zeroMargin
}

func (r *MobiReader) Kf8CoverUri() string {
	return r.kf8CoverUri
}

func (r *MobiReader) RegionMagni() string {
	return r.regionMagni
}

func (r *MobiReader) DictName() string {
	return r.dictName
}

func (r *MobiReader) Watermark() string {
	return r.watermark
}

func (r *MobiReader) DocType() string {
	return r.docType
}

func (r *MobiReader) LastUpdate() string {
	return r.lastUpdate
}

func (r *MobiReader) UpdatedTitle() string {
	return r.updatedTitle
}

func (r *MobiReader) Asin504() string {
	return r.asin504
}

func (r *MobiReader) TitleFileAs() string {
	return r.titleFileAs
}

func (r *MobiReader) CreatorFileAs() string {
	return r.creatorFileAs
}

func (r *MobiReader) PublisherFileAs() string {
	return r.publisherFileAs
}

func (r *MobiReader) Language() string {
	return r.language
}

func (r *MobiReader) Alignment() string {
	return r.alignment
}

func (r *MobiReader) PageDir() string {
	return r.pageDir
}

func (r *MobiReader) OverrideFonts() string {
	return r.overrideFonts
}

func (r *MobiReader) SourceDesc() string {
	return r.sourceDesc
}

func (r *MobiReader) DictLangIn() string {
	return r.dictLangIn
}

func (r *MobiReader) DictLangOut() string {
	return r.dictLangOut
}

