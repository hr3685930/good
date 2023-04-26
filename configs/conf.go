package configs

// ENV ENV
var ENV Conf

//Conf Conf
type Conf struct {
	App      App
	Database Database
	Queue    Queue
	Cache    Cache
	Trace    Trace
}
