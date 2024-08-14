package storage

type StorageType int

const (
    TypeUnknown StorageType = iota
    TypeMemory
    TypeFile
    TypeDatabase
)

func (st StorageType) String() string {
    return [...]string{"unknown", "memory", "file", "database"}[st]
}
