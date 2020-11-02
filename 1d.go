package cutlist

import (
	"fmt"
	"sort"
)

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
	stock := make([]Stock1D, len(cfg.Stock))
	for i, l := range cfg.Stock {
		stock[i].Waste = l
		stock[i].Length = l
	}

	cuts := make([]Length, len(_cuts))
	copy(cuts, _cuts)
	// always do longest first (in case we have existing long-stock)
	sort.Slice(cuts, func(i, j int) bool { return cuts[i] > cuts[j] })

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

	return stock, nil
}
