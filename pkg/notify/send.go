package notify

import (
	"os/exec"

	"libs.altipla.consulting/errors"
)

type Icon string

const (
	IconError = Icon("error")
	IconInfo  = Icon("info")
)

func Send(title, message string, icon Icon) error {
	cmd := exec.Command("notify-send", title, message, "-i", string(icon))
	return errors.Trace(cmd.Run())
}
