package helpers

import (
	"html/template"
)

// Content struct holds the parts to merge multiple templates.
// Contents are of type HTML to prevent escaping HTML.
type Content struct{
	ContainerHTML template.HTML 
}
