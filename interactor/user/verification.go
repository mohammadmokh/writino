package user

import (
	"bytes"
	"crypto/rand"
	"html/template"
	"math/big"
	"strconv"

	"gitlab.com/gocastsian/writino/entity"
)

func GenVerCode() (string, error) {

	var max int64 = 1000000
	var min int64 = 100000
	bg := big.NewInt(max - min)

	n, err := rand.Int(rand.Reader, bg)

	return strconv.FormatInt(n.Int64()+min, 10), err
}

func ParseVerificationTempl(user entity.User) (string, error) {

	t := template.New("email_tmpl.html")
	t, err := t.ParseFiles("./assets/verification_templ/email_tmpl.html")
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, user); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
