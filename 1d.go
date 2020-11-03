package cutlist

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

func init() {
	var b [8]byte
	_, err := crand.Read(b[:])
	if err != nil {
		// doesn't need to be secure, just random
		rand.Seed(time.Now().UnixNano())
		return
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

type Config1D struct {
	DefaultStock Length
	Kerf         Length

	Stock []Length
}

type Stock1D struct {
	Length Length
	Waste  Length
	Cuts   []Length
}

func (b *Stock1D) cut(kerf, cut Length) bool {
	if cut > b.Waste {
		return false
	}
	b.Cuts = append(b.Cuts, cut)
	b.Waste -= kerf + cut
	if b.Waste < 0 {
		b.Waste = 0
	}

	return true
}

func (cfg Config1D) Cutlist(_cuts []Length) ([]Stock1D, error) {
	if len(_cuts) == 0 {
		return nil, nil
	}
	baseStock := make([]Stock1D, len(cfg.Stock))
	for i, l := range cfg.Stock {
		baseStock[i].Waste = l
		baseStock[i].Length = l
	}

	baseCuts := make([]Length, len(_cuts))
	copy(baseCuts, _cuts)

	var totalWaste Length
	var logWaste float64
	var result []Stock1D
	for i := 0; i < 10000; i++ {
		stock := make([]Stock1D, len(baseStock))
		copy(stock, baseStock)
		cuts := make([]Length, len(baseCuts))
		copy(cuts, baseCuts)

		rand.Shuffle(len(stock), func(i, j int) { stock[i], stock[j] = stock[j], stock[i] })
		rand.Shuffle(len(cuts), func(i, j int) { cuts[i], cuts[j] = cuts[j], cuts[i] })

	cutLoop:
		for _, cut := range cuts {
			for i := range stock {
				if stock[i].cut(cfg.Kerf, cut) {
					continue cutLoop
				}
			}
			if cut > cfg.DefaultStock {
				return nil, fmt.Errorf("can not cut %s from new stock that is only %s", cut.String(), cfg.DefaultStock.String())
			}
			// add a board
			b := Stock1D{Waste: cfg.DefaultStock, Length: cfg.DefaultStock}
			b.cut(cfg.Kerf, cut)
			stock = append(stock, b)
		}

		var w Length
		var lw float64
		for _, s := range stock {
			w += s.Waste
			lw += math.Log(float64(s.Waste))
		}
		if result == nil || (w < totalWaste || (w == totalWaste && lw < logWaste)) {
			totalWaste = w
			logWaste = lw
			result = stock
		}
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Length > result[j].Length })
	for _, s := range result {
		sort.Slice(s.Cuts, func(i, j int) bool { return s.Cuts[i] > s.Cuts[j] })
	}

	return result, nil
}
