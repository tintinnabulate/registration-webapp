// Code generated by "strunger -type=EURYPAA_Country"; DO NOT EDIT.

package main

import "strconv"
import "strings"

const _EURYPAA_Country_name = "Albania_Andorra_Armenia_Austria_Azerbaijan_Belgium_Belarus_Bosnia_and_Herzegovina_Bulgaria_Croatia_Cyprus_Czech_Republic_Denmark_Estonia_Finland_France_Georgia_Germany_Greece_Hungary_Iceland_Ireland_Israel_Italy_Latvia_Liechtenstein_Lithuania_Luxembourg_Malta_Moldova_Republic_of_Monaco_Montenegro_Netherlands_North_Macedonia_Republic_of_Norway_Poland_Portugal_Russia_San_Marino_Serbia_Slovakia_Slovenia_Spain_Sweden_Turkey_United_Kingdom_"

var _EURYPAA_Country_index = [...]uint16{0, 8, 16, 24, 32, 43, 51, 59, 82, 91, 99, 106, 121, 129, 137, 145, 152, 160, 168, 175, 183, 191, 199, 206, 212, 219, 233, 243, 254, 260, 280, 287, 298, 310, 338, 345, 352, 361, 368, 379, 386, 395, 404, 410, 417, 424, 439}

func (i EURYPAA_Country) String() string {
	i -= 1
	if i < 0 || i >= EURYPAA_Country(len(_EURYPAA_Country_index)-1) {
		return "EURYPAA_Country(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	foo := _EURYPAA_Country_name[_EURYPAA_Country_index[i]:_EURYPAA_Country_index[i+1]]
	return strings.Replace(foo, "_", " ", -1)
}
