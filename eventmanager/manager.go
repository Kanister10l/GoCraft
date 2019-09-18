package eventmanager

import (
	"errors"
	"sync"

	"github.com/ztrue/tracerr"
)

var Master *Manager

type Manager struct {
	EventsMap map[string]*Event
}

type Event struct {
	Handlers []EventHandler
	Callback EventCallback
	Mutex    sync.Mutex
}

type EventHandler func()
type EventCallback func(interface{}) interface{}

func InitMasterManager() error {
	Master = &Manager{
		EventsMap: make(map[string]*Event),
	}

	return nil
}

func InitCustomManager() (*Manager, error) {
	return &Manager{
		EventsMap: make(map[string]*Event),
	}, nil
}

func (m *Manager) NewEvent(name string, callback EventCallback) error {
	var exist bool
	_, exist = m.EventsMap[name]
	if exist {
		return tracerr.Wrap(errors.New("event already exists"))
	}

	m.EventsMap[name] = &Event{
		Handlers: make([]EventHandler, 0, 8),
		Callback: callback,
		Mutex:    sync.Mutex{},
	}

	return nil
}

func (m *Manager) ExecEvent(name string) (interface{}, error) {
	var exist bool
	_, exist = m.EventsMap[name]
	if !exist {
		return nil, tracerr.Wrap(errors.New("no such event"))
	}

	wait := &sync.WaitGroup{}
	wait.Add(len(m.EventsMap[name].Handlers))
	for _, v := range m.EventsMap[name].Handlers {
		go func(f EventHandler) {
			f()
			wait.Done()
		}(v)
	}
	wait.Wait()

	return m.EventsMap[name].Callback(nil), nil
}

func (m *Manager) Register(event string, handler EventHandler) error {
	var exist bool
	_, exist = m.EventsMap[event]
	if !exist {
		return tracerr.Wrap(errors.New("no such event"))
	}

	m.EventsMap[event].Register(handler)
	return nil
}

func (e *Event) Register(handler EventHandler) {
	e.Mutex.Lock()
	e.Handlers = append(e.Handlers, handler)
	e.Mutex.Unlock()
}
