package utils

import (
	"github.com/astaxie/beego"
	"strings"
	"YYCMS/utils/YYLog"
	"net/smtp"

	"crypto/tls"
	"net"
)

type Email struct {
	To          string
	Subject     string
	Body        string
	ContentType EmailContentType
}

type EmailContentType string

const (
	EmailContentType_Html      EmailContentType = "Content-Type: text/html; charset=UTF-8"
	EmailContentType_TextPlain EmailContentType = "Content-Type: text/plain; charset=UTF-8"
)

type EmailServer struct {
	HostAndPort string
	ServerName  string
	ServerPwd   string
}

//type EmailHost string
//const (
//	EmailHost_QQ EmailHost = "smtp.exmail.qq.com:25"
//	EmailHost_Custon EmailHost = ""
//)

func EmailDefaltServer() EmailServer {
	user := beego.AppConfig.String("emailuser")
	host := beego.AppConfig.String("emailhost")
	pwd := beego.AppConfig.String("emailpwd")
	return EmailServer{
		HostAndPort:host,
		ServerName: user,
		ServerPwd:  pwd,
	}

}

func (es *EmailServer) SendEmailBySMTPUsingTLS(email *Email) error {
	hostport := strings.Split(es.HostAndPort, ":")
	msg := []byte("To: " + email.To + "\r\nFrom: " + es.ServerName + "\r\nSubject: " + email.Subject + "\r\n" + string(email.ContentType) + "\r\n\r\n" + email.Body)
	send_to := strings.Split(email.To, ";")
	auth := smtp.PlainAuth("", es.ServerName, es.ServerPwd, hostport[0])
	return SendMailUsingTLS(
		es.HostAndPort,
		auth,
		es.ServerName,
		send_to,
		msg,
	)
}

func (es *EmailServer) SendEmailBySMTP(email *Email) error {
	hostport := strings.Split(es.HostAndPort, ":")
	msg := []byte("To: " + email.To + "\r\nFrom: " + es.ServerName + "\r\nSubject: " + email.Subject + "\r\n" + string(email.ContentType) + "\r\n\r\n" + email.Body)
	send_to := strings.Split(email.To, ";")
	auth := smtp.PlainAuth("", es.ServerName, es.ServerPwd, hostport[0])
	err := smtp.SendMail(es.HostAndPort, auth, es.ServerName, send_to, msg)
	return err
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
to []string, msg []byte) (err error) {

	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		YYLog.Error("Create smpt client error:", err)
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				YYLog.Error("Error during AUTH", err)
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()

}

//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		YYLog.Error("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
