package server

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/infobloxopen/themis/pdp-service && protoc -I $GOPATH/src/github.com/infobloxopen/themis/proto/ $GOPATH/src/github.com/infobloxopen/themis/proto/service.proto --go_out=plugins=grpc:$GOPATH/src/github.com/infobloxopen/themis/pdp-service && ls $GOPATH/src/github.com/infobloxopen/themis/pdp-service"

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/infobloxopen/themis/pdp-control && protoc -I $GOPATH/src/github.com/infobloxopen/themis/proto/ $GOPATH/src/github.com/infobloxopen/themis/proto/control.proto --go_out=plugins=grpc:$GOPATH/src/github.com/infobloxopen/themis/pdp-control && ls $GOPATH/src/github.com/infobloxopen/themis/pdp-control"

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/infobloxopen/themis/pdp"
	pbc "github.com/infobloxopen/themis/pdp-control"
	pbs "github.com/infobloxopen/themis/pdp-service"
	"github.com/infobloxopen/themis/pdp/ast"
	"github.com/infobloxopen/themis/pdp/jcon"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	ot "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type transport struct {
	iface net.Listener
	proto *grpc.Server
}

// Option configures how we set up PDP server.
type Option func(*options)

// WithLogger returns a Option which sets logger
func WithLogger(logger *log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// WithPolicyParser returns a Option which sets policy parser
func WithPolicyParser(parser ast.Parser) Option {
	return func(o *options) {
		o.parser = parser
	}
}

// WithServiceAt returns a Option which sets service endpoint
func WithServiceAt(addr string) Option {
	return func(o *options) {
		o.service = addr
	}
}

// WithControlAt returns a Option which sets control endpoint
func WithControlAt(addr string) Option {
	return func(o *options) {
		o.control = addr
	}
}

// WithHealthAt returns a Option which sets healthcheck endpoint
func WithHealthAt(addr string) Option {
	return func(o *options) {
		o.health = addr
	}
}

// WithProfilerAt returns a Option which sets profiler endpoint
func WithProfilerAt(addr string) Option {
	return func(o *options) {
		o.profiler = addr
	}
}

// WithStorageAt returns a Option which sets storage endpoint
func WithStorageAt(addr string) Option {
	return func(o *options) {
		o.storage = addr
	}
}

// WithTracingAt returns a Option which sets tracing endpoint
func WithTracingAt(addr string) Option {
	return func(o *options) {
		o.tracing = addr
	}
}

// WithMaxGRPCStreams returns a Option which sets maximum gRPC streams count
func WithMaxGRPCStreams(limit uint32) Option {
	return func(o *options) {
		o.streams = limit
	}
}

// WithMemLimits returns a Option which sets memory limits
func WithMemLimits(limits MemLimits) Option {
	return func(o *options) {
		o.memLimits = &limits
	}
}

// WithMaxResponseSize creates an option which limits response size in bytes. Default is 10KB. In case if a response doesn't fit the constraint, PDP puts error message to response which indicates this fact. Buffer should be at least pdp.MinResponseSize long to accept the error.
func WithMaxResponseSize(size uint32) Option {
	return func(o *options) {
		o.maxResponseSize = size
	}
}

// WithMemStatsLogging returns a Option which enables regular runtime.MemStats logging. Path points to file where stats are logged as sequence of JSON objects splitted by new line. Each JSON object contains timestamp and output of runtime.ReadMemStats taken with given interval. Zero interval logs MemStats with minimum and maximum Alloc value between NumGC changes but not more than once a 100 ms. Negative interval disables logging.
func WithMemStatsLogging(path string, interval time.Duration) Option {
	return func(o *options) {
		o.memStatsLogPath = path
		o.memStatsLogInterval = interval
	}
}

// WithMemProfDumping returns a Option to enables memory profile dumps. Path points to a directory for memory profile files. File name is "mem-NNNNNNNNN.pprof" where NNNNNNNNN is value of runtime.MemStats.NumGC after which the profile is stored.Dumps are created at each numGC after first dump (negative or zero value disables dumping). If delay is above zero first dump is taken with the delay after request service port has been opened.
func WithMemProfDumping(path string, numGC uint32, delay time.Duration) Option {
	return func(o *options) {
		o.memProfDumpPath = path
		o.memProfNumGC = numGC
		o.memProfDelay = delay
	}
}

const memStatsCheckInterval = 100 * time.Millisecond

type options struct {
	logger    *log.Logger
	parser    ast.Parser
	service   string
	control   string
	health    string
	profiler  string
	storage   string
	tracing   string
	memLimits *MemLimits
	streams   uint32

	maxResponseSize uint32

	memStatsLogPath     string
	memStatsLogInterval time.Duration
	memProfDumpPath     string
	memProfNumGC        uint32
	memProfDelay        time.Duration
}

// Server structure is PDP server object
type Server struct {
	PDPService

	sync.RWMutex

	opts options

	startOnce sync.Once
	errCh     chan error

	requests    transport
	control     transport
	health      transport
	profiler    net.Listener
	storageCtrl net.Listener

	q *queue

	//	p *pdp.PolicyStorage
	//	c *pdp.LocalContentStorage

	softMemWarn *time.Time
	backMemWarn *time.Time
	fragMemWarn *time.Time

	memProfBaseDumpDone chan uint32
}

// NewServer returns new Server instance
func NewServer(opts ...Option) *Server {
	o := options{
		logger:              log.StandardLogger(),
		service:             ":5555",
		memStatsLogInterval: -1 * time.Second,
		maxResponseSize:     10240,
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.parser == nil {
		o.parser = ast.NewYAMLParser()
	}

	var memProfBaseDumpDone chan uint32
	if o.memProfNumGC > 0 && o.memProfDelay > 0 {
		memProfBaseDumpDone = make(chan uint32)
	}

	s := &Server{
		opts:                o,
		errCh:               make(chan error, 100),
		q:                   newQueue(),
		memProfBaseDumpDone: memProfBaseDumpDone,
	}
	s.PDPService.opts = o
	s.c = pdp.NewLocalContentStorage(nil)
	return s
}

// LoadPolicies loads policies from file
func (s *PDPService) LoadPolicies(path string) error {
	if len(path) <= 0 {
		return nil
	}

	s.opts.logger.WithField("policy", path).Info("Loading policy")
	pf, err := os.Open(path)
	if err != nil {
		s.opts.logger.WithFields(log.Fields{"policy": path, "error": err}).Error("Failed load policy")
		return err
	}

	s.opts.logger.WithField("policy", path).Info("Parsing policy")
	p, err := s.opts.parser.Unmarshal(pf, nil)
	if err != nil {
		s.opts.logger.WithFields(log.Fields{"policy": path, "error": err}).Error("Failed parse policy")
		return err
	}

	s.p = p

	return nil
}

// ReadPolicies reads policies with using io.Reader instance
func (s *Server) ReadPolicies(r io.Reader) error {
	if r == nil {
		return nil
	}

	s.opts.logger.Info("Parsing policy")
	p, err := s.opts.parser.Unmarshal(r, nil)
	if err != nil {
		s.opts.logger.WithError(err).Error("Failed parse policy")
		return err
	}

	s.p = p

	return nil
}

// LoadContent loads content from files
func (s *PDPService) LoadContent(paths []string) error {
	items := []*pdp.LocalContent{}
	for _, path := range paths {
		err := func() error {
			s.opts.logger.WithField("content", path).Info("Opening content")
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			defer f.Close()

			s.opts.logger.WithField("content", path).Info("Parsing content")
			item, err := jcon.Unmarshal(f, nil)
			if err != nil {
				return err
			}

			items = append(items, item)
			return nil
		}()
		if err != nil {
			return err
		}
	}

	s.c = pdp.NewLocalContentStorage(items)

	return nil
}

// InitializeSelectors initialize all registered selectors that are enabled
func (s *Server) InitializeSelectors() {
	pdp.InitializeSelectors()
}

// ReadContent reads content with using io.Reader instances
func (s *Server) ReadContent(readers ...io.Reader) error {
	items := []*pdp.LocalContent{}
	for _, r := range readers {
		s.opts.logger.Info("Parsing content")
		item, err := jcon.Unmarshal(r, nil)
		if err != nil {
			return err
		}

		items = append(items, item)
	}

	s.c = pdp.NewLocalContentStorage(items)

	return nil
}

func (s *Server) listenRequests() error {
	s.opts.logger.WithField("address", s.opts.service).Info("Opening service port")
	ln, err := net.Listen("tcp", s.opts.service)
	if err != nil {
		return err
	}

	s.requests.iface = ln
	return nil
}

func (s *Server) listenControl() error {
	if len(s.opts.control) <= 0 {
		return nil
	}

	s.opts.logger.WithField("address", s.opts.control).Info("Opening control port")
	ln, err := net.Listen("tcp", s.opts.control)
	if err != nil {
		return err
	}

	s.control.iface = ln
	return nil
}

func (s *Server) listenHealthCheck() error {
	if len(s.opts.health) <= 0 {
		return nil
	}

	s.opts.logger.WithField("address", s.opts.health).Info("Opening health check port")
	ln, err := net.Listen("tcp", s.opts.health)
	if err != nil {
		return err
	}

	s.health.iface = ln
	return nil
}

func (s *Server) listenProfiler() error {
	if len(s.opts.profiler) <= 0 {
		return nil
	}

	s.opts.logger.WithField("address", s.opts.profiler).Info("Opening profiler port")
	ln, err := net.Listen("tcp", s.opts.profiler)
	if err != nil {
		return err
	}

	s.profiler = ln
	return nil
}

func (s *Server) listenStorage() error {
	if len(s.opts.storage) <= 0 {
		return nil
	}

	s.opts.logger.WithField("address", s.opts.storage).Info("Opening storage port")
	ln, err := net.Listen("tcp", s.opts.storage)
	if err != nil {
		return err
	}

	s.storageCtrl = ln
	return nil
}

func (s *Server) configureRequests() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	if s.opts.streams > 0 {
		opts = append(opts, grpc.MaxConcurrentStreams(s.opts.streams))
	}

	if len(s.opts.tracing) > 0 {
		tracer, err := initTracing("zipkin", s.opts.tracing)
		if err != nil {
			s.opts.logger.WithFields(log.Fields{"err": err}).Warning("Cannot initialize tracing.")
		} else {
			onlyIfParent := func(parentSpanCtx ot.SpanContext, method string, req, resp interface{}) bool {
				return parentSpanCtx != nil
			}
			intercept := otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.IncludingSpans(onlyIfParent))
			opts = append(opts, grpc.UnaryInterceptor(intercept))
		}
	}

	return opts
}

