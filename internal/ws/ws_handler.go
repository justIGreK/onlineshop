package ws

import (
	"errors"
	"fmt"
	"net/http"
	"onlineshop/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type UserList interface {
	GetUserById(id int) (models.User, error)
}
type Handler struct {
	hub *Hub
	UserList
}

func NewHandler(h *Hub, userlist UserList) *Handler {
	return &Handler{
		hub:      h,
		UserList: userlist,
	}
}

type CreateRoomReq struct {
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, _, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	roomId := fmt.Sprintf("%v_%v", req.Name, userId)
	h.hub.Rooms[roomId] = &Room{
		ID:      roomId,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}
	fmt.Println(roomId)
	user := strconv.Itoa(userId)
	err = h.hub.RedisDB.RPush(user, req.Name).Err()
	if err != nil {
		fmt.Println("Error during getting data from redis:", err)
		return
	}
	err = h.hub.RedisDB.RPush("Chats", roomId).Err()
	if err != nil {
		fmt.Println("Error during getting data from redis:", err)
		return
	}
	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	AdministratorRole = "admin"
	CustomerRole      = "customer"
)

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, userRole, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	exists, err := h.hub.RedisDB.Exists("Chats", c.Query("room")).Result()
	if err != nil {
		fmt.Println("Error:", err)
	} else if exists == 0{
		c.JSON(http.StatusBadRequest, "No such room")
		return 
	}
	var roomID string
	if userRole == AdministratorRole {
		roomID = c.Query("room")
	} else if userRole == CustomerRole {
		roomID = fmt.Sprintf("%v_%v", c.Query("room"), userId)
	} else {
		c.JSON(http.StatusUnauthorized, "Something wrong with your role")
		return
	}
	clientID := string(user.Id)
	username := user.Login

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
	history := h.hub.RedisDB.LRange(roomID, 0, -1)
	c.JSON(http.StatusOK, history)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	var chats []string
	userId, userRole, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userRole == AdministratorRole {
		chats, err = h.hub.RedisDB.LRange("Chats", 0, -1).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else if userRole == CustomerRole {
		roomId := strconv.Itoa(userId)
		chats, err = h.hub.RedisDB.LRange(roomId, 0, -1).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, "Something wrong with your role")
		return
	}
	c.JSON(http.StatusOK, chats)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *gin.Context) {
	var clients []ClientRes
	roomId := c.Param("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}

func getUserInfo(c *gin.Context) (int, string, error) {
	role, ok := c.Get("userRole")
	if !ok {
		return 0, "", errors.New("user role not found")
	}
	userRole, ok := role.(string)
	if !ok {
		return 0, "", errors.New("user role is of invalid type")
	}
	id, ok := c.Get("userId")
	if !ok {
		return 0, "", errors.New("user id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		return 0, "", errors.New("user id is of invalid type")
	}
	return idInt, userRole, nil
}
func (h *Handler) GetHistory(c *gin.Context) {
	userId, userRole, err := getUserInfo(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userRole == AdministratorRole {
		roomId := c.Param("room")
		fmt.Println(roomId)
		messages, err := h.hub.RedisDB.LRange(roomId, 0, -1).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(messages)
		c.JSON(http.StatusOK, messages)
	}
	if userRole == CustomerRole {
		room := c.Query("room")
		fmt.Println(room)
		roomId := fmt.Sprintf("%v_%v", room, userId)
		fmt.Println(roomId)
		messages, err := h.hub.RedisDB.LRange(roomId, 0, -1).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(messages)
		c.JSON(http.StatusOK, messages)
	}

}

