/** layuiAdmin.std-v1.0.0 LPPL License By http://www.layui.com/admin/ */
;
layui.define(["table", "form"],
    function(e) {
        var t = layui.$,
            i = layui.table;
        layui.form;
        i.render({
            elem: "#LAY-module-manage",
            url: layui.setter.reqbaseurl + "api/v1/module/getall",
            method:'POST',
            contentType:"application/json",
            where:{
               "userName":"aidznc"
            },
            headers:{"":""},
            parseData: function(res){ //res 即为原始返回的数据
                return {
                    "code": res.code, //解析接口状态
                    "msg": res.msg, //解析提示文本
                    "count": res.data.total, //解析数据长度
                    "data": res.data //解析数据列表
                };
            },
            cols: [[{
                type: "checkbox",
                fixed: "left"
            },
                {
                    field: "id",
                    width: 100,
                    title: "ID",
                    sort: !0
                },
                {
                    field: "name",
                    title: "模块名",
                    minWidth: 100
                },
                {
                    field: "createtime",
                    title: "创建时间",
                    sort: !0
                },
                {
                    field: "updatetime",
                    title: "更新时间",
                    sort: !0
                },
                {
                    title: "操作",
                    width: 150,
                    align: "center",
                    fixed: "right",
                    toolbar: "#table-module-info"
                }]],
            page: !0,
            limit: 30,
            height: "full-220",
            text: "对不起，加载出现异常！"
        }),
            i.on("tool(LAY-module-manage)",
                function(e) {
                    e.data;
                    if ("del" === e.event) layer.prompt({
                            formType: 1,
                            title: "敏感操作，请验证口令"
                        },
                        function(t, i) {
                            layer.close(i),
                                layer.confirm("真的删除行么",
                                    function(t) {
                                        e.del(),
                                            layer.close(t)
                                    })
                        });
                    else if ("edit" === e.event) {
                        t(e.tr);
                        layer.open({
                            type: 2,
                            title: "编辑用户",
                            content: "../../../example/user/user/userform.html",
                            maxmin: !0,
                            area: ["500px", "450px"],
                            btn: ["确定", "取消"],
                            yes: function(e, t) {
                                var l = window["layui-layer-iframe" + e],
                                    r = "LAY-user-front-submit",
                                    n = t.find("iframe").contents().find("#" + r);
                                l.layui.form.on("submit(" + r + ")",
                                    function(t) {
                                        t.field;
                                        i.reload("LAY-user-front-submit"),
                                            layer.close(e)
                                    }),
                                    n.trigger("click")
                            },
                            success: function(e, t) {}
                        })
                    }
                }),
            e("module", {})
    });