package model

// ScriptDescriptor describes the meta data of a Nexus 3 script.
// More information can be found here: https://help.sonatype.com/display/NXRM3/Script+API
type ScriptDescriptor struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
}
