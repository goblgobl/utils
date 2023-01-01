package buffer

type Config struct {
	Count uint16 `json:"count"`
	Min   uint32 `json:"min"`
	Max   uint32 `json:"max"`
}
