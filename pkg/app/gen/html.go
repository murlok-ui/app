//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type tag struct {
	Name          string
	Type          tagType
	Doc           string
	Attrs         []attr
	EventHandlers []eventHandler
}

type tagType int

const (
	parent tagType = iota
	privateParent
	selfClosing
)

var tags = []tag{
	// A:
	{
		Name: "A",
		Doc:  "that creates a hyperlink, allowing navigation to other web pages or resources.",
		Attrs: withGlobalAttrs(attrsByNames(
			"download",
			"href",
			"hreflang",
			"media",
			"ping",
			"rel",
			"target",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Abbr",
		Doc:           "that represents an abbreviation or an acronym, providing a longer description or meaning of the content.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Address",
		Doc:           "that designates contact information for the author or owner of a document or web page.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Area",
		Type: selfClosing,
		Doc:  "that defines a clickable region within an image map, usually linking to another resource.",
		Attrs: withGlobalAttrs(attrsByNames(
			"alt",
			"coords",
			"download",
			"href",
			"hreflang",
			"media",
			"rel",
			"shape",
			"target",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Article",
		Doc:           "that marks a self-contained composition in a document, like a blog post or news story.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Aside",
		Doc:           "that represents content tangentially related to the main content, and can be considered separate.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Audio",
		Doc:  "that embeds an audio player for playing sound or music content.",
		Attrs: withGlobalAttrs(attrsByNames(
			"autoplay",
			"controls",
			"crossorigin",
			"loop",
			"muted",
			"preload",
			"src",
		)...),
		EventHandlers: withMediaEventHandlers(withGlobalEventHandlers()...),
	},

	// B:
	{
		Name:          "B",
		Doc:           "that applies bold styling to its content.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Base",
		Type: selfClosing,
		Doc:  "that specifies the base URL and target for all relative URLs in the document.",
		Attrs: withGlobalAttrs(attrsByNames(
			"href",
			"target",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Bdi",
		Doc:           "that isolates a section of text, allowing it to be formatted in a different direction than the surrounding content.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Bdo",
		Doc:           "that controls the text direction of its content, overriding other directional settings.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Blockquote",
		Doc:  "that represents a section of text quoted from another source.",
		Attrs: withGlobalAttrs(attrsByNames(
			"cite",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "Body",
		Type:  privateParent,
		Doc:   "that encloses the main content of the HTML document.",
		Attrs: withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"onafterprint",
			"onbeforeprint",
			"onbeforeunload",
			"onerror",
			"onhashchange",
			"onload",
			"onmessage",
			"onoffline",
			"ononline",
			"onpagehide",
			"onpageshow",
			"onpopstate",
			"onresize",
			"onstorage",
			"onunload",
		)...),
	},
	{
		Name:          "Br",
		Type:          selfClosing,
		Doc:           "that inserts a line break within inline content.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Button",
		Doc:  "that creates a clickable button, typically used for form submission or triggering interactions.",
		Attrs: withGlobalAttrs(attrsByNames(
			"autofocus",
			"disabled",
			"form",
			"formaction",
			"formenctype",
			"formmethod",
			"formnovalidate",
			"formtarget",
			"name",
			"type",
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// C:
	{
		Name: "Canvas",
		Doc:  "that provides a space where graphics can be rendered dynamically, such as 2D drawings or 3D visualizations.",
		Attrs: withGlobalAttrs(attrsByNames(
			"height",
			"width",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Caption",
		Doc:           "that represents the title or description of a table, usually appearing above or below the table.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Cite",
		Doc:           "that indicates the title or reference of a creative work, such as a book, film, or research paper.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Code",
		Doc:           "that displays a single line of code or a code snippet, preserving its formatting.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Col",
		Type: selfClosing,
		Doc:  "that defines the properties for a single column or a group of columns within a table, when nested within a `<colgroup>` element.",
		Attrs: withGlobalAttrs(attrsByNames(
			"span",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "ColGroup",
		Doc:  "that groups one or more `<col>` elements, providing a way to apply styles and attributes to multiple columns simultaneously.",
		Attrs: withGlobalAttrs(attrsByNames(
			"span",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// D:
	{
		Name: "Data",
		Doc:  "that pairs content with its machine-readable translation or value.",
		Attrs: withGlobalAttrs(attrsByNames(
			"value",
		)...),
	},
	{
		Name:          "DataList",
		Doc:           "that offers a predefined set of options for input controls.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Dd",
		Doc:           "that provides the description or value for a term in a description list.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Del",
		Doc:  "that denotes text segments that have been deleted or modified in the content.",
		Attrs: withGlobalAttrs(attrsByNames(
			"cite",
			"datetime",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Details",
		Doc:  "that encapsulates content users can toggle visibility for, such as additional information or context.",
		Attrs: withGlobalAttrs(attrsByNames(
			"open",
		)...),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"ontoggle",
		)...),
	},
	{
		Name:          "Dfn",
		Doc:           "that marks the defining occurrence or clarification of a term or phrase.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Dialog",
		Doc:  "that represents a popup dialog box or an interactive window overlay.",
		Attrs: withGlobalAttrs(attrsByNames(
			"open",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Div",
		Doc:           "that creates a generic container for flow content, usually combined with styles or scripts.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Dl",
		Doc:           "that structures a list of terms alongside their associated descriptions.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Dt",
		Doc:           "that specifies a term or name within a description list.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// E:
	{
		Name: "Elem",
		Doc:  "that is customizable.",
		Attrs: withGlobalAttrs(attrsByNames(
			"xmlns",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "ElemSelfClosing",
		Type: selfClosing,
		Doc:  "that is self-closing and customizable.",
		Attrs: withGlobalAttrs(attrsByNames(
			"xmlns",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Em",
		Doc:           "that marks text for emphasis, typically rendered as italicized text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Embed",
		Type: selfClosing,
		Doc:  "that offers a container for integrating non-HTML content or applications.",
		Attrs: withGlobalAttrs(attrsByNames(
			"height",
			"src",
			"type",
			"width",
		)...),
		EventHandlers: withMediaEventHandlers(withGlobalEventHandlers()...),
	},

	// F:
	{
		Name: "FieldSet",
		Doc:  "that clusters related input controls and labels within a form.",
		Attrs: withGlobalAttrs(attrsByNames(
			"disabled",
			"form",
			"name",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "FigCaption",
		Doc:           "that supplies a caption or explanation for content within the <figure> element.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Figure",
		Doc:           "that encapsulates media content or illustrations with an optional caption.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Footer",
		Doc:           "that denotes the footer of a section or the whole document, often containing metadata or author info.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Form",
		Doc:  "that constructs a user input form, allowing for various control elements and submission options.",
		Attrs: withGlobalAttrs(attrsByNames(
			"accept-charset",
			"action",
			"autocomplete",
			"enctype",
			"method",
			"name",
			"novalidate",
			"target",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// H:
	{
		Name:          "H1",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "H2",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "H3",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "H4",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "H5",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "H6",
		Doc:           "defines HTML heading.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "Head",
		Doc:   "defines information about the document.",
		Attrs: withGlobalAttrs(attrsByNames()...),
	},
	{
		Name:          "Header",
		Doc:           "defines a header for a document or section.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Hr",
		Type:          selfClosing,
		Doc:           "defines a thematic change in the content.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "Html",
		Type:  privateParent,
		Doc:   "defines the root of an HTML document.",
		Attrs: withGlobalAttrs(),
	},

	// I:
	{
		Name:          "I",
		Doc:           "defines a part of text in an alternate voice or mood.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "IFrame",
		Doc:  "defines an inline frame.",
		Attrs: withGlobalAttrs(attrsByNames(
			"allow",
			"allowfullscreen",
			"allowpaymentrequest",
			"height",
			"name",
			"referrerpolicy",
			"sandbox",
			"src",
			"srcdoc",
			"width",
			"loading",
		)...),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"onload",
		)...,
		),
	},
	{
		Name: "Img",
		Type: selfClosing,
		Doc:  "defines an image.",
		Attrs: withGlobalAttrs(attrsByNames(
			"alt",
			"crossorigin",
			"fetchpriority",
			"height",
			"ismap",
			"sizes",
			"src",
			"srcset",
			"usemap",
			"width",
		)...),
		EventHandlers: withMediaEventHandlers(withGlobalEventHandlers(
			eventHandlersByName(
				"onload",
			)...,
		)...),
	},
	{
		Name: "Input",
		Type: selfClosing,
		Doc:  "defines an input control.",
		Attrs: withGlobalAttrs(attrsByNames(
			"accept",
			"alt",
			"autocomplete",
			"autofocus",
			"capture",
			"checked",
			"dirname",
			"disabled",
			"form",
			"formaction",
			"formenctype",
			"formmethod",
			"formnovalidate",
			"formtarget",
			"height",
			"list",
			"max",
			"maxlength",
			"min",
			"multiple",
			"name",
			"pattern",
			"placeholder",
			"readonly",
			"required",
			"size",
			"src",
			"step",
			"type",
			"value",
			"width",
		)...),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"onload",
		)...,
		),
	},
	{
		Name:          "Ins",
		Doc:           "defines a text that has been inserted into a document.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// K:
	{
		Name:          "Kbd",
		Doc:           "defines keyboard input.",
		Attrs:         withGlobalAttrs(attrsByNames()...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// L:
	{
		Name: "Label",
		Doc:  "defines a label for an input element.",
		Attrs: withGlobalAttrs(attrsByNames(
			"for",
			"form",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Legend",
		Doc:           "defines a caption for a fieldset element.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Li",
		Doc:  "defines a list item.",
		Attrs: withGlobalAttrs(attrsByNames(
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Link",
		Type: selfClosing,
		Doc:  "defines the relationship between a document and an external resource (most used to link to style sheets).",
		Attrs: withGlobalAttrs(attrsByNames(
			"as",
			"crossorigin",
			"fetchpriority",
			"href",
			"hreflang",
			"media",
			"rel",
			"sizes",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"onload",
		)...),
	},

	// M:
	{
		Name:          "Main",
		Doc:           "specifies the main content of a document.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Map",
		Doc:  "defines a client-side image-map.",
		Attrs: withGlobalAttrs(attrsByNames(
			"name",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Mark",
		Doc:           "defines marked/highlighted text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Meta",
		Type: selfClosing,
		Doc:  ".",
		Attrs: withGlobalAttrs(attrsByNames(
			"charset",
			"content",
			"http-equiv",
			"name",
			"property",
		)...),
	},
	{
		Name: "Meter",
		Doc:  "defines a scalar measurement within a known range (a gauge).",
		Attrs: withGlobalAttrs(attrsByNames(
			"form",
			"high",
			"low",
			"max",
			"min",
			"optimum",
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// N:
	{
		Name:          "Nav",
		Doc:           "defines navigation links.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "NoScript",
		Doc:   "defines an alternate content for users that do not support client-side scripts.",
		Attrs: withGlobalAttrs(attrsByNames()...),
	},

	// O:
	{
		Name: "Object",
		Doc:  "defines an embedded object.",
		Attrs: withGlobalAttrs(attrsByNames(
			"data",
			"form",
			"height",
			"name",
			"type",
			"usemap",
			"width",
		)...),
		EventHandlers: withMediaEventHandlers(withGlobalEventHandlers()...),
	},
	{
		Name: "Ol",
		Doc:  "defines an ordered list.",
		Attrs: withGlobalAttrs(attrsByNames(
			"reversed",
			"start",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "OptGroup",
		Doc:  "defines a group of related options in a drop-down list.",
		Attrs: withGlobalAttrs(attrsByNames(
			"disabled",
			"label",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Option",
		Doc:  "defines an option in a drop-down list.",
		Attrs: withGlobalAttrs(attrsByNames(
			"disabled",
			"label",
			"selected",
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Output",
		Doc:  ".",
		Attrs: withGlobalAttrs(attrsByNames(
			"for",
			"form",
			"name",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// P:
	{
		Name:          "P",
		Doc:           "defines a paragraph.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Param",
		Type: selfClosing,
		Doc:  "defines a parameter for an object.",
		Attrs: withGlobalAttrs(attrsByNames(
			"name",
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Picture",
		Doc:           "defines a container for multiple image resources.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Pre",
		Doc:           "defines preformatted text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Progress",
		Doc:  "represents the progress of a task.",
		Attrs: withGlobalAttrs(attrsByNames(
			"max",
			"value",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// Q:
	{
		Name: "Q",
		Doc:  "defines a short quotation.",
		Attrs: withGlobalAttrs(attrsByNames(
			"cite",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},

	// R:
	{
		Name:          "Rp",
		Doc:           "defines what to show in browsers that do not support ruby annotations.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Rt",
		Doc:           "defines an explanation/pronunciation of characters (for East Asian typography).",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Ruby",
		Doc:           "defines a ruby annotation (for East Asian typography).",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// S:
	{
		Name:          "S",
		Doc:           "Defines text that is no longer correct.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Samp",
		Doc:           "defines sample output from a computer program.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Script",
		Doc:  "defines a client-side script.",
		Attrs: withGlobalAttrs(attrsByNames(
			"async",
			"charset",
			"crossorigin",
			"defer",
			"src",
			"type",
		)...),
		EventHandlers: eventHandlersByName("onload"),
	},
	{
		Name:          "Section",
		Doc:           "defines a section in a document.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Select",
		Doc:  "defines a drop-down list.",
		Attrs: withGlobalAttrs(attrsByNames(
			"autofocus",
			"disabled",
			"form",
			"multiple",
			"name",
			"required",
			"size",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Small",
		Doc:           "defines smaller text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Source",
		Type: selfClosing,
		Doc:  ".",
		Attrs: withGlobalAttrs(attrsByNames(
			"src",
			"srcset",
			"media",
			"sizes",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Span",
		Doc:           "defines a section in a document.",
		Attrs:         withGlobalAttrs(attrsByNames()...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Strong",
		Doc:           "defines important text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Style",
		Doc:  "defines style information for a document.",
		Attrs: withGlobalAttrs(attrsByNames(
			"media",
			"type",
		)...),
		EventHandlers: withGlobalEventHandlers(eventHandlersByName(
			"onload",
		)...),
	},
	{
		Name:          "Sub",
		Doc:           "defines subscripted text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Summary",
		Doc:           "defines a visible heading for a details element.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Sup",
		Doc:           "defines superscripted text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// T:
	{
		Name:          "Table",
		Doc:           "defines a table.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "TBody",
		Doc:           "groups the body content in a table.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Td",
		Doc:  "defines a cell in a table.",
		Attrs: withGlobalAttrs(attrsByNames(
			"colspan",
			"headers",
			"rowspan",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "Template",
		Doc:   "defines a template.",
		Attrs: withGlobalAttrs(attrsByNames()...),
	},
	{
		Name: "Textarea",
		Doc:  "defines a multiline input control (text area).",
		Attrs: withGlobalAttrs(attrsByNames(
			"autofocus",
			"cols",
			"dirname",
			"disabled",
			"form",
			"maxlength",
			"name",
			"placeholder",
			"readonly",
			"required",
			"rows",
			"wrap",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "TFoot",
		Doc:           "groups the footer content in a table.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Th",
		Doc:  "defines a header cell in a table.",
		Attrs: withGlobalAttrs(attrsByNames(
			"abbr",
			"colspan",
			"headers",
			"rowspan",
			"scope",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "THead",
		Doc:           "groups the header content in a table",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Time",
		Doc:  "defines a date/time.",
		Attrs: withGlobalAttrs(attrsByNames(
			"datetime",
		)...),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:  "Title",
		Doc:   "defines a title for the document.",
		Attrs: withGlobalAttrs(attrsByNames()...),
	},
	{
		Name:          "Tr",
		Doc:           "defines a row in a table.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// U:
	{
		Name:          "U",
		Doc:           "defines text that should be stylistically different from normal text.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name:          "Ul",
		Doc:           "defines an unordered list.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},

	// V:
	{
		Name:          "Var",
		Doc:           "defines a variable.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
	{
		Name: "Video",
		Doc:  "defines a video or movie.",
		Attrs: withGlobalAttrs(attrsByNames(
			"autoplay",
			"controls",
			"crossorigin",
			"height",
			"loop",
			"muted",
			"poster",
			"preload",
			"src",
			"width",
		)...),
		EventHandlers: withMediaEventHandlers(withGlobalEventHandlers()...),
	},
	{
		Name:          "Wbr",
		Doc:           "defines a possible line-break.",
		Attrs:         withGlobalAttrs(),
		EventHandlers: withGlobalEventHandlers(),
	},
}

type attr struct {
	Name         string
	NameOverride string
	Type         string
	Key          bool
	Doc          string
}

var attrs = map[string]attr{
	// A:
	"abbr": {
		Name: "Abbr",
		Type: "fmt",
		Doc:  "Denotes abbreviated content for header cells to provide clarity on shortened terms.",
	},
	"accept": {
		Name: "Accept",
		Type: "fmt",
		Doc:  "Restricts file types the server accepts, especially used for file input elements.",
	},
	"allow": {
		Name: "Allow",
		Type: "fmt",
		Doc:  "Sets a feature policy, enhancing security by controlling certain browser features. Allows multiple policies.",
	},
	"allowfullscreen": {
		Name: "AllowFullscreen",
		Type: "bool|force",
		Doc:  "Grants an iframe the capability to request fullscreen mode.",
	},
	"allowpaymentrequest": {
		Name: "AllowPaymentRequest",
		Type: "bool|force",
		Doc:  "Grants an iframe the permission to use the Payment Request API for smoother online transactions.",
	},
	"aria-*": {
		Name: "Aria",
		Type: "aria|value",
		Doc:  "Allocates ARIA roles and properties to the element to enhance accessibility for users with disabilities. Can be called multiple times to assign various roles and properties.",
	},
	"attribute": {
		Name: "Attr",
		Type: "attr|value",
		Doc:  "Sets an attribute with its associated value, allowing for flexible HTML customization.",
	},
	"accept-charset": {
		Name:         "AcceptCharset",
		NameOverride: "accept-charset",
		Type:         "fmt",
		Doc:          "Restricts the character encodings accepted for form submission, ensuring compatibility.",
	},
	"accesskey": {
		Name: "AccessKey",
		Type: "fmt",
		Doc:  "Assigns a keyboard shortcut for quick element activation or focus, enhancing user experience.",
	},
	"action": {
		Name: "Action",
		Type: "fmt",
		Doc:  "Specifies the server endpoint to which form-data should be sent upon submission.",
	},
	"alt": {
		Name: "Alt",
		Type: "fmt",
		Doc:  "Provides a text alternative for elements (often images) ensuring content is accessible when visuals can't be rendered.",
	},
	"as": {
		Name: "As",
		Type: "fmt",
		Doc:  "Hints the type of content to preload, optimizing loading for certain resources.",
	},
	"async": {
		Name: "Async",
		Type: "bool",
		Doc:  "Specifies that external scripts are executed asynchronously, preventing blocking of page rendering.",
	},
	"autocomplete": {
		Name: "AutoComplete",
		Type: "on/off",
		Doc:  "Toggles the browser's autocomplete feature, assisting users with common input values.",
	},
	"autofocus": {
		Name: "AutoFocus",
		Type: "bool",
		Doc:  "Instructs the browser to focus this element automatically when the page loads.",
	},
	"autoplay": {
		Name: "AutoPlay",
		Type: "bool",
		Doc:  "Automatically plays audio or video elements once they're ready, enhancing media responsiveness.",
	},

	// C:
	"capture": {
		Name: "Capture",
		Type: "fmt",
		Doc:  "Directs how media capture for file uploads should be handled, such as using the device's camera or microphone.",
	},
	"charset": {
		Name: "Charset",
		Type: "fmt",
		Doc:  "Specifies the character encoding for the linked document or external resource.",
	},
	"checked": {
		Name: "Checked",
		Type: "bool",
		Doc:  "Indicates that an input element (checkbox or radio) should start in a selected state upon page load.",
	},
	"cite": {
		Name: "Cite",
		Type: "fmt",
		Doc:  "Provides a reference or link to a source explaining quoted or modified content in the element.",
	},
	"class": {
		Name: "Class",
		Type: "string|class",
		Doc:  "Assigns one or more classnames to an element, linking it to styles defined in a stylesheet. Can be called multiple times to assign multiple classnames.",
	},
	"cols": {
		Name: "Cols",
		Type: "int",
		Doc:  "Defines the visible width, in character widths, of a text area element.",
	},
	"colspan": {
		Name: "ColSpan",
		Type: "int",
		Doc:  "Denotes how many columns a table cell should span across, allowing cells to occupy space of multiple columns.",
	},
	"content": {
		Name: "Content",
		Type: "fmt",
		Doc:  "Specifies metadata content for the `http-equiv` or `name` attributes, often used in meta tags.",
	},
	"contenteditable": {
		Name: "ContentEditable",
		Type: "bool",
		Doc:  "Determines if the content of an element is editable by the user, allowing for in-page content modification.",
	},
	"controls": {
		Name: "Controls",
		Type: "bool",
		Doc:  "Indicates the presence of user interface controls for audio or video elements, such as play or pause buttons.",
	},
	"coords": {
		Name: "Coords",
		Type: "fmt",
		Doc:  "Defines the coordinates for elements in an image map, establishing active regions for hyperlinks.",
	},
	"crossorigin": {
		Name: "CrossOrigin",
		Type: "fmt",
		Doc:  "Controls how cross-origin requests are managed for the element, supporting secure content integration from different origins.",
	},

	// D:
	"data": {
		Name: "Data",
		Type: "fmt",
		Doc:  "Specifies the URL of a resource associated with an embedded object, such as media or data.",
	},
	"data-*": {
		Name: "DataSet",
		Type: "data|value",
		Doc:  "Allows for storage of custom data specific to individual elements. Can be called multiple times to store multiple sets of data, often used for scripting purposes.",
	},
	"datasets": {
		Name: "DataSets",
		Type: "data|map",
		Doc:  "Denotes datasets linked to an element and can store multiple sets of data.",
	},
	"datetime": {
		Name: "DateTime",
		Type: "fmt",
		Doc:  "Represents the date and time, often used in context with machine-readable equivalents of time-related content.",
	},
	"default": {
		Name: "Default",
		Type: "bool",
		Doc:  "Indicates a default track for media elements and is selected unless the user or browser specifies otherwise.",
	},
	"defer": {
		Name: "Defer",
		Type: "bool",
		Doc:  "Delays the execution of a script until after the document is parsed, typically applied to external scripts.",
	},
	"dir": {
		Name: "Dir",
		Type: "fmt",
		Doc:  "Defines the text direction for the content within an element, such as 'ltr' (left-to-right) or 'rtl' (right-to-left).",
	},
	"dirname": {
		Name: "DirName",
		Type: "fmt",
		Doc:  "Instructs the browser to also submit the text direction of a form field when the form is submitted.",
	},
	"disabled": {
		Name: "Disabled",
		Type: "bool",
		Doc:  "Deactivates an element, rendering it uninteractive and visually distinct.",
	},
	"download": {
		Name: "Download",
		Type: "fmt",
		Doc:  "Hints the browser to download the linked resource, optionally providing a default filename.",
	},
	"draggable": {
		Name: "Draggable",
		Type: "bool",
		Doc:  "Specifies if an element can be dragged by the user, supporting drag-and-drop operations.",
	},

	// E:
	"enctype": {
		Name: "EncType",
		Type: "fmt",
		Doc:  "Describes how form data should be encoded upon submission, especially vital for forms submitting file uploads.",
	},

	// F:
	"fetchpriority": {
		Name: "FetchPriority",
		Type: "fmt",
		Doc:  "Provides a hint to the browser about how it should prioritize the fetch of the image in relation to other images.",
	},
	"for": {
		Name: "For",
		Type: "fmt",
		Doc:  "Associates a label or calculation with specific form element(s).",
	},
	"form": {
		Name: "Form",
		Type: "fmt",
		Doc:  "Identifies the form to which the element belongs.",
	},
	"formaction": {
		Name: "FormAction",
		Type: "fmt",
		Doc:  "Defines the URL to which form data should be sent upon submission. Applicable only to 'submit' type inputs.",
	},
	"formenctype": {
		Name: "FormEncType",
		Type: "fmt",
		Doc:  "Dictates the encoding method for form data prior to its submission to a server. Applicable only to 'submit' type inputs.",
	},
	"formmethod": {
		Name: "FormMethod",
		Type: "fmt",
		Doc:  "Determines the HTTP method for sending form data. Applicable only to 'submit' type inputs.",
	},
	"formnovalidate": {
		Name: "FormNoValidate",
		Type: "bool",
		Doc:  "Indicates that the form data should bypass validation upon submission. Applicable only to 'submit' type inputs.",
	},
	"formtarget": {
		Name: "FormTarget",
		Type: "fmt",
		Doc:  "Specifies where the server's response will be displayed after form submission. Applicable only to 'submit' type inputs.",
	},

	// H:
	"headers": {
		Name: "Headers",
		Type: "fmt",
		Doc:  "Designates one or more header cells to which a table cell is related.",
	},
	"height": {
		Name: "Height",
		Type: "int",
		Doc:  "Sets the height of the element, measured in pixels.",
	},
	"hidden": {
		Name: "Hidden",
		Type: "bool",
		Doc:  "Marks an element as currently irrelevant or not yet relevant.",
	},
	"high": {
		Name: "High",
		Type: "float64",
		Doc:  "Defines the value threshold considered as 'high' in a range context.",
	},
	"href": {
		Name: "Href",
		Type: "fmt",
		Doc:  "Points to the URL of the destination when the link is clicked.",
	},
	"hreflang": {
		Name: "HrefLang",
		Type: "fmt",
		Doc:  "Declares the language of the linked document's content.",
	},
	"http-equiv": {
		Name:         "HTTPEquiv",
		NameOverride: "http-equiv",
		Type:         "fmt",
		Doc:          "Supplies an HTTP header for the content attribute, often used for refresh rates or setting a default charset.",
	},

	// I:
	"id": {
		Name: "ID",
		Type: "fmt",
		Doc:  "Assigns a unique identifier to an element.",
	},
	"ismap": {
		Name: "IsMap",
		Type: "bool",
		Doc:  "Marks an image as a server-side image-map.",
	},

	// K:
	"kind": {
		Name: "Kind",
		Type: "fmt",
		Doc:  "Defines the type of text track for media elements.",
	},

	// L:
	"label": {
		Name: "Label",
		Type: "fmt",
		Doc:  "Provides a concise label for an option element.",
	},
	"lang": {
		Name: "Lang",
		Type: "fmt",
		Doc:  "Declares the language of the element's content.",
	},
	"list": {
		Name: "List",
		Type: "fmt",
		Doc:  "Links to a datalist element offering predefined options for an input element.",
	},
	"loading": {
		Name: "Loading",
		Type: "fmt",
		Doc:  "Determines the browser's loading behavior for the iframe ('eager' or 'lazy').",
	},
	"loop": {
		Name: "Loop",
		Type: "bool",
		Doc:  "Indicates that the audio or video should replay from the beginning upon reaching its end.",
	},
	"low": {
		Name: "Low",
		Type: "float64",
		Doc:  "Sets the value threshold regarded as 'low' in a range context.",
	},

	// M:
	"max": {
		Name: "Max",
		Type: "any",
		Doc:  "Establishes the maximum permissible value.",
	},
	"maxlength": {
		Name: "MaxLength",
		Type: "int",
		Doc:  "Defines the maximum number of characters permissible in an element.",
	},
	"media": {
		Name: "Media",
		Type: "fmt",
		Doc:  "Indicates the intended media or device for the linked document.",
	},
	"method": {
		Name: "Method",
		Type: "fmt",
		Doc:  "Determines the HTTP method for sending form data.",
	},
	"min": {
		Name: "Min",
		Type: "any",
		Doc:  "Establishes the minimum permissible value.",
	},
	"multiple": {
		Name: "Multiple",
		Type: "bool",
		Doc:  "Allows users to input multiple values.",
	},
	"muted": {
		Name: "Muted",
		Type: "bool",
		Doc:  "Ensures that the video's audio playback is muted.",
	},

	// N:
	"name": {
		Name: "Name",
		Type: "fmt",
		Doc:  "Assigns a name to the element.",
	},
	"novalidate": {
		Name: "NoValidate",
		Type: "bool",
		Doc:  "Indicates that the form should bypass validation upon submission.",
	},

	// O:
	"open": {
		Name: "Open",
		Type: "bool",
		Doc:  "Indicates that the details element is expanded and visible to the user.",
	},
	"optimum": {
		Name: "Optimum",
		Type: "float64",
		Doc:  "Sets the optimal numeric value for a gauge element.",
	},

	// P:
	"pattern": {
		Name: "Pattern",
		Type: "fmt",
		Doc:  "Establishes a regular expression against which an input element's value is validated.",
	},
	"ping": {
		Name: "Ping",
		Type: "fmt",
		Doc:  "Lists URLs to be notified when the user activates the hyperlink.",
	},
	"placeholder": {
		Name: "Placeholder",
		Type: "fmt",
		Doc:  "Provides a brief hint describing the expected value of the element.",
	},
	"poster": {
		Name: "Poster",
		Type: "fmt",
		Doc:  "Sets an image displayed before a video starts playing or while it's loading.",
	},
	"preload": {
		Name: "Preload",
		Type: "fmt",
		Doc:  "Indicates the preferred loading method for audio/video upon page load.",
	},
	"property": {
		Name: "Property",
		Type: "fmt",
		Doc:  "Defines the property name of the element.",
	},

	// R:
	"readonly": {
		Name: "ReadOnly",
		Type: "bool",
		Doc:  "Indicates that the element's value cannot be edited by the user.",
	},
	"referrerpolicy": {
		Name: "ReferrerPolicy",
		Type: "fmt",
		Doc:  "Determines the amount of referrer information sent when processing iframe attributes.",
	},
	"rel": {
		Name: "Rel",
		Type: "fmt",
		Doc:  "Describes the relationship between the current and linked documents.",
	},
	"required": {
		Name: "Required",
		Type: "bool",
		Doc:  "Indicates that the element must contain a value before form submission.",
	},
	"reversed": {
		Name: "Reversed",
		Type: "bool",
		Doc:  "States that the list items should be displayed in descending order.",
	},
	"role": {
		Name: "Role",
		Type: "fmt",
		Doc:  "Communicates the intended function or meaning of an element to assistive technologies.",
	},
	"rows": {
		Name: "Rows",
		Type: "int",
		Doc:  "Sets the number of visible lines in a textarea element.",
	},
	"rowspan": {
		Name: "Rowspan",
		Type: "int",
		Doc:  "Determines how many rows a table cell will span vertically.",
	},

	// S:
	"sandbox": {
		Name: "Sandbox",
		Type: "fmt",
		Doc:  "Applies extra security restrictions to content within an iframe.",
	},
	"scope": {
		Name: "Scope",
		Type: "fmt",
		Doc:  "Defines the set of cells a header cell provides header information for. Uses the given format and values.",
	},
	"selected": {
		Name: "Selected",
		Type: "bool",
		Doc:  "Indicates that an option should be pre-selected when the page loads.",
	},
	"shape": {
		Name: "Shape",
		Type: "fmt",
		Doc:  "Describes the shape of a clickable area within an image map. Uses the given format and values.",
	},
	"size": {
		Name: "Size",
		Type: "int",
		Doc:  "Indicates the width of the element, usually in characters for input elements.",
	},
	"sizes": {
		Name: "Sizes",
		Type: "fmt",
		Doc:  "Specifies sizes of icons and images for different page or screen scenarios. Uses the given format and values.",
	},
	"span": {
		Name: "Span",
		Type: "int",
		Doc:  "Defines how many columns or rows a cell should span.",
	},
	"spellcheck": {
		Name: "Spellcheck",
		Type: "bool|force",
		Doc:  "Indicates whether the element's content is subject to spell and grammar checks.",
	},
	"src": {
		Name: "Src",
		Type: "fmt",
		Doc:  "Provides the URL source of embedded content or media. Uses the given format and values.",
	},
	"srcdoc": {
		Name: "SrcDoc",
		Type: "fmt",
		Doc:  "Defines the HTML content to be displayed within an iframe. Uses the given format and values.",
	},
	"srclang": {
		Name: "SrcLang",
		Type: "fmt",
		Doc:  "Denotes the language of text track data (mandatory if kind = 'subtitles'). Uses the given format and values.",
	},
	"srcset": {
		Name: "SrcSet",
		Type: "fmt",
		Doc:  "Provides URLs of images to display in varied resolutions or viewport conditions. Uses the given format and values.",
	},
	"start": {
		Name: "Start",
		Type: "int",
		Doc:  "Determines the starting number for an ordered list.",
	},
	"step": {
		Name: "Step",
		Type: "float64",
		Doc:  "Specifies the interval between permissible values for an input field.",
	},
	"style": {
		Name: "Style",
		Type: "style",
		Doc:  "Assigns inline CSS styling to an element. Can be called multiple times to set multiple CSS styles.",
	},
	"styles": {
		Name: "Styles",
		Type: "style|map",
		Doc:  "Allocates multiple CSS styles to an element. Accepts multiple styling definitions.",
	},

	// T:
	"tabindex": {
		Name: "TabIndex",
		Type: "int",
		Doc:  "Determines the tabbing sequence of an element within the document navigation.",
	},
	"target": {
		Name: "Target",
		Type: "fmt",
		Doc:  "Indicates where to display the linked URL or where to submit the form. Can be called with various predefined values.",
	},
	"title": {
		Name: "Title",
		Type: "fmt",
		Doc:  "Provides additional information about an element, typically displayed as a tooltip. Can be called with the desired title format and content.",
	},
	"type": {
		Name: "Type",
		Type: "fmt",
		Doc:  "Designates the type of the element or its content. Can be called with specific format and values.",
	},

	// U:
	"usemap": {
		Name: "UseMap",
		Type: "fmt",
		Doc:  "Associates the element with a client-side image map. Can be called with the designated format and values.",
	},

	// V:
	"value": {
		Name: "Value",
		Type: "any",
		Doc:  "Assigns a value to the element.",
	},

	// W:
	"width": {
		Name: "Width",
		Type: "int",
		Doc:  "Sets the width of the element.",
	},
	"wrap": {
		Name: "Wrap",
		Type: "fmt",
		Doc:  "Determines how the text inside a text area is wrapped when submitted in a form. Can be called with specific format and values.",
	},
	"xmlns": {
		Name: "XMLNS",
		Type: "xmlns",
		Doc:  "Defines the XML namespace for the element.",
	},
}

func attrsByNames(names ...string) []attr {
	res := make([]attr, 0, len(names))
	for _, n := range names {
		attr, ok := attrs[n]
		if !ok {
			panic("unkowmn attr: " + n)
		}
		res = append(res, attr)
	}

	sort.Slice(res, func(i, j int) bool {
		return strings.Compare(res[i].Name, res[j].Name) <= 0
	})

	return res
}

func withGlobalAttrs(attrs ...attr) []attr {
	attrs = append(attrs, attrsByNames(
		"accesskey",
		"aria-*",
		"class",
		"contenteditable",
		"data-*",
		"datasets",
		"dir",
		"draggable",
		"hidden",
		"id",
		"lang",
		"role",
		"spellcheck",
		"style",
		"styles",
		"tabindex",
		"title",
		"attribute",
	)...)

	sort.Slice(attrs, func(i, j int) bool {
		return strings.Compare(attrs[i].Name, attrs[j].Name) <= 0
	})

	return attrs
}

type eventHandler struct {
	Name string
	Doc  string
}

var eventHandlers = map[string]eventHandler{
	// Window events:
	"onafterprint": {
		Name: "OnAfterPrint",
		Doc:  "Executes the given handler after the document has been printed.",
	},
	"onbeforeprint": {
		Name: "OnBeforePrint",
		Doc:  "Invokes the specified handler before the document gets printed.",
	},
	"onbeforeunload": {
		Name: "OnBeforeUnload",
		Doc:  "Triggers the specified handler when the document is about to be unloaded.",
	},
	"onerror": {
		Name: "OnError",
		Doc:  "Invokes the given handler when an error is encountered.",
	},
	"onhashchange": {
		Name: "OnHashChange",
		Doc:  "Triggers the specified handler when changes occur to the anchor part of the URL.",
	},
	"onload": {
		Name: "OnLoad",
		Doc:  "Executes the specified handler once the element has completely loaded.",
	},
	"onmessage": {
		Name: "OnMessage",
		Doc:  "Triggers the provided handler when a message event occurs.",
	},
	"onoffline": {
		Name: "OnOffline",
		Doc:  "Invokes the given handler when the browser transitions to offline mode.",
	},
	"ononline": {
		Name: "OnOnline",
		Doc:  "Executes the specified handler when the browser transitions to online mode.",
	},
	"onpagehide": {
		Name: "OnPageHide",
		Doc:  "Triggers the given handler when a user navigates away from the current page.",
	},
	"onpageshow": {
		Name: "OnPageShow",
		Doc:  "Invokes the specified handler when a user navigates to the page.",
	},
	"onpopstate": {
		Name: "OnPopState",
		Doc:  "Executes the provided handler when changes are made to the window's history.",
	},
	"onresize": {
		Name: "OnResize",
		Doc:  "Triggers the given handler upon resizing the browser window.",
	},
	"onstorage": {
		Name: "OnStorage",
		Doc:  "Invokes the specified handler when a Web Storage area undergoes updates.",
	},
	"onunload": {
		Name: "OnUnload",
		Doc:  "Executes the provided handler once the page has been unloaded or the browser window closes.",
	},

	// Form events:
	"onblur": {
		Name: "OnBlur",
		Doc:  "Executes the given handler when the element loses focus.",
	},
	"onchange": {
		Name: "OnChange",
		Doc:  "Triggers the specified handler when the element's value changes.",
	},
	"oncontextmenu": {
		Name: "OnContextMenu",
		Doc:  "Invokes the provided handler upon activation of a context menu.",
	},
	"onfocus": {
		Name: "OnFocus",
		Doc:  "Executes the given handler when the element receives focus.",
	},
	"oninput": {
		Name: "OnInput",
		Doc:  "Triggers the specified handler when the element receives user input.",
	},
	"oninvalid": {
		Name: "OnInvalid",
		Doc:  "Invokes the provided handler when the element is determined to be invalid.",
	},
	"onreset": {
		Name: "OnReset",
		Doc:  "Executes the given handler upon clicking the Reset button within a form.",
	},
	"onsearch": {
		Name: "OnSearch",
		Doc:  "Triggers the specified handler when input is provided in a search field.",
	},
	"onselect": {
		Name: "OnSelect",
		Doc:  "Invokes the provided handler after text within the element is selected.",
	},
	"onsubmit": {
		Name: "OnSubmit",
		Doc:  "Executes the given handler when the form undergoes submission.",
	},

	// Keyboard events:
	"onkeydown": {
		Name: "OnKeyDown",
		Doc:  "Executes the specified handler when a user starts pressing a key.",
	},
	"onkeypress": {
		Name: "OnKeyPress",
		Doc:  "Triggers the provided handler as a key is pressed by the user.",
	},
	"onkeyup": {
		Name: "OnKeyUp",
		Doc:  "Invokes the given handler when a user releases a key.",
	},

	// Mouse events:
	"onclick": {
		Name: "OnClick",
		Doc:  "Triggers the specified handler upon a mouse click on the element.",
	},
	"ondblclick": {
		Name: "OnDblClick",
		Doc:  "Executes the provided handler when the element is double-clicked by the mouse.",
	},
	"onmousedown": {
		Name: "OnMouseDown",
		Doc:  "Invokes the given handler as a mouse button is pressed on the element.",
	},
	"onmouseenter": {
		Name: "OnMouseEnter",
		Doc:  "Triggers the specified handler when the mouse pointer first enters the element's boundaries.",
	},
	"onmouseleave": {
		Name: "OnMouseLeave",
		Doc:  "Executes the provided handler when the mouse pointer leaves the element and its descendants.",
	},
	"onmousemove": {
		Name: "OnMouseMove",
		Doc:  "Invokes the given handler as the mouse pointer moves across the element.",
	},
	"onmouseout": {
		Name: "OnMouseOut",
		Doc:  "Triggers the specified handler when the mouse pointer exits the element.",
	},
	"onmouseover": {
		Name: "OnMouseOver",
		Doc:  "Executes the provided handler as the mouse pointer hovers over the element.",
	},
	"onmouseup": {
		Name: "OnMouseUp",
		Doc:  "Invokes the given handler when a mouse button is released above the element.",
	},
	"onwheel": {
		Name: "OnWheel",
		Doc:  "Triggers the specified handler as the mouse wheel scrolls over the element.",
	},

	// Drag events:
	"ondrag": {
		Name: "OnDrag",
		Doc:  "Executes the handler as an element is being dragged.",
	},
	"ondragend": {
		Name: "OnDragEnd",
		Doc:  "Invokes the handler at the conclusion of a drag operation.",
	},
	"ondragenter": {
		Name: "OnDragEnter",
		Doc:  "Triggers the handler when an element is dragged onto a valid drop target.",
	},
	"ondragleave": {
		Name: "OnDragLeave",
		Doc:  "Invokes the handler when an element exits a valid drop target.",
	},
	"ondragover": {
		Name: "OnDragOver",
		Doc:  "Executes the handler as an element is dragged over a valid drop target.",
	},
	"ondragstart": {
		Name: "OnDragStart",
		Doc:  "Triggers the handler at the initiation of a drag operation.",
	},
	"ondrop": {
		Name: "OnDrop",
		Doc:  "Invokes the handler when a dragged element is released onto a drop target.",
	},
	"onscroll": {
		Name: "OnScroll",
		Doc:  "Executes the handler as an element's scrollbar is scrolled.",
	},

	// Clipboard events:
	"oncopy": {
		Name: "OnCopy",
		Doc:  "Triggers the handler when content of an element is copied by the user.",
	},
	"oncut": {
		Name: "OnCut",
		Doc:  "Executes the handler when the user cuts content from an element.",
	},
	"onpaste": {
		Name: "OnPaste",
		Doc:  "Invokes the handler as content is pasted into an element by the user.",
	},

	// Media events:
	"onabort": {
		Name: "OnAbort",
		Doc:  "Triggers the handler when media loading is aborted.",
	},
	"oncanplay": {
		Name: "OnCanPlay",
		Doc:  "Executes the handler when media has buffered sufficiently to begin playback.",
	},
	"oncanplaythrough": {
		Name: "OnCanPlayThrough",
		Doc:  "Invokes the handler when media can be played through without buffering interruptions.",
	},
	"oncuechange": {
		Name: "OnCueChange",
		Doc:  "Triggers the handler upon cue changes within a track element.",
	},
	"ondurationchange": {
		Name: "OnDurationChange",
		Doc:  "Executes the handler when the media's duration changes.",
	},
	"onemptied": {
		Name: "OnEmptied",
		Doc:  "Invokes the handler when media unexpectedly becomes unavailable.",
	},
	"onended": {
		Name: "OnEnded",
		Doc:  "Triggers the handler when media playback reaches the end.",
	},
	"onloadeddata": {
		Name: "OnLoadedData",
		Doc:  "Executes the handler as media data finishes loading.",
	},
	"onloadedmetadata": {
		Name: "OnLoadedMetaData",
		Doc:  "Invokes the handler when metadata (like duration and dimensions) are fully loaded.",
	},
	"onloadstart": {
		Name: "OnLoadStart",
		Doc:  "Triggers the handler when media loading commences.",
	},
	"onpause": {
		Name: "OnPause",
		Doc:  "Executes the handler when media playback is paused.",
	},
	"onplay": {
		Name: "OnPlay",
		Doc:  "Invokes the handler when media starts its playback.",
	},
	"onplaying": {
		Name: "OnPlaying",
		Doc:  "Triggers the handler once the media has initiated playback.",
	},
	"onprogress": {
		Name: "OnProgress",
		Doc:  "Executes the handler while the browser fetches media data.",
	},
	"onratechange": {
		Name: "OnRateChange",
		Doc:  "Invokes the handler when playback rate changes (e.g., slow motion or fast forward).",
	},
	"onseeked": {
		Name: "OnSeeked",
		Doc:  "Triggers the handler post seeking completion.",
	},
	"onseeking": {
		Name: "OnSeeking",
		Doc:  "Executes the handler during the seeking process.",
	},
	"onstalled": {
		Name: "OnStalled",
		Doc:  "Invokes the handler when media data fetching stalls.",
	},
	"onsuspend": {
		Name: "OnSuspend",
		Doc:  "Triggers the handler when media data fetching is suspended.",
	},
	"ontimeupdate": {
		Name: "OnTimeUpdate",
		Doc:  "Executes the handler when the media's playback position changes.",
	},
	"onvolumechange": {
		Name: "OnVolumeChange",
		Doc:  "Invokes the handler upon volume changes or muting.",
	},
	"onwaiting": {
		Name: "OnWaiting",
		Doc:  "Triggers the handler when media pauses, awaiting further buffering.",
	},

	// Misc events:
	"ontoggle": {
		Name: "OnToggle",
		Doc:  "Executes the handler when the details element is toggled by the user.",
	},
}

func eventHandlersByName(names ...string) []eventHandler {
	res := make([]eventHandler, 0, len(names))
	for _, n := range names {
		h, ok := eventHandlers[n]
		if !ok {
			panic("unknown event handler: " + n)
		}
		res = append(res, h)
	}

	sort.Slice(res, func(i, j int) bool {
		return strings.Compare(res[i].Name, res[j].Name) <= 0
	})

	return res
}

func withGlobalEventHandlers(handlers ...eventHandler) []eventHandler {
	handlers = append(handlers, eventHandlersByName(
		"onblur",
		"onchange",
		"oncontextmenu",
		"onfocus",
		"oninput",
		"oninvalid",
		"onreset",
		"onsearch",
		"onselect",
		"onsubmit",

		"onkeydown",
		"onkeypress",
		"onkeyup",

		"onclick",
		"ondblclick",
		"onmousedown",
		"onmouseenter",
		"onmouseleave",
		"onmousemove",
		"onmouseout",
		"onmouseover",
		"onmouseup",
		"onwheel",

		"ondrag",
		"ondragend",
		"ondragenter",
		"ondragleave",
		"ondragover",
		"ondragstart",
		"ondrop",
		"onscroll",

		"oncopy",
		"oncut",
		"onpaste",
	)...)

	sort.Slice(handlers, func(i, j int) bool {
		return strings.Compare(handlers[i].Name, handlers[j].Name) <= 0
	})

	return handlers
}

func withMediaEventHandlers(handlers ...eventHandler) []eventHandler {
	handlers = append(handlers, eventHandlersByName(
		"onabort",
		"oncanplay",
		"oncanplaythrough",
		"oncuechange",
		"ondurationchange",
		"onemptied",
		"onended",
		"onerror",
		"onloadeddata",
		"onloadedmetadata",
		"onloadstart",
		"onpause",
		"onplay",
		"onplaying",
		"onprogress",
		"onratechange",
		"onseeked",
		"onseeking",
		"onstalled",
		"onsuspend",
		"ontimeupdate",
		"onvolumechange",
		"onwaiting",
	)...)

	sort.Slice(handlers, func(i, j int) bool {
		return strings.Compare(handlers[i].Name, handlers[j].Name) <= 0
	})

	return handlers
}

func main() {
	generateHTMLGo()
	generateHTMLTestGo()
}

func generateHTMLGo() {
	f, err := os.Create("html_gen.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "package app")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(f, `
import (
	"strings"
)
		`)

	for _, t := range tags {
		writeInterface(f, t)

		switch t.Name {
		case "Elem", "ElemSelfClosing":
			fmt.Fprintf(f, `
			// %s returns an HTML element that %s
			func %s(tag string) HTML%s {
				e := &html%s{
					htmlElement: htmlElement{
						tag: tag,
						isSelfClosing: %v,
					},
				}

				return e
			}
			`,
				t.Name,
				t.Doc,
				t.Name,
				t.Name,
				t.Name,
				t.Type == selfClosing,
			)

		default:
			fmt.Fprintf(f, `
			// %s returns an HTML element that %s
			func %s() HTML%s {
				e := &html%s{
					htmlElement: htmlElement{
						tag: "%s",
						isSelfClosing: %v,
					},
				}

				return e
			}
			`,
				t.Name,
				t.Doc,
				t.Name,
				t.Name,
				t.Name,
				strings.ToLower(t.Name),
				t.Type == selfClosing,
			)
		}

		fmt.Fprintln(f)
		fmt.Fprintln(f)
		writeStruct(f, t)
		fmt.Fprintln(f)
		fmt.Fprintln(f)
	}

}

func writeInterface(w io.Writer, t tag) {
	fmt.Fprintf(w, `
		// The interface that represents a "%s" HTML element.
		type HTML%s interface {
			UI
		`,
		strings.ToLower(t.Name),
		t.Name,
	)

	switch t.Type {
	case parent:
		fmt.Fprintf(w, `
			// Sets the content of the element.
			Body(elems ...UI) HTML%s 
		`, t.Name)

		fmt.Fprintf(w, `
			// Sets the content of the element with a text node containing the stringified given value.
			Text(v any) HTML%s
		`, t.Name)

		fmt.Fprintf(w, `
			// Sets the content of the element with a text node formatted according to a format specifier.
			Textf(format string, v ...any) HTML%s
		`, t.Name)

	case privateParent:
		fmt.Fprintf(w, `
			privateBody(elems ...UI) HTML%s 
		`, t.Name)
	}

	for _, a := range t.Attrs {
		fmt.Fprintln(w)
		fmt.Fprintln(w)

		fmt.Fprintf(w, "// %s\n", a.Doc)
		writeAttrFunction(w, a, t, true)
	}

	fmt.Fprintln(w)

	fmt.Fprintf(w, `
		// Invokes the specified handler when the corresponding event is triggered.
		On(event string, h EventHandler, scope ...any) HTML%s 
	`, t.Name)

	for _, e := range t.EventHandlers {
		fmt.Fprintln(w)
		fmt.Fprintln(w)

		fmt.Fprintf(w, "// %s\n", e.Doc)
		writeEventFunction(w, e, t, true)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "}")
}

func writeStruct(w io.Writer, t tag) {
	fmt.Fprintf(w, `type html%s struct {
			htmlElement
		}`, t.Name)

	switch t.Type {
	case parent:
		fmt.Fprintf(w, `
			func (e *html%s) Body(v ...UI) HTML%s {
				e.setChildren(v...)
				return e
			}
			`,
			t.Name,
			t.Name,
		)

		if t.Name == "Textarea" {
			fmt.Fprintf(w, `
			func (e *html%s) Text(v any) HTML%s {
				e.setAttr("value", v)
				return e
			}
			`,
				t.Name,
				t.Name,
			)
			fmt.Fprintf(w, `
			func (e *html%s) Textf(format string, v ...any) HTML%s {
				e.setAttr("value", FormatString(format, v...))
				return e
			}
			`,
				t.Name,
				t.Name,
			)
		} else {
			fmt.Fprintf(w, `
			func (e *html%s) Text(v any) HTML%s {
				return e.Body(Text(v))
			}
			`,
				t.Name,
				t.Name,
			)

			fmt.Fprintf(w, `
			func (e *html%s) Textf(format string, v ...any) HTML%s {
				return e.Body(Textf(format, v...))
			}
			`,
				t.Name,
				t.Name,
			)
		}

	case privateParent:
		fmt.Fprintf(w, `
			func (e *html%s) privateBody(v ...UI) HTML%s {
				e.setChildren(v...)
				return e
			}
			`,
			t.Name,
			t.Name,
		)
	}

	for _, a := range t.Attrs {
		fmt.Fprintln(w)
		fmt.Fprintln(w)

		writeAttrFunction(w, a, t, false)
	}

	fmt.Fprintln(w)

	fmt.Fprintf(w, `
		func (e *html%s) On(event string, h EventHandler, scope ...any)  HTML%s {
			e.setEventHandler(event, h, scope...)
			return e
		}
		`,
		t.Name,
		t.Name,
	)

	for _, e := range t.EventHandlers {
		fmt.Fprintln(w)
		fmt.Fprintln(w)

		writeEventFunction(w, e, t, false)
	}
}

func writeAttrFunction(w io.Writer, a attr, t tag, isInterface bool) {
	if !isInterface {
		fmt.Fprintf(w, "func (e *html%s)", t.Name)
	}

	var attrName string
	if a.NameOverride != "" {
		attrName = strings.ToLower(a.NameOverride)
	} else {
		attrName = strings.ToLower(a.Name)
	}

	switch a.Type {
	case "data|value":
		fmt.Fprintf(w, `%s(k string, v any) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("data-"+k, FormatString("%s", v))
				return e
			}`, "%v")
		}

	case "data|map":
		fmt.Fprintf(w, `%s(ds map[string]any) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				for k, v := range ds {
					e.DataSet(k, v)
				}
				return e
			}`)
		}

	case "attr|value":
		fmt.Fprintf(w, `%s(n string, v any) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr(n, v)
				return e
			}`)
		}

	case "aria|value":
		fmt.Fprintf(w, `%s(k string, v any) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("aria-"+k, FormatString("%s", v))
				return e
			}`, "%v")
		}

	case "style":
		fmt.Fprintf(w, `%s(k, format string, v ...any) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("style", k+":"+FormatString(format, v...))
				return e
			}`)
		}

	case "style|map":
		fmt.Fprintf(w, `%s(s map[string]string) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				for k, v := range s {
					e.Style(k, v)
				}
				return e
			}`)
		}

	case "on/off":
		fmt.Fprintf(w, `%s(v bool) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				s := "off"
				if (v) {
					s = "on"
				}
	
				e.setAttr("%s", s)
				return e
			}`, attrName)
		}

	case "bool|force":
		fmt.Fprintf(w, `%s(v bool) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				s := "false"
				if (v) {
					s = "true"
				}
	
				e.setAttr("%s", s)
				return e
			}`, attrName)
		}

	case "string|class":
		fmt.Fprintf(w, `%s(v ...string) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("%s", strings.Join(v, " "))
				return e
			}`, attrName)
		}

	case "xmlns":
		fmt.Fprintf(w, `%s(v string) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintln(w, `{
				e.xmlns = v
				return e
			}`)
		}

	case "fmt":
		fmt.Fprintf(w, `%s(format string, v ...any) HTML%s`, a.Name, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("%s", FormatString(format, v...))
				return e
			}`, attrName)
		}

	default:
		fmt.Fprintf(w, `%s(v %s) HTML%s`, a.Name, a.Type, t.Name)
		if !isInterface {
			fmt.Fprintf(w, `{
				e.setAttr("%s", v)
				return e
			}`, attrName)
		}
	}
}

func writeEventFunction(w io.Writer, e eventHandler, t tag, isInterface bool) {
	if !isInterface {
		fmt.Fprintf(w, `func (e *html%s)`, t.Name)
	}

	fmt.Fprintf(w, `%s (h EventHandler, scope ...any) HTML%s`, e.Name, t.Name)
	if !isInterface {
		fmt.Fprintf(w, `{
			e.setEventHandler("%s", h, scope...)
			return e
		}`, strings.TrimPrefix(strings.ToLower(e.Name), "on"))
	}
}

func generateHTMLTestGo() {
	f, err := os.Create("html_gen_test.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "package app")
	fmt.Fprintln(f)
	fmt.Fprintln(f, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintln(f, `
import (
	"testing"
)
		`)

	for _, t := range tags {
		fmt.Fprintln(f)
		fmt.Fprintf(f, `func Test%s(t *testing.T) {`, t.Name)
		fmt.Fprintln(f)

		switch t.Name {
		case "Elem", "ElemSelfClosing":
			fmt.Fprintf(f, `elem := %s("div")`, t.Name)

		default:
			fmt.Fprintf(f, `elem := %s()`, t.Name)
		}

		fmt.Fprintln(f)

		for _, a := range t.Attrs {
			fmt.Fprintf(f, `elem.%s(`, a.Name)

			switch a.Type {
			case "data|value", "aria|value", "attr|value":
				fmt.Fprintln(f, `"foo", "bar")`)

			case "data|map":
				fmt.Fprintln(f, `map[string]any{"foo": "bar"})`)

			case "style":
				fmt.Fprintln(f, `"margin", "%vpx", 42)`)

			case "style|map":
				fmt.Fprintln(f, `map[string]string{"color": "pink"})`)

			case "bool", "bool|force", "on/off":
				fmt.Fprintln(f, `true)`)
				fmt.Fprintf(f, `elem.%s(false)`, a.Name)
				fmt.Fprintln(f)

			case "int":
				fmt.Fprintln(f, `42)`)

			case "string":
				fmt.Fprintln(f, `"foo")`)

			case "url":
				fmt.Fprintln(f, `"http://foo.com")`)

			case "string|class":
				fmt.Fprintln(f, `"foo bar")`)

			case "xmlns":
				fmt.Fprintln(f, `"http://www.w3.org/2000/svg")`)

			case "fmt":
				fmt.Fprintln(f, `"hello %v", 42)`)

			default:
				fmt.Fprintln(f, `42)`)
			}
		}

		if len(t.EventHandlers) != 0 {
			fmt.Fprint(f, `
				h := func(ctx Context, e Event) {}
			`)
			fmt.Fprintf(f, `elem.On("click", h)`)
			fmt.Fprintln(f)
		}

		for _, e := range t.EventHandlers {
			fmt.Fprintf(f, `elem.%s(h)`, e.Name)
			fmt.Fprintln(f)
		}

		switch t.Type {
		case parent:
			fmt.Fprintln(f, `elem.Text("hello")`)
			fmt.Fprintln(f, `elem.Textf("hello %s", "Maxence")`)

		case privateParent:
			fmt.Fprintln(f, `elem.privateBody(Text("hello"))`)
		}

		fmt.Fprintln(f, "}")
	}
}
