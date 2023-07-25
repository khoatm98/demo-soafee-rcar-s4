package renesasrefappota

import (

    "encoding/json"
    "os"
    "path/filepath"
    "time"

    "github.com/aoscloud/aos_common/aoserrors"
    "github.com/aoscloud/aos_common/aostypes"
    "github.com/aoscloud/aos_common/image"
    "github.com/aoscloud/aos_updatemanager/updatehandler"
    log "github.com/sirupsen/logrus"

)

/***********************************************************************************************************************
 * Consts
 **********************************************************************************************************************/

const (
    otaCommandSyncCompose = 0
    otaCommandDownload	= 1
    otaCommandInstall 	= 2
    otaCommandActivate	= 3
    otaCommandRevert  	= 4
)

const (
    otaStatusSuccess = 0
    otaStatusFailed  = 1
)

const otaDefaultTimeout = 10 * time.Minute

const (
    idleState = iota
    preparedState
    updatedState
)

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// RenesasUpdateModule update components using Renesas OTA master.
type RenesasUpdateModule struct {
    id         	string
    config     	moduleConfig
    storage    	updatehandler.ModuleStorage
    State      	updateState `json:"state"`
    VendorVersion  string  	`json:"vendorVersion"`
    PendingVersion string  	`json:"pendingVersion"`
    rebooter   	Rebooter
}

type moduleConfig struct {
    TargetFile   	string        	`json:"targetFile"`
    Timeout      	aostypes.Duration `json:"timeout"`
}

type Rebooter interface {
    Reboot() (err error)
}
type updateState int

/***********************************************************************************************************************
 * Public
 **********************************************************************************************************************/

// New creates fs update module instance.
func New(id string, config json.RawMessage, storage updatehandler.ModuleStorage, rebooter Rebooter) (updatehandler.UpdateModule, error) {
    log.WithField("module", id).Debug("Create renesasupdate module")

    module := &RenesasUpdateModule{
   	 id:  	id,
   	 storage: storage,
   	 rebooter: rebooter,
   	 config: moduleConfig{
   		 Timeout: aostypes.Duration{Duration: otaDefaultTimeout},
   	 },
    }

    if err := json.Unmarshal(config, &module.config); err != nil {
   	 return nil, aoserrors.Wrap(err)
    }
	
    if module.config.TargetFile == "" {
   	 return nil, aoserrors.New("target file name should be configured")
    }

    state, err := storage.GetModuleState(id)
    if err != nil {
   	 return nil, aoserrors.Wrap(err)
    }

    if len(state) > 0 {
   	 if err := json.Unmarshal(state, module); err != nil {
   		 return nil, aoserrors.Wrap(err)
   	 }
    }
    return module, nil
}

// Close closes DualPartModule.
func (module *RenesasUpdateModule) Close() error {
    log.WithFields(log.Fields{"id": module.id}).Debug("Close renesasupdate module")

    return nil
}

// GetID returns module ID.
func (module *RenesasUpdateModule) GetID() string {
    return module.id
}

// Init initializes module.
func (module *RenesasUpdateModule) Init() error {
    return nil
}

// GetVendorVersion returns vendor version.
func (module *RenesasUpdateModule) GetVendorVersion() (string, error) {
    return module.VendorVersion, nil
}

