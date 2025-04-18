{{template "header" .}}
<style>
.login{width: 300px;margin-left:auto;margin-right:auto;padding-top:200px;}
</style>
<div class="login">
    <div class="form-horizontal" role="form">
        <div class="form-group" role="user_name">
            <label for="menu_name" class="col-sm-3 control-label">用户名</label>
            <div class="col-sm-9">
                <input value="" class="form-control" id="user_name" name="user_name" placeholder="user name" type="text">
            </div>
        </div>
        <div class="form-group" role="passwd">
            <label for="menu_name" class="col-sm-3 control-label">密&nbsp;&nbsp;&nbsp;&nbsp;码</label>
            <div class="col-sm-9">
                <input value="" class="form-control" id="passwd" name="passwd" placeholder="password" type="password">
            </div>
        </div>
        <div class="form-group" role="submit_btn">
            <div class="col-sm-offset-5 col-sm-8">
                <button id="login_submit_action" class="btn btn-primary">Login</button>
            </div>
        </div>
    </div>
</div>
<script>
$( function(){
    $("#login_submit_action").click( function(){
        var user_name = $("#user_name").val();
        var passwd = $("#passwd").val();
        comUtils.sendRequest( {
            url: "ajax/login",
            args: "user_name=" + user_name + "&passwd=" + passwd,
            onSuccess: function(){
                location.reload( true );
            }
        } );
    } );
} );
</script>
{{template "footer" .}}