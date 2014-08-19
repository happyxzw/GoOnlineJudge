package controller

import (
	"GoOnlineJudge/class"
	"GoOnlineJudge/config"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
)

type user struct {
	Uid string `json:"uid"bson:"uid"`
	Pwd string `json:"pwd"bson:"pwd"`

	Nick   string `json:"nick"bson:"nick"`
	Mail   string `json:"mail"bson:"mail"`
	School string `json:"school"bson:"school"`
	Motto  string `json:"motto"bson:"motto"`

	Privilege int `json:"privilege"bson:"privilege"`

	Solve  int `json:"solve"bson:"solve"`
	Submit int `json:"submit"bson:"submit"`

	Status int    `json:"status"bson:"status"`
	Create string `json:"create"bson:'create'`
}

type UserController struct {
	class.Controller
}

func (this *UserController) Signin(w http.ResponseWriter, r *http.Request) {
	class.Logger.Debug("User Login")
	this.Init(w, r)

	t := template.New("layout.tpl")
	t, err := t.ParseFiles("view/layout.tpl", "view/user_signin.tpl")
	if err != nil {
		http.Error(w, "tpl error", 500)
		return
	}

	this.Data["Title"] = "User Sign In"
	this.Data["IsUserSignIn"] = true
	err = t.Execute(w, this.Data)
	if err != nil {
		http.Error(w, "tpl error", 500)
		return
	}
}

func (this *UserController) Login(w http.ResponseWriter, r *http.Request) {
	class.Logger.Debug("User Login")
	this.Init(w, r)

	one := make(map[string]string)
	one["uid"] = r.FormValue("user[handle]")
	one["pwd"] = r.FormValue("user[password]")

	reader, err := this.PostReader(&one)
	if err != nil {
		http.Error(w, "read error", 500)
		return
	}

	response, err := http.Post(config.PostHost+"/user?login", "application/json", reader)
	if err != nil {
		http.Error(w, "post error", 500)
		return
	}
	defer response.Body.Close()

	var ret user
	err = this.LoadJson(response.Body, &ret)
	if err != nil {
		http.Error(w, "load error", 400)
		return
	}

	if response.StatusCode == 200 {
		if ret.Uid == "" {
			w.WriteHeader(400)
		} else {
			this.SetSession(w, r, "Uid", one["uid"])
			this.SetSession(w, r, "Privilege", strconv.Itoa(ret.Privilege))
			w.WriteHeader(200)
		}
		return
	} else {
		w.WriteHeader(response.StatusCode)
		return
	}
}

func (this *UserController) Signup(w http.ResponseWriter, r *http.Request) {
	class.Logger.Debug("User Sign Up")
	this.Init(w, r)

	this.Data["Title"] = "User Sign Up"
	this.Data["IsUserSignUp"] = true
	err := this.Execute(w, "view/layout.tpl", "view/user_signup.tpl")
	if err != nil {
		http.Error(w, "tpl error", 500)
		return
	}
}

func (this *UserController) Register(w http.ResponseWriter, r *http.Request) {
	class.Logger.Debug("User Register")
	this.Init(w, r)

	one := make(map[string]interface{})

	uid := r.FormValue("user[handle]")
	nick := r.FormValue("user[nick]")
	pwd := r.FormValue("user[password]")
	pwdConfirm := r.FormValue("user[confirmPassword]")
	one["mail"] = r.FormValue("user[mail]")
	one["school"] = r.FormValue("user[school]")
	one["motto"] = r.FormValue("user[motto]")

	ok := 1
	hint := make(map[string]string)
	response, err := http.Post(config.PostHost+"/user?list/uid?"+uid, "application/json", nil)
	if err != nil {
		http.Error(w, "post error", 500)
		return
	}
	defer response.Body.Close()

	if uid == "" {
		ok, hint["uid"] = 0, "Handle should not be empty."
	} else {
		ret := make(map[string][]*user)
		if response.StatusCode == 200 {
			err = this.LoadJson(response.Body, &ret)
			if err != nil {
				http.Error(w, "load error", 400)
				return
			}

			if len(ret["list"]) > 0 {
				ok, hint["uid"] = 0, "This handle is currently in use."
			}
		}
	}
	if nick == "" {
		ok, hint["nick"] = 0, "Nick should not be empty."
	}
	if len(pwd) < 6 {
		ok, hint["pwd"] = 0, "Password should contain at least six characters."
	}
	if pwd != pwdConfirm {
		ok, hint["pwdConfirm"] = 0, "Confirmation mismatched."
	}
	if ok == 1 {
		one["uid"] = uid
		one["nick"] = nick
		one["pwd"] = pwd
		one["pwdConfirm"] = pwdConfirm
		one["privilege"] = config.PrivilegePU
		//one["privilege"] = config.PrivilegeAD
		reader, err := this.PostReader(&one)
		if err != nil {
			http.Error(w, "read error", 500)
			return
		}

		response, err = http.Post(config.PostHost+"/user?insert", "application/json", reader)
		if err != nil {
			http.Error(w, "post error", 400)
			return
		}
		defer response.Body.Close()

		this.SetSession(w, r, "Uid", uid)
		this.SetSession(w, r, "Privilege", "1")
		w.WriteHeader(200)
	} else {
		b, err := json.Marshal(&hint)
		if err != nil {
			http.Error(w, "json error", 500)
			return
		}

		w.WriteHeader(400)
		w.Write(b)
	}
}

