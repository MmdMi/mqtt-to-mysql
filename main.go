package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
)

var broker_conf map[string]interface{}
var database_conf map[string]interface{}
var structure_conf map[string]interface{}
var database *sql.DB
var mqtt_client mqtt.Client

func main() {
	forever := make(chan bool)
	config_init()
	mqtt_init(broker_conf)
	database_init(database_conf)
	<-forever
}

func config_init() {
	fmt.Print("Read config...")
	byteValue, err := os.ReadFile("config/broker.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(byteValue, &broker_conf)
	//fmt.Println(broker_conf)

	byteValue, err = os.ReadFile("config/database.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(byteValue, &database_conf)
	//fmt.Println(database_conf)

	byteValue, err = os.ReadFile("config/structure.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(byteValue, &structure_conf)
	//fmt.Println(structure_conf)
	fmt.Println("OK")
}

func mqtt_init(conf map[string]interface{}) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", conf["host"].(string), conf["port"].(string)))
	opts.SetClientID(conf["client_id"].(string))
	opts.SetUsername(conf["username"].(string))
	opts.SetPassword(conf["password"].(string))
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	mqtt_client = mqtt.NewClient(opts)
	if token := mqtt_client.Connect(); token.Error() != nil {
		panic(token.Error())
	}
}

func mqtt_sub(client mqtt.Client, topic string) {
	_ = client.Subscribe(topic, 1, nil)
	fmt.Printf("Start subscribe [%s]\n", topic)
}

func database_init(conf map[string]interface{}) {
	var err error
	database, err = sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			conf["username"].(string),
			conf["password"].(string),
			conf["host"].(string),
			conf["port"].(string),
			conf["database"].(string),
		))

	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Database connection...OK")
}

func map_recv_data_to_database(byteValue []byte) {
	existKey := ""
	existValue := ""
	c := make(map[string]interface{})
	_ = json.Unmarshal(byteValue, &c)

	keys := reflect.ValueOf(c)

	for _, key := range keys.MapKeys() {
		if structure_conf[key.String()] != nil {
			if existKey != "" {
				existKey += ","
			}

			if existValue != "" {
				existValue += ","
			}

			existKey += structure_conf[key.String()].(string)
			existValue += "'" + c[key.String()].(string) + "'"

		}
	}
	// fmt.Println("exist key   : " + existKey)
	// fmt.Println("exist value : " + existValue)
	query := "INSERT INTO " + database_conf["table"].(string) + " (" + existKey + ") VALUES(" + existValue + ")"
	_, err := database.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s\n", msg.Payload())
	go map_recv_data_to_database(msg.Payload())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("MQTT connection...OK")
	mqtt_sub(mqtt_client, broker_conf["topic"].(string))
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("MQTT connection lost: %v", err)
}
