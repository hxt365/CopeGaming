package settings

var (
	AllowedOrigins   []string
	AllowedWSOrigins []string
)

func init() {
	AllowedOrigins = []string{"http://localhost:3000"}
	AllowedWSOrigins = []string{"*"}
}
