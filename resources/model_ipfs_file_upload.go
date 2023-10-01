/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type IpfsFileUpload struct {
	Key
	Attributes IpfsFileUploadAttributes `json:"attributes"`
}
type IpfsFileUploadResponse struct {
	Data     IpfsFileUpload `json:"data"`
	Included Included       `json:"included"`
}

type IpfsFileUploadListResponse struct {
	Data     []IpfsFileUpload `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
}

// MustIpfsFileUpload - returns IpfsFileUpload from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustIpfsFileUpload(key Key) *IpfsFileUpload {
	var ipfsFileUpload IpfsFileUpload
	if c.tryFindEntry(key, &ipfsFileUpload) {
		return &ipfsFileUpload
	}
	return nil
}
