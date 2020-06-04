package runtimes

type Archive struct {
	UUID      string
	Pid       int32
	BeginTime int64
	EndTime   int64
	Status    int32
	Signal    string
}

type AppArchive struct {
	UUID    string
	App     *App
	Archive *Archive
}

type JobArchive struct {
	UUID    string
	Job     *Job
	Archive *Archive
}

type TimingArchive struct {
	UUID    string
	Timing  *Timing
	Archive *Archive
}

func buildArchive(cmd *Command) *Archive {
	return &Archive{
		UUID:      cmd.UUID,
		Pid:       int32(cmd.Pid),
		BeginTime: cmd.Begin.Unix(),
		EndTime:   cmd.End.Unix(),
		Status:    int32(cmd.Status),
		Signal:    cmd.Signal,
	}
}
