// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package assets

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Chat() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"auth\" hx-swap-oob=\"true\" class=\"app max-w-screen-lg mx-auto text-white\"><script>\n     document.body.addEventListener('htmx:oobAfterSwap', function(evt) {\n        const form = document.querySelector(\"#form\");\n        form.reset();\n    });\n  </script><header class=\"bg-gray-900 px-6 py-4\"><h1 class=\"text-red-500 text-2xl font-bold text-center\">CHAT APPLICATION</h1></header><div id=\"messages\" class=\"message-container bg-black px-6 py-4\"><p class=\"text-blue-400 font-mono p-1\">Type <code>!join chatroom_name</code> to join a chatroom</p></div><footer class=\"bg-gray-900 px-6 py-4\"><div id=\"ws-form\" hx-ws=\"connect:/ws\"><form id=\"form\" hx-ws=\"send:submit\" class=\"flex\"><input name=\"message\" class=\"w-full rounded-l-md bg-gray-800 border-2 border-gray-700 focus:border-blue-500 focus:outline-none py-2 px-4\" placeholder=\"Type your message...\"> <button type=\"submit\" class=\"bg-transparent px-4 border border-solid border-green-500 rounded-r-md hover:bg-transparent focus:outline-none\"><svg class=\"text-green h-6 w-6\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M14 5l7 7m0 0l-7 7m7-7H3\"></path></svg></button></form></div></footer></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}