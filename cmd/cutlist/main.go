package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/mastercactapus/cutlist"
)

var (
	stockStr  = flag.String("stock", "8ft.", "Stock length.")
	kerfStr   = flag.String("kerf", "1/8in.", "Kerf (blade/cutting thickness).")
	file      = flag.String("f", "", "Read from a file instead of stdin.")
	writeFile = flag.String("o", "", "Write output to a file instead of stdout.")
	extraStr  = flag.String("extra", "", "Comma-delimited list of existing stock boards to attempt to use first.")

	multRx = regexp.MustCompile(`x[0-9]+$`)
)

func init() {
	flag.Parse()
}

func main() {
	log.SetFlags(log.Lshortfile)
	input := os.Stdin
	if *file != "" {
		fd, err := os.Open(*file)
		if err != nil {
			log.Fatalln("ERROR:", err)
		}
		defer fd.Close()
		input = fd
	}
	output := os.Stdout
	if *writeFile != "" {
		fd, err := os.Create(*writeFile)
		if err != nil {
			log.Fatalln("ERROR:", err)
		}
		defer fd.Close()
		output = fd
	}

	stock, err := cutlist.ParseLength(*stockStr)
	if err != nil {
		log.Fatalln("ERROR: invalid value for stock size:", err)
	}
	kerf, err := cutlist.ParseLength(*kerfStr)
	if err != nil {
		log.Fatalln("ERROR: invalid value for kerf size:", err)
	}

	r := bufio.NewScanner(input)
	var cuts []cutlist.Length
	for r.Scan() {
		line := strings.TrimSpace(r.Text())
		if line == "" {
			continue
		}

		times := 1
		mult := multRx.FindString(line)
		if mult != "" {
			line = line[:len(line)-len(mult)]
			mult = mult[1:]

			times, err = strconv.Atoi(mult)
			if err != nil {
				log.Fatalf("ERROR: invalid multiplier value '%s': %v", mult, err)
			}
		}

		cut, err := cutlist.ParseLength(line)
		if err != nil {
			log.Fatalln("ERROR: invalid dimension:", err)
		}
		for i := 0; i < times; i++ {
			cuts = append(cuts, cut)
		}
	}
	if r.Err() != nil {
		log.Fatalln("ERROR: read input:", err)
	}
	rand.Seed(2)
	rand.Shuffle(len(cuts), func(i, j int) { cuts[i], cuts[j] = cuts[j], cuts[i] })

	cfg := cutlist.Config1D{
		DefaultStock: stock,
		Kerf:         kerf,
	}
	if *extraStr != "" {
		extras := strings.Split(*extraStr, ",")
		for _, e := range extras {
			l, err := cutlist.ParseLength(e)
			if err != nil {
				log.Fatalln("ERROR: invalid dimension for extra board:", err)
			}
			cfg.Stock = append(cfg.Stock, l)
		}
	}

	toCut, err := cfg.Cutlist(cuts)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}

	var waste cutlist.Length
	for i, board := range toCut {
		waste += board.Waste
		fmt.Fprintf(output, "Board #%d (Length: %s, Waste: %s)\n", i+1, board.Length.String(), board.Waste.String())
		for _, cut := range board.Cuts {
			fmt.Fprintln(output, "  "+cut.String())
		}
		fmt.Fprintln(output)
	}
	fmt.Fprintf(output, "Total Boards (%s ea): %d\n", stock.String(), len(toCut))
	fmt.Fprintln(output, "Total Waste: "+waste.String())
}
