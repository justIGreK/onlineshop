package publisher

import (
	"encoding/json"
	"onlineshop/internal/models"
	"onlineshop/internal/storage"
	"onlineshop/pkg/util/logger"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type NATSMessageSender struct {
	nc   *nats.Conn
	conn GetConnections
}
type GetConnections interface {
	GetConnections(userID int) ([]storage.Connection, error)
}

func NewNATSMessageSender(nc *nats.Conn, conn *storage.UsersPostgres) *NATSMessageSender {
	return &NATSMessageSender{nc: nc, conn: conn}
}

func (sender *NATSMessageSender) SendMessage(userID int, orderCost float64) {
	conns, err := sender.conn.GetConnections(userID)
	if err != nil {
		logger.Logger.Info("error during getting connection: %v", zap.String("error", err.Error()))
		return
	}
	for _, conn := range conns {
		pl := &models.Payload{
			ServiceID: conn.ServiceID,
			Price:     orderCost,
		}
		data, _ := json.Marshal(pl)

		err := sender.nc.Publish(conn.ServiceName, data)
		if err != nil {
			logger.Logger.Info("error sending request: %v", zap.String("error", err.Error()))
		}
	}
	logger.Logger.Info("messages are sent")
}
