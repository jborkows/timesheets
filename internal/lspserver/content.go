package lspserver

type content struct {
	fileLines map[string][]string
}

func newContent() *content {
	return &content{
		fileLines: map[string][]string{},
	}
}

func (c *content) put(uri string, lines []string) {
	if len(lines) > 0 {
		if len(lines[len(lines)-1]) == 0 {
			lines = lines[:len(lines)-1]
		}
	}

	c.fileLines[uri] = lines
}

func (c *content) get(uri string) []string {
	return c.fileLines[uri]
}
