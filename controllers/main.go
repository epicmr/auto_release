package controllers

import (
	ms "auto_release/models/mysql"
    "encoding/hex"
    "encoding/json"
    "io/ioutil"
    "os"
    "strconv"
    "time"

    "github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
    "github.com/wumansgy/goEncrypt"
)

//MainController main controller
type MainController struct {
	beego.Controller
    JSONRetMsg
}

func (c *MainController) Prepare() {
	//c.Layout = "layout/app.tpl"
}

// Home goes to index.html file from ./../dist
func (c *MainController) Home() {
	//c.TplName = "main/index.tpl"
	c.TplName = "main/index.html"
}

func (c *MainController) Login() {
	c.TplName = "session/new.tpl"
}

func (c *MainController) Create() {
	phone := c.GetString("phone")
	password := c.GetString("password")
    if phone == "" || password == "" {
        logs.Error("参数非法，回到登录界面")
        c.Redirect("/session/login", 302)
        c.StopRun()
    }

	db, _ := ms.InitDb()
	var user ms.User
	db.Where("phone = ?", phone).Find(&user)
    
    //解密验证
    privateKey := readFile("private.pem")
    bytePasswd, err := hex.DecodeString(user.Password)
    if err != nil {
        logs.Error("DecodeString failed, error: [%v]", err)
    }
    plaintext, err := goEncrypt.RsaDecrypt(bytePasswd, privateKey)
    if err != nil {
        logs.Error("解密失败, Error: [%v]", err)
    }
	
    if string(plaintext) == password {
		logs.Info("验证成功，开始跳转")
		c.SetSession("current_user", phone)
		c.Redirect("/", 302)
		c.StopRun()
	}
	logs.Error("验证失败，回到登录界面")
	c.Redirect("/session/login", 302)
}

func (c *MainController) Logout() {
	c.DelSession("current_user")
	c.DestroySession()
	c.Redirect("/session/login", 302)
}

//Register new guy
func (c *MainController) Register() {
    var user    ms.User
    db ,_ :=ms.InitDb()

    //检验参数
    json.Unmarshal(c.Ctx.Input.RequestBody, &user)
    if user.Phone == "" || user.Password == "" {
        c.setError(2, "注册失败，缺少参数")
        logs.Error("注册失败，缺少参数。Phone: [%v], Passwprd: [%v]", user.Phone, user.Password)
        c.Data["json"] = c.GenRetJSON()
        c.ServeJSON()
    }

    if db.Where("phone = ?", user.Phone).First(&user).RecordNotFound() {
        phoneLast4, _ := strconv.Atoi(user.Phone[7:])
        user.UserID = uint64(time.Now().Unix() * 10000) + uint64(phoneLast4)

        //密码加密存储
        publicKey := readFile("public.pem")
        crypttext, err := goEncrypt.RsaEncrypt([]byte(user.Password), publicKey)
        if err != nil {
            logs.Error("加密失败，Error: [%v]", err)
        }
        user.Password = hex.EncodeToString(crypttext)
        
        db.Create(&user)
    } else {
        c.setError(2, "注册失败，用户已注册过")
        logs.Error("注册失败，[%v]已注册过", user.Phone)
    }
    c.Data["json"] = c.GenRetJSON()
    c.ServeJSON()
}

//简单封装读取文件
func readFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		logs.Error("openFile fail [%v]", err)
		return []byte{}
	}
	defer file.Close()

	result, err := ioutil.ReadAll(file)
	if err != nil {
		logs.Error("ReadAll fail [%v]", err)
		return []byte{}
	}

	return result
}
