# Store messages to database table by structure

## hints
  1. The message type in the broker must be the json object
  2. Only use one table for store data

## Configuration
You can set config for database, broker and table structure with config files
  
## Database configs
Set parameters in "config/database.json" for make a connection to mysql server and store data 
  ```json
  {
    "host":"127.0.0.1",
    "port":"3306",
    "username":"root",
    "password":"",
    "database":"",
    "table":""
  }
  ```
  
## Broker configs
Set parameters in "config/broker.json" for make a connection to mqtt broker
  ```json
  {
    "host":"127.0.0.1",
    "port":"1883",
    "username":"root",
    "password":"",
    "client_id":"MqttToMysql",
    "topic":"topic/test"
  }
  ```
  
## Structure configs
Make a structure for convert message to a row of table, you must set key as a key in message and value as a column name of table\
for examples:
  ```json
  {
    "usr": "user",
    "msg": "message"
  }
  ```
