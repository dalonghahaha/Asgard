package managers

import (
	"fmt"
	"sync"

	"github.com/dalonghahaha/avenger/components/logger"
	"github.com/dalonghahaha/avenger/tools/uuid"

	"Asgard/clients"
	"Asgard/runtimes"
)

type AppManager struct {
	lock         sync.Mutex
	apps         map[int64]*runtimes.App
	masterClient *clients.Master
	monitor      *runtimes.Monitor
}

func NewAppManager() *AppManager {
	manager := &AppManager{
		apps:    make(map[int64]*runtimes.App),
		monitor: runtimes.NewMonitor("app"),
	}
	return manager
}

func (m *AppManager) SetMaster(masterClient *clients.Master) {
	m.masterClient = masterClient
}

func (m *AppManager) StartMonitor() {
	go m.monitor.Start()
}

func (m *AppManager) NewApp(config map[string]interface{}) (*runtimes.App, error) {
	app := new(runtimes.App)
	err := app.Configure(config)
	if err != nil {
		return nil, err
	}
	autoRestart, ok := config["auto_restart"].(bool)
	if !ok {
		return nil, fmt.Errorf("config auto_restart type wrong")
	}
	app.AutoRestart = autoRestart
	return app, nil
}

func (m *AppManager) Register(id int64, config map[string]interface{}) error {
	app, err := m.NewApp(config)
	if err != nil {
		return err
	}
	app.ID = id
	app.Monitor = m.monitor
	app.ExceptionReport = func(message string) {
		appException := runtimes.AppException{
			UUID:  uuid.GenerateV4(),
			AppID: app.ID,
			Desc:  message,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(appException.UUID, 1)
			m.masterClient.AppExceptionChan <- appException
		}
	}
	app.MonitorReport = func(monitor *runtimes.MonitorInfo) {
		appMonitor := runtimes.AppMonitor{
			UUID:    uuid.GenerateV4(),
			App:     app,
			Monitor: monitor,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(appMonitor.UUID, 1)
			m.masterClient.AppMonitorChan <- appMonitor
		}
	}
	app.ArchiveReport = func(archive *runtimes.Archive) {
		appArchive := runtimes.AppArchive{
			UUID:    uuid.GenerateV4(),
			App:     app,
			Archive: archive,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(appArchive.UUID, 1)
			m.masterClient.AppArchiveChan <- appArchive
		}
	}
	m.lock.Lock()
	m.apps[id] = app
	m.lock.Unlock()
	return nil
}

func (m *AppManager) UnRegister(id int64) {
	m.lock.Lock()
	delete(m.apps, id)
	m.lock.Unlock()
}

func (m *AppManager) Count() int {
	return len(m.apps)
}

func (m *AppManager) GetList() []*runtimes.App {
	list := []*runtimes.App{}
	for _, app := range m.apps {
		list = append(list, app)
	}
	return list
}

func (m *AppManager) Get(id int64) *runtimes.App {
	app, ok := m.apps[id]
	if !ok {
		return nil
	}
	return app
}

func (m *AppManager) GetByName(name string) *runtimes.App {
	for _, app := range m.apps {
		if app.Name == name {
			return app
		}
	}
	return nil
}

func (m *AppManager) StartAll(monitor bool) {
	if monitor {
		m.StartMonitor()
	}
	for _, app := range m.apps {
		go app.Run()
	}
}

func (m *AppManager) StopAll() {
	for _, app := range m.apps {
		logger.Infof("killing app %s, runing: %v", app.Name, app.Running)
		if app.Running {
			app.Kill()
		}
	}
}

func (m *AppManager) Start(id int64) bool {
	app := m.Get(id)
	if app == nil {
		return true
	}
	if app.Running {
		return true
	}
	go app.Run()
	return true
}

func (m *AppManager) Stop(id int64) bool {
	app := m.Get(id)
	if app == nil {
		return true
	}
	if !app.Running {
		return true
	}
	app.AutoRestart = false
	go app.Kill()
	return true
}

func (m *AppManager) Add(id int64, config map[string]interface{}) error {
	app := m.Get(id)
	if app != nil {
		if !app.Running {
			ok := m.Start(id)
			if !ok {
				return fmt.Errorf("app start failed!")
			}
		}
		return nil
	}
	err := m.Register(id, config)
	if err != nil {
		return fmt.Errorf("app register failed:%s", err.Error())
	}
	ok := m.Start(id)
	if !ok {
		return fmt.Errorf("app start failed!")
	}
	return nil
}

func (m *AppManager) Update(id int64, config map[string]interface{}) error {
	err := m.Remove(id)
	if err != nil {
		return err
	}
	return m.Add(id, config)
}

func (m *AppManager) Remove(id int64) error {
	ok := m.Stop(id)
	if !ok {
		return fmt.Errorf("app stop failed!")
	}
	m.UnRegister(id)
	return nil
}
