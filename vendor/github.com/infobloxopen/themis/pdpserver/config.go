package main

import (
	"flag"
	"math"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/infobloxopen/themis/pdp"
	"github.com/infobloxopen/themis/pdp/ast"
	"github.com/infobloxopen/themis/pdpserver/server"
)

const (
	policyFormatNameYAML = "yaml"
	policyFormatNameJSON = "json"
)

var policyParsers = map[string]ast.Parser{
	policyFormatNameYAML: ast.NewYAMLParser(),
	policyFormatNameJSON: ast.NewJSONParser(),
}

type config struct {
	policy              string
	policyParser        ast.Parser
	content             stringSet
	serviceEP           string
	controlEP           string
	tracingEP           string
	healthEP            string
	profilerEP          string
	storageEP           string
	mem                 server.MemLimits
	maxStreams          uint
	maxResponseSize     uint
	memStatsLogPath     string
	memStatsLogInterval time.Duration
	memProfDumpPath     string
	memProfNumGC        uint
	memProfDelay        time.Duration
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
	verbose := flag.Int("v", 1, "log verbosity (0 - error, 1 - warn (default), 2 - info, 3 - debug)")
	flag.StringVar(&conf.policy, "p", "", "policy file to start with")
	policyFmt := flag.String("pfmt", policyFormatNameYAML, "policy data format \"yaml\" or \"json\"")
	flag.Var(&conf.content, "j", "JSON content files to start with")
	flag.StringVar(&conf.serviceEP, "l", ":5555", "listen for decision requests on this address:port")
	flag.StringVar(&conf.controlEP, "c", ":5554", "listen for policies on this address:port")
	flag.StringVar(&conf.tracingEP, "t", "", "OpenZipkin tracing endpoint")
	flag.StringVar(&conf.healthEP, "health", "", "health check endpoint")
	flag.StringVar(&conf.profilerEP, "pprof", "", "performance profiler endpoint")
	flag.StringVar(&conf.storageEP, "storage", ":5552", "storage control endpoint")
	limit := flag.Uint64("mem-limit", 0, "memory limit in megabytes")
	flag.UintVar(&conf.maxStreams, "max-streams", 0, "maximum number of parallel gRPC streams (0 - use gRPC default)")
	flag.UintVar(&conf.maxResponseSize, "max-response", 10240, "maximal response size")

	flag.StringVar(&conf.memStatsLogPath, "mem-stats-log", "mem-stats.log", "file to log memory allocator statistics")
	flag.DurationVar(&conf.memStatsLogInterval, "mem-stats-interval", -1,
		"interval for memory statistics logging. Zero interval logs maximum and minimum allocated values\n"+
			"\tbetween sequential GC calls but not more than once a 100 ms. Negative interval disables logging")

	flag.StringVar(&conf.memProfDumpPath, "mem-prof-path", "/tmp/mem-prof", "directory to dump memory profiles")
	flag.UintVar(&conf.memProfNumGC, "mem-prof-gc", 0, "dump at each given GC cycle (zero - no dumping)")
	flag.DurationVar(&conf.memProfDelay, "mem-prof-delay", 0,
		"delay after request serving start for first memory profile dump\n"+
			"(zero and below - dump from programm start)")

	flag.Parse()

	initLogging(*verbose)

	p, ok := policyParsers[strings.ToLower(*policyFmt)]
	if !ok {
		log.WithField("format", *policyFmt).Fatal("unknow policy format")
	}
	conf.policyParser = p

	mem, err := server.MakeMemLimits(*limit*1024*1024, 80, 70, 30, 30)
	if err != nil {
		log.WithError(err).Fatal("wrong memory limits")
	}
	conf.mem = mem

	if conf.maxStreams > math.MaxUint32 {
		log.WithFields(log.Fields{
			"max-streams": conf.maxStreams,
			"limit":       math.MaxUint32,
		}).Fatal("too big maximum number of parallel gRPC streams")
	}

	if conf.maxResponseSize > math.MaxUint32 {
		log.WithFields(log.Fields{
			"max-response": conf.maxResponseSize,
			"limit":        math.MaxUint32,
		}).Fatal("too big maximal response size")
	}

	if conf.maxResponseSize < pdp.MinResponseSize {
		log.WithFields(log.Fields{
			"max-response": conf.maxResponseSize,
			"limit":        pdp.MinResponseSize,
		}).Fatal("too tight response size limit")
	}

	if conf.memProfNumGC > math.MaxUint32 {
		log.WithFields(log.Fields{
			"mem-prof-gc": conf.memProfNumGC,
			"limit":       math.MaxUint32,
		}).Fatal("too big number of GC cycles")
	}
}
