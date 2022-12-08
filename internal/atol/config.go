package atol

type Config struct {
	ApiURL          string `toml:"api_url"`
	RequestTemplate string `toml:"request_template"`
}
