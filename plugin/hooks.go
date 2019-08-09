package plugin

type EventType string

const (
	EventTypeMessageCreated          EventType = "MessageCreated"          // payload is github.com/MixinNetwork/supergroup.mixin.one/models.Message
	EventTypeProhibitedStatusChanged EventType = "ProhibitedStatusChanged" // payload is bool
	EventTypeOrderPaid               EventType = "OrderPaid"               // payload is github.com/MixinNetwork/supergroup.mixin.one/models.Order
	EventTypeUserCreated             EventType = "UserCreated"             // payload is github.com/MixinNetwork/supergroup.mixin.one/models.User
	EventTypeTransactionCreated      EventType = "TransactionCreated"      // payload is github.com/NocturneLabs/supergroup.mixin.one-plugins/wallet.Transaction
)

var callbacks = map[EventType][]func(interface{}){}

// called by plugin implementations
func (*PluginContext) On(eventName EventType, fn func(interface{})) {
	mutex.RLock()
	defer mutex.RUnlock()

	cs, found := callbacks[eventName]
	if !found {
		cs = []func(interface{}){}
	}
	cs = append(cs, fn)
	callbacks[eventName] = cs
}

// pass hook from plugin
func (*PluginContext) Trigger(eventName EventType, obj interface{}) {
	Trigger(eventName, obj)
}

// called by main supergroup codebase
func Trigger(eventName EventType, obj interface{}) {
	mutex.RLock()
	defer mutex.RUnlock()

	cs, found := callbacks[eventName]
	if !found {
		return
	}

	for _, callback := range cs {
		callback(obj)
	}
}
