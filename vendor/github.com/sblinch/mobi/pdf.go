package mobi

type mobiPDF struct {
	DatabaseName       [32]byte `format:"string"` //Database name. This name is 0 terminated
	FileAttributes     uint16
	Version            uint16 //File version
	CreationTime       uint32 `format:"date"` //Timestamp, according to wiki it's supposed to be in Mac format, but Mobi files that I see use Unix. Not sure if it's important.
	ModificationTime   uint32 `format:"date"` //Timestamp
	BackupTime         uint32 `format:"date"` //Timestamp
	ModificationNumber uint32
	AppInfo            uint32
	SortInfo           uint32
	Type               [4]byte `format:"string"` //BOOK
	Creator            [4]byte `format:"string"` //MOBI
	UniqueIDSeed       uint32  //Used internally to identify record
	NextRecordList     uint32  //Only used when in-memory on Palm OS. Always set to zero in stored files.
	RecordsNum         uint16  //Number of records in the file. Records are stored as array starting with 0. RecordsNum is total count of records, not last ID.
}
