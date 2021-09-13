package gate

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/log"
	"Open_IM/src/utils"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WServer struct {
	wsAddr       string
	wsMaxConnNum int
	wsUpGrader   *websocket.Upgrader
	wsConnToUser map[*websocket.Conn]string
	wsUserToConn map[string]*websocket.Conn
}

func (ws *WServer) onInit(wsPort int) {
	ip := utils.ServerIP
	ws.wsAddr = ip + ":" + utils.IntToString(wsPort)
	ws.wsMaxConnNum = config.Config.LongConnSvr.WebsocketMaxConnNum
	ws.wsConnToUser = make(map[*websocket.Conn]string)
	ws.wsUserToConn = make(map[string]*websocket.Conn)
	ws.wsUpGrader = &websocket.Upgrader{
		HandshakeTimeout: time.Duration(config.Config.LongConnSvr.WebsocketTimeOut) * time.Second,
		ReadBufferSize:   config.Config.LongConnSvr.WebsocketMaxMsgLen,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
}

func (ws *WServer) run() {
	http.HandleFunc("/", ws.wsHandler)         //Get request from client to handle by wsHandler
	err := http.ListenAndServe(ws.wsAddr, nil) //Start listening
	if err != nil {
		log.ErrorByKv("Ws listening err", "", "err", err.Error())
	}
}

func (ws *WServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	if ws.headerCheck(w, r) {
		query := r.URL.Query()
		conn, err := ws.wsUpGrader.Upgrade(w, r, nil) //Conn is obtained through the upgraded escalator
		if err != nil {
			log.ErrorByKv("upgrade http conn err", "", "err", err)
			return
		} else {
			//Connection mapping relationship,
			//userID+" "+platformID->conn
			SendID := query["sendID"][0] + " " + utils.PlatformIDToName(int32(utils.StringToInt64(query["platformID"][0])))
			ws.addUserConn(SendID, conn)
			go ws.readMsg(conn)
		}
	}
}

func (ws *WServer) readMsg(conn *websocket.Conn) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.ErrorByKv("WS ReadMsg error", "", "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn), "error", err)
			ws.delUserConn(conn)
			return
		} else {
			log.ErrorByKv("test", "", "msgType", msgType, "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn))
		}
		ws.msgParse(conn, msg)
		//ws.writeMsg(conn, 1, chat)
	}

}

func (ws *WServer) writeMsg(conn *websocket.Conn, a int, msg []byte) error {
	rwLock.Lock()
	defer rwLock.Unlock()
	return conn.WriteMessage(a, msg)

}
func (ws *WServer) addUserConn(uid string, conn *websocket.Conn) {
	rwLock.Lock()
	defer rwLock.Unlock()
	if oldConn, ok := ws.wsUserToConn[uid]; ok {
		err := oldConn.Close()
		delete(ws.wsConnToUser, oldConn)
		if err != nil {
			log.ErrorByKv("close err", "", "uid", uid, "conn", conn)
		}
	} else {
		log.InfoByKv("this user is first login", "", "uid", uid)
	}
	ws.wsConnToUser[conn] = uid
	ws.wsUserToConn[uid] = conn
	log.WarnByKv("WS Add operation", "", "wsUser added", ws.wsUserToConn, "uid", uid, "online_num", len(ws.wsUserToConn))

}

func (ws *WServer) delUserConn(conn *websocket.Conn) {
	rwLock.Lock()
	defer rwLock.Unlock()
	var uidPlatform string
	if uid, ok := ws.wsConnToUser[conn]; ok {
		uidPlatform = uid
		if _, ok = ws.wsUserToConn[uid]; ok {
			delete(ws.wsUserToConn, uid)
			log.WarnByKv("WS delete operation", "", "wsUser deleted", ws.wsUserToConn, "uid", uid, "online_num", len(ws.wsUserToConn))
		} else {
			log.WarnByKv("uid not exist", "", "wsUser deleted", ws.wsUserToConn, "uid", uid, "online_num", len(ws.wsUserToConn))
		}
		delete(ws.wsConnToUser, conn)
	}
	err := conn.Close()
	if err != nil {
		log.ErrorByKv("close err", "", "uid", uidPlatform, "conn", conn)
	}

}

func (ws *WServer) getUserConn(uid string) *websocket.Conn {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if conn, ok := ws.wsUserToConn[uid]; ok {
		return conn
	}
	return nil
}
func (ws *WServer) getUserUid(conn *websocket.Conn) string {
	rwLock.RLock()
	defer rwLock.RUnlock()

	if uid, ok := ws.wsConnToUser[conn]; ok {
		return uid
	}
	return ""
}
func (ws *WServer) headerCheck(w http.ResponseWriter, r *http.Request) bool {
	status := http.StatusUnauthorized
	query := r.URL.Query()
	if len(query["token"]) != 0 && len(query["sendID"]) != 0 && len(query["platformID"]) != 0 {
		if !utils.VerifyToken(query["token"][0], query["sendID"][0]) {
			log.ErrorByKv("Token verify failed", "", "query", query)
			w.Header().Set("Sec-Websocket-Version", "13")
			http.Error(w, http.StatusText(status), status)
			return false
		} else {
			log.InfoByKv("Connection Authentication Success", "", "token", query["token"][0], "userID", query["sendID"][0])
			return true
		}
	} else {
		log.ErrorByKv("Args err", "", "query", query)
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(status), status)
		return false
	}

}
