package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

var TG_TOKEN string = ""

type SearchResults struct {
	ready   bool
	Query   string
	Results []Result
}

type Result struct {
	Name, Description, URL string
}

func urlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func telegramBot() {
	bot, err := tgbotapi.NewBotAPI(TG_TOKEN)
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {

			switch update.Message.Text {
			case "/start":

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, i'm a wikipedia bot")
				bot.Send(msg)

			default:
				language := "ru"
				ms, _ := urlEncoded(update.Message.Text)

				url := ms
				request := "https://" + language + ".wikipedia.org/w/api.php?action=opensearch&search=" + url + "&limit=3&origin=*&format=json"
				message := wikipediaAPI(request)

				for _, val := range message {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, val)
					bot.Send(msg)
				}
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Use the words for search.")
			bot.Send(msg)
		}
	}
}

func (sr *SearchResults) UnmarshalJSON(bs []byte) error {
	array := []interface{}{}
	if err := json.Unmarshal(bs, &array); err != nil {
		return err
	}
	sr.Query = array[0].(string)
	for i := range array[1].([]interface{}) {
		sr.Results = append(sr.Results, Result{
			array[1].([]interface{})[i].(string),
			array[2].([]interface{})[i].(string),
			array[3].([]interface{})[i].(string),
		})
	}
	return nil
}

func wikipediaAPI(request string) (answer []string) {

	s := make([]string, 3)

	if response, err := http.Get(request); err != nil {
		s[0] = "Wikipedia is not respond"
	} else {
		defer response.Body.Close()

		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		sr := &SearchResults{}
		if err = json.Unmarshal([]byte(contents), sr); err != nil {
			s[0] = "Something going wrong, try to change your question"
		}

		if !sr.ready {
			s[0] = "Something going wrong, try to change your question"
		}

		for i := range sr.Results {
			s[i] = sr.Results[i].URL
		}
	}

	return s
}

func main() {
	telegramBot()
}
