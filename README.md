# Cutlist Generator

This program can be used to generate single-dimension cutlists.

## Instal

`go get github.com/mastercactapus/cutlist/cmd/cutlist`

## Usage

```
Usage of cutlist:
  -extra string
        Comma-delimited list of existing stock boards to attempt to use first.
  -f string
        Read from a file instead of stdin.
  -kerf string
        Kerf (blade/cutting thickness). (default "1/8in.")
  -o string
        Write output to a file instead of stdout.
  -stock string
        Stock length. (default "8ft.")
```

### Example

```bash
# -stock indicates the length of new boards
# -extra indicates existing boards to be used up first if possible
cutlist -stock 8ft -extra 7ft,7ft,8ft,8ft -f cuts.txt
```

## Input

Input is one length per line, optionaly with a `xN` multiplier where `N` is the quantity needed of that size.

```
7.7cm x8
9cm x4
1m x4
1.8/2m x4
```

Multiple dimensions on the same line will be added:

- `1ft 3"` will be interpreted as `15 inches`/`38.1cm`.

You can use fractions with any unit:

- `1/4in`
- `1/2m`
- `1.8/2m`

## Output

It will attempt to minimize waste and use `extra` boards first.

Currently output cannot be rounded and units are fixed to `cm`.

```
Board #1 (Length: 243cm, Waste: 0.36cm)
  100cm
  90cm
  9cm
  9cm
  9cm
  7.7cm
  7.7cm
  7.7cm

Board #2 (Length: 243cm, Waste: 44.3475cm)
  100cm
  90cm
  7.7cm

Board #3 (Length: 213cm, Waste: 0.295cm)
  90cm
  90cm
  7.7cm
  7.7cm
  7.7cm
  7.7cm

Board #4 (Length: 213cm, Waste: 3.0475cm)
  100cm
  100cm
  9cm

Total Boards (243.84cm ea): 4
Total Waste: 48.05cm
```
