package mobi

const (
	INDX_TYPE_NORMAL     uint32 = 0
	INDX_TYPE_INFLECTION uint32 = 2
)

type mobiIndx struct {
	Identifier         [4]byte `format:"string"`
	HeaderLen          uint32
	Unk0               uint32
	Unk1               uint32 /* 1 when inflection is normal? */
	Indx_Type          uint32 /* 12: 0 - normal, 2 - inflection */
	Idxt_Offset        uint32 /* 20: IDXT offset */
	Idxt_Count         uint32 /* 24: entries count */
	Idxt_Encoding      uint32 /* 28: encoding */
	SetUnk2            uint32 //-1
	Idxt_Entry_Count   uint32 /* 36: total entries count */
	Ordt_Offset        uint32
	Ligt_Offset        uint32
	Ligt_Entries_Count uint32 /* 48: LIGT entries count */
	Cncx_Records_Count uint32 /* 52: CNCX entries count */
	Unk3               [108]byte
	Ordt_Type          uint32 /* 164: ORDT type */
	Ordt_Entries_Count uint32 /* 168: ORDT entries count */
	Ordt1_Offset       uint32 /* 172: ORDT1 offset */
	Ordt2_Offset       uint32 /* 176: ORDT2 offset */
	Tagx_Offset        uint32 /* 180: */
	Unk4               uint32 /* 184: */ /* ? Default index string offset ? */
	Unk5               uint32 /* 188: */ /* ? Default index string length ? */
}

type mobiIndxEntry struct {
	EntryID    uint8
	EntryValue uint32
}
