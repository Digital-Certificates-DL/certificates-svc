/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ResponseTemplate struct {
	Key
	Attributes ResponseTemplateAttributes `json:"attributes"`
}
type ResponseTemplateResponse struct {
	Data     ResponseTemplate `json:"data"`
	Included Included         `json:"included"`
}

type ResponseTemplateListResponse struct {
	Data     []ResponseTemplate `json:"data"`
	Included Included           `json:"included"`
	Links    *Links             `json:"links"`
}

// MustResponseTemplate - returns ResponseTemplate from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustResponseTemplate(key Key) *ResponseTemplate {
	var responseTemplate ResponseTemplate
	if c.tryFindEntry(key, &responseTemplate) {
		return &responseTemplate
	}
	return nil
}
