package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"text/template"
)

var svc *sns.SNS
var messageTemplate *template.Template
var subjectTemplate *template.Template

type Incoming struct {
	Receiver    string            `json:"receiver"`
	Status      string            `json:"status"`
	Alerts      Alerts            `json:"alerts"`
	ExternalURL string            `json:"externalURL"`
	GroupLabels map[string]string `json:"groupLabels"`
}

type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Status      string            `json:"status"`
}

func (a *Alert) Name() string {
	return a.Labels["alertname"]
}

type Alerts []Alert

func getSubject(incoming Incoming) string {
	var newMessage bytes.Buffer
	err := subjectTemplate.Execute(&newMessage, incoming)
	if err == nil {
		return newMessage.String()
	}
	log.Errorf("Got error executing template: %s", err)
	log.Error("Reverting to basic subject")
	if len(incoming.Alerts) == 1 {
		return fmt.Sprintf("%s, %s, %s", incoming.Status, incoming.Alerts[0].Labels["alertname"], incoming.Alerts[0].Annotations["summary"])
	} else {
		return fmt.Sprintf("%s, %d alerts", incoming.Status, len(incoming.Alerts))
	}
}

func getMessage(incoming Incoming) string {
	message := ""
	var newMessage bytes.Buffer
	err := messageTemplate.Execute(&newMessage, incoming)
	if err == nil {
		return newMessage.String()
	}
	log.Errorf("Got error executing template: %s", err)
	log.Error("Reverting to basic message")
	for _, alert := range incoming.Alerts {
		message += fmt.Sprintf("%s, %s, %s\n", alert.Status, alert.Name(), alert.Annotations["summary"])
	}
	return message
}

func recievePromAlert(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	incoming := Incoming{}
	if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	params := &sns.PublishInput{
		Message:  aws.String(getMessage(incoming)),
		Subject:  aws.String(getSubject(incoming)),
		TopicArn: aws.String(viper.GetString("sns.topicarn")),
	}
	log.Infof("Published Alert: %v", params)
	if !viper.GetBool("fakeMessage") {
		_, err := svc.Publish(params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func GetConfiguration() {
	viper.SetConfigName("prometheus-sns-webhook")
	viper.AddConfigPath("/etc/prometheus-sns-webhook/")
	viper.AddConfigPath("$HOME/.prometheus-sns-webhook")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.SetDefault("fakeMessage", false)
	viper.SetDefault("sns.region", "eu-west-1")
	viper.SetDefault("messageTemplate", "{{ range .Alerts }}{{ .Annotations.runbook_url }}\n{{ .Annotations.message }}\n{{ end }}")
	viper.SetDefault("subjectTemplate", "Alert: {{ .Status }}{{ .GroupLabels.alertname }}")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if viper.GetBool("fakeMessage") {
		log.Warn("fakeMessage is set, no messages will be published to sns")
	}
	if viper.GetString("sns.topicarn") == "" {
		log.Warn("sns.topicarn is not set, no messages will be published to sns")
	}
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	GetConfiguration()
	var err error
	messageTemplate, err = template.New("message").Parse(viper.GetString("messageTemplate"))
	if err != nil {
		log.Error("cannot read message template: %s", err)
		panic(err)
	}
	subjectTemplate, err = template.New("subject").Parse(viper.GetString("subjectTemplate"))
	if err != nil {
		log.Error("cannot read subject template: %s", err)
		panic(err)
	}
	svc = sns.New(session.New(
		&aws.Config{
			Region: aws.String(viper.GetString("sns.region")),
		}))
	router.HandleFunc("/alert", recievePromAlert)
	return router
}
