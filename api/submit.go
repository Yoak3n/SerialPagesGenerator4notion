package api


import(
	"b2n3/backend/model"
)


var datas []*model.Data



func SumbitVideo(){
	initVideoBody()
	

}



func initVideoBody(){
	parent := &model.Parent{
		Type: "database_id",
		DatabaseID: "b2n3",
	}


	for episode, episodeName := range video.Titles {
		properties := &model.Properties{
			Episode:  model.Episode{
				Number: episode + 1 ,
			},
			EpisodeName: *genEpisodeName(&episodeName),
			Name: model.Name{
				Select: struct{Name string "json:\"name\""}{
					Name: video.Name,
				},
			},
		}

		data := &model.Data{
			Parent: *parent,
			Properties: *properties,
		}

		datas = append(datas,data)
	}

}

func genEpisodeName(name *string)*model.EpisodeName{
	titles := make([]model.Title,0)
	title := &model.Title{
		Type: "text",
		Text: struct{Content string "json:\"content\""}{
			Content: *name,
		},
	}
	titles = append(titles,*title)

	episodeName := &model.EpisodeName{
		Title: titles,
	}
	
	return episodeName
}