package models

type ObjectsChunk struct {
	From          uint64
	To            uint64
	CountPerChunk uint64
	Objects       []interface{}
}
