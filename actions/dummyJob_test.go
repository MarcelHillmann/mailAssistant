package actions

import "testing"

func TestDummyJob(t *testing.T){
	var wg int32
	job := Job{}

	newDummy(job,&wg)
}
