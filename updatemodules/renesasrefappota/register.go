
package renesasrefappota

import (
    "encoding/json"
    "github.com/aoscloud/aos_common/aoserrors"
    "github.com/aoscloud/aos_updatemanager/updatehandler"
    "github.com/aoscloud/aos_updatemanager/updatemodules/partitions/rebooters/xenstorerebooter"
)

/*******************************************************************************
 * Init
 ******************************************************************************/

func init() {
    updatehandler.RegisterPlugin("renesasrefappota",func(id string, configJSON json.RawMessage,
   	 storage updatehandler.ModuleStorage,
   	 ) (module updatehandler.UpdateModule, err error) {
   		 if len(configJSON) == 0 {
   			 return nil, aoserrors.Errorf("config for %s module is required", id)
   		 }

   		 var config moduleConfig

   		 if err = json.Unmarshal(configJSON, &config); err != nil {
   			 return nil, aoserrors.Wrap(err)
   		 }

   		 if module, err = New(id, configJSON,
   			 storage, &xenstorerebooter.XenstoreRebooter{}); err != nil {
   			 return nil, aoserrors.Wrap(err)
   		 }

   		 return module, nil
   	 })
}
