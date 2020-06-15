package debug

import (
	"fmt"

	"github.com/dalonghahaha/avenger/components/mail"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"Asgard/constants"
	"Asgard/providers"
)

var mailCmd = &cobra.Command{
	Use:   "mail",
	Short: "debug send mail",
	PreRun: func(cmd *cobra.Command, args []string) {
		confPath := cmd.Flag("conf").Value.String()
		viper.SetConfigName("app")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(confPath)
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		err = mail.Register()
		if err != nil {
			panic(err)
		}
		mailUser := viper.GetString("component.mail." + constants.MAIL_NAME + ".user")
		if mailUser == "" {
			panic(fmt.Errorf("mail user can not be empty!"))
		}
		constants.MAIL_USER = mailUser
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("debug send mail")
		receiver := cmd.Flag("receiver").Value.String()
		if receiver == "" {
			fmt.Printf("receiver can not be empty!")
			return
		}
		subject := "Asgard Notice"
		body := "Asgard Message"
		err := providers.NoticeService.SendMail(receiver, subject, body)
		if err != nil {
			fmt.Printf("send mail failed:%+v\n", err)
			return
		}
		fmt.Println("send mail success!")
	},
}
