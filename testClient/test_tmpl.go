package main

import (
	"html/template"
	"os"
)

type Item struct {
	Ask     bool
	Content string
}

func main() {
	tpl, err := template.New("gptmail").Parse(mailTmpl)
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(os.Stderr, []Item{
		Item{
			Ask:     true,
			Content: "你是谁",
		},
		Item{
			Content: "我是                                        GPT",
		},
		Item{
			Ask:     true,
			Content: "GPT是啥",
		},
	})
	if err != nil {
		panic(err)
	}

}

var (
	mailTmpl = `<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<div>
	  <style>
		.tdl {
		  width: 40px;
		}
		.tdc {
		  background-color: rgb(0, 255, 38);
		}
	  </style>
	  <div>以下回复仅供参考，</div>
	  <table border="1px" style="border-collapse: collapse">
		<tr>
		  <td class="tdl">ASK</td>
		  <td>会话内容</td>
		  <td class="tdl">GPT</td>
		</tr>
		{{range .}}
		<tr>
		  {{ if .Ask }}
		  <td>问题</td>
		  <td>{{.Content}}</td>
		  <td></td>
		  {{else}}
		   <td></td>
		   <td class="tdc">{{.Content}}</td>
		   <td>回答</td>
		  {{end}}
		</tr>
		{{end}}
	  </table>
	</div>
	`
)
