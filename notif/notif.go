package notif

import (
	"github.com/godbus/dbus/v5"
)

const (
	dbusDest   = "org.freedesktop.Notifications"
	dbusPath   = "/org/freedesktop/Notifications"
	dbusIFace  = "org.freedesktop.Notifications"
	dbusMethod = "Notify"
	notifyApp  = "Wall Collage"
	notifyIcon = ""
)

type NotificationService struct {
	conn *dbus.Conn
	obj  dbus.BusObject
}

func NewNotificationService() (*NotificationService, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	return &NotificationService{
		conn: conn,
		obj:  conn.Object(dbusDest, dbusPath),
	}, nil
}

func (n *NotificationService) Close() error {
	return n.conn.Close()
}

func (n *NotificationService) Notify(title, body string) error {
	if n == nil {
		return nil
	}

	call := n.obj.Call(dbusIFace+"."+dbusMethod, 0,
		notifyApp,
		uint32(0),
		notifyIcon,
		title,
		body,
		[]string{},
		map[string]dbus.Variant{},
		int32(5000))

	if call.Err != nil {
		return call.Err
	}

	return nil
}
