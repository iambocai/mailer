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
		addr := strings.Split(r.RemoteAddr, ":")[0]
		allowList := strings.Split(g.Config().Http.WhiteList, ",")
		authorized := isValid(addr, allowList)

		if authorized == false {
			http.Error(w, "{\"status\":403,\"msg\":\"remote not in whitelist\"}", http.StatusBadRequest)
			return
		}

		if r.ContentLength == 0 {
			http.Error(w, "{\"status\":404,\"msg\":\"body is blank\"}", http.StatusBadRequest)
			return
		}

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, "{\"status\":400,\"msg\":\"param prase error\"}", http.StatusInternalServerError)
			return
		}

		if len(r.Form["tos"]) == 0 || len(r.Form["content"]) == 0 || len(r.Form["subject"]) == 0 {
			http.Error(w, "{\"status\":400,\"msg\":\"param tos,content or subject lost\"}", http.StatusBadRequest)
			return
		}

		e := email.NewEmail()
		e.From = g.Config().Smtp.User
		e.To = strings.Split(r.Form.Get("tos"), g.Config().Smtp.Spliter)
		e.Subject = r.Form.Get("subject")
		//消息体默认为text，如果format设置为html，则以html方式发送
		if len(r.Form.Get("format")) != 0 && r.Form.Get("format") == "html" {
			e.HTML = []byte(r.Form.Get("content"))

		} else {
			e.Text = []byte(r.Form.Get("content"))
		}

		if len(r.Form.Get("hasAttach")) != 0 && r.Form.Get("hasAttach") == "1" {
			m := r.MultipartForm

			//get the *fileheaders
			files := m.File["attachments"]
			for i, _ := range files {
				//for each fileheader, get a handle to the actual file
				file, err := files[i].Open()
				defer file.Close()
				if err != nil {
					http.Error(w, "{\"status\":503,\"msg\":\"open attach file stream error\"}", http.StatusInternalServerError)
					return
				}
				//attach each file to email
				e.Attach(file, files[i].Filename, "")
			}
		}

		hp := strings.Split(g.Config().Smtp.Addr, ":")
		auth := smtp.PlainAuth("", g.Config().Smtp.User, g.Config().Smtp.Pass, hp[0])
		error := e.Send(g.Config().Smtp.Addr, auth)
		if error != nil {
			log.Println(error)
			http.Error(w, "{\"status\":500,\"msg\":\"send mail error\"}", http.StatusBadRequest)
			return
		}

		w.Write([]byte("{\"status\":0,\"msg\":\"ok\"}"))
	})
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
				log.Println("ip list error")
			}
		} else {
			//log.Println("Contains ip ")
			if ipNet.Contains(rAddr) {
				authorized = true
			}
		}
	}
	log.Println("client", Addr, "auth", authorized)
	return authorized
}
