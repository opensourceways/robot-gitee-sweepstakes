package main

import (
	"errors"

	sdk "gitee.com/openeuler/go-gitee/gitee"
	libconfig "github.com/opensourceways/community-robot-lib/config"
	"github.com/opensourceways/community-robot-lib/giteeclient"
	libplugin "github.com/opensourceways/community-robot-lib/giteeplugin"
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

func (bot *robot) NewPluginConfig() libconfig.PluginConfig {
	return &configuration{}
}

func (bot *robot) getConfig(cfg libconfig.PluginConfig) (*configuration, error) {
	if c, ok := cfg.(*configuration); ok {
		return c, nil
	}
	return nil, errors.New("can't convert to configuration")
}

func (bot *robot) RegisterEventHandler(p libplugin.HandlerRegitster) {
	p.RegisterNoteEventHandler(bot.handleNoteEvent)
}

func (bot *robot) handleNoteEvent(e *sdk.NoteEvent, pc libconfig.PluginConfig, log *logrus.Entry) error {

	config, err := bot.getConfig(pc)
	if err != nil {
		return err
	}

	ne := giteeclient.NewIssueNoteEvent(e)

	cfg := config.configFor(ne.GetOrgRep())
	if cfg == nil {
		return nil
	}

	return bot.handleSweepstakes(ne, cfg)
}
