package ghastly

type Backend struct {
	Name                string
	Address             string
	Port                uint16
	UseSSL              bool
	ConnectTimeout      int
	FirstByteTimeout    int
	BetweenBytesTimeout int
	ErrorThreshold      int
	MaxConn             int
	Weight              int
	AutoLoadbalance     int
	RequestCondition    string
	Healthcheck         string
	version             *Version
}
