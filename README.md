# dedupe

![Go](https://github.com/glxxyz/dedupe/workflows/Go/badge.svg)

A tool for efficiently scanning over large numbers of files, to find and remove duplicates. Mac, Linux, and Windows binaries are [here](https://github.com/glxxyz/dedupe/tree/master/bin).

I wrote this for a bit of practice using Go, and to help me finally get around to cleaning up that mess of photo backups.

## Use Case
Here's an example use case for this tool:
I have spent a while organising most of my photos, by year etc. These are stored in `/photos`.
I have a copy in `/backup` from when I was part of the way through the re-organization.
I have a directory `/unsorted` with  photos that were downloaded from iCloud (in `/unsorted/icloud`), Amazon Photos (in `/unsorted/amazon`), Google Photos and others.
Some photos are only backed up in one place, some are duplicated, it's a big mess.

What I'd like to do is remove all duplicates from `/backup` and `/unsorted` so that I can organise the ones that aren't in `/photos` yet. I can do that with just one command:

    dedupe --trash=/trash /photos /backup /unsorted

I created a `/trash` directory to move the duplicates that are 'trashed'.
They won't be deleted, and I can undo the move of some/all of them by manually copying them back to where they were originally.
Duplicates will be moved to `/trash/backup`, `/trash/unsorted/amazon`, etc. based on where they started.

If I'm nervous about what could get moved I could leave off the `--trash` option, and it would just find duplicates and output what would be moved.

    dedupe /photos /backup /unsorted

It decides what to move based on the order of the directories- higher priority first.
If I have some directories within /backup that are better organised even than my main `/photos` directory I can make those higher priority:

    dedupe --trash=/trash /backup/2015move /backup/2020aprilfun /photos /backup /unsorted

## Downloads

If you don't want to build it yourself you can pull down a precompiled binary from here:

* [linux-amd64 latest build](https://github.com/glxxyz/dedupe/raw/master/bin/linux-amd64/dedupe)
* [macos-amd64 latest build](https://github.com/glxxyz/dedupe/raw/master/bin/macos-amd64/dedupe)
* [windows-amd64 latest build](https://github.com/glxxyz/dedupe/raw/master/bin/macos-amd64/dedupe.exe)

## Options

Here's the help text which the tool outputs, if it's unclear let me know or file an issue:
```
Usage: dedupe [OPTION] DIRECTORY...
       dedupe --trash=<trash> [OPTION]... DIRECTORY...

Search DIRECTORY(ies)... for duplicate files and optionally moves them to <trash> without user interaction.

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
```

## About

```
Copyright 2020 Alan Davies <alan@galax.xyz>
Distributed under the MIT license <https://opensource.org/licenses/MIT>.
See <https://github.com/glxxyz/dedupe> for documentation and help.
```
