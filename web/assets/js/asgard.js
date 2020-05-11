var Asgard = {
    "doAction": function (url, success_url) {
        $.getJSON(url, function (data) {
            if (data.code == 200) {
                alert("操作成功");
                if (success_url) {
                    window.location.href = success_url;
                } else {
                    window.location.reload(true)
                }
            } else {
                alert(data.message);
            }
        })
    },
    "SearchSubmit": function () {
        $("#main").find("form").submit();
    },
    "ActionShow": function () {
        $(this).parent().parent().next().toggle();
    },
    "Jump": function () {
        var mod = $(this).parents(".cmd-info").attr("data-mod");
        var action = $(this).attr("action");
        var url = "/" + mod + "/" + action;
        window.location.href = url;
    },
    "JumpWithID": function () {
        var mod = $(this).parents(".cmd-info").attr("data-mod");
        var action = $(this).attr("action");
        var id = $(this).parents(".cmd-info").attr("data-bind");
        var url = "/" + mod + "/" + action + "?id=" + id;
        window.location.href = url;
    },
    "runAction": function () {
        var mod = $(this).parents(".cmd-info").attr("data-mod");
        var action = $(this).attr("action");
        var id = $(this).parents(".cmd-info").attr("data-bind");
        var url = "/" + mod + "/" + action + "?id=" + id;
        Asgard.doAction(url);
    },
    "listDelegate": function (mod) {
        $("#main").delegate("button[action=reset]", "click", function () {
            window.location.href = "/" + mod + "/list";
        }).delegate("button[action=add]", "click", function () {
            window.location.href = "/" + mod + "/add";
        }).delegate("button[action=submit]", "click", Asgard.SearchSubmit)
            .delegate("button[action=show]", "click", Asgard.ActionShow)
            .delegate("button[action=edit]", "click", Asgard.JumpWithID)
            .delegate("button[action=monitor]", "click", Asgard.JumpWithID)
            .delegate("button[action=archive]", "click", Asgard.JumpWithID)
            .delegate("a[action=out_log]", "click", Asgard.JumpWithID)
            .delegate("a[action=err_log]", "click", Asgard.JumpWithID)
            .delegate("button[action=start]", "click", Asgard.runAction)
            .delegate("button[action=restart]", "click", Asgard.runAction)
            .delegate("button[action=pause]", "click", Asgard.runAction)
            .delegate("button[action=delete]", "click", Asgard.runAction)
            .delegate("button[action=copy]", "click", Asgard.runAction);
    }
}