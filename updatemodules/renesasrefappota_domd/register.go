package renesasrefappota_domd

import (
	"encoding/json"

	"github.com/aoscloud/aos_common/aoserrors"

	"github.com/aoscloud/aos_updatemanager/updatehandler"
	"github.com/aoscloud/aos_updatemanager/updatemodules/partitions/modules/overlaymodule"
	"github.com/aoscloud/aos_updatemanager/updatemodules/partitions/rebooters/xenstorerebooter"
	"github.com/aoscloud/aos_updatemanager/updatemodules/partitions/updatechecker/systemdchecker"
)

/*******************************************************************************
 * Types
 ******************************************************************************/

type moduleConfig struct {
	VersionFile    string                `json:"versionFile"`
	UpdateDir      string                `json:"updateDir"`
	SystemdChecker systemdchecker.Config `json:"systemdChecker"`
}

/*******************************************************************************
 * Init
 ******************************************************************************/

func init() {
	updatehandler.RegisterPlugin("renesasrefappota_domd",
		func(id string, configJSON json.RawMessage,
			storage updatehandler.ModuleStorage,
		) (module updatehandler.UpdateModule, err error) {
			if len(configJSON) == 0 {
				return nil, aoserrors.Errorf("config for %s module is required", id)
			}

			var config moduleConfig

			if err = json.Unmarshal(configJSON, &config); err != nil {
				return nil, aoserrors.Wrap(err)
			}

			if module, err = overlaymodule.New(id, config.VersionFile, config.UpdateDir,
				storage, &xenstorerebooter.XenstoreRebooter{}, systemdchecker.New(config.SystemdChecker)); err != nil {
				return nil, aoserrors.Wrap(err)
			}

			return module, nil
		})
}