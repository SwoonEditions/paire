package paire

type Release struct {
	Commit string
	Tag string
	Name string
	Packages []ReleasePackage
	Pushed bool
	PreRelease bool
	Id int
}
