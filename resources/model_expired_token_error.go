/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ExpiredTokenError struct {
	Key
	Attributes ExpiredTokenErrorAttributes `json:"attributes"`
}
type ExpiredTokenErrorResponse struct {
	Data     ExpiredTokenError `json:"data"`
	Included Included          `json:"included"`
}

type ExpiredTokenErrorListResponse struct {
	Data     []ExpiredTokenError `json:"data"`
	Included Included            `json:"included"`
	Links    *Links              `json:"links"`
}

// MustExpiredTokenError - returns ExpiredTokenError from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustExpiredTokenError(key Key) *ExpiredTokenError {
	var expiredTokenError ExpiredTokenError
	if c.tryFindEntry(key, &expiredTokenError) {
		return &expiredTokenError
	}
	return nil
}
