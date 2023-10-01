/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Settings struct {
	Key
	Attributes SettingsAttributes `json:"attributes"`
}
type SettingsResponse struct {
	Data     Settings `json:"data"`
	Included Included `json:"included"`
}

type SettingsListResponse struct {
	Data     []Settings `json:"data"`
	Included Included   `json:"included"`
	Links    *Links     `json:"links"`
}

// MustSettings - returns Settings from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSettings(key Key) *Settings {
	var settings Settings
	if c.tryFindEntry(key, &settings) {
		return &settings
	}
	return nil
}
