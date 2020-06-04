package managers

import (
	"fmt"
	"sync"

	"github.com/dalonghahaha/avenger/tools/uuid"

	"Asgard/applications"
	"Asgard/client"
)

type AppManager struct {
	lock           sync.Mutex
	apps           map[int64]*applications.App
	masterClient   *client.Master
	monitorManager *applications.MonitorMamager
}

func NewAppManager() (*AppManager, error) {
	manager := &AppManager{
		apps:           make(map[int64]*applications.App),
		monitorManager: applications.NewMonitorMamager(),
	}
	return manager, nil
}

func (m *AppManager) SetMaster(masterClient *client.Master) {
	m.masterClient = masterClient
}

func (m *AppManager) StartMonitor() {
	m.monitorManager.Start()
}

func (m *AppManager) StopMonitor() {
	m.monitorManager.Stop()
}

func (m *AppManager) Count() int {
	return len(m.apps)
}

func (m *AppManager) GetAppList() []*applications.App {
	list := []*applications.App{}
	for _, app := range m.apps {
		list = append(list, app)
	}
	return list
}

func (m *AppManager) GetApp(id int64) *applications.App {
	app, ok := m.apps[id]
	if !ok {
		return nil
	}
	return app
}

func (m *AppManager) StartAll() {
	for _, app := range m.apps {
		go app.Run()
	}
}

func (m *AppManager) StopAll() {
	for _, app := range m.apps {
		if app.Running {
			go app.Kill()
		}
	}
}

func (m *AppManager) Start(id int64) bool {
	app := m.GetApp(id)
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
	app := m.GetApp(id)
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

func (m *AppManager) NewApp(config map[string]interface{}) (*applications.App, error) {
	app := new(applications.App)
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
	app.MonitorMamager = m.monitorManager
	app.MonitorReport = func(monitor *applications.Monitor) {
		appMonitor := applications.AppMonitor{
			UUID:    uuid.GenerateV4(),
			App:     app,
			Monitor: monitor,
		}
		if m.masterClient != nil {
			m.masterClient.Reports.Store(appMonitor.UUID, 1)
			m.masterClient.AppMonitorChan <- appMonitor
		}
	}
	app.ArchiveReport = func(archive *applications.Archive) {
		appArchive := applications.AppArchive{
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
