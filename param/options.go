package param

type Options struct {
	trash       string
	doMove      bool
	modTime     bool
	name        bool
	size        bool
	hash        bool
	contents    bool
	minBytes    int64
	symLinks    bool
	verbose     bool
	scanBuffer  int
	scanners    int
	matchBuffer int
	matchers    int
	moveBuffer  int
	movers      int
	paths       []string
}

// dumb accessors that allow for encapsulation

func (options *Options) Trash() string {
	return options.trash
}

func (options *Options) DoMove() bool {
	return options.doMove
}

func (options *Options) ModTime() bool {
	return options.modTime
}

func (options *Options) Name() bool {
	return options.name
}

func (options *Options) Size() bool {
	return options.size
}

func (options *Options) Hash() bool {
	return options.hash
}

func (options *Options) Contents() bool {
	return options.contents
}

func (options *Options) MinBytes() int64 {
	return options.minBytes
}

func (options *Options) SymLinks() bool {
	return options.symLinks
}

func (options *Options) Verbose() bool {
	return options.verbose
}

func (options *Options) ScanBuffer() int {
	return options.scanBuffer
}

func (options *Options) Scanners() int {
	return options.scanners
}

func (options *Options) MatchBuffer() int {
	return options.matchBuffer
}

func (options *Options) Matchers() int {
	return options.matchers
}

func (options *Options) MoveBuffer() int {
	return options.moveBuffer
}

func (options *Options) Movers() int {
	return options.movers
}

func (options *Options) Paths() []string {
	return options.paths
}