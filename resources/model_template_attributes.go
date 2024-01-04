/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type TemplateAttributes struct {
	BackgroundImg     string          `json:"background_img"`
	IsCompleted       bool            `json:"is_completed"`
	IsDefaultTemplate bool            `json:"is_default_template"`
	Template          json.RawMessage `json:"template"`
	TemplateId        int64           `json:"template_id"`
	TemplateName      string          `json:"template_name"`
	TemplateShortName string          `json:"template_short_name"`
}
