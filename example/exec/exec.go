package exec

import (
	"strings"
	"testing"

	"github.com/edwarnicke/exechelper"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type Exec struct {
	canFail func(string) bool
	config  Config
	opts    []*exechelper.Option
	t       *testing.T
}

func New(t *testing.T) *Exec {
	writer := logrus.StandardLogger().Writer()

	result := &Exec{
		t: t,
		opts: []*exechelper.Option{
			exechelper.WithStderr(writer),
			exechelper.WithStdout(writer),
		},
	}

	err := envconfig.Usage("exec", &result.config)
	require.NoError(t, err)

	err = envconfig.Process("exec", &result.config)
	require.NoError(t, err)

	result.registerTimeoutCommand(func(s string) bool {
		return strings.HasPrefix(s, "kubectl wait")
	})

	return result
}

func (e *Exec) Run(cmd string) {
	require.Eventually(e.t, func() bool {
		err := exechelper.Run(cmd, e.opts...)
		if !e.canFail(cmd) {
			require.NoError(e.t, err)
		}
		return err == nil
	}, e.config.Timeout, e.config.Timeout/10)

}

// Sometimes some of the running commands can fail under exec on ci, but can not fail on the user manually run.
// These command should be registered here.
func (e *Exec) registerTimeoutCommand(filter func(string) bool) {
	if filter == nil {
		panic("filter cannot be nil")
	}
	old := e.canFail
	e.canFail = func(s string) bool {
		if filter(s) {
			return true
		}

		if old != nil {
			return old(s)
		}

		return false
	}
}
