package service

// type Authorization interface {
// 	CreateUser(login string, password string) (int, error)
// 	GenerateToken(login, password string) (string, error)
// 	ParseToken(token string) (int, string, error)
// }

// type UserList interface {
// 	GetUsersList() ([]models.User, error)
// 	GetUserById(id int) (models.User, error)
// 	ChangeBalance(id int, changeBalance float64) error
// 	DeleteAccount(id int, login string, password string) error
// 	LinkAccount(id int, login, password, srv string) error
// }

// type Product interface {
// 	CreateProduct(name string, cost float64, description string, amount int) (int, error)
// 	GetProductList() ([]models.Product, error)
// 	GetProductById(id int) (models.Product, error)
// 	DeleteProduct(id int) error
// 	UpdateProduct(id int, product models.UpdateProduct) error
// }

// type Cart interface {
// 	GetCart(id int) ([]models.GetCart, error)
// 	AddProductToCart(user_id int, product_id int, quantity int) error
// 	MakeOrder(user_id int) error
// 	randomDiscount(int64) (int, error)
// }

// type Order interface {
// 	GetOrderList(id int) ([]models.GetOrder, error)
// 	GetOrderDetails(userID int, orderID int) (models.GetOrder, []models.OrderItems, error)
// }

// type Service struct {
// 	Authorization
// 	UserList
// 	Product
// 	Cart
// 	Order
// }

// func NewHService(store *storage.Store) *HService{
// 	return &HService{}
// }

// func NewService(store *storage.Store, messageSender publisher.MessageSender) *Service {
// 	return &Service{
// 		Authorization: NewAuthService(store.Authorization, messageSender),
// 		Product:       NewProdService(store.Product, messageSender),
// 		UserList:      NewUserService(store.UserList, messageSender),
// 		Cart:          NewCartService(messageSender, store.Cart, store.Product, store.Order, store.UserList),
// 		Order:         NewOrderService(store.Order, messageSender),
// 	}
// }
