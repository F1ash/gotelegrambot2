// -*- Go -*-
/* ------------------------------------------------ */
/* Golang source                                    */
/* Author: Алексей Панов <a.panov@maximatelecom.ru> */
/* ------------------------------------------------ */

package main

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/telegram-bot-api.v4"
)

func commandsMainHandler(msg *tgbotapi.Message) {
	cmd := msg.Command()
	args := msg.CommandArguments()
	log.Debugf("Command from %s: `%s %s`", msg.From.String(), cmd, args)
	switch strings.ToLower(cmd) {
	case "start":
		go commandsStartHandler(msg)
	case "ban":
		go commandsBanHandler(msg)
	case "dnf", "yum":
		go commandsDNFHandler(msg)
	default:

	}
}

func commandsStartHandler(msg *tgbotapi.Message) {
	t := fmt.Sprintf("Привет %s!", msg.From.String())
	sendMessage(msg.Chat.ID, t, msg.MessageID)
	log.Debugf("Say hello to %s", msg.From.String())
}

func commandsBanHandler(msg *tgbotapi.Message) {
	if !msg.Chat.IsGroup() && !msg.Chat.IsSuperGroup() {
		sendMessage(msg.Chat.ID, "Кого будем банить в привате? 😂", msg.MessageID)
		log.Debugf("Command `ban` in private chat from %s", msg.From.String())
		return
	}

	if !isMeAdmin(msg.Chat) {
		sendMessage(msg.Chat.ID, "Бот не является администратором этого чата. Команда недоступна!", msg.MessageID)
		log.Warn("Command `ban` in chat with bot not admin from %s", msg.From.String())
		return
	} else {
		log.Debugf("Command `ban` in group or supergroup chat with bot admin from %s", msg.From.String())
	}

	if msg.CommandArguments() == "" {
		sendMessage(msg.Chat.ID, "Кого будем банить?", msg.MessageID)
		log.Debugf("Command `ban` without arguments from %s", msg.From.String())
		return
	}
}

func commandsDNFHandler(msg *tgbotapi.Message) {
	var err error
	args := strings.Replace(msg.CommandArguments(), "—", "--", -1)
	if args == "" {
		sendMessage(msg.Chat.ID, "Не знаю, что выполнять, ты же ничего не указал в аргументах", msg.MessageID)
		log.Debugf("Command `dnf` without arguments from %s", msg.From.String())
		return
	}

	arglist := strings.Split(args, " ")
	if arglist[0] == "info" || arglist[0] == "provides" || arglist[0] == "repolist" || arglist[0] == "repoquery" {
		if arglist[0] != "repolist" { //&& arglist[0] != "repoquery" {
			arglist = append(arglist, "-q")
		}
		cmd := exec.Command("/usr/bin/dnf", arglist...)
		if err = cmd.Start(); err != nil {
			log.Errorf("Unable to start command form %s: dnf %s: %s", msg.From.String(), strings.Join(arglist, " "), strings.Join(arglist, " "))
			sendMessage(msg.Chat.ID, "Ой. Что-то пошло не так!", msg.MessageID)
			return
		}
		if err = cmd.Wait(); err != nil {
			log.Errorf("Unable to wait command form %s: dnf %s: %s", msg.From.String(), strings.Join(arglist, " "), strings.Join(arglist, " "))
			sendMessage(msg.Chat.ID, "Ой. Что-то пошло не так!", msg.MessageID)
			return
		}

		f, _ := cmd.StdoutPipe()
		defer f.Close()

		if output, err := cmd.CombinedOutput(); err != nil {
		} else if len(output) == 0 {
			sendMessage(msg.Chat.ID, "А нечего выводить, вывод пустой", msg.MessageID)
			log.Warnf("Run command from %s: dnf %s with empty output", msg.From.String(), strings.Join(arglist, " "))
		} else {
			sendMessage(msg.Chat.ID, fmt.Sprintf("``` %s ```", output), msg.MessageID)
			log.Debugf("Run command from %s: dnf %s", msg.From.String(), strings.Join(arglist, " "))
		}
	}
}
