package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/fatih/color"
)

type LogcatLine struct {
	Type string
	Tag  string
	Pid  uint64
}

const (
	MATCH_V1 = `(\d\d.\d\d\s\d\d:\d\d:\d\d.\d\d\d)\s*(\d*)\s*(\d*)\s([VDIWEFS])\s(.*?)(:\s)(.*)`
	MATCH_V2 = `([VDIWEFS])/(.*?)\((\s*|)(\d*)\):.`
)

func matchV1(inputLine string) *LogcatLine {
	exp := regexp.MustCompile(MATCH_V1)
	m := exp.FindStringSubmatch(inputLine)
	if len(m) > 0 {
		var i uint64
		i, _ = strconv.ParseUint(m[2], 10, 64)
		return &LogcatLine{
			Type: m[4],
			Tag:  m[5],
			Pid:  i,
		}
	}
	return nil
}

func matchV2(inputLine string) *LogcatLine {
	exp := regexp.MustCompile(MATCH_V2)
	m := exp.FindStringSubmatch(inputLine)
	if len(m) > 0 {
		var i uint64
		i, _ = strconv.ParseUint(m[3], 10, 64)
		return &LogcatLine{
			Type: m[1],
			Tag:  m[2],
			Pid:  i,
		}
	}
	return nil
}

func detectFormat(inputLine string) *LogcatLine {
	if line := matchV1(inputLine); line != nil {
		return line
	}
	if line := matchV2(inputLine); line != nil {
		return line
	}

	return &LogcatLine{
		Type: "",
		Tag:  "",
		Pid:  0,
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		kind := detectFormat(text)
		switch kind.Type {
		case "D":
			color.Green(text)
			break
		case "W":
			color.Yellow(text)
			break
		case "E":
			color.Red(text)
			break
		case "I":
			color.Cyan(text)
		case "":
			fmt.Print(text)
		}
	}
}
