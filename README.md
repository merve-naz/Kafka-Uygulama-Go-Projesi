## BB Certificate Service Microservice  
This is the main core service for bb-Certficate Service. It is a microservice that provides certificateservice and manage features for BB platform.  


### Dev

- Run this command to golang swag with no env

```bash
# generate swag docs first! (MACOS)
$HOME/go/bin/swag init

# generate swag docs first! (LINUX)
go install github.com/swaggo/swag/cmd/swag
swag init

# run the server
go run .
```


// Prometheus metriclerinin event emitter içerisine taşınması
// kafka-trigger  uç noktasının kullanıcıdan parametreyi alacağı yapının kurulması ve test edilmesi.
// Kafka dead letter queue araştırılması