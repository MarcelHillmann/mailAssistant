package actions

func newDummy(job Job, _ *int32, result func(int)) {
	job.GetLogger().Debug(job)
	result(1)
}
