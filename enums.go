package main

// EURYPAA_Country : all the countries a EURYPAA can happen in
//go:generate strunger -type=EURYPAA_Country
type EURYPAA_Country int

const (
	Albania_ EURYPAA_Country = iota + 1
	Andorra_
	Armenia_
	Austria_
	Azerbaijan_
	Belgium_
	Belarus_
	Bosnia_and_Herzegovina_
	Bulgaria_
	Croatia_
	Cyprus_
	Czech_Republic_
	Denmark_
	Estonia_
	Finland_
	France_
	Georgia_
	Germany_
	Greece_
	Hungary_
	Iceland_
	Ireland_
	Israel_
	Italy_
	Latvia_
	Liechtenstein_
	Lithuania_
	Luxembourg_
	Malta_
	Moldova_Republic_of_
	Monaco_
	Montenegro_
	Netherlands_
	North_Macedonia_Republic_of_
	Norway_
	Poland_
	Portugal_
	Russia_
	San_Marino_
	Serbia_
	Slovakia_
	Slovenia_
	Spain_
	Sweden_
	Turkey_
	United_Kingdom_
	Online_
)

// Fellowship : which fellowship someone belongs to
//go:generate strunger -type=Fellowship
type Fellowship int

// AA : start of fellowship enumeration
const (
	AA Fellowship = iota + 1
	AlAnon
	Guest
)

// Fellowships : all the fellowships
var Fellowships = []Fellowship{
	AA,
	AlAnon,
	Guest,
}
