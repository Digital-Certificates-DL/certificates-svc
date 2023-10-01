/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ParseUsers struct {
	Key
	Attributes ParseUsersAttributes `json:"attributes"`
}
type ParseUsersResponse struct {
	Data     ParseUsers `json:"data"`
	Included Included   `json:"included"`
}

type ParseUsersListResponse struct {
	Data     []ParseUsers `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustParseUsers - returns ParseUsers from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustParseUsers(key Key) *ParseUsers {
	var parseUsers ParseUsers
	if c.tryFindEntry(key, &parseUsers) {
		return &parseUsers
	}
	return nil
}
