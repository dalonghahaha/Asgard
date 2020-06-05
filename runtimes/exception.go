package runtimes

type AppException struct {
	UUID  string
	AppID int64
	Desc  string
}

type JobException struct {
	UUID  string
	JobID int64
	Desc  string
}

type TimingException struct {
	UUID     string
	TimingID int64
	Desc     string
}
