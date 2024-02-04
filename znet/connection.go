package znet

import (
	"fmt"
	"net"
	"zinx-demo/ziface"
)

type Connection struct {
	// 当前连接
	Conn *net.TCPConn
	// 当前连接ID
	ConnId uint32
	// 当前连接状态
	IsClosed bool
	// 当前连接的业务处理方法
	HandleApi ziface.HandleFunc
	// 告知当前连接已经退出/停止的 channel
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connId uint32, handleApi ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ConnId:    connId,
		HandleApi: handleApi,
		IsClosed:  false,
		ExitChan:  make(chan bool, 1),
	}
}

func (conn *Connection) StartReader() {
	fmt.Println("[Client]Read message, connId =", conn.ConnId)
	defer fmt.Println("[Client]Read message finished, connId =", conn.ConnId, " RemoteAddr =", conn.RemoteAddr().String())
	defer conn.Stop()

	for {
		buf := make([]byte, 512)
		read, err := conn.Conn.Read(buf)
		if err != nil {
			fmt.Println("[Client]Read message error", err)
			continue
		}

		// 调用handleApi处理数据
		if err := conn.HandleApi(conn.Conn, buf, read); err != nil {
			fmt.Println("[Client]Handle data error", err, " connId =", conn.ConnId)
			break
		}
	}
}

func (conn *Connection) Start() {
	fmt.Println("Connection start... connId =", conn.ConnId)

	// 处理数据
	go conn.StartReader()
}

func (conn *Connection) Stop() {
	fmt.Println("Connection stop... connId =", conn.ConnId)

	if conn.IsClosed {
		return
	}

	conn.IsClosed = true
	conn.Conn.Close()
	close(conn.ExitChan)
}

func (conn *Connection) GetTcpConnection() *net.TCPConn {
	return conn.Conn
}

func (conn *Connection) GetConnId() uint32 {
	return conn.ConnId
}

func (conn *Connection) RemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}

func (conn *Connection) Send(data []byte) error {
	_, err := conn.Conn.Write(data)
	return err
}