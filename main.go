package main

import(
	"github.com/marcusolsson/tui-go"
	"time"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type hitRow struct{
	hits,url  *tui.Entry

}

func main(){
	hits := tui.NewTable(0,0)
	hits.SetColumnStretch(0, 1)
	hits.SetColumnStretch(1, 1)
	hits.SetColumnStretch(2, 4)
	hits.AppendRow(
		tui.NewLabel("HITS"),
		tui.NewLabel("URL"),
	)
	rows := make([]hitRow, 10)
	for i:=0;i < 10;i++{
		r := hitRow{tui.NewEntry(),tui.NewEntry()}
		r.hits.SetSizePolicy(tui.Minimum,tui.Expanding)
		r.url.SetSizePolicy(tui.Minimum,tui.Expanding)
		hits.AppendRow(r.hits,r.url)
		rows[i] = r
	}
	hits.SetSizePolicy(tui.Minimum,tui.Expanding)

	root := tui.NewVBox(hits)
	ui, err := tui.New(root)
	if err != nil{
		log.Fatal(err)
	}
	ui.SetKeybinding("Esc", func() { ui.Quit() })
	ui.SetKeybinding("q", func() { ui.Quit() })
	//table.
	go loop(rows, ui)
	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}


type jsonStruct struct{
	Hits int `json:"hits"`
	Target string `json:"target"`
	Key string `json:"key"`
}
func loop(rows []hitRow, ui tui.UI){
	for{
		res, err := http.Get("http://46.101.97.8/api/v1/list")
		if err != nil{
			return
		}
		body, err := ioutil.ReadAll(res.Body)

		var hits []jsonStruct
		err = json.Unmarshal(body,&hits)
		if err != nil{
			return
		}

		for i, h := range hits{
			rows[i].hits.SetText(fmt.Sprint(h.Hits))
			rows[i].url.SetText(fmt.Sprint(h.Target))
		}
		ui.Update(func(){})
		time.Sleep(time.Second)
	}
}
