// connections.go 

package main

import (
	"github.com/gorilla/websocket"
    "sync"

    "github.com/mriusd/game-contracts/account"
    "github.com/mriusd/game-contracts/fighters"
)

type Connection struct {
    AccountId int
    Session *account.Session
    Fighter *fighters.Fighter
    WSConn *websocket.Conn
    sync.RWMutex 
}

func (i *Connection) GetFighter() *fighters.Fighter {
    i.RLock()
    defer i.RUnlock()

    return i.Fighter
}


func (i *Connection) GetAccountID() int {
    i.RLock()
    defer i.RUnlock()

    return i.AccountId
}

func (i *Connection) GetSession() *account.Session {
    i.RLock()
    defer i.RUnlock()

    return i.Session
}

type SafeConnectionsMap struct {
    Map map[*websocket.Conn]*Connection
    sync.RWMutex
}

var ConnectionsMap = &SafeConnectionsMap{Map: make(map[*websocket.Conn]*Connection)}


func (i *SafeConnectionsMap) GetMap() map[*websocket.Conn]*Connection {
    i.RLock()
    defer i.RUnlock()

    copy := make(map[*websocket.Conn]*Connection, len(i.Map))
    for key, val := range i.Map {
        copy[key] = val
    }
    return copy
}

func (i *SafeConnectionsMap) Find(conn *websocket.Conn) *Connection {
    i.RLock()
    defer i.RUnlock()

    connection, ok := i.Map[conn]
    if !ok {
        return nil
    }

    return connection
}



func (i *SafeConnectionsMap) Remove(conn *websocket.Conn) {
    connection := i.Find(conn)

    if connection != nil && connection.Fighter != nil {
        unauthFighter(connection.Fighter)
    }    

    i.Lock()
    delete(i.Map, conn)
    i.Unlock()
}

func (i *SafeConnectionsMap) Add(conn *websocket.Conn, accountId int, session *account.Session) *Connection {
    i.Lock()
    defer i.Unlock()

    i.Map[conn] = &Connection{
        AccountId: accountId,
        Session: session,
        WSConn: conn,
    }   

    return i.Map[conn] 
}

func (i *SafeConnectionsMap) AddWithValues(c *websocket.Conn, fighter *fighters.Fighter) *Connection {
    i.Lock()
    defer i.Unlock()


    newConn := &Connection{
        Fighter: fighter,
        WSConn: c,
    }   

    i.Map[c] = newConn

    return newConn 
}


func GetConnection(conn *websocket.Conn) *Connection {
    ConnectionsMap.RLock()
    defer ConnectionsMap.RUnlock()

    connection, ok := ConnectionsMap.Map[conn]
    if !ok {
        return nil
    }

    return connection
}


func findConnectionByFighter(fighter *fighters.Fighter) (*websocket.Conn, *Connection) {
    ConnectionsMap.RLock() // Only lock once at the start
    defer ConnectionsMap.RUnlock()

    for conn, connection := range ConnectionsMap.Map {
        // Directly access the Fighter field since we're just comparing pointers
        // Ensure this does not introduce race conditions in your specific case
        if connection.Fighter == fighter {
            return conn, connection
        }
    }

    return nil, nil // Return nil instead of an empty Connection object when not found
}
