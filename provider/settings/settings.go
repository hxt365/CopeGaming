package settings

type Range struct {
	Min uint16
	Max uint16
}

var (
	SinglePort                 int
	PortRange                  Range
	IceIpMap                   string
	DisableDefaultInterceptors bool

	VideoCodec string

	CoordinatorAddr string
)

func init() {
	SinglePort = 8443
	DisableDefaultInterceptors = false

	VideoCodec = "vpx"

	CoordinatorAddr = "localhost:8080"
}
