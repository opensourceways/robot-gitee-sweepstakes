package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"k8s.io/apimachinery/pkg/util/sets"
)

func (bot *robot) handleSweepstakes(ne *sdk.NoteEvent, cfg *botConfig) error {
	if !ne.IsIssue() || !ne.IsCreatingCommentEvent() {
		return nil
	}

	n := parseCmd(ne.GetComment().GetBody())
	if n <= 0 {
		return nil
	}

	org, repo := ne.GetOrgRepo()

	comments, err := bot.cli.ListIssueComments(org, repo, ne.GetIssueNumber())
	if err != nil {
		return err
	}

	commenter := ne.GetCommenter()
	candidates := sets.NewString()
	for i := range comments {
		v := comments[i].User.Login
		if v != bot.botName && v != commenter {
			candidates.Insert(v)
		}
	}

	max := candidates.Len()
	if max == 0 {
		return nil
	}

	v := candidates.UnsortedList()
	topN := genRandomNumber(max, n)

	r := make([]string, len(topN))
	for i := range topN {
		r[i] = v[topN[i]]
	}

	return bot.cli.CreateIssueComment(
		org, repo, ne.GetIssueNumber(),
		fmt.Sprintf(cfg.Congratulation, "@"+strings.Join(r, ", @")),
	)
}

func genRandomNumber(max, n int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	v := r.Perm(max)

	if max <= n {
		return v
	}
	return v[:n]
}
