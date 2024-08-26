package publisher

// import (
// 	"encoding/json"
// 	"onlineshop/internal/models"
// 	"onlineshop/pkg/util/logger"

// 	"github.com/nats-io/nats.go"
// 	"go.uber.org/zap"
// )

// type KafkaMessageSender struct {
// 	nc *nats.Conn
// }

// func NewKafkaMessageSender(nc *nats.Conn) *KafkaMessageSender {
// 	return &KafkaMessageSender{nc: nc}
// }

// func (sender *KafkaMessageSender) SendMessage(serviceID int, orderCost float64) {
// 	pl := &models.Payload{
// 		ServiceID: serviceID,
// 		Price:     orderCost,
// 	}
// 	data, _ := json.Marshal(pl)

// 	err := sender.nc.Publish("joker", data)
// 	if err != nil {
// 		logger.Logger.Info("error sending request: %v", zap.String("error", err.Error()))
// 	}
// 	logger.Logger.Info("messages are sent")
// }
