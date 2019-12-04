package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	ms "auto_release/models/mysql"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

//ThirdController for backstage api
type ThirdController struct {
	beego.Controller
	JSONRetMsg
}

//GetEnvs return envs
func (c *ThirdController) DeleteUser() {
	phone, _ := c.GetInt("phone")
	var stderr, stdout bytes.Buffer
	s := fmt.Sprintf("cd /home/tools/del_user_wechat/;./repair.py %d", phone)
	logs.Info(s)
	cmd := exec.Command("ssh", "STG-ALL-83", s)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if nil != err {
		c.setError(1, fmt.Sprintf("STG-ALL-83 exec[%s] failed. ", s))
		logs.Error("exec[%s] failed. Error:[%s]", s, stderr.String())
		goto end
	}

	logs.Info("DeleteUser %d passed", phone)

end:
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}

func (c *ThirdController) CasUser() {
	employee_id := c.GetString("employee_id")

	flt := &ms.CasUser{
		EmployeeId: employee_id,
	}

	client := &http.Client{}
	url := "https://cas.360gst.com/api/getUsers"

	tmp, _ := json.Marshal(&flt)
	logs.Info(string(tmp))
	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(tmp))
	r.Header.Set("Content-Type", "application/json")

	rsp, _ := client.Do(r)
	defer rsp.Body.Close()

	body, _ := ioutil.ReadAll(rsp.Body)

	logs.Info(string(body))
	var resp ms.CasUserResp
	json.Unmarshal(body, &resp)
	logs.Info(resp)

	for _, v := range resp.Users {
		if v.EmployeeId == employee_id {
			c.setData(v)
			goto end
		}
	}

	c.setError(1, "未找到cas用户")
	logs.Error("employee_id:[%s] cant find casUser", employee_id)

end:
	c.Data["json"] = c.GenRetJSON()
	c.ServeJSON()
}
