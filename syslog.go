// +build !windows,!nacl,!plan9

package logrus_syslog

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"time"

	syslog "github.com/RackSec/srslog"
	"github.com/sirupsen/logrus"
)

const (
	SecureProto = "tcp+tls"
)

// SyslogHook to send logs via syslog.
type SyslogHook struct {
	Writer *syslog.Writer
}

// Creates a hook to be added to an instance of logger. This is called with
// `hook, err := NewSyslogHook("udp", "localhost:514", syslog.LOG_DEBUG, "")`
// `if err == nil { log.Hooks.Add(hook) }`
func NewSyslogHook(network, raddr string, priority syslog.Priority, tag string) (*SyslogHook, error) {
	w, err := syslog.Dial(network, raddr, priority, tag)
	return &SyslogHook{w}, err
}

func NewSyslogHookTls(raddr string, priority syslog.Priority, tag string, tlsConfig *tls.Config) (*SyslogHook, error) {
	dial := func(network, addr string) (net.Conn, error) {
		return tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", addr, tlsConfig)
	}
	w, err := syslog.DialWithCustomDialer("custom", raddr, priority, tag, dial)
	return &SyslogHook{w}, err
}

func (hook *SyslogHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel:
		return hook.Writer.Crit(line)
	case logrus.FatalLevel:
		return hook.Writer.Crit(line)
	case logrus.ErrorLevel:
		return hook.Writer.Err(line)
	case logrus.WarnLevel:
		return hook.Writer.Warning(line)
	case logrus.InfoLevel:
		return hook.Writer.Info(line)
	case logrus.DebugLevel:
		return hook.Writer.Debug(line)
	default:
		return nil
	}
}

func (hook *SyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
