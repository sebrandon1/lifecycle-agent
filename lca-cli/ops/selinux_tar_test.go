package ops

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/openshift-kni/lifecycle-agent/internal/common"
)

func TestPrepareSELinuxTarCleansUpOnChconFailure(t *testing.T) {
	if _, err := os.Stat("/usr/bin/tar"); err != nil {
		t.Skip("requires /usr/bin/tar")
	}

	mockExec := NewMockExecute(gomock.NewController(t))
	o := newTestOps(t, 0)
	o.hostCommandsExecutor = mockExec
	var tarCopy string
	mockExec.EXPECT().Execute("chcon", "-t", "install_exec_t", gomock.Any()).
		DoAndReturn(func(_ string, args ...string) (string, error) {
			tarCopy = args[2]
			return "", errors.New("operation not permitted")
		})

	_, err := o.prepareSELinuxTar()
	assert.ErrorContains(t, err, "failed to set SELinux context install_exec_t")
	if assert.NotEmpty(t, tarCopy) {
		assert.NoFileExists(t, common.PathOutsideChroot(tarCopy), "the temporary tar copy must be removed")
	}
}
