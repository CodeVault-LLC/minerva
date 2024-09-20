package cache

import (
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

var SessionManager *scs.SessionManager

func InitSessionManager() *scs.SessionManager {
	SessionManager = scs.New()
	SessionManager.Store = memstore.New()

	return SessionManager
}
