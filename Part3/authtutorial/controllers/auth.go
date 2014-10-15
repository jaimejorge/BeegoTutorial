package controllers

import (
	"authtutorial/models"
	"authtutorial/utils"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	_ "github.com/go-sql-driver/mysql"
)

var sessionName = beego.AppConfig.String("SessionName")

type LoginController struct {
	beego.Controller
}

func (this *LoginController) LoginView() {
	this.TplNames = "login.html"
}

func (this *LoginController) Login() {
	username := this.GetString("username")
	password := this.GetString("password")

	var user models.User
	if VerifyUser(&user, username, password) {
		v := this.GetSession(sessionName)
		if v == nil {
			this.SetSession(sessionName, user.Id)
		}
		this.Redirect("/secret", 302)

	} else {
		this.Redirect("/register", 302)
	}

}

func (this *LoginController) Logout() {
	this.DelSession(sessionName)
	this.Redirect("/login", 302)
}

func (this *LoginController) RegisterView() {
	this.TplNames = "register.html"
}

func (this *LoginController) Register() {
	username := this.GetString("username")
	password := this.GetString("password")
	passwordre := this.GetString("passwordre")
	test := models.RegisterForm{Username: username, Password: password, PasswordRe: passwordre}

	valid := validation.Validation{}
	b, err := valid.Valid(&test)
	if err != nil {
	}
	if !b {
		for _, err := range valid.Errors {
			fmt.Println(err.Key, err.Message)
		}
	} else {
		salt := utils.GetRandomString(10)
		encodedPwd := salt + "$" + utils.EncodePassword(password, salt)

		o := orm.NewOrm()
		o.Using("default")

		user := new(models.User)
		user.Username = username
		user.Password = encodedPwd
		user.Rands = salt

		o.Insert(user)

		this.Redirect("/", 302)

	}
	this.TplNames = "register.html"
}

func (this *LoginController) SecretView() {
	this.TplNames = "secret.html"
}

func HasUser(user *models.User, username string) bool {
	var err error
	qs := orm.NewOrm()
	user.Username = username
	err = qs.Read(user, "Username")
	if err == nil {
		return true
	}
	return false
}

func VerifyPassword(rawPwd, encodedPwd string) bool {
	// split
	var salt, encoded string
	salt = encodedPwd[:10]
	encoded = encodedPwd[11:]

	return utils.EncodePassword(rawPwd, salt) == encoded
}

func VerifyUser(user *models.User, username, password string) (success bool) {
	// search user by username or email
	if HasUser(user, username) == false {
		return
	}
	if VerifyPassword(password, user.Password) {
		// success
		success = true
	}
	return
}

var FilterUser = func(ctx *context.Context) {
	_, ok := ctx.Input.Session(sessionName).(int)
	if !ok && ctx.Input.Uri() != "/login" && ctx.Input.Uri() != "/register" {
		ctx.Redirect(302, "/login")
	}
}
