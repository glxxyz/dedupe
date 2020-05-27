module github.com/glxxyz/dedupe

go 1.14

replace (
    github.com/glxxyz/dedupe/param v0.0.0 => ./src/param
    github.com/glxxyz/dedupe/repo v0.0.0 => ./src/repo
)

require (
	github.com/glxxyz/dedupe/param v0.0.0
	github.com/glxxyz/dedupe/repo v0.0.0
)
