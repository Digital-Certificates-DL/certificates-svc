/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Container struct {
	Key
	Attributes ContainerAttributes `json:"attributes"`
}
type ContainerResponse struct {
	Data     Container `json:"data"`
	Included Included  `json:"included"`
}

type ContainerListResponse struct {
	Data     []Container `json:"data"`
	Included Included    `json:"included"`
	Links    *Links      `json:"links"`
}

// MustContainer - returns Container from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustContainer(key Key) *Container {
	var container Container
	if c.tryFindEntry(key, &container) {
		return &container
	}
	return nil
}
