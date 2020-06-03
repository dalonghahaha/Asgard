package services

import (
	"fmt"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/components/mail"

	"Asgard/constants"
	"Asgard/models"
	"Asgard/rpc"
)

type NoticeService struct {
}

func NewNoticeService() *NoticeService {
	return &NoticeService{}
}

func (s *NoticeService) SendMail(address, subject, body string) error {
	message := mail.BuildMessage(constants.MAIL_USER, []string{address}, subject, body)
	dialer := mail.Get(constants.MAIL_NAME)
	return mail.Send(dialer, message)
}

func (s *NoticeService) AppUnsuccessNotify(app *models.App, agent *models.Agent, archive *rpc.AppArchive) {
	subject := "app finished unsuccess notice"
	body := fmt.Sprintf("<p>agent info:%s(%s:%s)</p>", agent.Alias, agent.IP, agent.Port)
	body += fmt.Sprintf("<p>app name:%s</p>", app.Name)
	body += fmt.Sprintf("<p>app dir:%s</p>", app.Dir)
	body += fmt.Sprintf("<p>app cmd:%s %s</p>", app.Program, app.Args)
	body += fmt.Sprintf("<p>app std:%s</p>", app.StdErr)
	body += fmt.Sprintf("<p>app err:%s</p>", app.StdOut)
	body += fmt.Sprintf("<p>app uuid:%s</p>", archive.GetArchive().GetUuid())
	body += fmt.Sprintf("<p>app exit status:%d</p>", archive.GetArchive().GetStatus())
	body += fmt.Sprintf("<p>app exit singal:%s</p>", archive.GetArchive().GetSignal())
	err := s.SendMail(constants.MASTER_RECEIVER, subject, body)
	if err != nil {
		logger.Errorf("AppUnsuccessNotify Failed:%+v", err)
	}
}

func (s *NoticeService) JobUnsuccessNotify(job *models.Job, agent *models.Agent, archive *rpc.JobArchive) {
	subject := "job finished unsuccess notice"
	body := fmt.Sprintf("<p>agent info:%s(%s:%s)</p>", agent.Alias, agent.IP, agent.Port)
	body += fmt.Sprintf("<p>job name:%s</p>", job.Name)
	body += fmt.Sprintf("<p>job dir:%s</p>", job.Dir)
	body += fmt.Sprintf("<p>job cmd:%s %s</p>", job.Program, job.Args)
	body += fmt.Sprintf("<p>job std:%s</p>", job.StdErr)
	body += fmt.Sprintf("<p>job err:%s</p>", job.StdOut)
	body += fmt.Sprintf("<p>job uuid:%s</p>", archive.GetArchive().GetUuid())
	body += fmt.Sprintf("<p>job exit status:%d</p>", archive.GetArchive().GetStatus())
	body += fmt.Sprintf("<p>job exit singal:%s</p>", archive.GetArchive().GetSignal())
	err := s.SendMail(constants.MASTER_RECEIVER, subject, body)
	if err != nil {
		logger.Errorf("JobUnsuccessNotify Failed:%+v", err)
	}
}

func (s *NoticeService) TimingUnsuccessNotify(timing *models.Timing, agent *models.Agent, archive *rpc.TimingArchive) {
	subject := "timing finished unsuccess notice"
	body := fmt.Sprintf("agent info:%s(%s:%s)</p>", agent.Alias, agent.IP, agent.Port)
	body += fmt.Sprintf("<p>timing name:%s</p>", timing.Name)
	body += fmt.Sprintf("<p>timing dir:%s</p>", timing.Dir)
	body += fmt.Sprintf("<p>timing cmd:%s %s</p>", timing.Program, timing.Args)
	body += fmt.Sprintf("<p>timing std:%s</p>", timing.StdErr)
	body += fmt.Sprintf("<p>timing err:%s</p>", timing.StdOut)
	body += fmt.Sprintf("<p>timing uuid:%s</p>", archive.GetArchive().GetUuid())
	body += fmt.Sprintf("<p>timing exit status:%d</p>", archive.GetArchive().GetStatus())
	body += fmt.Sprintf("<p>timing exit singal:%s</p>", archive.GetArchive().GetSignal())
	err := s.SendMail(constants.MASTER_RECEIVER, subject, body)
	if err != nil {
		logger.Errorf("TimingUnsuccessNotify Failed:%+v", err)
	}
}
