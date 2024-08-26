package publisher

type MessageSender interface {
	SendMessage(userID int, orderCost float64)
}
