package models

type Error struct {
	Message string
}

func NewError(error error) Error {
	err := Error{
		Message: error.Error(),
	}
	return err
}
