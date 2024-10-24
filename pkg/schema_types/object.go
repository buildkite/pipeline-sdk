package schema_types

import (
	"fmt"
	"strings"

	"github.com/buildkite/pipeline-sdk/pkg/utils"
)

type SchemaObject struct {
	Name       AttributeName
	Properties []Field
}

func (SchemaObject) IsUnion() bool {
	return false
}

func (s SchemaObject) Fields() []Field {
	return s.Properties
}

func (s SchemaObject) TypeScriptType() string {
	tsInterface := utils.CodeBlock{
		fmt.Sprintf("export interface %s {", s.Name.TitleCase()),
	}

	var properties []string
	for _, p := range s.Properties {
		optionalMarker := ""
		if !p.required {
			optionalMarker = "?"
		}

		if p.fieldref != nil {
			properties = append(properties, utils.CodeBlock{
				utils.CodeGen.Comment.TypeScript(p.description, 0),
				fmt.Sprintf("%s%s: %s;", p.name.CamelCase(), optionalMarker, p.fieldref.name.TitleCase()),
			}.DisplayIndent(4))
			continue
		}

		propType := p.typ.TypeScriptType()

		properties = append(properties, utils.CodeBlock{
			utils.CodeGen.Comment.TypeScript(p.description, 0),
			fmt.Sprintf("%s%s: %s;", p.name.CamelCase(), optionalMarker, propType),
		}.DisplayIndent(4))
	}

	tsInterface = append(tsInterface, properties...)
	tsInterface = append(tsInterface, "}")

	return tsInterface.Display()
}

func (s SchemaObject) GoType() string {
	var properties []string
	for _, prop := range s.Properties {
		if prop.fieldref != nil {
			properties = append(properties, utils.CodeBlock{
				utils.CodeGen.Comment.Go(prop.description, 0),
				fmt.Sprintf("%s %s `json:\"%s,omitempty\"`", prop.name.TitleCase(), prop.fieldref.name.TitleCase(), prop.name),
			}.DisplayIndent(4))
			continue
		}

		propType := prop.typ.GoType()
		switch prop.typ.(type) {
		case SchemaUnion:
			properties = append(properties, utils.CodeBlock{
				utils.CodeGen.Comment.Go(prop.description, 0),
				fmt.Sprintf("%s %s `json:\"%s,omitempty\"`", prop.name.TitleCase(), prop.name.TitleCase(), prop.name),
			}.DisplayIndent(4))
		default:
			properties = append(properties, fmt.Sprintf("%s\n%s\n",
				utils.CodeGen.Comment.Go(prop.description, 4),
				fmt.Sprintf("    %s %s `json:\"%s,omitempty\"`", prop.name.TitleCase(), propType, prop.name),
			))
		}
	}

	return utils.CodeBlock{
		fmt.Sprintf("type %s struct {", s.Name.TitleCase()),
		strings.Join(properties, "\n"),
		"}",
	}.Display()
}

type object struct{}

func (object) New(name string, props []Field) SchemaObject {
	return SchemaObject{
		Name:       AttributeName(name),
		Properties: props,
	}
}

var Object = object{}