// Prepare preparing image.
func (module *RenesasUpdateModule) Prepare(imagePath string, vendorVersion string, annotations json.RawMessage) error {
    log.WithFields(log.Fields{
   	 "id":        	module.id,
   	 "imagePath": 	imagePath,
   	 "vendorVersion": vendorVersion,
    }).Debug("Prepare renesasupdate module")

    if module.State == preparedState {
   	 return nil
    }

    if err := os.MkdirAll(filepath.Dir(module.config.TargetFile), 0o700); err != nil {
   	 return aoserrors.Wrap(err)
    }

    file, err := os.Create(module.config.TargetFile)
    if err != nil {
   	 return aoserrors.Wrap(err)
    }
    file.Close()

    if _, err := image.CopyFromGzipArchive(module.config.TargetFile, imagePath); err != nil {
   	 return aoserrors.Wrap(err)
    }

    module.PendingVersion = vendorVersion
    //Create flag for updating request
    log.WithFields(log.Fields{"id": module.id}).Debug("Make dowload done flag")
    if err := os.MkdirAll("/var/aos/status/"+module.id+"/downloadedFlag", 0o700); err != nil {
   	 return aoserrors.Wrap(err)
    }
    if err := module.setState(preparedState); err != nil {
   	 return err
    }
	log.WithFields(log.Fields{"id": module.id}).Debug("Waiting for updates ...")
	for {
		// condition to terminate the loop
		if _, err := os.Stat("/var/aos/status/"+module.id+"/downloadedFlag"); os.IsNotExist(err) {
			log.WithFields(log.Fields{"id": module.id}).Debug("Make update flag")
			if err := os.MkdirAll("/var/aos/status/"+module.id+"/updateFlag", 0o700); err != nil {
				return aoserrors.Wrap(err)
			}
			break
		}
	}
    return nil
}

// Update updates module.
func (module *RenesasUpdateModule) Update() (rebootRequired bool, err error) {
    log.WithFields(log.Fields{"id": module.id}).Debug("Update renesas update module")

    if module.State == updatedState {
   	 return false, nil
    }
    
    log.WithFields(log.Fields{"id": module.id}).Debug("Check update flag ...")
	if _, err := os.Stat("/var/aos/status/"+module.id+"/updateFlag"); !os.IsNotExist(err) {
		log.WithFields(log.Fields{"id": module.id}).Debug("On updating process...")
	}

	for {
		// condition to terminate the loop
		if _, err := os.Stat("/var/aos/status/"+module.id+"/updateFlag"); os.IsNotExist(err) {
			break
		}
	}

    module.VendorVersion, module.PendingVersion = module.PendingVersion, module.VendorVersion
    log.WithFields(log.Fields{"id": module.id}).Debug("Done updating")
    if err := module.setState(updatedState); err != nil {
   	 return false, err
    }

    return false, nil
}

// Revert reverts update.
func (module *RenesasUpdateModule) Revert() (rebootRequired bool, err error) {
    log.WithFields(log.Fields{"id": module.id}).Debug("Revert renesas update module")

    if module.State == idleState {
   	 return false, nil
    }

    if module.State == updatedState {
   	 module.VendorVersion, module.PendingVersion = module.PendingVersion, module.VendorVersion
    }

    if err := module.setState(idleState); err != nil {
   	 return false, err
    }

    return true, nil
}

// Apply applies update.
func (module *RenesasUpdateModule) Apply() (rebootRequired bool, err error) {
    log.WithFields(log.Fields{"id": module.id}).Debug("Apply renesas update module")
    
    if module.State == idleState {
   	 return false, nil
    }
    
    if err := module.setState(idleState); err != nil {
   	 return false, err
    }

    return true, nil
}

// Reboot performs module reboot.
func (module *RenesasUpdateModule) Reboot() (err error) {
    if module.rebooter != nil {
   	 log.WithFields(log.Fields{"id": module.id}).Debug("Reboot renesas update module")

   	 if err = module.rebooter.Reboot(); err != nil {
   		 return aoserrors.Wrap(err)
   	 }
    }

    return nil
}

func (state updateState) String() string {
    return []string{"idle", "prepared", "updated"}[state]
}

/***********************************************************************************************************************
 * Private
 **********************************************************************************************************************/

func (module *RenesasUpdateModule) setState(state updateState) error {
    log.WithFields(log.Fields{"id": module.id, "state": state}).Debugf("State changed")

    module.State = state

    data, err := json.Marshal(module)
    if err != nil {
   	 return aoserrors.Wrap(err)
    }

    if err = module.storage.SetModuleState(module.id, data); err != nil {
   	 return aoserrors.Wrap(err)
    }

    return nil
}
