package mobi

const (
	TagEntry_END                uint8 = 0
	TagEntry_Pos                      = 1  // NCX | Position offset for the beginning of NCX record (filepos) Ex: Beginning of a chapter
	TagEntry_Len                      = 2  // NCX | Record lenght. Ex: Chapter lenght
	TagEntry_NameOffset               = 3  // NCX | Label text offset in CNCX
	TagEntry_DepthLvl                 = 4  // NCX | Depth/Level of CNCX
	TagEntry_KOffs                    = 5  // NCX | kind CNCX offset
	TagEntry_PosFid                   = 6  // NCX | pos:fid
	TagEntry_Parent                   = 21 // NCX | Parent
	TagEntry_Child1                   = 22 // NCX | First child
	TagEntry_ChildN                   = 23 // NCX | Last child
	TagEntry_ImageIndex               = 69
	TagEntry_DescOffset               = 70 // Description offset in cncx
	TagEntry_AuthorOffset             = 71 // Author offset in cncx
	TagEntry_ImageCaptionOffset       = 72 // Image caption offset in cncx
	TagEntry_ImgAttrOffset            = 73 // Image attribution offset in cncx
)

var tagEntryMap = map[uint8]string{
	TagEntry_Pos:                "Offset",
	TagEntry_Len:                "Lenght",
	TagEntry_NameOffset:         "Label",
	TagEntry_DepthLvl:           "Depth",
	TagEntry_KOffs:              "Kind",
	TagEntry_PosFid:             "Pos:Fid",
	TagEntry_Parent:             "Parent",
	TagEntry_Child1:             "First Child",
	TagEntry_ChildN:             "Last Child",
	TagEntry_ImageIndex:         "Image Index",
	TagEntry_DescOffset:         "Description",
	TagEntry_AuthorOffset:       "Author",
	TagEntry_ImageCaptionOffset: "Image Caption Offset",
	TagEntry_ImgAttrOffset:      "Image Attr Offset"}

type mobiPTagx struct {
	Tag             uint8
	Tag_Value_Count uint8
	Value_Count     uint32
	Value_Bytes     uint32
}
