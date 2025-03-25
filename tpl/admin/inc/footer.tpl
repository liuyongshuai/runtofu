{{define "footer"}}
        {{if gt .userInfo.Uid 0}}
            </div>
            </div>
            {{/******展开和隐藏左侧的菜单栏******/}}
            <span class="glyphicon glyphicon-move tofu-expand-icon" isShow="1"></span>
            </div>

            {{/******返回顶部的小图标******/}}
            <div class="goto-top">
                <a href="#" title="返回顶部" data-toggle="tooltip" data-placement="top">
                    <span class="glyphicon glyphicon-home"></span>
                </a>
            </div>

            {{/******修改密码******/}}
            <div id="modify_user_passwd_action_dialog" class="modal" tabindex="-1" role="dialog">
                <div class="modal-dialog modal-sm">
                    <div class="modal-content">
                        <div class="modal-header">
                            <button type="button" class="close" data-dismiss="modal">
                                <span aria-hidden="true">×</span>
                            </button>
                            <h4 class="modal-title">修改密码</h4>
                        </div>
                        <div class="modal-body">
                            <div class="form-horizontal">
                                <div class="form-group">
                                    <label class="col-sm-2 control-label">原密码</label>
                                    <div class="col-sm-9">
                                        <input value="" class="form-control" id="user_old_passwd" placeholder="输入当前密码" type="password">
                                    </div>
                                </div>
                                  <div class="form-group">
                                      <label class="col-sm-2 control-label">新密码</label>
                                      <div class="col-sm-9">
                                          <input value="" class="form-control" id="user_new_passwd" placeholder="输入新密码" type="password">
                                      </div>
                                  </div>
                                <div class="form-group">
                                    <label class="col-sm-2 control-label">新密码</label>
                                    <div class="col-sm-9">
                                        <input value="" class="form-control" id="user_new_passwd_again" placeholder="再次输入新密码" type="password">
                                    </div>
                                </div>
                                <div class="form-group">
                                    <div class="col-sm-offset-5 col-sm-8">
                                        <button id="modify_user_passwd_action_submit_btn" class="btn btn-default">确定</button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

        {{end}}
    </body>
    </html>
{{ end }}