// Code generated by "strunger -type=HelpOutreach"; DO NOT EDIT.

package main

import "strconv"
import "strings"

const _HelpOutreach_name = "Yes_Help_OutreachNo_Thanks"

var _HelpOutreach_index = [...]uint8{0, 17, 26}

func (i HelpOutreach) String() string {
	i -= 1
	if i < 0 || i >= HelpOutreach(len(_HelpOutreach_index)-1) {
		return "HelpOutreach(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	foo := _HelpOutreach_name[_HelpOutreach_index[i]:_HelpOutreach_index[i+1]]
	return strings.Replace(foo, "_", " ", -1)
}
