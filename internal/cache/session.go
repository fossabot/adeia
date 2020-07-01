package cache

// SetSession stores a sessionID:value pair, with a TTL of expiry seconds.
func (r *RedisCache) SetSession(sessionID, value string, expiry int) error {
	sessKey := buildSessionKey(sessionID)
	return r.SetWithExpiry(sessKey, value, expiry)
}

// GetAndRefreshSession returns the value and resets the TTL for the provided sessID.
func (r *RedisCache) GetAndRefreshSession(sessionID string, expiry int) (string, error) {
	var value string
	sessKey := buildSessionKey(sessionID)
	if err := r.Get(&value, sessKey); err != nil {
		return "", err
	}

	// refresh expiry
	if err := r.Expire(sessKey, expiry); err != nil {
		return "", err
	}
	return value, nil
}

// ExpireSession instants expires the session identified by the sessionID.
func (r *RedisCache) ExpireSession(sessionID string) error {
	sessKey := buildSessionKey(sessionID)
	return r.Expire(sessKey, 0)
}

func buildSessionKey(sessID string) string {
	return buildKey("session", sessID)
}
