package constants

import (
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

var SessionManager *scs.SessionManager

func InitSessionManager() {
	SessionManager = scs.New()
	SessionManager.Store = memstore.New()
}
