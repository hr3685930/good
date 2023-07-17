package queue

// ChannelDrive ChannelDrive
type ChannelDrive struct {}

// Connect Connect
func (m ChannelDrive) Connect(key string) error {
	c := NewChannel()
	QueueStore.Store(key, c)
	return nil
}

// Register Register
func (ChannelDrive) Register(key string) {
	MQ = GetQueueDrive(key)
}
