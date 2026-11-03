package app

import (
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"net/http"
	"os"
)

func (a *App) voiceShares(w http.ResponseWriter, r *http.Request) {
	var owner string
	if a.H.DB.QueryRow(`SELECT owner FROM voices WHERE id=?`, r.PathValue("id")).Scan(&owner) != nil || owner != uid(r) {
		apphost.Problem(w, 404, "Voice not found")
		return
	}
	rows, e := a.H.DB.Query(`SELECT g.user,u.email FROM grants g JOIN host_users u ON u.id=g.user WHERE g.voice=?`, r.PathValue("id"))
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	defer rows.Close()
	out := []map[string]string{}
	for rows.Next() {
		var user, email string
		if rows.Scan(&user, &email) == nil {
			out = append(out, map[string]string{"user": user, "email": email})
		}
	}
	apphost.JSON(w, 200, out)
}

func (a *App) favoriteVoice(w http.ResponseWriter, r *http.Request) {
	var p struct {
		Favorite bool `json:"favorite"`
	}
	if !apphost.Decode(w, r, &p) {
		return
	}
	v, e := a.findVoice(r.PathValue("id"), uid(r))
	if e != nil {
		apphost.Problem(w, 404, "Voice not found")
		return
	}
	if v.Preset {
		_, e = a.H.DB.Exec(`INSERT INTO preferences(owner,key,body)VALUES(?,?,?) ON CONFLICT(owner,key) DO UPDATE SET body=excluded.body`, uid(r), "favorite-"+v.ID, jsonText(p.Favorite))
	} else if v.Owner == uid(r) {
		v.Favorite = p.Favorite
		_, e = a.H.DB.Exec(`UPDATE voices SET body=? WHERE id=? AND owner=?`, jsonText(v), v.ID, uid(r))
	} else {
		apphost.Problem(w, 403, "Only the owner may edit this voice")
		return
	}
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	apphost.JSON(w, 200, map[string]bool{"favorite": p.Favorite})
}
func (a *App) deleteHistory(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	j, e := a.readJob(r.PathValue("id"), uid(r))
	if e != nil {
		apphost.Problem(w, 404, "Recording not found")
		return
	}
	if j.Status == "running" || j.Status == "queued" {
		apphost.Problem(w, 409, "Cancel this recording before removing it")
		return
	}
	if !safeID.MatchString(j.ID) {
		apphost.Problem(w, 400, "Invalid recording")
		return
	}
	// Reusable audio is materialized into each job; deleting a parent cannot
	// remove completed descendants. Queued descendants fail safely if necessary.
	if _, e = a.H.DB.Exec(`DELETE FROM jobs WHERE id=? AND owner=?`, j.ID, uid(r)); e != nil {
		a.H.DBError(w, e)
		return
	}
	if e = os.RemoveAll(a.work(j.ID)); e != nil {
		apphost.Problem(w, 500, "Metadata was removed but some audio could not be deleted; check disk access")
		return
	}
	a.H.Audit(uid(r), "recording-delete", j.ID)
	apphost.JSON(w, 200, json.RawMessage(`{"deleted":true}`))
}
