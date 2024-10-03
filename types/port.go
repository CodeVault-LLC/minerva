package types

// Port is a struct that represents a port along with what its commonly used as: (e.g. 80 -> http)
type Port struct {
	Port uint16
	Name string
	State string
}
