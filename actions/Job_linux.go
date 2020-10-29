//+build !windows

package actions

func saveTo(j *Job) string {
	if j.saveTo != "" {
		return j.saveTo
	}
	saveTo := j.GetString("saveTo")
	j.saveTo = saveTo
	return j.saveTo
}
