/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type IpfsFile struct {
	Key
	Attributes IpfsFileAttributes `json:"attributes"`
}
type IpfsFileResponse struct {
	Data     IpfsFile `json:"data"`
	Included Included `json:"included"`
}

type IpfsFileListResponse struct {
	Data     []IpfsFile `json:"data"`
	Included Included   `json:"included"`
	Links    *Links     `json:"links"`
}

// MustIpfsFile - returns IpfsFile from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustIpfsFile(key Key) *IpfsFile {
	var ipfsFile IpfsFile
	if c.tryFindEntry(key, &ipfsFile) {
		return &ipfsFile
	}
	return nil
}
