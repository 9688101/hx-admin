package message

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/9688101/hx-admin/core/logger"
	"github.com/9688101/hx-admin/global"
)

func shouldAuth() bool {
	return global.SMTPAccount != "" || global.SMTPToken != ""
}

func SendEmail(subject string, receiver string, content string) error {
	if receiver == "" {
		return fmt.Errorf("receiver is empty")
	}
	if global.SMTPFrom == "" { // for compatibility
		global.SMTPFrom = global.SMTPAccount
	}
	encodedSubject := fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))

	// Extract domain from SMTPFrom
	parts := strings.Split(global.SMTPFrom, "@")
	var domain string
	if len(parts) > 1 {
		domain = parts[1]
	}
	// Generate a unique Message-ID
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		return err
	}
	messageId := fmt.Sprintf("<%x@%s>", buf, domain)

	mail := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s<%s>\r\n"+
		"Subject: %s\r\n"+
		"Message-ID: %s\r\n"+ // add Message-ID header to avoid being treated as spam, RFC 5322
		"Date: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n",
		receiver, global.SystemName, global.SMTPFrom, encodedSubject, messageId, time.Now().Format(time.RFC1123Z), content))

	auth := smtp.PlainAuth("", global.SMTPAccount, global.SMTPToken, global.SMTPServer)
	addr := fmt.Sprintf("%s:%d", global.SMTPServer, global.SMTPPort)
	to := strings.Split(receiver, ";")

	if global.SMTPPort == 465 || !shouldAuth() {
		// need advanced client
		var conn net.Conn
		var err error
		if global.SMTPPort == 465 {
			tlsglobal := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         global.SMTPServer,
			}
			conn, err = tls.Dial("tcp", fmt.Sprintf("%s:%d", global.SMTPServer, global.SMTPPort), tlsglobal)
		} else {
			conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", global.SMTPServer, global.SMTPPort))
		}
		if err != nil {
			return err
		}
		client, err := smtp.NewClient(conn, global.SMTPServer)
		if err != nil {
			return err
		}
		defer client.Close()
		if shouldAuth() {
			if err = client.Auth(auth); err != nil {
				return err
			}
		}
		if err = client.Mail(global.SMTPFrom); err != nil {
			return err
		}
		receiverEmails := strings.Split(receiver, ";")
		for _, receiver := range receiverEmails {
			if err = client.Rcpt(receiver); err != nil {
				return err
			}
		}
		w, err := client.Data()
		if err != nil {
			return err
		}
		_, err = w.Write(mail)
		if err != nil {
			return err
		}
		err = w.Close()
		if err != nil {
			return err
		}
		return nil
	}
	err = smtp.SendMail(addr, auth, global.SMTPAccount, to, mail)
	if err != nil && strings.Contains(err.Error(), "short response") { // 部分提供商返回该错误，但实际上邮件已经发送成功
		logger.SysWarnf("short response from SMTP server, return nil instead of error: %s", err.Error())
		return nil
	}
	return err
}
