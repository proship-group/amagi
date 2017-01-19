package mailerBackend

import (
	"fmt"
	"net/smtp"
	"time"

	utils "github.com/b-eee/amagi"
)

// PostfixSendEmail send email to postfix
func PostfixSendEmail() error {
	s := time.Now()
	if err := smtp.SendMail("104.198.115.53:4425", nil, "j.soliva@b-eee.com", []string{"jeanepaul@gmail.com"}, []byte("testing!")); err != nil {
		utils.Error(fmt.Sprintf("error PostfixSendEmail %v", err))
		return err
	}

	// c, err := smtp.Dial("104.198.115.53:4425")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// if err := c.Hello("localhost"); err != nil {
	// 	return err
	// }

	utils.Info(fmt.Sprintf("PostfixSendEmail took: %v", time.Since(s)))
	return nil
}
