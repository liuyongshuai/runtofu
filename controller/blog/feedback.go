/**
 * @author      Liu Yongshuai
 * @package     blog
 * @date        2018-03-08 20:47
 */
package blog

type FeedBackController struct {
	RunToFuBaseController
}

//提取首页的文章列表
func (bc *FeedBackController) Run() {
	bc.TplName = "feedback.tpl"
	bc.RenderHtml()
}
