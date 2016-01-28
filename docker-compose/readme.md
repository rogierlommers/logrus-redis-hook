# Redis Hook for [Logrus](https://github.com/Sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Running Redis and Kibana with docker-compose

By using the docker-compose.yml file, you can easily run a complete testing environment. The network settings are `net: "host"`, which means that the redis server will be listening on your localhost, default port ().

## Usage
- Start by typing `docker-compose up`.
- After a while, Redis is listening on localhost, port 6379.
- Every message sent to Redis with key `my_redis_key` will be processed by Logstash.

## Change configuration
If you want to change the key, simply edit the `logstash.conf` file, there is an entry `key`. By default, Logstash is setup to process V0 messages, if you want to change this to V1, edit the entry `codec` from `oldlogstashjson` to `json`.

## Screenshot after compose-up
<img src="http://i.imgur.com/eFu3v5X.png"/>
