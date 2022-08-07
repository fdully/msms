package smsgw

import (
	"testing"
)

func TestModifyMessage(t *testing.T) {
	t.Parallel()

	msg := `Сведения о вашей учетной записи:
	Имя пользователя: 197644
	Пароль: 1`

	if modifyMessage(msg) != "197644" {
		t.Fail()
	}

	// changing value
	msg = `<pre class="code-java" style="margin: 0px; padding: 0px; max-height: 30em; overflow: auto; white-space: pre-wrap; word-wrap: normal; color: rgb(51, 51, 51); font-size: 12px; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: normal; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; background-color: rgb(245, 245, 245); text-decoration-style: initial; text-decoration-color: initial;">Your guest account details: Username: 351182 Password: 1</pre>`

	if modifyMessage(msg) != "351182" {
		t.Fail()
	}

	// changing ones again
	msg = `got nothing but digits 456765`
	if modifyMessage(msg) != "" {
		t.Fail()
	}

}
