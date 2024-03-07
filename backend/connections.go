// connections.go 

package main

import (
	"github.com/gorilla/websocket"
	"github.com/ethereum/go-ethereum/common"
    "sync"
    "fmt"
)

type Connection struct {
    Fighter *Fighter
    OwnerAddress common.Address
    sync.RWMutex 
}

func (i *Connection) gFighter() *Fighter {
    i.RLock()
    i.RUnlock()

    return i.Fighter
}

func (i *Connection) gOwnerAddress() common.Address {
    i.RLock()
    i.RUnlock()

    return i.OwnerAddress
}

type SafeConnectionsMap struct {
    Map map[*websocket.Conn]*Connection
    sync.RWMutex
}

var ConnectionsMap = &SafeConnectionsMap{Map: make(map[*websocket.Conn]*Connection)}


func (i *SafeConnectionsMap) gMap() map[*websocket.Conn]*Connection {
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
    i.Lock()
    defer i.Unlock()

    delete(i.Map, conn)
}

func (i *SafeConnectionsMap) Add(conn *websocket.Conn) {
    i.Lock()
    defer i.Unlock()

    i.Map[conn] = &Connection{}    
}

func (i *SafeConnectionsMap) AddWithValues(conn *websocket.Conn, fighter *Fighter, ownerAddress common.Address) {
    i.Lock()
    defer i.Unlock()

    i.Map[conn] = &Connection{
    	Fighter: fighter,
    	OwnerAddress: ownerAddress,
    }    
}

func (i *SafeConnectionsMap) SetConnectionOwnerAddress(conn *websocket.Conn, ownerAddress common.Address) {
    i.Lock()
    defer i.Unlock()

    i.Map[conn].OwnerAddress = ownerAddress
}



func GetConnection(conn *websocket.Conn) *Connection {
    ConnectionsMap.RLock()
    ConnectionsMap.RUnlock()

    connection, ok := ConnectionsMap.Map[conn]
    if !ok {
        return nil
    }

    return connection
}


func getOwnerAddressByConn(conn *websocket.Conn) (common.Address, error) {
    connection := ConnectionsMap.Find(conn)

    if connection == nil {
        return common.Address{}, fmt.Errorf("[getOwnerAddressByConn] Connection not found")
    }

    return connection.gOwnerAddress(), nil
}


func findConnectionByFighter(fighter *Fighter) (*websocket.Conn, *Connection) {
    for conn, connection := range ConnectionsMap.gMap() {
        if connection.gFighter() == fighter {
            return conn, connection
        }
    }

    return nil, &Connection{}
}
