package main

import (
	_ "net/http/pprof"
	"runtime"

	log "github.com/sirupsen/logrus"

	_ "github.com/infobloxopen/themis/pdp/selector"
	"github.com/infobloxopen/themis/pdpserver/server"
)

func main() {
	logger := log.StandardLogger()
	logger.Info("Starting PDP server")

	pdp := server.NewServer(
		server.WithLogger(logger),
		server.WithPolicyParser(conf.policyParser),
		server.WithServiceAt(conf.serviceEP),
		server.WithControlAt(conf.controlEP),
		server.WithHealthAt(conf.healthEP),
		server.WithProfilerAt(conf.profilerEP),
		server.WithStorageAt(conf.storageEP),
		server.WithTracingAt(conf.tracingEP),
		server.WithMemLimits(conf.mem),
		server.WithMaxGRPCStreams(uint32(conf.maxStreams)),
		server.WithAutoResponseSize(conf.autoResponseSize),
		server.WithMaxResponseSize(uint32(conf.maxResponseSize)),
		server.WithMemStatsLogging(
			conf.memStatsLogPath,
			conf.memStatsLogInterval,
		),
		server.WithMemProfDumping(
			conf.memProfDumpPath,
			uint32(conf.memProfNumGC),
			conf.memProfDelay,
		),
	)

	pdp.InitializeSelectors()

	err := pdp.LoadPolicies(conf.policy)
	if err != nil {
		logger.WithFields(
			log.Fields{
				"policy": conf.policy,
				"err":    err,
			},
		).Error("Failed to load policy. Continue with no policy...")
	}

	err = pdp.LoadContent(conf.content)
	if err != nil {
		logger.WithField("err", err).Error("Failed to load content. Continue with no content...")
	}

	runtime.GC()

	err = pdp.Serve()
	if err != nil {
		logger.WithError(err).Error("Failed to run server")
	}
}
