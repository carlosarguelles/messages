// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func CreateChat() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html lang=\"en\"><head><meta charset=\"UTF-8\"><title>Chat</title><script src=\"https://cdn.tailwindcss.com\"></script><script src=\"https://unpkg.com/htmx.org@1.9.10\"></script></head><body class=\"h-screen antialiased\"><div class=\"container mx-auto max-w-4xl\"><h1 class=\"text-xl text-center p-4 font-bold\">Create Chat</h1><form class=\"m-0\" method=\"POST\"><div class=\"flex items-center gap-3\"><label for=\"name\">Chat Name</label> <input id=\"name\" type=\"text\" name=\"name\" class=\"flex-1 border bg-gray-100 py-1 px-2 border border-gray-400 rounded-lg\"> <button type=\"submit\" class=\"bg-green-500 font-medium rounded-lg py-1 px-2 text-white border-green-700\">Create Chat</button></div></form></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
