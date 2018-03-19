package model

// JsonStruct is a Type interface for json marshalling
type JsonStruct interface {
	Unmarshal([]byte) (JsonStruct, error)
}
