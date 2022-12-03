package telegram

const msgHelp = `I can save and keep your pages. Also I can offer you some of them to read

In order to save the page just send me a link.

In order to get a random page from your list send me command "/rnd"
Caution! After that your link will be deleted from your list`

const msgHello = "Hi there! \n\n" + msgHelp

const (
	msgUnknownCommand = "Unknown command received"
	msgNoSavedPages   = "You have no saved pages!"
	msgSaved          = "Saved!"
	msgAlreadyExists  = "You have already have this page in your list"
)
