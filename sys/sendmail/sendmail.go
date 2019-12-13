package sendmail

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/lj-team/go-generic/encode/pack"
	"github.com/lj-team/go-generic/log"
)

var Util string = "/usr/sbin/sendmail"
var Args string = "-t"

type Attachment struct {
	Name    string
	Content []byte
	IsFile  bool
}

type Mail struct {
	From        string
	To          []string
	Subject     string
	Message     string
	ContentType string
	Attachments []*Attachment
}

func Send(mail *Mail) {
	list := strings.Join(mail.To, ", ")

	num := 0
	boundary := make_boundary(num)

	msg := ""
	if mail.From != "" {
		msg += "From: " + mail.From + "\n"
	}

	msg += "To: " + list + "\n"
	msg += "Subject: " + mail.Subject + "\n"
	msg += "Mime-Version: 1.0\n"

	msg += "Content-Type: multipart/mixed; boundary=\"" + boundary + "\"\n"
	msg += "\n"
	msg += "--" + boundary + "\n"

	if mail.ContentType == "" {
		msg += "Content-Type: text/plain; charset=utf-8\n"
	} else {
		msg += "Content-Type: " + mail.ContentType + "\n"
	}

	msg += "\n"
	msg += mail.Message + "\n"

	re := regexp.MustCompile("[^/]+$")

	for _, attach := range mail.Attachments {
		name := attach.Name
		if attach.IsFile {
			list := re.FindAllStringSubmatch(name, 1)
			name = list[0][0]
			var err error
			if attach.Content, err = ioutil.ReadFile(attach.Name); err != nil {
				panic("read file " + attach.Name + " error")
			}
		}

		encoded := base64.StdEncoding.EncodeToString(attach.Content)

		msg += "--" + boundary + "\n"
		msg += "Content-Type: application/octet-stream; name=\"" + name + "\"\n"
		msg += "Content-Transfer-Encoding: base64\n"
		msg += "Content-Disposition: attachment; filename=\"" + name + "\"\n"
		msg += "\n"
		msg += encoded + "\n"
	}

	msg += "--" + boundary + "--\n"
	msg += ".\n"

	log.Trace("send mail:\n" + msg)

	util := exec.Command(Util, Args)
	stdin, err := util.StdinPipe()

	if err != nil {
		log.Error("sendmail error: " + err.Error())
		return
	}

	if err = util.Start(); err != nil {
		log.Error("sendmail error: " + err.Error())
		return
	}

	if _, err = stdin.Write([]byte(msg)); err != nil {
		log.Error("sendmail error: " + err.Error())
		return
	}

	stdin.Close()
	util.Wait()
}

func make_boundary(num int) string {
	return fmt.Sprintf("%x%d", md5.Sum(pack.Encode(time.Now().UnixNano())), num)
}
