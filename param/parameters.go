package param

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var errLog = log.New(os.Stderr, "", 0)

var versionMessage = `
dedupe version 1.02
Copyright 2020 Alan Davies <alan@galax.xyz>
Distributed under the MIT license <https://opensource.org/licenses/MIT>.
See <https://github.com/glxxyz/dupes> for documentation and help.
`

var usageMessage = `
Usage: dedupe [OPTION] DIRECTORY...
       dedupe --trash=<trash> [OPTION]... DIRECTORY...

Search DIRECTORY(ies)... for duplicate files and optionally moves them to to <trash> without user interaction.

DIRECTORY order is used for priority, highest first. Higher priority files are left untouched and lower priority files are moved. 

Mandatory parameters:

Options:
        --trash             root directory for moved duplicates, (default: files not moved)
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

func ParseParameters() (*Options, error) {

	if len(os.Args) < 2 {
		fmt.Print(usageMessage)
		os.Exit(0)
	}

	trash := flag.String("trash", "", "directory for 'trashed' files")
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
		return nil, nil
	}

	if !(*modTime || *name || *size || *hash || *contents) {
		return nil, errors.New("at least one compare- option must be true")
	}

	if *contents && !*hash {
		return nil, errors.New("when compare-contents=true then compare-hash=true must also be set")
	}

	if *hash && !*size {
		return nil, errors.New("when compare-hash=true then compare-size=true must also be set")
	}

	if len(flag.Args()) < 1 {
		return nil, errors.New("at least one directory to scan must be passed in")
	}

	minBytes, err := parseHumanReadableSize(*minSize)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse miniumum size: %w", err)
	}

	if *verbose {
		fmt.Printf("minimum file size in bytes: %v\n", minBytes)
	}

	var absoluteTrash string
	if *trash != "" {
		if absolute, err := filepath.Abs(*trash); err == nil {
			absoluteTrash = absolute
		} else {
			return nil, fmt.Errorf("failed to get an absolute path for %q: %w", *trash, err)
		}
		if _, err := os.Stat(absoluteTrash); os.IsNotExist(err) {
			return nil, fmt.Errorf("trash path does not exist: %s\n", *trash)
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
			return nil,fmt.Errorf("failed to get an absolute path for %q: %w", path, err)
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
		doMove:      *trash != "",
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
	}, nil
}
