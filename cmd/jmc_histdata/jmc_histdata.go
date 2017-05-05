package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"encoding/csv"

	"github.com/gleroi/histdata"
)

func main() {
	help := flag.Bool("h", false, "Display help message")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] files...\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "  convertes JMC histdata files to csv\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
	}

	if flag.NArg() <= 0 {
		flag.Usage()
		return
	}

	inputs := flag.Args()

	for _, input := range inputs {
		filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %s: %s\n", path, err)
					return nil
				}
				r, err := histdata.NewReader(file)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %s: %s\n", path, err)
					return nil
				}
				output, err := os.Create(filepath.Base(path) + ".csv")
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %s: %s\n", path, err)
					return nil
				}
				writer := csv.NewWriter(output)
				writer.Comma = ';'
				for entry, err := r.Read(); err == nil; entry, err = r.Read() {
					writer.Write([]string{
						entry.Attribute,
						fmt.Sprintf("%s", entry.Date),
						fmt.Sprintf("%f", entry.Value),
					})
				}
				writer.Flush()
				output.Close()
			}
			return nil
		})
	}
}
