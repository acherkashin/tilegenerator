package geo

// BaseGeometry is a geometry structure
type BaseGeometry struct {
	ID     int
	TypeID int
	Value  string
}

type BaseAttribute struct {
	Value    string
	Code     string
	ObjectID int
}
