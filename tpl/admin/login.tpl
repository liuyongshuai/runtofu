{{template "header" .}}
<div class="login">
    <div class="form-horizontal" role="form">
        <div class="form-group">
            <label for="menu_name" class="col-sm-3 control-label">用户名</label>
            <div class="col-sm-9">
                <input value="" class="form-control" id="name" name="name" placeholder="user name" type="text">
            </div>
        </div>
        <div class="form-group">
            <label for="menu_name" class="col-sm-3 control-label">密&nbsp;&nbsp;&nbsp;&nbsp;码</label>
            <div class="col-sm-9">
                <input value="" class="form-control" id="passwd" name="passwd" placeholder="password" type="password">
            </div>
        </div>
        <div class="form-group">
            <div class="col-sm-offset-5 col-sm-8">
                <button id="login_submit_action" class="btn btn-primary">Login</button>
            </div>
        </div>
    </div>
</div>
{{template "footer" .}}