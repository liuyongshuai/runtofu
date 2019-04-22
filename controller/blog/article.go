/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     controller
 * @date        2018-02-02 21:50
 */
package blog

type ArticleController struct {
	RunToFuBaseController
}

//运行主逻辑
func (bc *ArticleController) Run() {
	bc.TplName = "article_info.tpl"
	aid, _ := bc.GetParam("aid", 0).ToInt64()
	bc.AddTplData("aid", aid)
	bc.RenderHtml()
}
