package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	m_rep = iota
)

type micro struct {
	// micro type
	Type int

	// replace the micro with Puts
	Puts string
}

func main() {
	//  open

	wd, e := os.Getwd()
	if e != nil {
		panic(e)
	}

	gomiF, e := os.Open(wd + "\\" + os.Args[2])
	if e != nil {
		panic(e)
	}
	defer gomiF.Close()

	var resFName string
	if strings.HasSuffix(os.Args[2], "_mi.go") {
		resFName = strings.Replace(os.Args[2], "_mi.go", ".go", 1)
	} else if strings.HasSuffix(os.Args[2], ".gomi") {
		resFName = strings.Replace(os.Args[2], ".gomi", ".go", 1)
	} else {
		panic("Warning: gomi accepts only files that end with `_mi.go` or `.gomi`")
	}

	resF, e := os.Create(resFName)
	if e != nil {
		panic(nil)
	}
	defer func() {
		if e := resF.Close(); e != nil {
			panic(e)
		}
	}()

	//  search for micros

	// micro map
	// micros[key] -> micro
	micros := make(map[string]micro)

	shout_out := "panic(V)"

	micro_chunk := true

	// line segment
	var ln string

	sc := bufio.NewScanner(gomiF)
	for sc.Scan() {
		ln = sc.Text()
		if micro_chunk {
			if strings.HasPrefix(ln, "#mi ") {
				t := strings.SplitN(ln, " ", 3)

				// t[1] - text to replace
				// t[2] - text to replace with
				micros[t[1]] = micro{Type: m_rep, Puts: t[2]}

				ln = "//" + ln
			} else if strings.HasPrefix(ln, "#shout ") {
				shout_out = strings.Trim(ln, "#shout ")
				ln = "//" + ln
			} else if strings.HasPrefix(ln, "import") {
				micro_chunk = false
			}
		} else {
			if strings.Contains(ln, "shout ") {
				ln = convert_shout(&ln, &shout_out)
			}
			for k, v := range micros {
				if strings.Contains(ln, k) {
					ln = convert_micro(&ln, &k, &v)
				}
			}
		}
		if _, e := resF.WriteString(ln + "\n"); e != nil {
			panic(e)
		}
	}
	if e := sc.Err(); e != nil {
		panic(e)
	}

	//  check os.Args[1], generate will only generate the file, anything else will be run in console
	if !strings.Contains(os.Args[1], "gen") {
		cmd := exec.Command("go", strings.Join(os.Args[1:], " "))
		if e := cmd.Run(); e != nil {
			fmt.Printf("\nSomething went wrong while executing go compiler: %s", e.Error())
		}
	}
}

func convert_micro(ln *string, mi_declaration *string, mi *micro) string {
	switch mi.Type {
	case m_rep:
		return strings.ReplaceAll(*ln, *mi_declaration, mi.Puts)
	default:
		return ""
	}
}

func convert_shout(ln *string, shout_out *string) string {
	lns := strings.SplitN(*ln, "shout ", 2)
	err_name := strings.SplitN(lns[1], " ", 2)[0]

	declared := " != nil"
	if strings.Contains(lns[1], ":=") {
		declared = fmt.Sprintf("; %s != nil", err_name)
	}

	return fmt.Sprintf(strings.ReplaceAll("Tif %s%s {\n\tT%s\nT}", "T", lns[0]),
		lns[1], declared, strings.Replace(*shout_out, "V", err_name, 1))
}
