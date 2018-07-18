package main

import (
	"flag"
	"strings"
	"time"
)

type config struct {
	policy    string
	content   string
	addresses stringSet
	timeout   time.Duration
	chunkSize int
	contentID string
	fromTag   string
	toTag     string
}

type stringSet []string

func (s *stringSet) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSet) Set(v string) error {
	*s = append(*s, v)
	return nil
}

var conf config

func init() {
	flag.StringVar(&conf.policy, "p", "", "policy file to upload")
	flag.StringVar(&conf.content, "j", "", "JSON content to upload")
	flag.Var(&conf.addresses, "s", "server(s) to upload policy to")
	flag.DurationVar(&conf.timeout, "t", 5*time.Second, "connection timeout")
	flag.IntVar(&conf.chunkSize, "c", 64*1024, "size of chunk for splitting uploads")
	flag.StringVar(&conf.contentID, "id", "", "id of content to upload")
	flag.StringVar(&conf.fromTag, "vf", "", "tag to update from (if not specified data to upload is full snapshot)")
	flag.StringVar(&conf.toTag, "vt", "", "new tag to set (if not specified data to upload is not updateable)")

	flag.Parse()
}
