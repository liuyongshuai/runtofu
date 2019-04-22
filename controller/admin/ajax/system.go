/**
 * @author      Liu Yongshuai
 * @package     ajax
 * @date        2018-02-16 22:33
 */
package ajax

import (
	"bufio"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/liuyongshuai/goUtils"
	"github.com/liuyongshuai/runtofu/configer"
	"github.com/liuyongshuai/runtofu/model"
	"strings"
	"time"
)

type AdminAjaxSystemController struct {
	AdminAjaxBaseController
}

//返回数据信息
func (bc *AdminAjaxSystemController) Run() {
	action := bc.GetParam("action", "").ToString()
	fn := make(map[string]func())
	fn["delMenu"] = bc.delMenu
	fn["getMenuInfo"] = bc.getMenuInfo
	fn["modifyMenuInfo"] = bc.modifyMenuInfo
	fn["upMenu"] = bc.upMenu
	fn["downMenu"] = bc.downMenu
	fn["uploadImage"] = bc.uploadImage
	if f, ok := fn[action]; ok {
		f()
	} else {
		bc.Notice(nil, 100100, "非法操作")
	}
}

//删除菜单操作
func (bc *AdminAjaxSystemController) delMenu() {
	menuId, _ := bc.GetParam("menu_id", 0).ToInt()
	_, e := model.MAdminMenu.DeleteMenuInfo(menuId)
	if e != nil {
		bc.Notice(nil, 100100, e.Error())
	} else {
		bc.Notice(nil)
	}
}

//获取菜单详细信息
func (bc *AdminAjaxSystemController) getMenuInfo() {
	menuId, _ := bc.GetParam("menu_id", 0).ToInt()
	mInfo, err := model.MAdminMenu.GetAdminMenuInfo(menuId)
	if err != nil {
		bc.Notice(nil, 100100, err.Error())
	} else {
		if mInfo.MenuId > 0 {
			bc.Notice(mInfo)
			return
		}
	}
	bc.Notice(nil)
}

//获取菜单详细信息
func (bc *AdminAjaxSystemController) modifyMenuInfo() {
	menuId, _ := bc.GetParam("menu_id", 0).ToInt()
	name := bc.GetParam("name", "").ToString()
	path := bc.GetParam("path", "").ToString()
	icon := bc.GetParam("icon", "").ToString()
	iconColor := bc.GetParam("icon_color", "").ToString()
	parentId, _ := bc.GetParam("parent_menu_id", -1).ToInt()

	//先检查当前菜单信息
	menuInfo, _ := model.MAdminMenu.GetAdminMenuInfo(menuId)
	if menuInfo.MenuId > 0 { //修改菜单信息
		data := make(map[string]interface{})
		if len(name) > 0 {
			data["menu_name"] = name
		}
		if len(path) > 0 {
			data["menu_path"] = path
		}
		if len(icon) > 0 {
			data["icon_name"] = icon
		}
		if len(iconColor) > 0 {
			data["icon_color"] = iconColor
		}
		if parentId >= 0 {
			data["parent_menu_id"] = parentId
		}
		_, e := model.MAdminMenu.UpdateMenuInfo(menuId, data)
		if e != nil {
			bc.Notice(nil, 100100, e.Error())
		} else {
			bc.Notice(nil)
		}
	} else { //添加菜单信息
		if parentId < 0 {
			parentId = 0
		}
		_, _, e := model.MAdminMenu.AddMenuInfo(model.AdminMenuInfo{
			MenuName:     name,
			MenuPath:     path,
			IconName:     icon,
			IconColor:    iconColor,
			ParentMenuId: parentId,
		})
		if e != nil {
			bc.Notice(nil, 100100, e.Error())
		} else {
			bc.Notice(nil)
		}
	}
}

//上调菜单顺序
func (bc *AdminAjaxSystemController) upMenu() {
	menuId, _ := bc.GetParam("menu_id", 0).ToInt()
	menuInfo, err := model.MAdminMenu.GetAdminMenuInfo(menuId)
	if err != nil || menuInfo.MenuId <= 0 {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		bc.Notice(nil, 100100, "提取菜单信息失败 "+errMsg)
		return
	}
	cond := make(map[string]interface{})
	cond["parent_menu_id"] = menuInfo.ParentMenuId
	total := model.MAdminMenu.GetAdminMenuTotal(cond)
	mList := model.MAdminMenu.GetAdminMenuList(cond, 1, int(total))
	if total <= 0 || len(mList) <= 0 {
		bc.Notice(nil, 100100, "提取同级菜单信息失败："+menuInfo.MenuName)
		return
	}
	if mList[0].MenuId == menuId {
		bc.Notice(nil, 100100, "当前菜单已经是第一个了："+menuInfo.MenuName)
		return
	}
	curIndex := 0
	for i, mInfo := range mList {
		if mInfo.MenuId == menuId {
			curIndex = i
			break
		}
	}
	mList[curIndex-1], mList[curIndex] = mList[curIndex], mList[curIndex-1]
	for i, mInfo := range mList {
		data := make(map[string]interface{})
		data["menu_sort"] = int(total) - i
		model.MAdminMenu.UpdateMenuInfo(mInfo.MenuId, data)
	}
	bc.Notice(nil)
}

