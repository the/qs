package process

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
)

type JSONURL struct {
	URL      string            `json:"url"`
	Scheme   string            `json:"scheme,omitempty"`
	Host     string            `json:"host,omitempty"`
	Path     string            `json:"path,omitempty"`
	Query    map[string]string `json:"query"`
	Fragment string            `json:"fragment,omitempty"`
}

func JSON(in io.Reader, out io.Writer, paramNames []string) {
	encoder := json.NewEncoder(out)
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		u, err := url.Parse(clean(scanner.Text()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		params := u.Query()
		host := u.Host
		if u.Opaque != "" {
			host = u.Opaque
		}
		j := JSONURL{
			URL:      u.String(),
			Scheme:   u.Scheme,
			Host:     host,
			Path:     u.Path,
			Fragment: u.Fragment,
			Query:    map[string]string{},
		}
		if len(paramNames) > 0 {
			for _, param := range paramNames {
				j.Query[param] = params.Get(param)
			}
		} else {
			for param := range params {
				j.Query[param] = params.Get(param)
			}
		}
		encoder.Encode(j)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