func (s *Server) serveRequests() error {
	err := s.listenRequests()
	if err != nil {
		return err
	}

	if s.opts.memProfNumGC > 0 && s.opts.memProfDelay > 0 {
		time.AfterFunc(s.opts.memProfDelay, s.memProfBaseDump)
	}

	s.opts.logger.Info("Serving decision requests")
	if err := s.requests.proto.Serve(s.requests.iface); err != nil {
		return err
	}

	return nil
}

func (s *Server) flushErrors() {
	for len(s.errCh) > 0 {
		select {
		default:
			return
		case <-s.errCh:
		}
	}
}

// Serve starts PDP server service
func (s *Server) Serve() error {
	memStatsLoggingDone := make(chan struct{})
	go s.memStatsLogging(memStatsLoggingDone)
	defer close(memStatsLoggingDone)

	memProfDumpingDone := make(chan struct{})
	go s.memProfDumping(memProfDumpingDone)
	defer close(memProfDumpingDone)

	s.flushErrors()

	if err := s.listenControl(); err != nil {
		return err
	}
	if err := s.listenHealthCheck(); err != nil {
		return err
	}
	if err := s.listenProfiler(); err != nil {
		return err
	}

	if err := s.listenStorage(); err != nil {
		return err
	}

	if s.health.iface != nil {
		healthMux := http.NewServeMux()
		healthMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			s.opts.logger.Info("Health check responding with OK")
			io.WriteString(w, "OK")
		})

		healthServer := &http.Server{Handler: healthMux}
		defer func() {
			s.health.iface.Close()
			s.health.iface = nil
		}()

		go func(l net.Listener) {
			s.errCh <- healthServer.Serve(l)
		}(s.health.iface)
	}

	if s.profiler != nil {
		profilerServer := &http.Server{}
		defer func() {
			s.profiler.Close()
			s.profiler = nil
		}()

		go func(l net.Listener) {
			s.errCh <- profilerServer.Serve(l)
		}(s.profiler)
	}

	if s.storageCtrl != nil {
		storageServer := &http.Server{Handler: &storageHandler{s}}
		defer func() {
			s.storageCtrl.Close()
			s.storageCtrl = nil
		}()

		go func(l net.Listener) {
			s.errCh <- storageServer.Serve(l)
		}(s.storageCtrl)
	}

	s.opts.logger.Info("Creating service protocol handler")
	s.requests.proto = grpc.NewServer(s.configureRequests()...)
	pbs.RegisterPDPServer(s.requests.proto, s)
	defer s.requests.proto.Stop()

	go s.memoryChecker()

	if s.p != nil {
		// We already have policy info applied; supplied from local files,
		// pointed to by CLI options.
		go s.startOnce.Do(func() {
			s.errCh <- s.serveRequests()
		})
	} else {
		if s.control.iface == nil {
			return fmt.Errorf("nothing to server - no policies provided and no control endpoint specified")
		}

		// serveRequests() will be executed by external request.
		s.opts.logger.Info("Waiting for policies to be applied.")
	}

	if s.control.iface != nil {
		s.opts.logger.Info("Creating control protocol handler")
		s.control.proto = grpc.NewServer()
		pbc.RegisterPDPControlServer(s.control.proto, s)
		defer s.control.proto.Stop()

		go func() {
			s.opts.logger.Info("Serving control requests")
			s.errCh <- s.control.proto.Serve(s.control.iface)
		}()
	}

	err := <-s.errCh
	s.flushErrors()
	return err
}

// Stop stops PDP server service
func (s *Server) Stop() error {
	if s.control.proto != nil {
		s.control.proto.Stop()
		return nil
	}

	s.RLock()
	p := s.p
	s.RUnlock()

	if p != nil && s.requests.proto != nil {
		s.requests.proto.Stop()
		return nil
	}

	return fmt.Errorf("server hasn't been started")
}
