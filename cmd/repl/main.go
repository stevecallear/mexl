package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/stevecallear/mexl"
)

func main() {
	f := flag.String("f", "", "specifies the input file path")
	flag.Parse()

	var env map[string]any
	if f != nil && *f != "" {
		var err error
		env, err = loadEnv(*f)
		if err != nil {
			log.Fatal(err)
		}
	}

	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")

		s, err := r.ReadString('\n')
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		out, err := mexl.Eval(s, env)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		fmt.Printf("= %v\n", out)
	}
}

func loadEnv(filepath string) (map[string]any, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var env map[string]any
	if err = json.NewDecoder(f).Decode(&env); err != nil {
		return nil, err
	}

	return env, nil
}
