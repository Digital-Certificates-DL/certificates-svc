/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type PrepareCertificates struct {
	Key
	Attributes PrepareCertificatesAttributes `json:"attributes"`
}
type PrepareCertificatesResponse struct {
	Data     PrepareCertificates `json:"data"`
	Included Included            `json:"included"`
}

type PrepareCertificatesListResponse struct {
	Data     []PrepareCertificates `json:"data"`
	Included Included              `json:"included"`
	Links    *Links                `json:"links"`
}

// MustPrepareCertificates - returns PrepareCertificates from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustPrepareCertificates(key Key) *PrepareCertificates {
	var prepareCertificates PrepareCertificates
	if c.tryFindEntry(key, &prepareCertificates) {
		return &prepareCertificates
	}
	return nil
}