//下调菜单顺序
func (bc *AdminAjaxSystemController) downMenu() {
	menuId, _ := bc.GetParam("menu_id", 0).ToInt()
	menuInfo, err := model.MAdminMenu.GetAdminMenuInfo(menuId)
	if err != nil || menuInfo.MenuId <= 0 {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		bc.Notice(nil, 100100, "提取菜单信息失败 "+errMsg)
		return
	}
	cond := make(map[string]interface{})
	cond["parent_menu_id"] = menuInfo.ParentMenuId
	total := model.MAdminMenu.GetAdminMenuTotal(cond)
	mList := model.MAdminMenu.GetAdminMenuList(cond, 1, int(total))
	if total <= 0 || len(mList) <= 0 {
		bc.Notice(nil, 100100, "提取同级菜单信息失败："+menuInfo.MenuName)
		return
	}
	if mList[len(mList)-1].MenuId == menuId {
		bc.Notice(nil, 100100, "当前菜单已经是最后一个了："+menuInfo.MenuName)
		return
	}
	curIndex := 0
	for i, mInfo := range mList {
		if mInfo.MenuId == menuId {
			curIndex = i
			break
		}
	}
	mList[curIndex+1], mList[curIndex] = mList[curIndex], mList[curIndex+1]
	for i, mInfo := range mList {
		data := make(map[string]interface{})
		data["menu_sort"] = int(total) - i
		model.MAdminMenu.UpdateMenuInfo(mInfo.MenuId, data)
	}
	bc.Notice(nil)
}

//editor.md要求返回的数据格式
type UpImgResp struct {
	Success int64  `json:"success"`
	Message string `json:"message"`
	Url     string `json:"url"`
}

//上传图片信息
func (bc *AdminAjaxSystemController) uploadImage() {
	allImageExt := []string{"jpg", "jpeg", "gif", "png", "bmp", "PNG", "JPG", "JPEG", "GIF"}
	resp := UpImgResp{Success: 1}
	mf := bc.Ctx.Request.MultipartForm
	if mf == nil {
		resp.Success = 0
		resp.Message = "上传图片失败：未接收到上传的文件"
		bc.RenderJson(resp)
		return
	}
	form := mf.File
	if len(form) <= 0 {
		resp.Success = 0
		resp.Message = "上传图片失败：未接收到上传的文件"
		bc.RenderJson(resp)
		return
	}

	//form里就没这个字段的话
	imgs, ok := form["editormd-image-file"]
	if !ok {
		resp.Success = 0
		resp.Message = "上传图片失败：未包含指定字段名"
		bc.RenderJson(resp)
		return
	}
	//至少有一个
	if len(imgs) <= 0 {
		resp.Success = 0
		resp.Message = "上传图片失败：上传文件为空"
		bc.RenderJson(resp)
		return
	}
	//对第一个图片进行解析
	img := imgs[0]
	fname := img.Filename
	//提取mime头信息，判断上传的是否为图片
	mimeHeader := img.Header
	contentType := mimeHeader.Get("Content-Type")
	if !strings.Contains(contentType, "image") && !strings.Contains(contentType, "IMAGE") {
		resp.Success = 0
		resp.Message = "上传图片失败：上传文件非图片，其mime为 " + contentType
		bc.RenderJson(resp)
		return
	}
	size := img.Size
	//取其后缀，包括点，如".jpg"
	ext := ""
	pos := strings.LastIndex(fname, ".")
	if pos > 0 {
		ext = fname[pos:]
	}
	isAllow := false
	for _, e := range allImageExt {
		if "."+e == ext {
			isAllow = true
			break
		}
	}
	if !isAllow {
		resp.Success = 0
		resp.Message = "上传图片失败：允许的文件后缀有 " + strings.Join(allImageExt, ",")
		bc.RenderJson(resp)
		return
	}
	//打开这个文件
	fp, err := img.Open()
	if err != nil {
		resp.Success = 0
		resp.Message = "上传图片失败 " + err.Error()
		bc.RenderJson(resp)
		return
	}
	defer fp.Close()
	data := make([]byte, size)
	_, err = fp.Read(data)
	if err != nil {
		resp.Success = 0
		resp.Message = "上传图片失败 " + err.Error()
		bc.RenderJson(resp)
		return
	}
	md5 := goUtils.MD5(string(data))
	key := md5 + ext
	fp.Seek(0, 0)
	rder := bufio.NewReader(fp)
	maxAge := 63072000
	cacheControl := oss.CacheControl(fmt.Sprintf("max-age=%d", maxAge))
	expire := oss.Expires(time.Now().Add(time.Duration(maxAge) * time.Second))
	err = model.AliyunOSSBucket.PutObject(key, rder, cacheControl, expire)
	if err != nil {
		resp.Success = 0
		resp.Message = "上传图片失败 " + err.Error()
		bc.RenderJson(resp)
		return
	}
	resp.Success = 1
	resp.Url = configer.GetConfiger().Common.ImagePrefix + "/" + key
	bc.RenderJson(resp)
}
