package main

import "sync"



type NotificationManager struct {
	clients map[string]map[chan string]bool //! map -> key string , val of type --> map[whose key is string] bool
	//* channels graphy *
	// $ reading from channel into  x <- chan
	//# directing output to channel Y <- "str"
	goRoutinelockSafetyMutex sync.RWMutex //* prevent multi go routines to crash by locking current r/w to be finished first before executing another go routine
}

// @interface for notification manager type
type NotificationManagerIface interface {
	AddClient(key string,client chan string)
	RemoveClient(key string,client chan string)
	Notify(key,message string)
}

// type that stores interface which has all the methods that belongs to the type Notification
type NotificationManagerStore struct {
	NotiManagerIface NotificationManagerIface
} 

// func that returns type that stores interface of type *NotifictionStore
func NewNotificationManagerStore(notificationManager *NotificationManager) NotificationManagerStore {
	return NotificationManagerStore{
		NotiManagerIface: notificationManager ,
	}
}

// func that returns the instance of type 'this'
func NewNotiFManager() *NotificationManager {
	return &NotificationManager{
		clients: make(map[string]map[chan string]bool) ,
	}
}

// add client
func(nm *NotificationManager) AddClient(key string,client chan string) {
	// locking first fired goR with mutex
	nm.goRoutinelockSafetyMutex.Lock()
	defer nm.goRoutinelockSafetyMutex.Unlock() //! deferred it to unlock the goR at the end when sorrounding things are invoked

	// checking if already any client is listening
	if nm.clients[key] == nil { // checking if it has map or not
		nm.clients[key] = make(map[chan string]bool) // setting map 
	}
	nm.clients[key][client] = true //! setting it to true --> nm.clients[key][innermap key ] = setting value
}

// remove client
func(nm *NotificationManager) RemoveClient(key string,client chan string) {
	// locking process
	nm.goRoutinelockSafetyMutex.Lock()
	defer nm.goRoutinelockSafetyMutex.Unlock()

	clients := nm.clients[key]
	if clients != nil {
		delete(clients,client)
		if len(clients) == 0 {
			delete(nm.clients,key)
		}
		close(client)
	}
}

// notifies the client if there is any --> nm.client[thisKey]{...chan-> val as bool}
func(nm *NotificationManager) Notify(key,message string) {
	nm.goRoutinelockSafetyMutex.RLock()
	defer nm.goRoutinelockSafetyMutex.RUnlock()

	for client := range nm.clients[key]{
		select {
			// one of block executes safely as we read in chan
			case client <- message: // * if there exists a chan and we are sending this output to that chan
			default:
		}
	}
}