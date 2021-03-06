package update

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/config"
)

func Check() error {
	// Jenkins y el entorno de desarrollo de actools no deben comprobar la versión
	if config.Development() || config.Jenkins() {
		return nil
	}

	lastUpdateFilename := filepath.Join(config.Home(), ".actools", "last-update-check.txt")

	if err := os.MkdirAll(filepath.Dir(lastUpdateFilename), 0700); err != nil {
		return errors.Trace(err)
	}

	lastUpdate := time.Time{}
	lastUpdateContent, err := ioutil.ReadFile(lastUpdateFilename)
	if err != nil && !os.IsNotExist(err) {
		return errors.Trace(err)
	} else if err == nil {
		if err := lastUpdate.UnmarshalText(lastUpdateContent); err != nil {
			return errors.Trace(err)
		}
	}

	if time.Now().Sub(lastUpdate) > 1*time.Hour {
		reply, err := http.Get("https://tools.altipla.consulting/version-manifest/actools")
		if err != nil {
			return errors.Trace(err)
		}
		defer reply.Body.Close()

		content, err := ioutil.ReadAll(reply.Body)
		if err != nil {
			return errors.Trace(err)
		}
		version := strings.TrimSpace(string(content))

		if version != config.Version {
			log.WithFields(log.Fields{"current": config.Version, "latest": version}).Error("actools is not updated")

			log.Warning()
			log.Warning("Run the following command to install the latest version:")
			log.Warning()
			log.Warning("\tcurl https://tools.altipla.consulting/install/actools | bash")
			log.Warning()

			os.Exit(2)

			return nil
		}

		lastUpdateContent, err = time.Now().MarshalText()
		if err != nil {
			return errors.Trace(err)
		}
		if err := ioutil.WriteFile(lastUpdateFilename, lastUpdateContent, 0600); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}
