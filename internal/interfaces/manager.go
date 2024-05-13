package interfaces

type IManager interface {
	Service() IService
	Repository() IRepository
	Processor() IProcessor
	Server() IServer
}
