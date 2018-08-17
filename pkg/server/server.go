package server

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/model"
	"github.com/spf13/viper"
	"net/http"
)

var svc *sns.SNS

type Incoming struct {
	Receiver string       `json:"receiver"`
	Status   string       `json:"status"`
	Alerts   model.Alerts `json:"alerts"`

	ExternalURL string `json:"externalURL"`
}

func getSubject(incoming Incoming) string {
	if len(incoming.Alerts) == 1 {
		return fmt.Sprintf("%s, %s, %s", incoming.Status, incoming.Alerts[0].Labels["alertname"], incoming.Alerts[0].Annotations["summary"])
	} else {
		return fmt.Sprintf("%s, %d alerts", incoming.Status, len(incoming.Alerts))
	}
}

func getMessage(incoming Incoming) string {
	message := ""
	for _, alert := range incoming.Alerts {
		message += fmt.Sprintf("%s, %s, %s\n", alert.Status(), alert.Name(), alert.Annotations["summary"])
	}
	return message
}

func recievePromAlert(w http.ResponseWriter, r *http.Request) {
	fmt.Println("called")
	defer r.Body.Close()
	incoming := Incoming{}
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	params := &sns.PublishInput{
		Message:  aws.String(getMessage(incoming)),
		Subject:  aws.String(getSubject(incoming)),
		TopicArn: aws.String(viper.GetString("sns.topicarn")),
	}
	fmt.Println(params)
	if !viper.GetBool("fakeMessage") {
		_, err := svc.Publish(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println("done")
}

func GetConfiguration() {
	viper.SetConfigName("prometheus-sns-webhook")
	viper.AddConfigPath("/etc/prometheus-sns-webhook/")
	viper.AddConfigPath("$HOME/.prometheus-sns-webhook")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.SetDefault("fakeMessage", false)
	viper.SetDefault("sns.region", "eu-west-1")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if viper.GetBool("fakeMessage") {
		fmt.Println("WARNING: fakeMessage is set, no messages will be published to sns")
	}
	if viper.GetString("sns.topicarn") == "" {
		fmt.Println("WARNING: sns.topicarn is not set, no messages will be published to sns")
	}
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	GetConfiguration()
	svc = sns.New(session.New(
		&aws.Config{
			Region: aws.String(viper.GetString("sns.region")),
		}))
	router.HandleFunc("/alert", recievePromAlert)
	return router
}
