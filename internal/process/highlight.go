package process

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
)

const (
	formatReset      = "\033[0m"
	formatBold       = "\033[1m"
	formatUnderscore = "\033[4m"
	formatReverse    = "\033[7m"
)

const (
	colorBlack  = "\033[30m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

var colors = []string{colorCyan, colorPurple, colorRed, colorYellow}

func Highlight(in io.Reader, out io.Writer, paramNames []string) {
	matcher := matcherForParams(paramNames)
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		u, err := url.Parse(clean(scanner.Text()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		formatter := matcher.Match(u.RawQuery)
		formatter.Format(u, out)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

type Matcher struct {
	Param string
	Color string
	Regex *regexp.Regexp
}

type ParamMatcher []Matcher

type Match struct {
	Start int
	End   int
	Color string
}

type ParamFormatter []Match

func matcherForParams(params []string) ParamMatcher {
	var matcher ParamMatcher

	for i, param := range params {
		re, err := regexp.Compile(fmt.Sprintf(`\b%s=[^&]+`, param))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
		matcher = append(matcher, Matcher{
			Param: param,
			Color: colors[i%len(colors)],
			Regex: re,
		})
	}

	return matcher
}

func (pm ParamMatcher) Match(qs string) ParamFormatter {
	var formatter ParamFormatter

	for _, matcher := range pm {
		matches := matcher.Regex.FindAllStringIndex(qs, -1)
		for _, match := range matches {
			formatter = append(formatter, Match{
				Start: match[0],
				End:   match[1],
				Color: matcher.Color,
			})
		}
	}

	sort.Sort(formatter)
	return formatter
}

func (f ParamFormatter) Format(u *url.URL, out io.Writer) {
	len := len(f)
	if len > 0 && out == os.Stdout {
		if u.Scheme != "" {
			if u.Opaque == "" {
				fmt.Fprint(out, u.Scheme, "://")
			} else {
				fmt.Fprint(out, u.Scheme, ":")
			}
		}
		fmt.Fprint(out, u.Opaque, u.Host, u.Path, "?")

		for i, match := range f {
			if i == 0 && match.Start > 0 {
				fmt.Fprint(out, u.RawQuery[0:match.Start])
			}
			equals := match.Start + strings.Index(u.RawQuery[match.Start:match.End], "=")
			fmt.Fprint(out, formatUnderscore, match.Color, u.RawQuery[match.Start:equals], formatReset)
			fmt.Fprint(out, match.Color, u.RawQuery[equals:match.End], formatReset)
			if i < len-1 {
				nextMatch := f[i+1]
				fmt.Fprint(out, u.RawQuery[match.End:nextMatch.Start])
			} else if i == len-1 {
				fmt.Fprint(out, u.RawQuery[match.End:])
			}
		}

		if u.Fragment != "" {
			fmt.Fprint(out, "#", u.Fragment)
		}
		fmt.Fprintln(out)
	} else {
		fmt.Fprintln(out, u.String())
	}
}

func (f ParamFormatter) Len() int           { return len(f) }
func (f ParamFormatter) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f ParamFormatter) Less(i, j int) bool { return f[i].Start < f[j].Start }
