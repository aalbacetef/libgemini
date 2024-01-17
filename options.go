package libgemini

import (
	_ "embed"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aalbacetef/libgemini/tofu"
)

type Options struct {
	Timeout         time.Duration
	FollowRedirects bool
	StorePath       string
	RCFilepath      string
	DumpHeaders     string
	Trace           string
	Insecure        bool
}

const (
	DefaultTimeout         = time.Second * 30
	DefaultFollowRedirects = false
	DefaultInsecure        = false
	InMemoryStoreVal       = ":memory:"
)

func defaultOpts() Options {
	return Options{
		Timeout:         DefaultTimeout,
		FollowRedirects: DefaultFollowRedirects,
		Insecure:        DefaultInsecure,
	}
}

const (
	EnvRC              = "LIBGEMINI_RC"
	EnvFollowRedirects = "LIBGEMINI_FOLLOW_REDIRECTS"
	EnvStorePath       = "LIBGEMINI_STORE_PATH"
	EnvDumpHeaders     = "LIBGEMINI_DUMP_HEADERS"
	EnvTrace           = "LIBGEMINI_TRACE"
	EnvInsecure        = "LIBGEMINI_INSECURE"
	KeyRC              = "RC"
	KeyFollowRedirects = "FollowRedirects"
	KeyStorePath       = "StorePath"
	KeyDumpHeaders     = "DumpHeaders"
	KeyTrace           = "Trace"
	KeyInsecure        = "Insecure"
)

type strOrBool struct {
	s string
	b bool
}

// envOpts will lookup the environment variables one by one,
// using only the ones that are set. It will check for
// truthy values (on, 1, true), defaulting to false for anything else.
func envOpts() map[string]strOrBool {
	opts := make(map[string]strOrBool)

	if v, set := os.LookupEnv(EnvRC); set {
		opts[KeyRC] = strOrBool{s: v}
	}

	if v, set := os.LookupEnv(EnvFollowRedirects); set {
		opts[KeyFollowRedirects] = strOrBool{b: toBool(v)}
	}

	if v, set := os.LookupEnv(EnvStorePath); set {
		opts[KeyStorePath] = strOrBool{s: v}
	}

	if v, set := os.LookupEnv(EnvDumpHeaders); set {
		opts[KeyDumpHeaders] = strOrBool{s: v}
	}

	if v, set := os.LookupEnv(EnvTrace); set {
		opts[KeyTrace] = strOrBool{s: v}
	}

	if v, set := os.LookupEnv(EnvInsecure); set {
		opts[KeyInsecure] = strOrBool{b: toBool(v)}
	}

	return opts
}

const (
	ConfigFollowRedirects = "follow"
	ConfigStore           = "store"
	ConfigDumpHeaders     = "dump-headers"
	ConfigTrace           = "trace"
	ConfigInsecure        = "insecure"
)

func configOpts(contents string) map[string]strOrBool {
	opts := make(map[string]strOrBool)

	if contents == "" {
		return opts
	}

	for _, line := range strings.Split(contents, "\n") {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}

		if !strings.HasPrefix(l, "--") {
			continue
		}

		parts := strings.Split(l, " ")
		optName := parts[0][2:]

		switch optName {
		case ConfigFollowRedirects, ConfigInsecure:
			opts[optName] = strOrBool{b: true}
		case ConfigStore, ConfigDumpHeaders, ConfigTrace:
			val := strings.TrimSpace(strings.Join(parts[1:], " "))
			if val == "" {
				continue
			}
			// NOTE: maybe add a strings.Split(val, "#")[0] to allow comments on the same line?
			opts[optName] = strOrBool{s: val}
		}
	}

	return opts
}

func mergeOpts(base Options, applyOpts ...map[string]strOrBool) Options {
	if len(applyOpts) == 0 {
		return base
	}

	for _, opts := range applyOpts {
		for key, val := range opts {
			switch key {
			case KeyRC:
				base.RCFilepath = val.s
			case KeyFollowRedirects:
				base.FollowRedirects = val.b
			case KeyStorePath:
				base.StorePath = val.s
			case KeyDumpHeaders:
				base.DumpHeaders = val.s
			case KeyTrace:
				base.Trace = val.s
			case KeyInsecure:
				base.Insecure = val.b
			}
		}
	}

	return base
}

func toBool(s string) bool {
	switch s {
	case "on", "1", "true":
		return true
	default:
		return false
	}
}

type OptsFn func(*Options)

func WithInMemoryStore() OptsFn {
	return func(opts *Options) {
		opts.StorePath = InMemoryStoreVal
	}
}

func WithStore(val string) OptsFn {
	return func(opts *Options) {
		opts.StorePath = val
	}
}

func WithInsecure() OptsFn {
	return func(opts *Options) {
		opts.Insecure = true
	}
}

//go:embed data/geminirc
var stubRCFile []byte

// @TODO: log errors
func resolveConfigFile() string {
	if val, set := os.LookupEnv(EnvRC); set {
		writeIfNotExists(val, stubRCFile)
		return val
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	libgeminiDir := filepath.Join(homeDir, ".config", "libgemini")
	if err := os.MkdirAll(libgeminiDir, fs.FileMode(0o700)); err != nil {
		return ""
	}

	cfgFile := filepath.Join(libgeminiDir, "geminirc")
	writeIfNotExists(cfgFile, stubRCFile)

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return ""
	}

	return string(data)
}

func writeIfNotExists(fpath string, file []byte) {
	if _, err := os.Stat(fpath); errors.Is(err, fs.ErrNotExist) {
		os.WriteFile(fpath, file, fs.FileMode(0o644))
	}
}

// @TODO: log errors
func resolveStore(storeOpt string) tofu.Store {
	if storeOpt == InMemoryStoreVal {
		return tofu.NewInMemoryStore()
	}

	if storeOpt == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return tofu.NewInMemoryStore()
		}

		libgeminiDir := filepath.Join(homeDir, ".config", "libgemini")
		if err := os.MkdirAll(libgeminiDir, fs.FileMode(0o700)); err != nil {
			return tofu.NewInMemoryStore()
		}

		store, err := tofu.NewFileStore(filepath.Join(libgeminiDir, "known_hosts"))
		if err != nil {
			return tofu.NewInMemoryStore()
		}

		return store
	}

	expanded := os.ExpandEnv(storeOpt)
	if expanded == "" {
		return tofu.NewInMemoryStore()
	}

	store, err := tofu.NewFileStore(expanded)
	if err != nil {
		return tofu.NewInMemoryStore()
	}

	return store
}
