package v8go_win

//#include "v8go.h"
import "C"
import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var globalInspectorIdGenerator int32 = 0
var globalInspectorMap sync.Map

type InspectorInfo struct {
	Description          string `json:"description"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	Url                  string `json:"url"`
	WebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
}
type InspectorVersion struct {
	Browser         string `json:"Browser"`
	ProtocolVersion string `json:"Protocol-Version"`
}

type InspectorServer struct {
	InspectorClientPtr      C.RawInspectorClientPtr
	InspectorId             int32
	Isolate                 *Isolate
	Context                 *Context
	CurrentClientConnection *ClientConnection
	Lock                    sync.Mutex
	JSONVersion             *InspectorVersion
	JSONLIST                []*InspectorInfo
	Server                  *http.Server
}

func (i *InspectorServer) WaitDebugger() {
	ticker := time.NewTicker(time.Millisecond * 100)
	for range ticker.C {
		if i.DebuggerConnected() {
			break
		}
	}
}

func (i *InspectorServer) DebuggerConnected() bool {
	if i.Server == nil {
		return false
	}
	return i.CurrentClientConnection != nil
}

func (i *InspectorServer) WebSocketServe(conn *websocket.Conn) {
	if i.CurrentClientConnection == nil {
		i.CurrentClientConnection = NewClientConnection(conn, i.ProcessClientCall, i.PostConnectionClosed)
		i.CurrentClientConnection.Run()
	} else {
		fmt.Printf("Debugger Already In Use \n")
	}
}

func (i *InspectorServer) ProcessClientCall(connection *ClientConnection, data []byte) {
	mes := string(data)
	//直接调用C的函数
	//fmt.Printf("Receive Message %s  \n", mes)
	C.OnReceiveMessage(i.InspectorClientPtr, C.CString(mes))
}

func (i *InspectorServer) PostConnectionClosed(connection *ClientConnection) {
	i.CurrentClientConnection = nil
}

func (i *InspectorServer) StartWebSocketServer(port uint32) {
	handler := &websocket.Server{
		Config:    websocket.Config{},
		Handshake: nil,
		Handler:   i.WebSocketServe,
	}
	address := fmt.Sprintf("127.0.0.1:%d", port)
	i.JSONVersion = &InspectorVersion{
		Browser:         "Puerts/v1.0.0",
		ProtocolVersion: "1.1",
	}
	i.JSONLIST = append(i.JSONLIST, &InspectorInfo{
		Description:          "go instance",
		Title:                "debug tools for V8",
		Type:                 "node",
		Url:                  "",
		WebSocketDebuggerUrl: fmt.Sprintf("ws://%s", address),
	})
	mux := &http.ServeMux{}
	mux.Handle("/json/version", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		s, _ := json.Marshal(i.JSONVersion)
		writer.Write(s)
	}))
	mux.Handle("/json/list", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		s, _ := json.Marshal(i.JSONLIST)
		writer.Write(s)
	}))
	mux.Handle("/json", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		s, _ := json.Marshal(i.JSONLIST)
		writer.Write(s)
	}))
	mux.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
	}))
	server := &http.Server{Addr: address, Handler: mux}
	i.Server = server
	go server.ListenAndServe()
}

func (i *InspectorServer) Destroy() {
	globalInspectorMap.Delete(i.InspectorId)
	i.Server.Shutdown(context.TODO())
}

func NewClient(iso *Isolate, ctx *Context, port uint32) *InspectorServer {
	inspectorId := atomic.AddInt32(&globalInspectorIdGenerator, 1)
	cl := &InspectorServer{
		InspectorId:        inspectorId,
		InspectorClientPtr: C.NewInspectorClient(ctx.ptr, C.int32_t(inspectorId)),
		Isolate:            iso,
		Context:            ctx,
	}
	globalInspectorMap.Store(inspectorId, cl)
	C.BindMessageSendFuncToClient(cl.InspectorClientPtr)
	C.BindTickFuncToClient(cl.InspectorClientPtr)
	cl.StartWebSocketServer(port)
	return cl
}

//export goSendMessage
func goSendMessage(inspectorId C.int32_t, message C.RawCharPtr) {
	id := int32(inspectorId)
	inspectorClientInterface, ok := globalInspectorMap.Load(id)
	if !ok {
		return
	}
	inspectorClient := inspectorClientInterface.(*InspectorServer)
	msg := C.GoString(message)
	//fmt.Printf("Send Message %s  \n", msg)
	websocket.Message.Send(inspectorClient.CurrentClientConnection.Connection, msg)
}

//export goTick
func goTick(inspectorId C.int32_t) {
	id := int32(inspectorId)
	inspectorClientInterface, ok := globalInspectorMap.Load(id)
	if !ok {
		return
	}
	inspectorClient := inspectorClientInterface.(*InspectorServer)
	//一直接受网络请求
	data2 := make([]byte, 4096*100)
	connection := inspectorClient.CurrentClientConnection
	err := websocket.Message.Receive(connection.Connection, &data2)
	if err != nil && err != io.EOF {
		panic(err)
	}
	connection.ProcessFun(connection, data2)
}
