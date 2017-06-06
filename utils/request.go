package utils

type IRequest interface {
	POST() (error)
}

type Request struct {

}
