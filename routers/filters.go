package routers

import (
    "strings"
    "github.com/astaxie/beego/context"
    "github.com/astaxie/beego/logs"

)

var filterLoggedInUser = func(ctx *context.Context) {
    url := ctx.Input.URL()
    logs.Info("url: [%v], 开始判断路由", url)
    if strings.HasPrefix(url, "/session") {
        logs.Info("[%v]路由，无需判断", url)
        return
    }
    sessionId, ok := ctx.Input.Session("current_user").(string)
    logs.Info("==========sessionID: [%v]", sessionId)
    //判断特殊路由 api release third
    if ok && (strings.HasPrefix(url, "/api") || strings.HasPrefix(url, "/release") || strings.HasPrefix(url, "/third")) {
        return
    }
    if !ok {
        logs.Info("没有写入current_user, 跳回登录界面")
        ctx.Redirect(302, "/session/login")
    } else if url != "/"{
        ctx.Redirect(302, "/")
    }
}
