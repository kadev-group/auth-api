package interfaces

type IServer interface {
	REST() IRESTServer
}

type IRESTServer interface {
	Run()
}
