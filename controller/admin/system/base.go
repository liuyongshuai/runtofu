/**
 * @author      Liu Yongshuai
 * @package     system
 * @date        2018-02-15 16:49
 */
package system

import (
	"fmt"
	"github.com/liuyongshuai/runtofu/controller/admin"
	"github.com/liuyongshuai/runtofu/model"
)

//system层的基类
type AdminSystemBaseController struct {
	admin.AdminBaseController
}

//校验是否为超级用户
func (bc *AdminSystemBaseController) Prepare() error {
	bc.UserInfo = bc.CheckLogin(true, func() {
		bc.Ctx.Redirect("/login")
	})
	if bc.UserInfo.Type != model.ADMIN_USER_TYPE_SUPER {
		fmt.Println(bc.UserInfo)
		return fmt.Errorf("只有超级用户可以访问")
	}
	bc.SetLeftMenu()
	return nil
}