func (this *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	class.Logger.Debug("User Logout")
	this.Init(w, r)

	this.DeleteSession(w, r)
	w.WriteHeader(200)
}

func (this *UserController) Detail(w http.ResponseWriter, r *http.Request) {
	class.Logger.Debug("User Detail")
	this.Init(w, r)

	args := this.ParseURL(r.URL.String())
	uid := args["uid"]
	response, err := http.Post(config.PostHost+"/user?detail/uid?"+uid, "application/json", nil)
	if err != nil {
		http.Error(w, "post error", 500)
		return
	}
	defer response.Body.Close()

	var one user
	if response.StatusCode == 200 {
		err = this.LoadJson(response.Body, &one)
		if err != nil {
			http.Error(w, "load error", 400)
			return
		}
		this.Data["Detail"] = one
	}

	response, err = http.Post(config.PostHost+"/solution?achieve/uid?"+uid, "application/json", nil)
	if err != nil {
		http.Error(w, "post error", 500)
		return
	}
	defer response.Body.Close()

	solvedList := make(map[string][]int)
	if response.StatusCode == 200 {
		err = this.LoadJson(response.Body, &solvedList)
		if err != nil {
			http.Error(w, "load error", 400)
			return
		}
		this.Data["List"] = solvedList["list"]
	}
	//class.Logger.Debug(solvedList["list"])
	t := template.New("layout.tpl")
	t, err = t.ParseFiles("view/layout.tpl", "view/user_detail.tpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	this.Data["Title"] = "User Detail"
	if uid != "" && uid == this.Uid {
		this.Data["IsSettings"] = true
		this.Data["IsSettingsDetail"] = true
	}

	err = t.Execute(w, this.Data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (this *UserController) Settings(w http.ResponseWriter, r *http.Request) {
	class.Logger.Debug("User Settings")
	this.Init(w, r)

	if this.Privilege == config.PrivilegeNA {
		this.Data["Title"] = "Warning"
		this.Data["Info"] = "You must login!"
		t := template.New("layout.tpl")
		t, err := t.ParseFiles("view/layout.tpl", "view/400.tpl")
		if err != nil {
			http.Error(w, "tpl error", 500)
			return
		}
		err = t.Execute(w, this.Data)
		if err != nil {
			http.Error(w, "tpl error", 500)
			return
		}
		return
	}

	response, err := http.Post(config.PostHost+"/user?detail/uid?"+this.Uid, "application/json", nil)
	if err != nil {
		http.Error(w, "post error", 500)
		return
	}
	defer response.Body.Close()

	var one user
	if response.StatusCode == 200 {
		err = this.LoadJson(response.Body, &one)
		if err != nil {
			http.Error(w, "load error", 400)
			return
		}
		this.Data["Detail"] = one
	}

	response, err = http.Post(config.PostHost+"/solution?achieve/uid?"+this.Uid, "application/json", nil)
	if err != nil {
		http.Error(w, "post error", 500)
		return
	}
	defer response.Body.Close()

	solvedList := make(map[string][]int)
	if response.StatusCode == 200 {
		err = this.LoadJson(response.Body, &solvedList)
		if err != nil {
			http.Error(w, "load error", 400)
			return
		}
		this.Data["List"] = solvedList["list"]
	}

	t := template.New("layout.tpl")
	t, err = t.ParseFiles("view/layout.tpl", "view/user_detail.tpl")
	if err != nil {
		http.Error(w, "tpl error", 500)
		return
	}

	this.Data["Title"] = "User Settings"
	this.Data["IsSettings"] = true
	this.Data["IsSettingsDetail"] = true

	err = t.Execute(w, this.Data)
	if err != nil {
		http.Error(w, "tpl error", 500)
		return
	}
}

func (this *UserController) Edit(w http.ResponseWriter, r *http.Request) {
	class.Logger.Debug("User Edit")
	this.Init(w, r)

	if this.Privilege == config.PrivilegeNA {
		this.Data["Title"] = "Warning"
		this.Data["Info"] = "You must login!"
		t := template.New("layout.tpl")
		t, err := t.ParseFiles("view/layout.tpl", "view/400.tpl")
		if err != nil {
			http.Error(w, "tpl error", 500)
			return
		}
		err = t.Execute(w, this.Data)
		if err != nil {
			http.Error(w, "tpl error", 500)
			return
		}
		return
	}

	uid := this.Uid
	response, err := http.Post(config.PostHost+"/user?detail/uid?"+uid, "application/json", nil)
	if err != nil {
		http.Error(w, "post error", 500)
		return
	}
	defer response.Body.Close()

	var one user
	if response.StatusCode == 200 {
		err = this.LoadJson(response.Body, &one)
		if err != nil {
			http.Error(w, "load error", 400)
			return
		}
		this.Data["Detail"] = one
	}

	t := template.New("layout.tpl")
	t, err = t.ParseFiles("view/layout.tpl", "view/user_edit.tpl")
	if err != nil {
		http.Error(w, "tpl error", 500)
		return
	}

	this.Data["Title"] = "User Edit"
	this.Data["IsSettings"] = true
	this.Data["IsSettingsEdit"] = true

	err = t.Execute(w, this.Data)
	if err != nil {
		http.Error(w, "tpl error", 500)
		return
	}
}

func (this *UserController) Update(w http.ResponseWriter, r *http.Request) {
	class.Logger.Debug("User Update")
	this.Init(w, r)

	ok := 1
	hint := make(map[string]string)
	hint["uid"] = this.Uid

	one := make(map[string]interface{})
	one["nick"] = r.FormValue("user[nick]")
	one["mail"] = r.FormValue("user[mail]")
	one["school"] = r.FormValue("user[school]")
	one["motto"] = r.FormValue("user[motto]")

	if one["nick"] == "" {
		ok, hint["nick"] = 0, "Nick should not be empty."
	}

	if ok == 1 {
		reader, err := this.PostReader(&one)
		if err != nil {
			http.Error(w, "read error", 500)
			return
		}

		response, err := http.Post(config.PostHost+"/user?update/uid?"+this.Uid, "application/json", reader)
		if err != nil {
			http.Error(w, "post error", 500)
			return
		}
		defer response.Body.Close()

		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}

	b, err := json.Marshal(&hint)
	if err != nil {
		http.Error(w, "json error", 400)
		return
	}
	w.Write(b)
}

func (this *UserController) Pagepassword(w http.ResponseWriter, r *http.Request) {
	class.Logger.Debug("User Password Page")
	this.Init(w, r)

	if this.Privilege == config.PrivilegeNA {
		this.Data["Title"] = "Warning"
		this.Data["Info"] = "You must login!"
		t := template.New("layout.tpl")
		t, err := t.ParseFiles("view/layout.tpl", "view/400.tpl")
		if err != nil {
			http.Error(w, "tpl error", 500)
			return
		}
		err = t.Execute(w, this.Data)
		if err != nil {
			http.Error(w, "tpl error", 500)
			return
		}
		return
	}

	var err error
	t := template.New("layout.tpl")
	t, err = t.ParseFiles("view/layout.tpl", "view/user_password.tpl")
	if err != nil {
		http.Error(w, "tpl error", 500)
		return
	}

	this.Data["Title"] = "User Password"
	this.Data["IsSettings"] = true
	this.Data["IsSettingsPassword"] = true

	err = t.Execute(w, this.Data)
	if err != nil {
		http.Error(w, "tpl error", 400)
		return
	}
}

func (this *UserController) Password(w http.ResponseWriter, r *http.Request) {
	class.Logger.Debug("User Password")
	this.Init(w, r)

	ok := 1
	hint := make(map[string]string)
	hint["uid"] = this.Uid

	data := make(map[string]string)
	data["oldPassword"] = r.FormValue("user[oldPassword]")
	data["newPassword"] = r.FormValue("user[newPassword]")
	data["confirmPassword"] = r.FormValue("user[confirmPassword]")

	one := make(map[string]interface{})
	one["uid"] = this.Uid
	one["pwd"] = data["oldPassword"]

	reader, err := this.PostReader(&one)
	if err != nil {
		http.Error(w, "read error", 500)
		return
	}

	response, err := http.Post(config.PostHost+"/user?login", "application/json", reader)
	if err != nil {
		http.Error(w, "post error", 500)
		return
	}
	defer response.Body.Close()

	var ret user
	if response.StatusCode == 200 {
		err = this.LoadJson(response.Body, &ret)
		if err != nil {
			http.Error(w, "load error", 400)
			return
		}
	}

	if ret.Uid == "" {
		ok, hint["oldPassword"] = 0, "Old Password is Incorrect."
	}
	if len(data["newPassword"]) < 6 {
		ok, hint["newPassword"] = 0, "Password should contain at least six characters."
	}
	if data["newPassword"] != data["confirmPassword"] {
		ok, hint["confirmPassword"] = 0, "Confirmation mismatched."
	}

	if ok == 1 {
		one["pwd"] = data["newPassword"]
		reader, err = this.PostReader(&one)
		if err != nil {
			http.Error(w, "read error", 500)
			return
		}

		response, err = http.Post(config.PostHost+"/user?password/uid?"+this.Uid, "application/json", reader)
		if err != nil {
			http.Error(w, "post error", 400)
			return
		}
		defer response.Body.Close()

		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
	b, err := json.Marshal(&hint)
	if err != nil {
		http.Error(w, "json error", 400)
		return
	}

	w.Write(b)
}
