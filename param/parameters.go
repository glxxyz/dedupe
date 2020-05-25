package param

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

var errLog = log.New(os.Stderr, "", 0)

var versionMessage = `
dedupe version 1.01
Copyright 2020 Alan Davies <alan@galax.xyz>
Distributed under the MIT license <https://opensource.org/licenses/MIT>.
See <https://github.com/glxxyz/dupes> for documentation and help.
`

var usageMessage = `
Usage: dedupe [OPTION] PATH...
       dedupe --move-files --trash=<trash> [OPTION] PATH...

Find PATH... for duplicate files and moves them to to <trash> without user interaction. If multiple PATHs are specified their order is used for priority, highest first.

Mandatory parameters:

Options:
        --move-files        actually move files, otherwise just output (default: false)
        --trash             root directory for 'trashed' files, mandatory if moving files
        --compare-time      compare file modification time (default: false)
        --compare-name      compare file name (default: false)
        --compare-size      compare file size (default: true)
        --compare-hash      compare file hash (default: true)
        --compare-contents  compare whole file contents (default: false)
        --min-size          minimum file size, bytes or human readable e.g. 4M, 5G (default 1)
		--follow-symlinks   follow symbolic links, false ignores them (default false)
        --verbose           emit verbose information (default: false)
        --version           output version and license information and exit

Advanced options:
        --scan-buffer       size of the scan buffer (default: 100)
        --scanners          number of scanner coroutines (default: 2)
        --match-buffer      size of the match buffer (default: 100)
        --matchers          number of matcher coroutines (default: 4)
        --move-buffer       size of the move buffer (default: 100)
        --movers            number of mover coroutines (default: 2)
        --max-cpus          maximum CPUs to use (default: system setting)

See <https://github.com/glxxyz/dedupe> for documentation and help.
`

