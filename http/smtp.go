package http

import (
	"log"
	"net"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/iambocai/mailer/g"
	"github.com/jordan-wright/email"
)

func configSmtpRoutes() {
	http.HandleFunc("/api/mail", func(w http.ResponseWriter, r *http.Request) {
		SendMailBySmtp(w, r, false)
	})
	http.HandleFunc("/api/attachmail", func(w http.ResponseWriter, r *http.Request) {
		SendMailBySmtp(w, r, true)
	})

}

func SendMailBySmtp(w http.ResponseWriter, r *http.Request, hasAttach bool) {
	var (
		err    error
		server string
		user   string
		passwd string
	)

	addr := strings.Split(r.RemoteAddr, ":")[0]
	allowList := strings.Split(g.Config().Http.WhiteList, ",")
	authorized := isValid(addr, allowList)

	//授权检查
	if authorized == false {
		http.Error(w, "remote not in whitelist", http.StatusBadRequest)
		return
	}

	if r.ContentLength == 0 {
		http.Error(w, "body is blank", http.StatusBadRequest)
		return
	}

	//根据是否有附件，使用不同的ParseForm方法
	if hasAttach == true {
		err = r.ParseMultipartForm(g.Config().Smtp.MaxBytes)
	} else {
		err = r.ParseForm()
	}
	if err != nil {
		http.Error(w, "param prase error", http.StatusBadRequest)
		return
	}

	if len(r.Form["tos"]) == 0 || len(r.Form["content"]) == 0 || len(r.Form["subject"]) == 0 {
		http.Error(w, "param tos,content or subject lost", http.StatusBadRequest)
		return
	}

	e := email.NewEmail()
	//收件人
	e.To = strings.Split(r.Form.Get("tos"), g.Config().Smtp.Spliter)
	//主题
	e.Subject = r.Form.Get("subject")

	//抄送，可选
	if len(r.Form.Get("cc")) != 0 {
		e.Cc = strings.Split(r.Form.Get("cc"), g.Config().Smtp.Spliter)
	}

	//密送，可选
	if len(r.Form.Get("bcc")) != 0 {
		e.Bcc = strings.Split(r.Form.Get("bcc"), g.Config().Smtp.Spliter)
	}

	//消息体格式，可选。默认为text，如果format设置为html，则以html方式发送
	if len(r.Form.Get("format")) != 0 && r.Form.Get("format") == "html" {
		e.HTML = []byte(r.Form.Get("content"))

	} else {
		e.Text = []byte(r.Form.Get("content"))
	}

	//附件，可选
	if hasAttach == true {
		m := r.MultipartForm

		//get the *fileheaders
		files := m.File["attachments"]
		for i, _ := range files {
			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, "open attach file stream error", http.StatusInternalServerError)
				return
			}
			//attach each file to email
			e.Attach(file, files[i].Filename, "")
		}
	}

	//smtp服务器，如果没有则使用默认配置
	if len(r.Form.Get("server")) != 0 && len(r.Form.Get("user")) != 0 && len(r.Form.Get("passwd")) != 0 {
		server = r.Form.Get("server")
		user = r.Form.Get("user")
		passwd = r.Form.Get("passwd")
	} else {
		server = g.Config().Smtp.Addr
		user = g.Config().Smtp.User
		passwd = g.Config().Smtp.Pass
	}

	//发件人，可以设置成 San Zhang <zhangsan@example.com> 这种形式(注意不能有非ascii字符)，如果不设置，默认为登陆使用的账号
	if len(r.Form.Get("from")) != 0 {
		e.From = r.Form.Get("from")
	} else {
		e.From = user
	}

	hp := strings.Split(server, ":")
	auth := smtp.PlainAuth("", user, passwd, hp[0])
	//暂时不支持TLS/StartTLS等加密认证
	error := e.Send(server, auth)
	if error != nil {
		log.Println("[ERROR]", addr, e.From, e.To, e.Subject, error)
		http.Error(w, error.Error(), http.StatusBadRequest)
		return
	}
	log.Println("[INFO]", addr, e.From, e.To, e.Subject)
	w.Write([]byte("{\"status\":0,\"msg\":\"ok\"}"))
}

//检查调用IP合法性函数
func isValid(Addr string, allowList []string) (authorized bool) {
	rAddr := net.ParseIP(Addr)
	for v := range allowList {
		_, ipNet, err := net.ParseCIDR(allowList[v])
		if err != nil {
			//log.Println("parse ip net error")
			ipHost := net.ParseIP(allowList[v])
			if ipHost != nil {
				if ipHost.Equal(rAddr) {
					authorized = true
				}
			} else {
				log.Fatalln("ip list error")
			}
		} else {
			//log.Println("Contains ip ")
			if ipNet.Contains(rAddr) {
				authorized = true
			}
		}
	}
	return authorized
}
