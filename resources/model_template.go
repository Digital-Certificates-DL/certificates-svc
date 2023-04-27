/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Template struct {
	Key
	Attributes TemplateAttributes `json:"attributes"`
}
type TemplateResponse struct {
	Data     Template `json:"data"`
	Included Included `json:"included"`
}

type TemplateListResponse struct {
	Data     []Template `json:"data"`
	Included Included   `json:"included"`
	Links    *Links     `json:"links"`
}

// MustTemplate - returns Template from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustTemplate(key Key) *Template {
	var template Template
	if c.tryFindEntry(key, &template) {
		return &template
	}
	return nil
}
