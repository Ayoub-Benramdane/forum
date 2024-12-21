package database

import (
    structs "forum/Data"
    "time"
    "net/http"
)

func CreateSession(w http.ResponseWriter, username string, id int64) error {
    _, err := DB.Exec("DELETE FROM session WHERE user_id = ?", id)
    if err != nil {
        return err
    }

    cookie := &http.Cookie{
        Name:     "session_id",
        Value:    username, 
        Expires:  time.Now().Add(3 * time.Minute),
    }
    http.SetCookie (w, cookie)

    _, err = DB.Exec(`
        INSERT INTO session (username, user_id, statut, last_activity) 
        VALUES (?, ?, ?, ?)`,
        username, id, "Connected", time.Now())

    return err
}

func GetUserConnected(r *http.Request) *structs.Session {
    cookie, err := r.Cookie("session_id")
    if err != nil {
        return nil
    }

    var session structs.Session
    var lastActivity time.Time

    err = DB.QueryRow(`
        SELECT id, username, user_id, statut, last_activity 
        FROM session 
        WHERE username = ?`,
        cookie.Value).Scan(
            &session.ID,
            &session.Username,
            &session.UserID,
            &session.Status,
            &lastActivity)

    if err != nil {
        return nil
    }

    if time.Since(lastActivity) > 3*time.Minute {
        DeleteSession(session.Username)
        return nil
    }

    DB.Exec("UPDATE session SET last_activity = ? WHERE username = ?",
        time.Now(), session.Username)

    return &session
}

func DeleteSession(username string) error {
    _, err := DB.Exec("DELETE FROM session WHERE username = ?", username)
    return err
}