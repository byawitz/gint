package configurator

type Config struct {
	preset  string
	rules   []Rule
	exclude []string
	notName []string
	notPath []string
}
