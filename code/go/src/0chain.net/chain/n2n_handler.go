package chain

import (
	"net/http"
)

func (c *Chain) SetupNodeHandlers() {
	http.HandleFunc("/_nh/status", c.StatusHandler)
	http.HandleFunc("/_nh/list/m", c.GetMinersHandler)
	http.HandleFunc("/_nh/list/s", c.GetShardersHandler)
	http.HandleFunc("/_nh/list/b", c.GetBlobbersHandler)
}