function err(msg, f) {
  layer.alert(msg, {icon: 2, closeBtn: 0, title: '出错了'}, f);
}

function succ(msg, f) {
  layer.msg(msg, {
      icon: 1,
      time: 1000
  }, f);
}

function handle_json(json, f) {
    if (json.msg.length > 0) {
        err(json.msg);
    } else {
        succ('恭喜，操作成功：）', f);
    }
}

function my_confirm(msg, btns, yes_func, no_func) {
    layer.confirm(msg, {
        btn: btns
    }, yes_func, no_func);
}

function show_profile_page() {
  $.get("/me.json", {}, function(json){
    if (json.msg.length>0){
      err(json.msg);
      return
    }
    $("#username").val(json.data.name);
    $("#email").val(json.data.email);
    $("#phone").val(json.data.phone);
    $("#qq").val(json.data.qq);
    $("#wechat").val(json.data.wechat);
  });
  $('#profile_div').modal();
}

function show_chpwd_page() {
  $('#chpwd_div').modal();
}

function htmlEncode(value){
  return $('<div/>').text(value).html();
}


function del_strategy(item_id){
    my_confirm("确定删除此策略？", [ '确定', '取消' ], function() {
        $.post('/strategy/'+item_id+'/delete', {}, function(json) {
            handle_json(json, function (){location.reload()})
        });
    });
}

function get_strategy(id){

  url = "/strategy/" + id
  $.post(url, {}, function(res){
    $("#url").val(res.data.url)
    $("#expect_code").val(res.data.expect_code)
    $("#timeout").val(res.data.timeout)
    $("#data").val(res.data.data)
    $("#keywords").val(res.data.keywords)
    $("#monitor_idc").val(res.data.monitor_idc)
    $("#note").val(res.data.note)
    $("#times").val(res.data.times)
    $("#max_step").val(res.data.max_step)
    $("#tags").val(res.data.tag.replace(/,/g,"\n"))
  }, "json");
};

function update_strategy(id){
    $.post('/strategy/' + id + '/edit', {
      "url": $('#url').val(),
      "expect_code": $('#expect_code').val(),
      "timeout": $('#timeout').val(),
      "times": $('#times').val(),
      "teams": $('#teams').val(),
      "max_step": $('#max_step').val(),
      "tags": $('#tags').val(),
      "note": $('#note').val(),
      "keywords": $('#keywords').val(),
      "monitor_idc": $('#monitor_idc').val(),
      "data": $('#data').val()
    }, function(json) {
      handle_json(json, function (){location.reload()})
    });
}

function add_strategy() {
    $.post("/strategy/add", {
      "url": $('#url').val(),
      "expect_code": $('#expect_code').val(),
      "timeout": $('#timeout').val(),
      "times": $('#times').val(),
      "teams": $('#teams').val(),
      "max_step": $('#max_step').val(),
      "tags": $('#tags').val(),
      "note": $('#note').val(),
      "keywords": $('#keywords').val(),
      "monitor_idc": $('#monitor_idc').val(),
      "data": $('#data').val()
    }, function(json){
        handle_json(json, function(){
            location.href=$("/").val();
        });
    });
}

function register() {
    $.post("/auth/register", {
        "username": $("#username").val(),
        "password": $("#password").val(),
        "repeat": $("#repeat").val()
    }, function(json){
        handle_json(json, function(){
            location.href=$("#callback").val();
        });
    });
}

function login() {
    $.post("/auth/login", {
        "username": $("#username").val(),
        "password": $("#password").val()
    }, function(json){
        handle_json(json, function(){
            location.href=$("#callback").val();
        });
    });
}

function update_profile() {
    $.post("/me/profile", {
        "email": $("#email").val(),
        "phone": $("#phone").val(),
        "wechat": $("#wechat").val()
    }, function(json){
        handle_json(json);
    });
}

function chpwd() {
    $.post("/me/chpwd", {
        "old_password": $("#old_password").val(),
        "new_password": $("#new_password").val(),
        "repeat": $("#repeat_password").val()
    }, function(json){
        handle_json(json, function() {location.href="/auth/login"});
    });
}