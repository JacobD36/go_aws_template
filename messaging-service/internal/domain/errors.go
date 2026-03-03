package domain

import "errors"

var (
	// ErrMessageSendFailed indica que el envío del mensaje falló
	ErrMessageSendFailed = errors.New("failed to send message")

	// ErrInvalidMessage indica que el mensaje no es válido
	ErrInvalidMessage = errors.New("invalid message")

	// ErrEventProcessing indica que hubo un error procesando el evento
	ErrEventProcessing = errors.New("failed to process event")
)
