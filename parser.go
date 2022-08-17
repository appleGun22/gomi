package main

import (
	"bufio"
	"fmt"
	"log"
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
	if len(os.Args) <= 2 {
		log.Fatal("go micros\nyou can find more about the usage of gomi in ")
	}

	wd, e := os.Getwd()
	if e != nil {
		log.Fatal(e)
	}

	gomiF, e := os.Open(wd + "\\" + os.Args[2])
	if e != nil {
		log.Fatal(e)
	}
	defer gomiF.Close()

	var resFName string
	if strings.HasSuffix(os.Args[2], ".gomi") {
		resFName = strings.Replace(os.Args[2], ".gomi", ".go", 1)
	} else {
		log.Fatal("Warning: gomi only accepts files that end with `.gomi`")
	}

	resF, e := os.Create(resFName)
	if e != nil {
		log.Fatal(e)
	}

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
			log.Fatal(e)
		}
	}
	if e := sc.Err(); e != nil {
		log.Fatal(e)
	}

	if e := resF.Close(); e != nil {
		log.Fatal(e)
	}

	//  check os.Args[1], generate will only generate the file, anything else will be run in console
	if !strings.Contains(os.Args[1], "gen") {
		cmd := exec.Command("go", os.Args[1:]...)
		out, e := cmd.CombinedOutput()
		if e != nil {
			fmt.Printf("\nSomething went wrong while refering to the go compiler\n%s: %s", e.Error(), string(out))
		} else {
			fmt.Printf("%s", string(out))
		}
	} else {
		fmt.Printf("`%s` successfully generated", resFName)
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
