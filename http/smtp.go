package http

import (
	"log"
	"mailer/g"
	"net"
	"net/http"
	"net/smtp"
	"strings"
)

func configSmtpRoutes() {
	http.HandleFunc("/api/mail", func(w http.ResponseWriter, req *http.Request) {
		addr := strings.Split(req.RemoteAddr, ":")[0]
		allowList := strings.Split(g.Config().Http.WhiteList, ",")
		authorized := isValid(addr, allowList)

		if authorized == false {
			http.Error(w, "remote not in whitelist", http.StatusBadRequest)
			return
		}

		if req.ContentLength == 0 {
			http.Error(w, "body is blank", http.StatusBadRequest)
			return
		}

		req.ParseForm()
		if len(req.Form["tos"]) == 0 || len(req.Form["content"]) == 0 || len(req.Form["subject"]) == 0 {
			http.Error(w, "connot decode body", http.StatusBadRequest)
			return
		}

		err := SendMailBySmtp(req.Form.Get("tos"), req.Form.Get("subject"), req.Form.Get("content"))
		if err != nil {
			log.Println(err)
			http.Error(w, "send mail error", http.StatusBadRequest)
			return
		}

		w.Write([]byte("success"))
	})
}

//下面两个函数放g里老是调用出错，先放这里吧

// 发送邮件函数
func SendMailBySmtp(to string, subject string, content string) error {
	var cfg = g.Config()

	hp := strings.Split(cfg.Smtp.Addr, ":")
	auth := smtp.PlainAuth("", cfg.Smtp.User, cfg.Smtp.Pass, hp[0])

	var content_type string = "Content-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + to + "\r\nFrom: " + cfg.Smtp.User + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + content)
	send_to := strings.Split(to, cfg.Smtp.Spliter)

	err := smtp.SendMail(cfg.Smtp.Addr, auth, cfg.Smtp.User, send_to, msg)
	return err
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
