package routes

import (
	"github.com/jbockle/captivated/server/routes/events_handler"
	"github.com/jbockle/captivated/server/routes/home_handler"
	"github.com/jbockle/captivated/server/routes/new_event_handler"
)

var NewEventHandler = new_event_handler.NewEventHandler
var HomeHandler = home_handler.HomeHandler
var EventsHandler = events_handler.EventsHandler
