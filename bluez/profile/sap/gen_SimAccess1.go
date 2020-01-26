package sap



import (
   "sync"
   "github.com/muka/go-bluetooth/bluez"
   "github.com/muka/go-bluetooth/util"
   "github.com/muka/go-bluetooth/props"
   "github.com/godbus/dbus"
)

var SimAccess1Interface = "org.bluez.SimAccess1"


// NewSimAccess1 create a new instance of SimAccess1
//
// Args:
// - objectPath: [variable prefix]/{hci0,hci1,...}
func NewSimAccess1(objectPath dbus.ObjectPath) (*SimAccess1, error) {
	a := new(SimAccess1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: SimAccess1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(SimAccess1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


/*
SimAccess1 Sim Access Profile hierarchy

*/
type SimAccess1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*SimAccess1Properties
	watchPropertiesChannel chan *dbus.Signal
}

// SimAccess1Properties contains the exposed properties of an interface
type SimAccess1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	/*
	Connected Indicates if SAP client is connected to the server.
	*/
	Connected bool

}

//Lock access to properties
func (p *SimAccess1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *SimAccess1Properties) Unlock() {
	p.lock.Unlock()
}






// GetConnected get Connected value
func (a *SimAccess1) GetConnected() (bool, error) {
	v, err := a.GetProperty("Connected")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}



// Close the connection
func (a *SimAccess1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return SimAccess1 object path
func (a *SimAccess1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Client return SimAccess1 dbus client
func (a *SimAccess1) Client() *bluez.Client {
	return a.client
}

// Interface return SimAccess1 interface
func (a *SimAccess1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *SimAccess1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

	if a.objectManagerSignal == nil {
		if a.objectManager == nil {
			om, err := bluez.GetObjectManager()
			if err != nil {
				return nil, nil, err
			}
			a.objectManager = om
		}

		s, err := a.objectManager.Register()
		if err != nil {
			return nil, nil, err
		}
		a.objectManagerSignal = s
	}

	cancel := func() {
		if a.objectManagerSignal == nil {
			return
		}
		a.objectManagerSignal <- nil
		a.objectManager.Unregister(a.objectManagerSignal)
		a.objectManagerSignal = nil
	}

	return a.objectManagerSignal, cancel, nil
}


// ToMap convert a SimAccess1Properties to map
func (a *SimAccess1Properties) ToMap() (map[string]interface{}, error) {
	return props.ToMap(a), nil
}

// FromMap convert a map to an SimAccess1Properties
func (a *SimAccess1Properties) FromMap(props map[string]interface{}) (*SimAccess1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an SimAccess1Properties
func (a *SimAccess1Properties) FromDBusMap(props map[string]dbus.Variant) (*SimAccess1Properties, error) {
	s := new(SimAccess1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// ToProps return the properties interface
func (a *SimAccess1) ToProps() bluez.Properties {
	return a.Properties
}

// GetWatchPropertiesChannel return the dbus channel to receive properties interface
func (a *SimAccess1) GetWatchPropertiesChannel() chan *dbus.Signal {
	return a.watchPropertiesChannel
}

// SetWatchPropertiesChannel set the dbus channel to receive properties interface
func (a *SimAccess1) SetWatchPropertiesChannel(c chan *dbus.Signal) {
	a.watchPropertiesChannel = c
}

// GetProperties load all available properties
func (a *SimAccess1) GetProperties() (*SimAccess1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *SimAccess1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *SimAccess1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *SimAccess1) GetPropertiesSignal() (chan *dbus.Signal, error) {

	if a.propertiesSignal == nil {
		s, err := a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
		if err != nil {
			return nil, err
		}
		a.propertiesSignal = s
	}

	return a.propertiesSignal, nil
}

// Unregister for changes signalling
func (a *SimAccess1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *SimAccess1) WatchProperties() (chan *bluez.PropertyChanged, error) {
	return bluez.WatchProperties(a)
}

func (a *SimAccess1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	return bluez.UnwatchProperties(a, ch)
}




/*
Disconnect 
			Disconnects SAP client from the server.

			Possible errors: org.bluez.Error.Failed


*/
func (a *SimAccess1) Disconnect() error {
	
	return a.client.Call("Disconnect", 0, ).Store()
	
}

