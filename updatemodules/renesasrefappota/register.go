package renesasrefappota

import (
	"github.com/aoscloud/aos_updatemanager/updatehandler"
	"github.com/aoscloud/aos_updatemanager/updatemodules/partitions/rebooters/xenstorerebooter"
	"github.com/aoscloud/aos_updatemanager/updatemodules/partitions/updatechecker/systemdchecker"
	"github.com/khoatm98/demo-soafee-rcar-s4/updatemodules/renesasrefappota"
)

/*******************************************************************************
 * Init
 ******************************************************************************/

func init() {
	updatehandler.RegisterPlugin("renesasrefappota"func(id string, configJSON json.RawMessage,
		storage updatehandler.ModuleStorage,
		) (module updatehandler.UpdateModule, err error) {
			if len(configJSON) == 0 {
				return nil, aoserrors.Errorf("config for %s module is required", id)
			}

			var config moduleConfig

			if err = json.Unmarshal(configJSON, &config); err != nil {
				return nil, aoserrors.Wrap(err)
			}

			if module, err = renesasrefappota.New(id, config.VersionFile, config.UpdateDir,
				storage, &xenstorerebooter.XenstoreRebooter{}, systemdchecker.New(config.SystemdChecker)); err != nil {
				return nil, aoserrors.Wrap(err)
			}

			return module, nil
		})
}
