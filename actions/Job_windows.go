//+build windows

package actions

import "runtime"

func saveTo(j *Job) string {
	if j.saveTo != "" {
		return j.saveTo
	}
	saveTo := ""
	if runtime.GOOS == "windows" {
		saveTo = j.GetString("saveToWin")
	}
	if saveTo == "" {
		saveTo = j.GetString("saveTo")
	}

	j.saveTo = saveTo
	return j.saveTo
}