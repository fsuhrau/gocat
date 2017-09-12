package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/fatih/color"
)

type Process struct {
	Tag string
	Pid uint64
}

type LogcatLine struct {
	Type string
	Tag  string
	Pid  uint64
}

const (
	MATCH_V1         = `(\d\d.\d\d\s\d\d:\d\d:\d\d.\d\d\d)\s*(\d*)\s*(\d*)\s([VDIWEFS])\s(.*?)(:\s)(.*)`
	MATCH_V2         = `([VDIWEFS])/(.*?)\((\s*|)(\d*)\):.`
	MATCH_PROCESS_V1 = `\d\d.\d\d\s\d\d:\d\d:\d\d.\d\d\d\s*\d*\s*\d*\s[VDIWEFS]\sActivityManager:\sStart\sproc\s(\d*):(.*?)/.*`
	MATCH_PROCESS_V2 = `.*Start\sproc\s(.*)\sfor\sactivity\s.*:\spid=(\d*)\s.*`
)

func matchProcessV1(inputLine string) *Process {
	exp := regexp.MustCompile(MATCH_PROCESS_V1)
	m := exp.FindStringSubmatch(inputLine)
	if len(m) > 0 {
		var i uint64
		i, _ = strconv.ParseUint(m[1], 10, 64)
		return &Process{
			Tag: m[2],
			Pid: i,
		}
	}
	return nil
}

func matchProcessV2(inputLine string) *Process {
	exp := regexp.MustCompile(MATCH_PROCESS_V2)
	m := exp.FindStringSubmatch(inputLine)
	if len(m) > 0 {
		var i uint64
		i, _ = strconv.ParseUint(m[2], 10, 64)
		return &Process{
			Tag: m[1],
			Pid: i,
		}
	}
	return nil
}

func process(inputLine string) *Process {
	if proc := matchProcessV1(inputLine); proc != nil {
		return proc
	}
	if proc := matchProcessV2(inputLine); proc != nil {
		return proc
	}
	return nil
}

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
		i, _ = strconv.ParseUint(m[4], 10, 64)
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

var tag string

func init() {
	flag.StringVar(&tag, "tag", "", "-tag com.yourcompany.yourapp")
}

func main() {
	flag.Parse()
	reader := bufio.NewReader(os.Stdin)
	var currentPid uint64
	currentPid = 0
	filterProcess := false
	if tag != "" {
		filterProcess = true
	}

	fmt.Println("Filter by tag:" + tag)

	for {
		text, _ := reader.ReadString('\n')
		kind := detectFormat(text)

		if filterProcess {
			process := process(text)
			if process != nil && process.Tag == tag {
				currentPid = process.Pid
			}
			if kind.Pid != currentPid {
				continue
			}
		}

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
