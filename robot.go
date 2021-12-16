package main

import (
	"errors"

	"github.com/opensourceways/community-robot-lib/config"
	"github.com/opensourceways/community-robot-lib/robot-gitee-framework"
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/sirupsen/logrus"
)

const botName = "sweepstakes"

type iClient interface {
	ListIssueComments(org, repo, number string) ([]sdk.Note, error)
	CreateIssueComment(org, repo string, number string, comment string) error
}

func newRobot(cli iClient, botName string) *robot {
	return &robot{cli, botName}
}

type robot struct {
	cli     iClient
	botName string
}

func (bot *robot) NewConfig() config.Config {
	return &configuration{}
}

func (bot *robot) getConfig(cfg config.Config) (*configuration, error) {
	if c, ok := cfg.(*configuration); ok {
		return c, nil
	}
	return nil, errors.New("can't convert to configuration")
}

func (bot *robot) RegisterEventHandler(f framework.HandlerRegitster) {
	f.RegisterNoteEventHandler(bot.handleNoteEvent)
}

func (bot *robot) handleNoteEvent(e *sdk.NoteEvent, c config.Config, log *logrus.Entry) error {
	config, err := bot.getConfig(c)
	if err != nil {
		return err
	}

	cfg := config.configFor(e.GetOrgRepo())
	if cfg == nil {
		return nil
	}

	return bot.handleSweepstakes(e, cfg)
}