func ParseParameters() *Options {

	if len(os.Args) < 2 {
		fmt.Print(usageMessage)
		os.Exit(0)
	}

	trash := flag.String("trash", "", "directory for 'trashed' files")
	doMove := flag.Bool("move-files", false, "if set will move files, otherwise just output")
	modTime := flag.Bool("compare-time", false, "compare file modification time")
	name := flag.Bool("compare-name", false, "compare file name")
	size := flag.Bool("compare-size", true, "compare file size")
	hash := flag.Bool("compare-hash", true, "compare file hash")
	contents := flag.Bool("compare-contents", false, "compare file contents")
	minSize := flag.String("min-size", "0", "minimum file size, bytes or human readable e.g. 4M, 5G")
	symLinks := flag.Bool("follow-symlinks", false, "follow symbolic links, false ignores them")
	verbose := flag.Bool("verbose", false, "emit verbose information")
	version := flag.Bool("version", false, "output version and license information and exit")
	scanBuffer := flag.Int("scan-buffer", 100, "size of the scan buffer")
	scanners := flag.Int("scanners", 2, " number of scanner coroutines")
	matchBuffer := flag.Int("match-buffer", 100, "size of the match buffer")
	matchers := flag.Int("matchers", 4, " number of matcher coroutines")
	moveBuffer := flag.Int("move-buffer", 100, "size of the move buffer")
	movers := flag.Int("movers", 2, "number of mover coroutines")
	maxProcs := flag.Int("max-cpus", 0, "maximum CPUs to use ")

	flag.Parse()

	if *version {
		fmt.Print(versionMessage)
		os.Exit(0)
	}

	if *trash == "" && *doMove {
		errLog.Print("when moving files a trash path must be specified with --trash=")
		os.Exit(1)
	}

	if *trash != "" && !*doMove {
		errLog.Print("When specifying a trash directory you must also set --move-files")
		os.Exit(1)
	}

	if !(*modTime || *name || *size || *hash || *contents) {
		errLog.Print("at least one compare- option must be true")
		os.Exit(1)
	}

	if *contents && !*hash {
		errLog.Print("when compare-contents=true then compare-hash=true must also be set")
		os.Exit(1)
	}

	if *hash && !*size {
		errLog.Print("when compare-hash=true then compare-size=true must also be set")
		os.Exit(1)
	}

	if len(flag.Args()) < 1 {
		errLog.Print("at least one path to scan must be passed in")
		os.Exit(1)
	}

	minBytes := parseMinSize(*minSize)
	if *verbose {
		fmt.Printf("minimum file size in bytes: %v\n", minBytes)
	}

	var absoluteTrash string
	if *trash != "" {
		if absolute, err := filepath.Abs(*trash); err == nil {
			absoluteTrash = absolute
		} else {
			errLog.Printf("failed to get an absolute path for %q: %v\n", *trash, err)
			panic(err)
		}
		if _, err := os.Stat(absoluteTrash); os.IsNotExist(err) {
			errLog.Printf("trash path does not exist: %s\n", *trash)
			panic(err)
		}
	}

	absolutePaths := make([]string, len(flag.Args()))
	for i, path := range flag.Args() {
		if absolute, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absolute); os.IsNotExist(err) {
				errLog.Printf("path does not exist: %s\n", *trash)
			}
			absolutePaths[i] = absolute
		} else {
			errLog.Printf("failed to get an absolute path for %q: %v\n", path, err)
			panic(err)
		}
	}

	if *verbose {
		fmt.Printf("System default is %d CPUs\n", runtime.NumCPU())
	}

	if *maxProcs > 0 {
		if *verbose {
			fmt.Printf("Setting GOMAXPROCS to %d CPUs\n", *maxProcs)
		}
		previous := runtime.GOMAXPROCS(*maxProcs)
		if *verbose {
			fmt.Printf("Previous GOMAXPROCS was %d CPUs\n", previous)
		}
	}

	return &Options{
		trash:       absoluteTrash,
		doMove:      *doMove,
		modTime:     *modTime,
		name:        *name,
		size:        *size,
		hash:        *hash,
		contents:    *contents,
		minBytes:    minBytes,
		symLinks:    *symLinks,
		verbose:     *verbose,
		scanBuffer:  *scanBuffer,
		scanners:    *scanners,
		matchBuffer: *matchBuffer,
		matchers:    *matchers,
		moveBuffer:  *moveBuffer,
		movers:      *movers,
		paths:       absolutePaths,
	}
}

func parseMinSize(minSize string) int64 {
	scaleFactor := 1.0
	switch minSize[len(minSize)-1] {
	case 'K': // kibibyte
		scaleFactor = 1024
		minSize = minSize[:len(minSize)-1]
	case 'M': // mebibyte
		scaleFactor *= math.Pow(1024, 2)
		minSize = minSize[:len(minSize)-1]
	case 'G': // gibibyte
		scaleFactor *= math.Pow(1024, 3)
		minSize = minSize[:len(minSize)-1]
	case 'T': // tebibyte
		scaleFactor *= math.Pow(1024, 4)
		minSize = minSize[:len(minSize)-1]
	case 'P': // pebibyte
		scaleFactor *= math.Pow(1024, 5)
		minSize = minSize[:len(minSize)-1]
	case 'E': // exbibyte
		scaleFactor *= math.Pow(1024, 6)
		minSize = minSize[:len(minSize)-1]
	case 'Z': // zebibyte
		scaleFactor *= math.Pow(1024, 7)
		minSize = minSize[:len(minSize)-1]
	case 'Y': // yobibyte
		scaleFactor *= math.Pow(1024, 8)
		minSize = minSize[:len(minSize)-1]
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		break
	default:
		errLog.Print("min-size can be an integer number of bytes or a floating point value with a human readable size e.g. 1.5K, 3M, 4G, 4.2T")
		os.Exit(1)
	}
	floatVal, err := strconv.ParseFloat(minSize, 64)
	if err != nil {
		errLog.Printf("failed to parse min size: %v\n", err)
		panic(err)
	}
	return int64(math.Round(floatVal * scaleFactor))
}
