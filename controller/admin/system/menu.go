/**
 * @author      Liu Yongshuai
 * @package     system
 * @date        2018-02-14 16:43
 */
package system

// system层的基类
type AdminSystemMenuController struct {
	AdminSystemBaseController
}

// 开始
func (c *AdminSystemMenuController) Run() {
	c.TplName = "menu.tpl"
	c.RenderHtml()
}
