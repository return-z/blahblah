// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Login() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"auth\" class=\"h-screen flex justify-center items-center\"><div class=\"w-full flex flex-col justify-center items-center px-4\"><div class=\"w-full max-w-md bg-gray-900 p-8 rounded-lg shadow-md\"><h2 class=\"text-red-500 text-3xl font-bold mb-4 text-center\">Chat Application</h2><div id=\"form-div\"><form id=\"form\" hx-post=\"/auth\" class=\"flex flex-col items-center\"><div class=\"relative w-full text-white\"><input type=\"text\" id=\"username\" name=\"username\" class=\"bg-black w-full pl-2 pr-10 py-2 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-blink\" placeholder=\"Join room as\" autocomplete=\"off\" required></div><button type=\"submit\" class=\"w-full bg-green-600 mt-6 py-2 rounded text-white hover:bg-green-700 focus:outline-none focus:bg-green-700\">Submit</button></form><button hx-get=\"/register\" type=\"submit\" class=\"w-full bg-green-600 mt-4 py-2 rounded text-white hover:bg-green-700 focus:outline-none focus:bg-green-700\">Register</button></div></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
