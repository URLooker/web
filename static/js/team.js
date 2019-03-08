function create_team() {
    $.post('/team/create', {
        'name' : $("#name").val(),
        'resume' : $("#resume").val(),
        'emails' : $("#emails").val()
    }, function(json) {
        handle_json(json, function(){
            location.href = '/teams';
        });
    });
}

function edit_team(team_id) {
    $.post('/team/edit?tid='+team_id, {
        'resume' : $("#resume").val(),
        'emails' : $("#emails").val(),
    }, function(json) {
        handle_json(json, function(){
            location.href = '/teams';
        });
    });
}

function del_team(team_id){
    my_confirm("确定删除此团队？", [ '确定', '取消' ], function() {
        $.post('/team/delete?tid='+team_id, {}, function(json) {
            handle_json(json, function (){location.reload()})
        });
    });
}