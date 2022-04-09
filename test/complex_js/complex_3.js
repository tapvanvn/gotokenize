var design ={
	quizz :{
		timeout : 60 ,
		game_duration : 30 ,
		popup :{
			panel :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 - 25 ,
				img_name : "quizz.popup",
				scale :{
					x : 1.0 ,
					y : 2.0
				}
			},
			panel_glow :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 ,
				img_name : "quizz.popup_glow",
			},
			tutorial :{
				panel_content :{
					x : config.design_w / 2 ,
					y : config.design_h / 2 - 75 ,
					img_name : "quizz.content_tutorial",
				},
				panel_btn_close :{
					x : config.design_w / 2 + 225 ,
					y : config.design_h / 2 - 335 ,
					img_name : "quizz.btn_close"
				}
			}
		},
		btn_play :{
			x : config.design_w / 2 ,
			y : config.design_h - 175 ,
			img_name : "quizz.btn_play",
		},
		btn_tutorial :{
			x : config.design_w / 2 ,
			y : config.design_h - 75 ,
			img_name : "quizz.btn_tutorial",
			scale :{
				x : 1.4 ,
				y : 1.4
			}
		},
		avatar :{
			x : 150 ,
			y : 90 ,
			img_name : "quizz.avatar"
		},
		avatar_img :{
			x : 57 ,
			y : 90 ,
			w : 85 ,
			h : 85 ,
			img_name : "avatar"
		},
		timeline_bg :{
			x : config.design_w / 2 ,
			y : timeline_base_y ,
			img_name : "quizz.timeline_2"
		},
		timeline :{
			x : config.design_w / 2 ,
			y : timeline_base_y ,
			img_name : "quizz.timeline_1"
		},
		time :{
			x : config.design_w - 40 ,
			y : timeline_base_y + 30
		},
		question :{
			x : 50 ,
			vid_h : 350 ,
			y : timeline_base_y + 75 ,
			w : config.design_w - 100
		}
	}
}