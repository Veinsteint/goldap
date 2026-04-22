package tools

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"goldap-server/config"

	"github.com/patrickmn/go-cache"
	"gopkg.in/gomail.v2"
)

var VerificationCodeCache = cache.New(24*time.Hour, 48*time.Hour)

func email(mailTo []string, subject string, body string) error {
	mailConn := map[string]string{
		"user": config.Conf.Email.User,
		"pass": config.Conf.Email.Pass,
		"host": config.Conf.Email.Host,
		"port": config.Conf.Email.Port,
	}
	port, _ := strconv.Atoi(mailConn["port"])

	newmail := gomail.NewMessage()
	newmail.SetHeader("From", newmail.FormatAddress(mailConn["user"], config.Conf.Email.From))
	newmail.SetHeader("To", mailTo...)
	newmail.SetHeader("Subject", subject)
	newmail.SetBody("text/html", body)

	do := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	if config.Conf.Email.IsSSL || port == 465 {
		do.SSL = true
	}
	return do.DialAndSend(newmail)
}

// SendMail sends password reset success notification
func SendMail(sendto []string, pass string) error {
	subject := "LDAP Password Reset Success"
	body := fmt.Sprintf(`<div>
        <div>Dear User,</div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>Your password has been reset successfully. New password: <b>%s</b></p>
            <p>Please keep your account information secure.</p>
        </div>
        <div><p>This is an automated message. Please do not reply.</p></div>
    </div>`, pass)
	return email(sendto, subject, body)
}

// SendCode sends verification code for password reset
func SendCode(sendto []string) error {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	VerificationCodeCache.Set(sendto[0], vcode, time.Minute*5)
	subject := "Verification Code - Password Reset"
	body := fmt.Sprintf(`<div>
        <div>Dear User,</div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>Your verification code is: %s (valid for 5 minutes)</p>
            <p>Please ensure this is your own action and do not share with others.</p>
        </div>
        <div><p>This is an automated message. Please do not reply.</p></div>
    </div>`, vcode)
	return email(sendto, subject, body)
}

// SendUserCreationNotification sends account creation notification
func SendUserCreationNotification(username, nickname, mail, password string) error {
	subject := "LDAP Account Created"
	body := fmt.Sprintf(`<div>
        <div>Dear %s,</div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>Your LDAP account has been created:</p>
            <ul>
                <li>Username: %s</li>
                <li>Nickname: %s</li>
                <li>Initial Password: %s</li>
            </ul>
            <p style="color: #ff6600;">Please change your password after first login.</p>
        </div>
        <div><p>This is an automated message. Please do not reply.</p></div>
    </div>`, nickname, username, nickname, password)
	return email([]string{mail}, subject, body)
}

// SendPasswordResetNotification sends password reset notification
func SendPasswordResetNotification(username, nickname, mail, newPassword string) error {
	subject := "LDAP Password Reset Notification"
	body := fmt.Sprintf(`<div>
        <div>Dear %s,</div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>Your LDAP password has been reset:</p>
            <ul>
                <li>Username: %s</li>
                <li>New Password: %s</li>
            </ul>
            <p style="color: #ff6600;">Please login and change your password immediately.</p>
        </div>
        <div><p>This is an automated message. Please do not reply.</p></div>
    </div>`, nickname, username, newPassword)
	return email([]string{mail}, subject, body)
}

// SendUserRegistrationPendingNotification sends registration pending notification
func SendUserRegistrationPendingNotification(username, nickname, mail string) error {
	subject := "LDAP Registration Submitted"
	body := fmt.Sprintf(`<div>
        <div>Dear %s,</div>
        <div style="padding: 8px 40px 8px 50px;">
            <p>Your LDAP registration has been submitted:</p>
            <ul>
                <li>Username: %s</li>
                <li>Name: %s</li>
                <li>Email: %s</li>
            </ul>
            <p style="color: #ff6600;">Your registration is pending admin approval.</p>
        </div>
        <div><p>This is an automated message. Please do not reply.</p></div>
    </div>`, nickname, username, nickname, mail)
	return email([]string{mail}, subject, body)
}

// SendUserRegistrationApprovedNotification sends registration approval notification
func SendUserRegistrationApprovedNotification(username, nickname, mail string) error {
	subject := "LDAP Registration Approved"
	body := fmt.Sprintf(`<div>
        <div>Dear %s,</div>
        <div style="padding: 8px 40px 8px 50px;">
            <p style="color: #67C23A; font-weight: bold;">Your LDAP registration has been approved!</p>
            <ul>
                <li>Username: %s</li>
                <li>Name: %s</li>
            </ul>
            <p style="color: #ff6600;">You can now login with the password you set during registration.</p>
        </div>
        <div><p>This is an automated message. Please do not reply.</p></div>
    </div>`, nickname, username, nickname)
	return email([]string{mail}, subject, body)
}

// SendUserRegistrationRejectedNotification sends registration rejection notification
func SendUserRegistrationRejectedNotification(username, nickname, mail, reviewRemark string) error {
	subject := "LDAP Registration Rejected"
	remarkHTML := ""
	if reviewRemark != "" {
		remarkHTML = fmt.Sprintf("<p><b>Reason:</b> %s</p>", reviewRemark)
	}
	body := fmt.Sprintf(`<div>
        <div>Dear %s,</div>
        <div style="padding: 8px 40px 8px 50px;">
            <p style="color: #F56C6C; font-weight: bold;">Your LDAP registration has been rejected.</p>
            <ul>
                <li>Username: %s</li>
                <li>Name: %s</li>
            </ul>
            %s
            <p>Please contact the administrator for more information.</p>
        </div>
        <div><p>This is an automated message. Please do not reply.</p></div>
    </div>`, nickname, username, nickname, remarkHTML)
	return email([]string{mail}, subject, body)
}
