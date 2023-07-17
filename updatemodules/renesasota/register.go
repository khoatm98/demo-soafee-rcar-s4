package renesasrefappota

import (
	"github.com/aoscloud/aos_updatemanager/updatehandler"
)

/*******************************************************************************
 * Init
 ******************************************************************************/

func init() {
	updatehandler.RegisterPlugin("renesasrefappota", New)
}
