
var config ={
	design_w : 720 ,
	design_h : 1480 ,
	design_pc_w_scale : 1 ,
}
var design ={
	shaking :{
		timeout : 60 ,
		shake_duration : 11 ,
		anim_scale : 1.2 ,
		popup :{
			panel :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 - 25 ,
				img_name : "shaking.popup",
				scale :{
					x : 1.0 ,
					y : 2.0
				}
			},
			panel_glow :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 ,
				img_name : "shaking.popup_glow",
			},
			tutorial :{
				panel_content :{
					x : config.design_w / 2 ,
					y : config.design_h / 2 - 75 ,
					img_name : "shaking.content_tutorial",
				},
				panel_btn_close :{
					x : config.design_w / 2 + 225 ,
					y : config.design_h / 2 - 335 ,
					img_name : "shaking.btn_close"
				}
			}
		},
		tutorial :{
			avatar :{
				x : 150 ,
				y : 90 ,
				img_name : "shaking.avatar"
			},
			avatar_img :{
				x : 57 ,
				y : 90 ,
				w : 85 ,
				h : 85 ,
				img_name : "avatar"
			},
			title :{
				x : config.design_w / 2 ,
				y : 250 ,
				img_name : "shaking.title_shaking"
			},
			count :{
				x : config.design_w / 2 ,
				y : config.design_h - 225 ,
				img_name : "shaking.count3.3"
			},
			count_s :{
				org_x : config.design_w / 2 + 125 ,
				x : config.design_w / 2 + 125 ,
				y : config.design_h - 185 ,
				img_name : "shaking.count10.s"
			},
			anim :{
				x : config.design_w / 2 ,
				y : config.design_h
			},
			anim_shaking_now :{
				x : config.design_w / 2 ,
				y : config.design_h - 300
			},
			anim_narrow :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 + 30
			},
			panel :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 - 25 ,
				img_name : "shaking.popup",
				scale :{
					x : 1.0 ,
					y : 1.3
				}
			},
			panel_content :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 - 75 ,
				img_name : "shaking.content_tutorial",
			},
			panel_btn_close :{
				x : config.design_w / 2 + 225 ,
				y : config.design_h / 2 - 335 ,
				img_name : "shaking.btn_close"
			},
			panel_btn_start :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 + 225 ,
				img_name : "shaking.btn_tutorial_start"
			},
			popup_return :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 ,
				img_name : "shaking.popup_return"
			},
			btn_start :{
				x : config.design_w / 2 ,
				y : config.design_h - 175 ,
				img_name : "shaking.btn_start"
			},
			num_play :{
				x : config.design_w / 2 ,
				y : config.design_h - 105 ,
			},
			btn_tutorial :{
				x : config.design_w / 2 ,
				y : config.design_h - 50 ,
				img_name : "shaking.title_tutorial"
			}
		},
		play :{
			glow :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 ,
				img_name : "shaking.glow"
			},
			title_shaking :{
				x : config.design_w / 2 ,
				y : 230 ,
				img_name : "shaking.title_shaking"
			},
			gift_box :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 ,
				img_name : "shaking.gift_box"
			},
			partical_1 :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 ,
				img_name : "shaking.partical_1"
			},
			partical_2 :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 ,
				img_name : "shaking.partical_2"
			},
			btn_shaking :{
				x : config.design_w / 2 ,
				y : config.design_h / 2 + 200 ,
				img_name : "shaking.btn_shaking"
			},
		}
	}
}
var background_src = "res/shaking/background.jpg"
var bg_music = "shaking.bg_music"
var game_id = 1
var game_state ={
	title : "Họ tên",
	num_play : 3 ,
	avatar : "res/shaking/fake_avatar.jpg"
}
var activeScene = null
var activePopup = null
var is_start_journey = false
var is_start_shaking = false
var last_lock_touch = 0
var session = "748b7b8f-096e-4701-ac92-d13278406188"
var session_type = "unknown"
function changeScene( scene ){
	touchCleanAll ()
	if( activeScene != null && typeof activeScene.touchHelper == 'object'){
		last_lock_touch = activeScene.touchHelper.lock_time + 200
		activeScene.touchHelper.clean ()
	}
	if( scene != null && typeof scene.touchHelper == 'object'){
		scene.touchHelper.lock_time = last_lock_touch
		scene.touchHelper.clean ()
	}
	if( scene != null && typeof scene.name == 'string'){
		console.log ( "CHANGE SCENE", scene.name )
	}
	else{
		console.log ( "CHANGE SCENE untitle")
	}
	activeScene = scene
	if( scene != null && typeof scene.begin == 'function'){
		scene.begin ()
	}
}
function popupScene( scene ){
	touchCleanAll ()
	if( activePopup != null && typeof activePopup.touchHelper == 'object'){
		last_lock_touch = activePopup.touchHelper.lock_time + 200
		activePopup.touchHelper.clean ()
	}
	if( scene != null && typeof scene.touchHelper == 'object'){
		scene.touchHelper.lock_time = last_lock_touch
		scene.touchHelper.clean ()
	}
	if( scene != null && typeof scene.name == 'string'){
		console.log ( "POPUP SCENE", scene.name )
	}
	else{
		console.log ( "POPUP SCENE untitle")
	}
	if( typeof registry.gui_container != 'undefined'){
		registry.gui_container.className = 'disable'
	}
	activePopup = scene
	if( scene != null && typeof scene.begin == 'function'){
		scene.begin ()
	}
}
function closePopup(){
	if( activePopup != null && typeof activePopup.touchHelper == 'object'){
		last_lock_touch = activePopup.touchHelper.lock_time + 200
		activePopup.touchHelper.clean ()
	}
	touchCleanAll ()
	if( typeof registry.gui_container != 'undefined'){
		registry.gui_container.className = ''
	}
	activePopup = null
	if( typeof registry.gui_container != 'undefined'){
		registry.gui_container.style.display = 'block'
	}
}
function popupError(){
	Popup.hasClose = true
	setPopupContentFunc ( 400 , 200 ,
	( x , y )=>{
		var font = 36
		var ctx = registry.context
		ctx.fillStyle = '#666666'
		ctx.font = font + "pt fontAlbert"
		drawTextBox ( "Có lỗi xảy ra vui lòng quay lại để làm mới", x - 250 , y , 500 , font + 15 ,{
			halign : 'center',
		})
	}
	,[
	{
		img_name : "popup.btn_back",
		callback :()=>{
			closePopup ()
			goBack ()
		}
	}
	])
	popupScene ( Popup )
}
function login( callback ){
	if( session == "dev"){
		callback ({
			id : "test",
			profile :{
				full_name : "Họ và tên",
				avatar : "https://interactive-examples.mdn.mozilla.net/media/cc0-images/grapefruit-slice-332-332.jpg"
			},
			game_turns :[
			{
				game_id : 1 ,
				num_turn : 1
			},
			{
				game_id : 2 ,
				num_turn : 1
			}
			]
		})
		return
	}
	var xmlhttp = new XMLHttpRequest ()
	xmlhttp.timeout = 2000
	var path = window.server_url + "/api/v1/games/user-info?type="+ session_type
	xmlhttp.onreadystatechange = function (){
		if( this.readyState == 4 ){
			if( typeof callback === "function"){
				if( this.status == 200 || this.status == 201 ){
					try {
						callback ( JSON.parse ( this.responseText ))
					} catch ( ex ){
						callback ( null )
					}
				}
				else{
					callback ( null )
				}
			}
		}
	}
	xmlhttp.ontimeout = function ( e ){
		callback ( null )
	}
	try {
		xmlhttp.open ( "GET", path , true )
	} catch ( ex ){
		if( typeof callback === "function"){
			callback ( null )
		}
	}
	xmlhttp.setRequestHeader ( "session", session )
	xmlhttp.send ()
}
function updateResult( game_id , success , callback ){
	if( session == "dev"){
		callback (
		{
			point : 100 ,
			game_turns :[
			{
				game_id : 1 ,
				num_turn : game_state.num_play - 1
			},
			{
				game_id : 2 ,
				num_turn : game_state.num_play - 1
			}
			]
		}
		)
		return
	}
	var xmlhttp = new XMLHttpRequest ()
	xmlhttp.timeout = 2000
	var path = window.server_url + "/api/v1/games/update-result"
	xmlhttp.onreadystatechange = function (){
		if( this.readyState == 4 ){
			if( typeof callback === "function"){
				if( this.status == 200 || this.status == 201 ){
					try {
						callback ( JSON.parse ( this.responseText ))
					} catch ( ex ){
						callback ( null )
					}
				}
				else{
					callback ( null )
				}
			}
		}
	}
	xmlhttp.ontimeout = function ( e ){
		callback ( null )
	}
	try {
		xmlhttp.open ( "POST", path , true )
	} catch ( ex ){
		if( typeof callback === "function"){
			callback ( null )
		}
	}
	xmlhttp.setRequestHeader ( "session", session )
	xmlhttp.setRequestHeader ( "Content-Type", "application/json")
	var data ={
		game_id : game_id ,
		result : success ? "Win": "Lose",
		type : session_type
	}
	xmlhttp.send ( JSON.stringify ( data ))
}
var __spine_res = new spine.AssetManager ( "")
var __spines ={
}
var __init_skeleton = false
var __spine_renderer = null
function spineCalculateSetupPoseBounds( skeleton ){
	skeleton.setToSetupPose ()
	skeleton.updateWorldTransform ()
	var offset = new spine.Vector2 ()
	var size = new spine.Vector2 ()
	skeleton.getBounds ( offset , size ,[])
	return {
	offset : offset , size : size}
}
var SpineAnimation = class {
	constructor ( skeleton_data , animationName , loop ){
		this.skin = 'default'
		this.skeleton = new spine.Skeleton ( skeleton_data )
		this.skeleton.setSkinByName ( this.skin )
		this.skeleton.scaleY =- 1
		this.bounds = spineCalculateSetupPoseBounds ( this.skeleton )
		this.data = new spine.AnimationStateData ( this.skeleton.data )
		this.state = new spine.AnimationState ( this.data )
		this.state.setAnimation ( 0 , animationName , loop )
	}
	setAnimation ( name , loop ){
		this.state.setAnimation ( 0 , name , loop )
	}
	scale ( x , y ){
		this.skeleton.scaleX = x
		this.skeleton.scaleY =- y
	}
	position ( x , y ){
		this.skeleton.x = x
		this.skeleton.y = y
	}
	angle ( angle ){
		this.skeleton.getRootBone ().rotation = angle
	}
}
var SpineRes = class {
	constructor ( res_name , skel_path ){
		var last_slash = skel_path.lastIndexOf ( "/")
		var container_path = ""
		var name = skel_path
		if( last_slash >= 0 ){
			name = skel_path.substring ( last_slash + 1 )
			container_path = skel_path.substring ( 0 , last_slash )
		}
		name = name.replace ( /(\-ess|\-pro)*\.(skel|json)/g, "")
		console.log ( "spine:", name )
		var ext = skel_path.split ( '.').pop ()
		if( ext != "skel"&& ext != 'json'){
			throw "not support extension"+ ext
		}
		__spine_res.loadTextureAtlas ( container_path + "/"+ name + ".atlas")
		if( ext == "skel"){
			__spine_res.loadBinary ( skel_path )
		}
		else{
			__spine_res.loadText ( skel_path )
		}
		__spines [ res_name ]= this
		this.premultipliedAlpha = false
		this.initialAnimation = ""
		this.name = name
		this.path = skel_path
		this.containerPath = container_path
		this.atlas = null
		this.isBinary =( ext == "skel")
		this.skeletonData = null
		this.loaded = false
	}
	loadSkeleton (){
		var atlas_path = this.containerPath + "/"+ this.name +( this.premultipliedAlpha ? "-pma": "")+ ".atlas"
		this.atlas = __spine_res.require ( atlas_path )
		var atlasLoader = new spine.AtlasAttachmentLoader ( this.atlas )
		var skeletonLoader = this.isBinary ? new spine.SkeletonBinary ( atlasLoader ): new spine.SkeletonJson ( atlasLoader )
		this.skeletonData = skeletonLoader.readSkeletonData ( __spine_res.require ( this.path ))
		this.loaded = true
	}
	newAnimation ( animationName , loop ){
		if(! this.loaded ){
			throw ( "spine skeleton is not loaded yet")
		}
		return new SpineAnimation ( this.skeletonData , animationName , loop )
	}
}
var spineLoadingComplete =
()=>{
	var complete = false
	if( __spine_res.isLoadingComplete ()){
		if(! __init_skeleton ){
			__spine_renderer = new spine.SkeletonRenderer ( registry.context )
			__init_skeleton = true
		}
		else{
			complete = true
			Object.values ( __spines ).forEach ( spineRes =>{
				if(! spineRes.loaded ){
					spineRes.loadSkeleton ()
					complete = false
				}
				else{
					console.log ( "loaded", spineRes )
				}
			})
		}
	}
	return complete
}
var spineUpdate =
( animation , delta )=>{
	animation.state.update ( delta )
}
var spineDraw =
( animation )=>{
	animation.state.apply ( animation.skeleton )
	animation.skeleton.updateWorldTransform ()
	__spine_renderer.draw ( animation.skeleton )
}
function drawRect( rect , color ){
	var ctx = registry.context
	var bk_color = ctx.strokeStyle
	ctx.beginPath ()
	ctx.strokeStyle = color
	ctx.rect ( rect.x , rect.y , rect.w , rect.h )
	ctx.stroke ()
	ctx.strokeStyle = bk_color
}
function fillRect( rect , color ){
	var ctx = registry.context
	var bk_color = ctx.fillStyle
	ctx.beginPath ()
	ctx.fillStyle = color
	ctx.rect ( rect.x , rect.y , rect.w , rect.h )
	ctx.fill ()
	ctx.fillStyle = bk_color
}
function fillAlphaRect( rect , color , alpha ){
	var ctx = registry.context
	var bk_color = ctx.fillStyle
	var bk_alpha = ctx.globalAlpha
	ctx.beginPath ()
	ctx.fillStyle = color
	ctx.globalAlpha = alpha
	ctx.rect ( rect.x , rect.y , rect.w , rect.h )
	ctx.fill ()
	ctx.globalAlpha = bk_alpha
	ctx.fillStyle = bk_color
}
function drawRoundRect( rect , radius , color ){
	if( typeof radius === 'undefined'){
		radius = 5
	}
	if( typeof radius === 'number'){
		radius ={
		tl : radius , tr : radius , br : radius , bl : radius}
	}
	else{
		var defaultRadius ={
		tl : 0 , tr : 0 , br : 0 , bl : 0}
		for( var side in defaultRadius ){
			radius [ side ]= radius [ side ]|| defaultRadius [ side ]
		}
	}
	var ctx = registry.context
	var bk_color = ctx.strokeStyle
	ctx.beginPath ()
	ctx.moveTo ( rect.x + radius.tl , rect.y )
	ctx.lineTo ( rect.x + rect.w - radius.tr , rect.y )
	ctx.quadraticCurveTo ( rect.x + rect.w , rect.y , rect.x + rect.w , rect.y + radius.tr )
	ctx.lineTo ( rect.x + rect.w , rect.y + rect.h - radius.br )
	ctx.quadraticCurveTo ( rect.x + rect.w , rect.y + rect.h , rect.x + rect.w - radius.br , rect.y + rect.h )
	ctx.lineTo ( rect.x + radius.bl , rect.y + rect.h )
	ctx.quadraticCurveTo ( rect.x , rect.y + rect.h , rect.x , rect.y + rect.h - radius.bl )
	ctx.lineTo ( rect.x , rect.y + radius.tl )
	ctx.quadraticCurveTo ( rect.x , rect.y , rect.x + radius.tl , rect.y )
	ctx.closePath ()
	ctx.strokeStyle = color
	ctx.stroke ()
	ctx.strokeStyle = bk_color
}
function fillRoundRect( rect , radius , color ){
	if( typeof radius === 'undefined'){
		radius = 5
	}
	if( typeof radius === 'number'){
		radius ={
		tl : radius , tr : radius , br : radius , bl : radius}
	}
	else{
		var defaultRadius ={
		tl : 0 , tr : 0 , br : 0 , bl : 0}
		for( var side in defaultRadius ){
			radius [ side ]= radius [ side ]|| defaultRadius [ side ]
		}
	}
	var ctx = registry.context
	var bk_color = ctx.strokeStyle
	ctx.beginPath ()
	ctx.moveTo ( rect.x + radius.tl , rect.y )
	ctx.lineTo ( rect.x + rect.w - radius.tr , rect.y )
	ctx.quadraticCurveTo ( rect.x + rect.w , rect.y , rect.x + rect.w , rect.y + radius.tr )
	ctx.lineTo ( rect.x + rect.w , rect.y + rect.h - radius.br )
	ctx.quadraticCurveTo ( rect.x + rect.w , rect.y + rect.h , rect.x + rect.w - radius.br , rect.y + rect.h )
	ctx.lineTo ( rect.x + radius.bl , rect.y + rect.h )
	ctx.quadraticCurveTo ( rect.x , rect.y + rect.h , rect.x , rect.y + rect.h - radius.bl )
	ctx.lineTo ( rect.x , rect.y + radius.tl )
	ctx.quadraticCurveTo ( rect.x , rect.y , rect.x + radius.tl , rect.y )
	ctx.closePath ()
	ctx.fillStyle = color
	ctx.fill ()
	ctx.fillStyle = bk_color
}
function drawCircle( pos , radius , color ){
	var ctx = registry.context
	var bkcolor = ctx.strokeStyle
	ctx.beginPath ()
	ctx.arc ( pos.x , pos.y , radius , 0 , 2 * Math.PI , false )
	ctx.strokeStyle = color
	ctx.stroke ()
	ctx.strokeStyle = bkcolor
}
function fillCircle( pos , radius , color ){
	var ctx = registry.context
	var bkcolor = ctx.fillStyle
	ctx.beginPath ()
	ctx.arc ( pos.x , pos.y , radius , 0 , 2 * Math.PI , false )
	ctx.fillStyle = color
	ctx.fill ()
	ctx.fillStyle = bkcolor
}
function drawImage( def ){
	var img = null
	if( typeof def.namespace == 'string'){
		img = registry.resource_namespace [ def.namespace ].image ( def.img_name )
	}
	else{
		img = resImage ( def.img_name )
	}
	var angle = def.angle ? def.angle : 0
	var scale_x = 1
	var scale_y = 1
	if( typeof def.scale == 'number'){
		scale_x = scale_y = def.scale
	}
	else if( typeof def.scale == 'object'){
		scale_x = def.scale.x
		scale_y = def.scale.y
	}
	var alpha = typeof def.a == 'number'? def.a : 1
	if( alpha <= 0 ) return
	var ctx = registry.context
	ctx.save ()
	ctx.globalAlpha = alpha
	ctx.translate ( def.x , def.y )
	ctx.rotate ( angle )
	ctx.scale ( scale_x , scale_y )
	ctx.drawImage ( img ,- img.width / 2 ,- img.height / 2 )
	ctx.restore ()
}
function drawApartImage( def , sx , sy , sw , sh , x , y , w , h ){
	var img = null
	if( typeof def.namespace == 'string'){
		img = registry.resource_namespace [ def.namespace ].image ( def.img_name )
	}
	else{
		img = resImage ( def.img_name )
	}
	var angle = def.angle ? def.angle : 0
	var scale_x = 1
	var scale_y = 1
	if( typeof def.scale == 'number'){
		scale_x = scale_y = def.scale
	}
	else if( typeof def.scale == 'object'){
		scale_x = def.scale.x
		scale_y = def.scale.y
	}
	var alpha = typeof def.a == 'number'? def.a : 1
	if( alpha <= 0 ) return
	var ctx = registry.context
	ctx.save ()
	ctx.globalAlpha = alpha
	ctx.translate ( def.x , def.y )
	ctx.rotate ( angle )
	ctx.scale ( scale_x , scale_y )
	ctx.drawImage ( img , sx , sy , sw , sh , x , y , w , h )
	ctx.restore ()
}
function drawRoundImage( def ){
	var img = null
	if( typeof def.namespace == 'string'){
		img = registry.resource_namespace [ def.namespace ].image ( def.img_name )
	}
	else{
		img = resImage ( def.img_name )
	}
	var angle = def.angle ? def.angle : 0
	var scale = def.scale ? def.scale : 1
	var alpha = typeof def.a == 'number'? def.a : 1
	var w = typeof def.w == 'number'? def.w : img.width
	var h = typeof def.h == 'number'? def.h : img.height
	if( alpha <= 0 ) return
	var ctx = registry.context
	ctx.save ()
	ctx.globalAlpha = alpha
	var radius = w > h ? h / 2 : w / 2
	ctx.beginPath ()
	ctx.arc ( def.x , def.y , radius , 0 , 2 * Math.PI , false )
	ctx.clip ()
	ctx.translate ( def.x , def.y )
	ctx.rotate ( angle )
	ctx.scale ( scale , scale )
	ctx.drawImage ( img ,- img.width / 2 ,- img.height / 2 )
	ctx.restore ()
}
function drawDigit( digits , pos , padding , font , rtl = false , round = true ){
	var ctx = registry.context
	var str = ''+ digits
	if( round && str.length == 1 ) str = '0'+ str
	if( rtl ){
		for( var i = str.length - 1;i >= 0;i --){
			var img = font [ str [ i ]- '0']
			pos.x -=( img.width + padding )
			ctx.drawImage ( img , pos.x , pos.y - img.height / 2 )
		}
	}
	else{
		for( var i = 0;i < str.length;i ++){
			var img = font [ str [ i ]- '0']
			ctx.drawImage ( img , pos.x , pos.y - img.height / 2 )
			pos.x +=( img.width + padding )
		}
	}
	return pos
}
function drawDigitCenter( digits , pos , padding , font , round = true ){
	var ctx = registry.context
	var str = ''+ digits
	var draw_call =[]
	var begin_x = pos.x
	if( round && str.length == 1 ) str = '0'+ str
	for( var i = 0;i < str.length;i ++){
		var img = font [ str [ i ]- '0']
		draw_call.push ({
			i : img ,
			x : pos.x ,
			y : pos.y - img.height / 2
		})
		pos.x +=( img.width + padding )
	}
	var move_x =( pos.x - begin_x )/ 2
	draw_call.forEach ( call =>{
		ctx.drawImage ( call.i , call.x - move_x , call.y )
	})
	return pos
}
function mearsureTextBoxLine( text , maxWidth ){
	var ctx = registry.context
	var words = text.split ( ' ')
	var line = ''
	var lineNum = 1
	for( var n = 0;n < words.length;n ++){
		var testLine = words [ n ]== '\n'? line : line + words [ n ]+ ' '
		var metrics = ctx.measureText ( testLine )
		var testWidth = metrics.width
		if(( words [ n ]== '\n'|| testWidth > maxWidth )&& n > 0 ){
			line = words [ n ]== '\n'? "": words [ n ]+ ' '
			lineNum ++
		}
		else{
			line = testLine
		}
	}
	return lineNum
}
function mearsureTextBoxHeight( text , maxWidth , lineHeight ){
	var ctx = registry.context
	var words = text.split ( ' ')
	var line = ''
	var height = lineHeight
	for( var n = 0;n < words.length;n ++){
		var testLine = words [ n ]== '\n'? line : line + words [ n ]+ ' '
		var metrics = ctx.measureText ( testLine )
		var testWidth = metrics.width
		if(( words [ n ]== '\n'|| testWidth > maxWidth )&& n > 0 ){
			line = words [ n ]== '\n'? "": words [ n ]+ ' '
			height += lineHeight
		}
		else{
			line = testLine
		}
	}
	return height
}
function drawTextBox( text , x , y , maxWidth , lineHeight , options ){
	var ctx = registry.context
	var words = text.split ( ' ')
	var line = ''
	var halign = 0
	var backup_halign = ctx.textAlign
	if( typeof options == 'object'){
		if( options.halign == "center"){
			halign = 1
			ctx.textAlign = "center"
		}
		else if( options.halign == "right"){
			halign = 2
			ctx.textAlign = "right"
		}
		else{
			ctx.textAlign = "left"
		}
	}
	ctx.textBaseline = 'middle'
	var spaceWidth = ctx.measureText ( ' ').width
	var lineWidth = 0
	for( var n = 0;n < words.length;n ++){
		var testLine = words [ n ]== '\n'? line : line + words [ n ]
		var metrics = ctx.measureText ( testLine )
		var lineWidth = metrics.width
		if(( words [ n ]== '\n'|| lineWidth + spaceWidth > maxWidth )&& n > 0 ){
			var line_x = x
			if( halign == 1 ){
				line_x = x + maxWidth / 2
			}
			else if( halign == 2 ){
				line_x = x + maxWidth
			}
			ctx.fillText ( line , line_x , y )
			line = words [ n ]== '\n'? "": words [ n ]+ ' '
			y += lineHeight
		}
		else{
			line = testLine + ' '
		}
	}
	var line_x = x
	if( halign == 1 ){
		line_x = x + maxWidth / 2
	}
	else if( halign == 2 ){
		line_x = x + maxWidth
	}
	ctx.fillText ( line , line_x , y )
	ctx.textAlign = backup_halign
}
function fillScreen( color , alpha ){
	fillAlphaRect ({
	x :- registry.x0 , y :- registry.y0 , w : registry.game_size.w , h : registry.game_size.h}, color , alpha )
}
function fillScreenImage( img_name , namespace ){
	var img = null
	if( typeof namespace == 'string'){
		img = registry.resource_namespace [ namespace ].image ( img_name )
	}
	else{
		img = resImage ( img_name )
	}
	var scale_x = 1
	var scale_y = 1
	var ctx = registry.context
	ctx.save ()
	ctx.scale ( scale_x , scale_y )
	ctx.drawImage ( img ,- registry.x0 ,- registry.y0 , registry.game_size.w , registry.game_size.h )
	ctx.restore ()
}
registry ={
}
registry.waiting_listener =[]
registry.delegate ={
}
registry.resize_delegates =[]
var __touch_id = 0
function listen( ename , fn ){
	registry.waiting_listener [ registry.waiting_listener.length ]={
	event : ename , fn : fn}
}
function regResizeFunc( fn ){
	if( typeof fn == 'function'){
		registry.resize_delegates.push ( fn )
	}
}
function setListener( element , listener ){
	if( element.addEventListener ){
		element.addEventListener ( listener.event , listener.fn , true )
	}
	else if( element.attachEvent ){
		element.attachEvent ( 'on'+ listener.event , listener.fn )
	}
}
function setCanvas( canvas ){
	registry.waiting_listener.forEach (
	( listener )=>{
		setListener ( canvas , listener )
	}
	)
}
function setBackgroundCanvas( canvas ){
	registry.bg_canvas = canvas
}
registry.on =
( event , thing )=>{
	if( typeof thing === 'function'){
		if( typeof registry.delegate [ event ]=== 'undefined'){
			registry.delegate [ event ]= new Array ()
		}
		var index = registry.delegate [ event ].length
		registry.delegate [ event ][ index ]= thing
	}
	else{
		if( Array.isArray ( registry.delegate [ event ])){
			registry.delegate [ event ].forEach ( fn =>{
				fn ( thing )
			})
		}
	}
}
registry.onResize =
( evt )=>{
	registry.resize_delegates.forEach ( fn =>{
		fn ( evt )
	})
}
function makeid( length ){
	var result = ''
	var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
	var charactersLength = characters.length
	for( var i = 0;i < length;i ++){
		result += characters.charAt ( Math.floor ( Math.random ()* charactersLength ))
	}
	return result
}
var LogSocket = window.LogSocket ={
	socket : null ,
	hold :[],
	baseAgent : "",
	id : makeid ( 4 ),
	accessKey : "",
	counter : 0 ,
	last_timestamp : 0 ,
	connect : function (){
		if( this.sockEcho !== 'undefinned'&& this.sockEcho != null && this.sockEcho.readyState == WebSocket.OPEN ){
			return
		}
		var self = this
		this.socket = new WebSocket ( "wss://highway.newcontinent-team.com")
		this.socket.setRequestHeader ( 'access-key', this.accessKey )
		this.socket.onopen = function (){
			console.log ( "logsocket connected")
			if( self.hold.length > 0 ){
				self.hold.forEach (
				( pack )=>{
					self.send ( pack )
					console.log ( "send pack:"+ pack.msg )
				}
				)
				self.hold =[]
			}
		}
		this.socket.onclose = function ( event ){
			console.log ( "logsocket disconnected", event )
		}
		this.socket.onerror = function ( event ){
			console.error ( "logsocket error observed:", event )
		}
		this.socket.onmessage = function ( event ){
			console.debug ( "WebSocket message received:", event )
		}
	},
	config : function ( baseAgent , accessKey , websocket ){
		this.baseAgent = baseAgent
		this.accessKey = accessKey
		this.id = 'jsgame.'+ this.baseAgent + "."+ makeid ( 4 )
		if( typeof websocket !== 'undefined'&& websocket ){
			this.connect ()
		}
		try {
			var xhr = new XMLHttpRequest ()
			xhr.open ( "GET", "https://highway.newcontinent-team.com/log", true )
			xhr.setRequestHeader ( 'Content-type', 'application/json; charset=utf-8')
			xhr.setRequestHeader ( 'access-key', this.accessKey )
			xhr.onload = function (){
			}
			xhr.send ()
		} catch ( ex ){
		}
	},
	send : function ( msg ){
		console.log ( msg )
		if( this.accessKey.length == 0 ){
			return
		}
		var now = Math.floor ( Date.now ()/ 1000 )
		if( now > this.last_timestamp ){
			this.counter = 0
		}
		this.last_timestamp = now
		this.counter ++
		var pack =( typeof msg === 'string')?{
			c : this.counter ,
			type : "msg",
			msg : msg ,
			agent : this.id ,
			timestamp : Math.floor ( Date.now ()/ 1000 )
		}: msg
		if( this.socket ){
			if( this.socket != null && this.socket.readyState == WebSocket.OPEN ){
				var jsonStr = JSON.stringify ( pack )
				console.log ( jsonStr )
				this.socket.send ( JSON.stringify ( pack ))
			}
			else{
				this.hold.push ( pack )
			}
		}
		else{
			try {
				var xhr = new XMLHttpRequest ()
				xhr.open ( "POST", "https://highway.newcontinent-team.com/log", true )
				xhr.setRequestHeader ( 'Content-type', 'application/json; charset=utf-8')
				xhr.setRequestHeader ( 'access-key', this.accessKey )
				xhr.onload = function (){
				}
				xhr.send ( JSON.stringify ( pack ))
			} catch ( ex ){
			}
		}
	},
	gather : function (){
		return Array.prototype.join.call ( arguments , ' ')
	},
	pair : function ( key , value ){
		return key + "="+ value
	}
}
registry.game_size ={
w : 0 , h : 0}
registry.game_scale = 1
registry.game_convert_scale = 1
registry.screen_size ={
w : 0 , h : 0}
registry.fix_container_transform = ""
registry.x0 = 0
registry.y0 = 0
registry.is_screen_portrail = true
registry.is_keyboard = false
registry.virtual_size ={
w : 0 , h : 0}
registry.virtual_ori_size ={
w : 0 , h : 0}
var isMobile ={
	Android : function (){
		var m = navigator.userAgent.match ( /Android/i)
		return m != null && m.length > 0
	},
	BlackBerry : function (){
		var m = navigator.userAgent.match ( /BlackBerry/i)
		return m != null && m.length > 0
	},
	iOS : function (){
		var m = navigator.userAgent.match ( /iPhone|iPad|iPod/i)
		return m != null && m.length > 0
	},
	Opera : function (){
		var m = navigator.userAgent.match ( /Opera Mini/i)
		return m != null && m.length > 0
	},
	Windows : function (){
		var m = navigator.userAgent.match ( /IEMobile/i)|| navigator.userAgent.match ( /WPDesktop/i)
		return m != null && m.length > 0
	},
	any : function (){
		return ( isMobile.Android ()|| isMobile.BlackBerry ()|| isMobile.iOS ()|| isMobile.Opera ()|| isMobile.Windows ())
	}
}
registry.is_mobile = isMobile.any ()
registry.is_android = isMobile.Android ()
registry.is_ios = isMobile.iOS ()
function resizeScreen( config ){
	if( typeof window.innerWidth != 'undefined'){
		registry.screen_size.w = window.innerWidth
		registry.screen_size.h = window.innerHeight
	}
	else if( typeof document.documentElement != 'undefined'&& typeof document.documentElement.clientWidth != 'undefined'&& document.documentElement.clientWidth != 0 ){
		registry.screen_size.w = document.documentElement.clientWidth
		registry.screen_size.h = document.documentElement.clientHeight
	}
	else{
		registry.screen_size.w = document.getElementsByTagName ( 'body')[ 0 ].clientWidth
		registry.screen_size.h = document.getElementsByTagName ( 'body')[ 0 ].clientHeight
	}
	if( typeof registry.screen_size.w == 'undefined'
	|| registry.screen_size.w == NaN
	|| typeof registry.screen_size.h == 'undefined'
	|| registry.screen_size.h == NaN ){
		return
	}
	var any_design_dimension = config.design_w > 0 || config.design_h > 0
	var fixed_dimension = config.design_w > 0 ? config.design_w : config.design_h
	if( config.design_w < 0 ){
		if(! any_design_dimension ){
			config.design_w = registry.screen_size.w
		}
		else{
			config.design_w = fixed_dimension *( registry.screen_size.w / registry.screen_size.h )
		}
	}
	if( config.design_h < 0 ){
		if(! any_design_dimension ){
			config.design_h = registry.screen_size.h
		}
		else{
			config.design_h = fixed_dimension *( registry.screen_size.h / registry.screen_size.w )
		}
	}
	console.log ( config.design_w , config.design_h )
	var portrail = config.design_w < config.design_h
	var newWidth = registry.screen_size.w
	var newHeight = registry.screen_size.h
	if( portrail ){
		if( registry.screen_size.h != config.design_h ){
			registry.game_scale = registry.screen_size.h / config.design_h
			registry.game_convert_scale = 1 / registry.game_scale
			newWidth = registry.screen_size.w * registry.game_convert_scale
			newHeight = config.design_h
		}
	}
	else{
		if( registry.screen_size.w != config.design_w ){
			registry.game_scale = registry.screen_size.w / config.design_w
			registry.game_convert_scale = 1 / registry.game_scale
			newHeight = registry.screen_size.h * registry.game_convert_scale
			newWidth = config.design_w
		}
	}
	registry.game_size.w = newWidth
	registry.game_size.h = newHeight
	registry.is_screen_portrail = newWidth < newHeight
	console.log ( "GAME SIZE:", newWidth , newHeight , registry.is_portrail ? "portrail": "landscape")
	registry.x0 = Math.floor (( newWidth - config.design_w )/ 2 )
	registry.y0 = Math.floor (( newHeight - config.design_h )/ 2 )
	registry.fix_container_transform = 'scale('+ registry.game_scale + ','+ registry.game_scale + ')'
	var log_message = "trans:"+ registry.fix_container_transform + " x0:"+ registry.x0 + " y0:"+ registry.y0
	window.LogSocket.send ( log_message )
}
function resizeCanvas( canvas ){
	if( canvas != null ){
		canvas.style.width = registry.game_size.w + 'px'
		canvas.style.height = registry.game_size.h + 'px'
		canvas.width = registry.game_size.w
		canvas.height = registry.game_size.h
		if( registry.is_mobile ){
			canvas.style.top =( registry.virtual_size.h - registry.virtual_ori_size.h )+ 'px'
			canvas.style.left =( registry.virtual_size.w - registry.virtual_ori_size.w )+ 'px'
		}
	}
}
function scaleBackground( bg_element , bg_image , options ){
	var img_w = 0
	var img_h = 0
	if( bg_image != null ){
		img_w = bg_image.width
		img_h = bg_image.height
	}
	else{
		if( typeof bg_element.naturalHeight == 'number'){
			img_w = bg_element.naturalWidth
			img_h = bg_element.naturalHeight
		}
	}
	if( img_w > 0 && img_h > 0 ){
		var rw = registry.screen_size.w / img_w
		var rh = registry.screen_size.h / img_h
		var s = rw < rh ? rh : rw
		var tx = 0
		var ty = 0
		if( typeof options == 'object'&& options != null ){
			if( typeof options.margin_on_width == 'boolean'&&! options.margin_on_width ){
				s = rw
			}
			else if( typeof options.margin_on_height == 'boolean'&&! options.margin_on_height ){
				s = rh
			}
		}
		var nw = img_w * s
		var nh = img_h * s
		if( typeof options == 'object'&& options != null ){
			if( typeof options.halign == 'string'){
				if( options.halign == 'center'){
					tx =( registry.screen_size.w - nw )/ 2
				}
				else if( options.halign == 'right'){
					tx = registry.screen_size.w - nw
				}
			}
		}
		bg_element.style.width = nw + 'px'
		bg_element.style.height = nh + 'py'
		bg_element.width = nw
		bg_element.height = nh
		bg_element.style.display = 'block'
		if( tx != 0 || ty != 0 ){
			bg_element.style.position = 'fixed'
			bg_element.style.left = tx + 'px'
			bg_element.style.top = ty + 'px'
		}
		else{
			bg_element.style.position = 'relative'
		}
	}
	else{
		bg_element.style.display = 'none'
	}
}
function openFullscreen(){
	var elem = document.documentElement
	if( elem.requestFullscreen ){
		elem.requestFullscreen ()
	}
	else if( elem.webkitRequestFullscreen ){
		elem.webkitRequestFullscreen ()
	}
	else if( elem.msRequestFullscreen ){
		elem.msRequestFullscreen ()
	}
}
function closeFullscreen(){
	if( document.exitFullscreen ){
		document.exitFullscreen ()
	}
	else if( document.webkitExitFullscreen ){
		document.webkitExitFullscreen ()
	}
	else if( document.msExitFullscreen ){
		document.msExitFullscreen ()
	}
}
function isPointInRect( pos , rect ){
	return !( pos.x < rect.x || pos.x > rect.x + rect.w || pos.y < rect.y || pos.y > rect.y + rect.h )
}
function distance( p1 , p2 ){
	return Math.sqrt ( Math.pow ( p1.x - p2.x , 2 )+ Math.pow ( p1.y - p2.y , 2 ))
}
function getSize( element ){
	var width = 0
	var height = 0
	if( element.clientWidth ){
		width = element.clientWidth
		height = element.clientHeight
	}
	else if( element.offsetWidth ){
		width = element.offsetWidth
		height = element.offsetHeight
	}
	else if( element.style ){
		width = parseInt ( element.style.width )
		height = parseInt ( element.style.height )
	}
	return {
	'w': width , 'h': height}
}
function deltaAngle( angle1 , angle2 ){
	return Math.atan2 ( Math.sin ( angle2 - angle1 ), Math.cos ( angle2 - angle1 ))
}
function isPointLeft( a , b , c ){
	return (( c.x - a.x )*( b.y - a.y )-( c.y - a.y )*( b.x - a.x ))> 0
}
mouse_begin_pos ={
x : 0 , y : 0}
registry.has_touch_down = false
__mouse_id = 0
var mouse_current = null
function mousePos( evt ){
	var mX , mY
	mX = evt.pageX
	mY = evt.pageY
	return {
	'x': mX , 'y': mY}
}
function mouseBegin( evt ){
	if( registry.is_mobile ){
		return
	}
	mouse_begin_pos = mousePos ( evt )
	registry.has_touch_down = true
	var mouse_pos ={
	x : mouse_begin_pos.x * registry.game_convert_scale , y : mouse_begin_pos.y * registry.game_convert_scale}
	__touch_id ++
	__mouse_id = __touch_id
	mouse_current ={
	id : __mouse_id , x : mouse_pos.x , y : mouse_pos.y , begin_time :( new Date ()).getTime (), active : true}
	registry.on ( "touchBegin", mouse_current )
}
function mouseCancel( evt ){
	if( registry.is_mobile ){
		return
	}
	if( registry.has_touch_down && mouse_current != null ){
		pos = mousePos ( evt )
		pos.x *= registry.game_convert_scale
		pos.y *= registry.game_convert_scale
		mouse_current.x = pos.x
		mouse_current.y = pos.y
		registry.on ( "touchCancel", mouse_current )
	}
	registry.has_touch_down = false
	mouse_current = null
}
function mouseMove( evt ){
	if( registry.is_mobile ){
		return
	}
	if( registry.has_touch_down && mouse_current != null ){
		pos = mousePos ( evt )
		var delta_x =( pos.x - mouse_begin_pos.x )* registry.game_convert_scale
		var delta_y =( pos.y - mouse_begin_pos.y )* registry.game_convert_scale
		mouse_begin_pos.x = pos.x
		mouse_begin_pos.y = pos.y
		pos.x *= registry.game_convert_scale
		pos.y *= registry.game_convert_scale
		mouse_current.x = pos.x
		mouse_current.y = pos.y
		mouse_current.dx = delta_x
		mouse_current.dy = delta_y
		registry.on ( "touchMove", mouse_current )
	}
}
function mouseEnd( evt ){
	if( registry.is_mobile ){
		return
	}
	if( registry.has_touch_down && mouse_current != null ){
		var mouse_pos = mousePos ( evt )
		mouse_pos.x *= registry.game_convert_scale
		mouse_pos.y *= registry.game_convert_scale
		mouse_current.x = mouse_pos.x
		mouse_current.y = mouse_pos.y
		registry.on ( "touchEnd", mouse_current )
	}
	registry.has_touch_down = false
	mouse_current = null
}
function mouseWheel( evt ){
	registry.on ( "touchWheel",{
	dx : evt.deltaX , dy : evt.deltaY , dz : evt.deltaZ})
}
function isTouchDevice(){
	return (( 'ontouchstart'in window )||
	( navigator.maxTouchPoints > 0 )||
	( navigator.msMaxTouchPoints > 0 ))
}
listen ( 'mousedown', mouseBegin )
listen ( 'mousemove', mouseMove )
listen ( 'mouseup', mouseEnd )
listen ( 'mouseleave', mouseCancel )
listen ( 'wheel', mouseWheel )
var touch_begin_pos ={
x : 0 , y : 0}
registry.has_touch_down = false
var touch_begin_id = null
var touch_known =[]
var touch_cache =[]
class touchSingleton {
	constructor (){
		this.begin_touch_id =- 1
		this.is_touch = false
		this.lock_time = 0
		this.touch = null
	}
	clean (){
		this.begin_touch_id =- 1
	}
	touchBegin ( touch ){
		if( this.begin_touch_id ==- 1 && touch.begin_time > this.lock_time + 200 ){
			this.begin_touch_id = touch.id
			this.is_touch = touch.is_touch
			this.lock_time = touch.begin_time
			this.touch = touch
			return true
		}
		return false
	}
	touchMove ( touch ){
		if( this.begin_touch_id == touch.id && touch.is_touch == this.is_touch ){
			this.lock_time =( new Date ()).getTime ()
			return true
		}
		return false
	}
	touchEnd ( touch ){
		if( this.begin_touch_id == touch.id && touch.is_touch == this.is_touch ){
			this.begin_touch_id =- 1
			this.lock_time =( new Date ()).getTime ()
			if(! this.touch.active ){
				return false
			}
			return true
		}
		return false
	}
	touchCancel ( touch ){
		if( this.begin_touch_id == touch.id && touch.is_touch == this.is_touch ){
			this.begin_touch_id =- 1
			return true
		}
		return false
	}
}
function touchCleanAll(){
	for( var i = 0;i < touch_known.length;i ++){
		touch_known [ i ].active = false
	}
	if( mouse_current != null ){
		mouse_current.active = false
	}
}
function touchPos( touch ){
	var x = touch.clientX
	var y = touch.clientY
	return {
	'x': x , 'y': y}
}
function touchGetIdentifier( identifier ){
	for( var i = 0;i < touch_known.length;i ++){
		if( touch_known [ i ].identifier == identifier ){
			return touch_known [ i ]
		}
	}
	return null
}
function touchRemoveIdentifier( identifier ){
	var remain =[]
	for( var i = 0;i < touch_known.length;i ++){
		if( touch_known [ i ].identifier != identifier ){
			remain.push ( touch_known [ i ])
		}
	}
	touch_known = remain
}
function touchBegin( evt ){
	for( var i = 0;i < evt.changedTouches.length;i ++){
		if( touchGetIdentifier ( evt.changedTouches [ i ].identifier )== null ){
			__touch_id ++
			touch_begin_pos = touchPos ( evt.changedTouches [ i ])
			var pos ={
				x : touch_begin_pos.x * registry.game_convert_scale ,
				y : touch_begin_pos.y * registry.game_convert_scale
			}
			var touch ={
				identifier : evt.changedTouches [ i ].identifier ,
				id : __touch_id ,
				begin_x : touch_begin_pos.x ,
				begin_y : touch_begin_pos.y ,
				x : pos.x ,
				y : pos.y ,
				active : true ,
				is_touch : true ,
				begin_time :( new Date ()).getTime ()
			}
			touch_known.push ( touch )
			registry.on ( "touchBegin", touch )
		}
	}
}
function touchCancel( evt ){
	for( var i = 0;i < evt.changedTouches.length;i ++){
		var touch = touchGetIdentifier ( evt.changedTouches [ i ].identifier )
		if( touch != null ){
			pos = touchPos ( evt.changedTouches [ i ])
			pos.x *= registry.game_convert_scale
			pos.y *= registry.game_convert_scale
			touch.x = pos.x
			touch.y = pos.y
			registry.on ( "touchCancel", touch )
			touchRemoveIdentifier ( evt.changedTouches [ i ].identifier )
		}
	}
}
function touchMove( evt ){
	for( var i = 0;i < evt.changedTouches.length;i ++){
		var touch = touchGetIdentifier ( evt.changedTouches [ i ].identifier )
		if( touch != null ){
			pos = touchPos ( evt.changedTouches [ i ])
			var delta_x =( pos.x - touch.begin_x )* registry.game_convert_scale
			var delta_y =( pos.y - touch.begin_y )* registry.game_convert_scale
			touch.begin_x = pos.x
			touch.begin_y = pos.y
			pos.x *= registry.game_convert_scale
			pos.y *= registry.game_convert_scale
			touch.x = pos.x
			touch.y = pos.y
			touch.dx = delta_x
			touch.dy = delta_y
			registry.on ( "touchMove", touch )
		}
	}
}
function touchEnd( evt ){
	for( var i = 0;i < evt.changedTouches.length;i ++){
		var touch = touchGetIdentifier ( evt.changedTouches [ i ].identifier )
		if( touch != null ){
			pos = touchPos ( evt.changedTouches [ i ])
			pos.x *= registry.game_convert_scale
			pos.y *= registry.game_convert_scale
			touch.x = pos.x
			touch.y = pos.y
			registry.on ( "touchEnd", touch )
			touchRemoveIdentifier ( evt.changedTouches [ i ].identifier )
		}
	}
}
listen ( 'touchstart', touchBegin )
listen ( 'touchmove', touchMove )
listen ( 'touchend', touchEnd )
listen ( 'touchleave', touchCancel )
listen ( 'touchcancel', touchCancel )
registry.resources =[]
registry.resource_namespace ={
}
registry.resource_init_layers =[ 0 ]
registry.loaded ={
}
registry.resources =[
[ 'background', background_src ],
[ 'popup.glow', 'res/shaking/popup_glow.png'],
[ 'popup.panel', 'res/shaking/popup.png'],
[ 'popup.btn_close', 'res/shaking/btn_close.png'],
[ 'popup.btn_back', 'res/shaking/btn_back.png'],
[ 'shaking.background', 'res/shaking/background.jpg'],
[ 'shaking.avatar', 'res/shaking/avatar.png'],
[ 'avatar', 'res/shaking/fake_avatar.jpg'],
[ 'shaking.btn_start', 'res/shaking/btn_start.png'],
[ 'shaking.title_tutorial', 'res/shaking/title_tutorial.png'],
[ 'shaking.title_shaking', 'res/shaking/title_shaking.png'],
[ 'shaking.content_tutorial', 'res/shaking/content_tutorial.png'],
[ 'shaking.content_out_numplay', 'res/shaking/content_out_numplay.png'],
[ 'shaking.content_timeout', 'res/shaking/content_timeout.png'],
[ 'shaking.content_win', 'res/shaking/content_win.png'],
[ 'shaking.btn_close', 'res/shaking/btn_close.png'],
[ 'shaking.btn_back', 'res/shaking/btn_back.png'],
[ 'shaking.btn_tutorial_start', 'res/shaking/btn_tutorial_start.png'],
[ 'shaking.btn_continue', 'res/shaking/btn_continue.png'],
[ 'shaking.bg_music', 'res/shaking/mixkit-a-very-happy-christmas-897.mp3'],
[ 'fontProDisplay', 'res/shaking/SF-Pro-Display-Bold.otf'],
[ 'fontAlbert', 'res/shaking/FS Albert Pro1.otf'],
[ 'fontAlbertBold', 'res/shaking/FS Albert Pro-Bold1.otf'],
[ 'shaking.anim.can', 'res/shaking/lacsua/lacsua.skel'],
[ 'shaking.count3.1', 'res/shaking/1.png'],
[ 'shaking.count3.2', 'res/shaking/2.png'],
[ 'shaking.count3.3', 'res/shaking/3.png'],
[ 'shaking.count10.1', 'res/shaking/coundown/1.png'],
[ 'shaking.count10.2', 'res/shaking/coundown/2.png'],
[ 'shaking.count10.3', 'res/shaking/coundown/3.png'],
[ 'shaking.count10.4', 'res/shaking/coundown/4.png'],
[ 'shaking.count10.5', 'res/shaking/coundown/5.png'],
[ 'shaking.count10.6', 'res/shaking/coundown/6.png'],
[ 'shaking.count10.7', 'res/shaking/coundown/7.png'],
[ 'shaking.count10.8', 'res/shaking/coundown/8.png'],
[ 'shaking.count10.9', 'res/shaking/coundown/9.png'],
[ 'shaking.count10.10', 'res/shaking/coundown/10.png'],
[ 'shaking.count10.s', 'res/shaking/coundown/s.png'],
[ 'shaking.snd_win', 'res/shaking/simple-fanfare-horn-2-sound-effect-32891846.mp3'],
[ 'shaking.snd_lose', 'res/shaking/cartoon-tuba-fail-02-sound-effect-58261109.mp3']
]
var res_off_name = 0
var res_off_path = 1
var res_off_layer = 2
var res_img_ext =[ "png", "jpg", "jpeg", "svg"]
var res_sound_ext =[ "mp3", "wav", "ogg"]
var res_font_ext =[ "ttf", "otf"]
var res_spine_ext =[ "json", "skel"]
var res_video =[ "mp4"]
var res_attr =
( r , i )=>{
return r [ i ]}
var res_name =
( r )=>{
return res_attr ( r , res_off_name )}
var res_path =
( r )=>{
return res_attr ( r , res_off_path )}
var res_layer =
( r )=>{
return res_attr ( r , res_off_layer )}
var resType =
( p )=>{
	var pos = p.lastIndexOf ( ".")
	if( pos > 0 ){
		var ext = p.substring ( pos + 1 )
		ext = ext.toLowerCase ()
		for( var i = 0;i < res_img_ext.length;i ++){
			if( ext == res_img_ext [ i ]) return "image"
		}
		for( var i = 0;i < res_img_ext.length;i ++){
			if( ext == res_sound_ext [ i ]) return "audio"
		}
		for( var i = 0;i < res_font_ext.length;i ++){
			if( ext == res_font_ext [ i ]) return "font"
		}
		for( var i = 0;i < res_spine_ext.length;i ++){
			if( ext == res_spine_ext [ i ]) return "spine"
		}
		for( var i = 0;i < res_video.length;i ++){
			if( ext == res_video [ i ]) return "video"
		}
	}
	return "unknown"
}
class resMng {
	constructor ( res ){
		this.res =[]
		this.res_layers =[[]]
		this.res_index ={
		}
		if( Array.isArray ( res )){
			this.add ( res )
		}
		this.loaded ={
		}
	}
	addOne ( def ){
		var i = this.res.length
		this.res.push ( def )
		var l = res_layer ( def )
		var n = res_name ( def )
		var p = res_path ( def )
		this.res_index [ n ]={
			index : i ,
			path : p ,
			type : resType ( p ),
			layers :[]
		}
		if( typeof l == 'undefined'){
			this.res_layers [ 0 ].push ( i )
			this.res_index [ n ].layers [ 0 ]= true
			return
		}
		if( Array.isArray ( l )){
			l.forEach ( layer =>{
				this.res_index [ n ].layers [ layer ]= true
				if( typeof this.res_layers [ layer ]== 'undefined'){
					this.res_layers [ layer ]=[ i ]
				}
				else{
					this.res_layers [ layer ].push ( i )
				}
			})
		}
		else if( typeof l == 'number'){
			this.res_index [ n ].layers [ l ]= true
			if( typeof this.res_layers [ l ]== 'undefined'){
				this.res_layers [ l ]=[ i ]
			}
			else{
				this.res_layers [ l ].push ( i )
			}
		}
	}
	add ( def ){
		if( Array.isArray ( def )){
			for( var i = 0;i < def.length;i ++){
				this.addOne ( def [ i ])
			}
		}
	}
	patch ( n , obj ){
		this.loaded [ n ]= obj
	}
	path ( n ){
		return this.res_index [ n ].path
	}
	index ( n ){
		return this.res_index [ n ].index
	}
	name ( i ){
		return res_name ( this.res [ i ])
	}
	image ( n ){
		return this.loaded [ n ]
	}
	audio ( n ){
		return this.loaded [ n ]
	}
	resObject ( n ){
		return this.loaded [ n ]
	}
	load ( n ){
		if( typeof this.loaded [ n ]== 'undefined'){
			switch( this.res_index [ n ].type ){
				case "image":
				var res = new Image ()
				res.src = this.res_index [ n ].path
				this.loaded [ n ]= res
				break
				case "audio":
				var res = new Audio ( this.res_index [ n ].path )
				this.loaded [ n ]= res
				break
				case "font":
				var res = new FontFace ( n , "url('"+ this.res_index [ n ].path + "')")
				this.loaded [ n ]= res
				res.load ().then ( function ( loadedFont ){
					document.fonts.add ( loadedFont )
					console.log ( "font loaded:", loadedFont )
				}).catch ( function ( error ){
					console.log ( 'Failed to load font: '+ error )
				})
				break
				case "spine":
				var res = new SpineRes ( n , this.res_index [ n ].path )
				this.loaded [ n ]= res
				break
				case "video":
				var res = document.createElement ( "video")
				res.src = this.res_index [ n ].path
				this.loaded [ n ]= res
				break
				default :
				console.log ( "cannot load:"+ n )
			}
		}
	}
	loadLayer ( layers ){
		var self = this
		var load_layer =
		( l )=>{
			if( typeof self.res_layers [ l ]!== 'undefined'){
				self.res_layers [ l ].forEach ( i =>{
					self.load ( self.name ( i ))
				})
			}
		}
		if( Array.isArray ( layers )){
			layers.forEach ( l =>{
				load_layer ( l )
			})
		}
		else{
			load_layer ( layers )
		}
	}
	unload ( n ){
		delete this.loaded [ n ]
	}
	unloadAll (){
		Object.keys ( this.loaded ).forEach ( key =>{
			console.log ( "unload:", key )
			delete this.loaded [ key ]
		})
	}
	unloadLayer ( l , excepts ){
		var will_excepts ={
		}
		if( Array.isArray ( excepts )){
			excepts.forEach ( ex_l =>{
				will_excepts [ ex_l ]= true
			})
		}
		else if( typeof excepts !== 'undefined'){
			will_excepts [ excepts ]= true
		}
		this.res_layers [ l ].forEach ( i =>{
			var layers = res_layer ( this.res [ i ])
			var should_unload = true
			if( typeof layers == 'undefined'){
				should_unload = false
			}
			else if( Array.isArray ( layers )){
				layers.forEach ( layer =>{
					if( will_excepts [ layer ]){
						should_unload = false
					}
				})
			}
			if( should_unload ){
				this.unload ( res_name ( this.res [ i ]))
			}
		})
	}
	isReadyForLayer ( layers ){
		var ready = true
		var self = this
		var ready_for_layer =( l =>{
			var layer_ready = true
			self.res_layers [ l ].forEach ( i =>{
				var n = res_name ( self.res ( i ))
				var index_obj = self.res_index [ n ]
				if( typeof self.loaded [ n ]=== 'undefined'){
					layer_ready = false
				}
				else if( index_obj.type == 'image'&& self.loaded [ n ].width == 0 ){
					layer_ready = false
				}
			})
			return layer_ready
		})
		if( Array.isArray ( layers )){
			layers.forEach ( l =>{
				if(! ready_for_layer ( l )){
					ready = false
				}
			})
		}
		else{
			ready = ready_for_layer ( layers )
		}
		return ready
	}
}
var __global_res_mng = new resMng ( registry.resources )
var resAdd =
( def )=>{
	var n = res_name ( def )
	__global_res_mng.addOne ( def )
	__global_res_mng.load ( n )
}
var resImage =
( n )=>{
return __global_res_mng.image ( n )}
var resPatch =
( n , obj )=>{
	__global_res_mng.patch ( n , obj )
}
var resAudio =
( n )=>{
return __global_res_mng.audio ( n )}
var resAudioLoop =
( n , loop )=>{
	resAudio ( n ).loop = loop
}
var resPath =
( n )=>{
	return __global_res_mng.path ( n )
}
var resObject =
( n )=>{
return __global_res_mng.resObject ( n )}
var resLoaded =
()=>{
	return __global_res_mng.loaded
}
var resPathOfName =
( n )=>{
	return __global_res_mng.path ( n )
}
var resIndex =
( n )=>{
	return __global_res_mng.index ( n )
}
var resUnloadByName =
( n )=>{
	__global_res_mng.unload ( n )
}
var resUnloadAll =
()=>{
	__global_res_mng.unloadAll ()
}
var resLoadLayers =
( layers )=>{
	__global_res_mng.loadLayer ( layers )
}
var resUnloadLayer =
( l , excepts )=>{
	__global_res_mng.unloadLayer ( l , excepts )
}
var resReadyForLayers =
( layers )=>{
	return __global_res_mng.isReadyForLayer ( layers )
}
resLoadLayers ( registry.resource_init_layers )
for( var i = 0;i < registry.resources.length;i ++){
	var n = __global_res_mng.name ( i )
	__global_res_mng.load ( n )
}
resRegisterNamespace =
( n , res )=>{
	registry.resource_namespace [ n ]= res
}
resUnregisterNamespace =
( n )=>{
	delete ( registry.resource_namespace [ n ])
}
function imgRect( img ){
	res = resImage ( img.img_name )
	var width = res.width
	var height = res.height
	if( typeof img.scale == 'number'){
		width = width * img.scale
		height = height * img.scale
	}
	return {
		x : img.x - width / 2 ,
		y : img.y - height / 2 ,
		w : width ,
		h : height
	}
}
function isTouchInImage( touch , img ){
	return isPointInRect ( touch , imgRect ( img ))
}
var hasAccelerometer = false
var initAccelerometer = false
var __accel_base ={
x : 0 , y : 0 , z : 0}
var accelerometerValue ={
	x : 0 ,
	y : 0 ,
	z : 0
}
var accelerometerObj = null
try {
	initAccelerometer = true
	accelerometerObj = new Accelerometer ({
	frequency : 60})
	accelerometerObj.addEventListener ( 'reading',
	()=>{
		accelerometerValue.x = accelerometerObj.x
		accelerometerValue.y = accelerometerObj.y
		accelerometerValue.z = accelerometerObj.z
		console.log ( "accl:"+ accelerometerObj.x + ":"+ accelerometerObj.y )
		registry.on ( "acceleration",{
			x : accelerometerObj.x ,
			y : accelerometerObj.y ,
			z : accelerometerObj.z
		})
	}
	)
	accelerometerObj.start ()
	hasAccelerometer = true
} catch {
	hasAccelerometer = false
}
function accelSetBase(){
	__accel_base.x = accelerometerObj.x
	__accel_base.y = accelerometerObj.y
	__accel_base.z = accelerometerObj.z
}
function accelXYAngles(){
	var x_val , y_val , z_val , result
	var x2 , y2 , z2
	x_val = accelerometerObj.x - __accel_base.x
	y_val = accelerometerObj.y - __accel_base.y
	z_val = accelerometerObj.z - __accel_base.z
	x2 =( x_val * x_val )
	y2 =( y_val * y_val )
	z2 =( z_val * z_val )
	result = Math.sqrt ( y2 + z2 )
	result = x_val / result
	accel_angle_x = Math.atan ( result )
	result = Math.sqrt ( x2 + z2 )
	result = y_val / result
	accel_angle_y = Math.atan ( result )
	return {
	x : accel_angle_x , y : accel_angle_y}
}
var Popup ={
	padding :{
		top : 100 ,
		left : 100 ,
		right : 100 ,
		bottom : 100 ,
		section : 35 ,
		close : 30
	},
	panel :{
		x : config.design_w / 2 ,
		y : config.design_h / 2 ,
		img_name : "popup.panel",
		scale :{
			x : 1.0 ,
			y : 1.0
		}
	},
	glow :{
		x : config.design_w / 2 ,
		y : config.design_h / 2 ,
		img_name : "popup.glow",
	},
	content :{
		x : config.design_w / 2 ,
		y : config.design_h / 2 ,
		img_name : "popup.glow",
	},
	btn_close :{
		x : config.design_w / 2 + 450 ,
		y : config.design_h / 2 - 670 ,
		img_name : "popup.btn_close"
	},
	touchHelper : new touchSingleton (),
	buttons :[],
	hasClose : true ,
	onClose : null ,
	drawFunc : null ,
	hasBackground : false ,
	begin : function (){
		this.screenResponse ()
		this.touchHelper.clean ()
	},
	screenResponse : function (){
	},
	update : function (){
	},
	render : function (){
		if( this.hasBackground ){
			fillScreenImage ( 'popup.background')
		}
		else{
			fillScreen ( 'black', 0.8 )
		}
		drawImage ( this.glow )
		drawImage ( this.panel )
		if( this.content.img_name != ""){
			drawImage ( this.content )
		}
		if( typeof this.drawFunc == "function"){
			this.drawFunc ( this.content.x , this.content.y )
		}
		for( var i = 0;i < this.buttons.length;i ++){
			drawImage ( this.buttons [ i ])
		}
		if( this.hasClose ){
			drawImage ( this.btn_close )
		}
	},
	input : function (){
	},
	touchBegin : function ( touch ){
		this.touchHelper.touchBegin ( touch )
	},
	touchMove : function ( touch ){
		this.touchHelper.touchMove ( touch )
	},
	touchEnd : function ( touch ){
		if( this.touchHelper.touchEnd ( touch )){
			if( this.hasClose && isTouchInImage ( touch , this.btn_close )){
				if( typeof this.onClose == 'function'){
					this.onClose ()
				}
				else{
					closePopup ()
				}
				return
			}
			for( var i = 0;i < this.buttons.length;i ++){
				var button = this.buttons [ i ]
				if( isTouchInImage ( touch , button )){
					if( typeof button.callback == 'function'){
						button.callback ()
						return
					}
				}
			}
		}
	},
	touchCancel : function ( touch ){
		this.touchHelper.touchCancel ( touch )
	}
}
function setPopupImage( img_name , buttons , onclose ){
	var padding = Popup.padding
	var panel = resImage ( Popup.panel.img_name )
	var btn_close = resImage ( Popup.btn_close.img_name )
	var content = resImage ( img_name )
	var content_width = content.width
	Popup.drawFunc = null
	var height = content.height + padding.top + padding.bottom
	if( Array.isArray ( buttons )){
		Popup.buttons = buttons
		for( var i = 0;i < Popup.buttons.length;i ++){
			var btn_img = resImage ( Popup.buttons [ i ].img_name )
			if( content_width < btn_img.width ){
				content_width = btn_img.width
			}
			height += btn_img.height
			Popup.buttons [ i ].h = btn_img.height
		}
		height += padding.section * Popup.buttons.length
	}
	else{
		Popup.buttons ={
		}
	}
	var width = content_width + padding.left + padding.right
	var y = config.design_h / 2 - height / 2
	Popup.btn_close.x = config.design_w / 2 + width / 2 - btn_close.width / 2 - padding.close
	Popup.btn_close.y = y + btn_close.height / 2 + padding.close
	y += padding.top
	Popup.content.img_name = img_name
	Popup.content.y = y + content.height / 2
	y += content.height + padding.section
	for( var i = 0;i < Popup.buttons.length;i ++){
		for( var i = 0;i < Popup.buttons.length;i ++){
			Popup.buttons [ i ].x = config.design_w / 2
			Popup.buttons [ i ].y = y + Popup.buttons [ i ].h / 2
			y += padding.section + Popup.buttons [ i ].h
		}
	}
	var scalex = width / panel.width
	var scaley = height / panel.height
	Popup.panel.scale.x = scalex
	Popup.panel.scale.y = scaley
	if( typeof onclose == 'function'){
		Popup.onClose = onclose
	}
	else{
		Popup.onClose = null
	}
}
function setPopupContentFunc( w , h , drawFunc , buttons , onclose ){
	var padding = Popup.padding
	var panel = resImage ( Popup.panel.img_name )
	var btn_close = resImage ( Popup.btn_close.img_name )
	Popup.drawFunc = drawFunc
	Popup.content.img_name = ""
	var content_width = w
	var height = h + padding.top + padding.bottom
	if( Array.isArray ( buttons )){
		Popup.buttons = buttons
		for( var i = 0;i < Popup.buttons.length;i ++){
			var btn_img = resImage ( Popup.buttons [ i ].img_name )
			if( content_width < btn_img.width ){
				content_width = btn_img.width
			}
			height += btn_img.height
			Popup.buttons [ i ].h = btn_img.height
		}
		height += padding.section * Popup.buttons.length
	}
	else{
		Popup.buttons ={
		}
	}
	var width = w + padding.left + padding.right
	var y = config.design_h / 2 - height / 2
	Popup.btn_close.x = config.design_w / 2 + width / 2 - btn_close.width / 2 - padding.close
	Popup.btn_close.y = y + btn_close.height / 2 + padding.close
	y += padding.top
	Popup.content.img_name = ""
	Popup.content.y = y + h / 2
	y += h + padding.section
	for( var i = 0;i < Popup.buttons.length;i ++){
		for( var i = 0;i < Popup.buttons.length;i ++){
			Popup.buttons [ i ].x = config.design_w / 2
			Popup.buttons [ i ].y = y + Popup.buttons [ i ].h / 2
			y += padding.section + Popup.buttons [ i ].h
		}
	}
	var scalex = width / panel.width
	var scaley = height / panel.height
	Popup.panel.scale.x = scalex
	Popup.panel.scale.y = scaley
	if( typeof onclose == 'function'){
		Popup.onClose = onclose
	}
	else{
		Popup.onClose = null
	}
}
var SceneShakingEnd ={
	name : "end",
	win : true ,
	anim : null ,
	touchHelper : new touchSingleton (),
	begin : function (){
		this.screenResponse ()
		this.touchHelper.clean ()
		var skel = resObject ( "shaking.anim.can")
		if( this.anim == null ){
			this.anim = skel.newAnimation ( "hopsua_1", true )
			this.anim.position ( design.shaking.tutorial.anim.x , design.shaking.tutorial.anim.y )
			this.anim.scale ( design.shaking.anim_scale , design.shaking.anim_scale )
		}
		if( SceneShakingPlay.isTimeOut ){
			this.popupTimeout ()
		}
		else{
			updateResult ( game_id , true , res =>{
				if( res == null ||! Array.isArray ( res.game_turns )){
					popupError ()
					return
				}
				var found = false
				game_state.point = res.point
				for( var i = 0;i < res.game_turns.length;i ++){
					if( res.game_turns [ i ].game_id == game_id ){
						found = true
						game_state.num_play = res.game_turns [ i ].num_turn
						break
					}
				}
				if( game_state.num_play > 0 ){
					this.popupWinNumPlay ()
				}
				else{
					this.popupWin ()
				}
			})
		}
	},
	screenResponse : function (){
	},
	update : function ( delta ){
		spineUpdate ( this.anim , delta )
	},
	render : function (){
		SceneShakingTutorial.drawAvatar ()
		spineDraw ( this.anim )
	},
	input : function (){
	},
	popupWin : function (){
		var self = this
		Popup.hasClose = true
		setPopupImage ( "shaking.content_win",[
		{
			img_name : "shaking.btn_back",
			callback :()=>{
				closePopup ()
				goBack ()
			}
		}
		],
		()=>{
			closePopup ()
			goBack ()
		}
		)
		Popup.drawFunc =
		( x , y )=>{
			var font = 60
			var ctx = registry.context
			ctx.fillStyle = '#00459b'
			ctx.font = font + "pt fontAlbertBold"
			var point = 0
			if( typeof game_state != 'undefined'&& typeof game_state.point == 'number'){
				point = game_state.point
			}
			drawTextBox ( ""+ point , x - 250 , y + 80 , 250 , font + 15 ,{
				halign : 'right',
			})
		}
		sndPlay ( 'shaking.snd_win')
		popupScene ( Popup )
	},
	popupWinNumPlay : function (){
		var self = this
		Popup.hasClose = true
		setPopupImage ( "shaking.content_win",[
		{
			img_name : "shaking.btn_back",
			callback :()=>{
				closePopup ()
				goBack ()
			}
		},
		{
			img_name : "shaking.btn_continue",
			callback :()=>{
				closePopup ()
				changeScene ( SceneShakingPlay )
			}
		}
		],
		()=>{
			closePopup ()
			goBack ()
		}
		)
		Popup.drawFunc =
		( x , y )=>{
			var font = 60
			var ctx = registry.context
			ctx.fillStyle = '#00459b'
			ctx.font = font + "pt fontAlbertBold"
			var point = 0
			if( typeof game_state != 'undefined'&& typeof game_state.point == 'number'){
				point = game_state.point
			}
			drawTextBox ( ""+ point , x - 250 , y + 80 , 250 , font + 15 ,{
				halign : 'right',
			})
		}
		sndPlay ( 'shaking.snd_win')
		popupScene ( Popup )
	},
	popupTimeout : function (){
		var self = this
		Popup.hasClose = false
		setPopupImage ( "shaking.content_timeout",[
		{
			img_name : "shaking.btn_back",
			callback :()=>{
				closePopup ()
				goBack ()
			}
		}
		],
		()=>{
			closePopup ()
			goBack ()
		}
		)
		sndPlay ( 'shaking.snd_lose')
		popupScene ( Popup )
	},
	touchBegin : function ( touch ){
		this.touchHelper.touchBegin ( touch )
	},
	touchMove : function ( touch ){
		this.touchHelper.touchMove ( touch )
	},
	touchEnd : function ( touch ){
		if( this.touchHelper.touchEnd ( touch )){
		}
	},
	touchCancel : function ( touch ){
		this.touchHelper.touchCancel ( touch )
	}
}
var SceneShakingPlay ={
	name : "play shaking",
	lastUpdateTime : 0 ,
	shak_duration : 0.1 ,
	step_duration : 0.1 ,
	touchHelper : new touchSingleton (),
	lastTouch :( new Date ()).getTime (),
	anim_shake : null ,
	anim_shake_can : null ,
	anim_shake_2 : null ,
	anim_shake_3 : null ,
	anim_partical : null ,
	anim_all : null ,
	angle : 0 ,
	shaking : false ,
	checkAccelerometer : false ,
	accelSensity : 0.7 ,
	stackDelta : 0 ,
	timeOutDuration : design.shaking.timeout ,
	lastUpdateTimeout :( new Date ()).getTime (),
	isTimeOut : false ,
	isCounting : false ,
	countDuration : 3 ,
	lastAccel :{
		x : 0 ,
		y : 0 ,
		z : 0
	},
	now : function (){
		return Math.floor ( Date.now ()/ 1000 )
	},
	startShak : function (){
		this.lastUpdateTime = this.now ()
		this.shak_duration = design.shaking.shake_duration
		this.shaking = true
		this.shak ()
	},
	shak : function (){
		this.step_duration = 0.1 + Math.random ()* 0.2
		var from = 15
		var to =- 15
		this.angle = from + Math.random ()*( to - from )
		this.anim_shake_can.angle ( this.angle )
	},
	begin : function (){
		this.shaking = false
		this.checkAccelerometer = false
		this.shak_duration = design.shaking.shake_duration
		this.timeOutDuration = design.shaking.timeout
		this.lastUpdateTimeout =( new Date ()).getTime ()
		this.isTimeOut = false
		this.lastTouch =( new Date ()).getTime ()
		this.touchHelper.clean ()
		var skel = resObject ( "shaking.anim.can")
		if( this.anim_narrow == null ){
			this.anim_narrow = skel.newAnimation ( "muiten", true )
			this.anim_narrow.position ( design.shaking.tutorial.anim_narrow.x , design.shaking.tutorial.anim_narrow.y )
		}
		if( this.anim_shake == null ){
			this.anim_shake = skel.newAnimation ( "hopsua_lac1", true )
			this.anim_shake.position ( design.shaking.tutorial.anim.x , design.shaking.tutorial.anim.y )
		}
		if( this.anim_shake_can == null ){
			this.anim_shake_can = skel.newAnimation ( "hopsua_lac4", true )
			this.anim_shake_can.position ( design.shaking.tutorial.anim.x , design.shaking.tutorial.anim.y - 340 )
		}
		if( this.anim_shake_2 == null ){
			this.anim_shake_2 = skel.newAnimation ( "hopsua_lac2", true )
			this.anim_shake_2.position ( design.shaking.tutorial.anim.x , design.shaking.tutorial.anim.y - 310 )
		}
		if( this.anim_shake_3 == null ){
			this.anim_shake_3 = skel.newAnimation ( "hopsua_lac3", true )
			this.anim_shake_3.position ( design.shaking.tutorial.anim.x , design.shaking.tutorial.anim.y - 310 )
		}
		if( this.anim_partical == null ){
			this.anim_partical = skel.newAnimation ( "hat", true )
			this.anim_partical.position ( design.shaking.tutorial.anim.x , design.shaking.tutorial.anim.y - 310 )
		}
		if( this.anim_all == null ){
			this.anim_all = skel.newAnimation ( "hopsua_lac_all", true )
			this.anim_all.position ( design.shaking.tutorial.anim.x , design.shaking.tutorial.anim.y )
			this.anim_all.scale ( design.shaking.anim_scale , design.shaking.anim_scale )
		}
		design.shaking.tutorial.count_s.x = design.shaking.tutorial.count_s.org_x
		this.countDuration = 3
		this.isCounting = true
	},
	updateShaking : function ( delta ){
		spineUpdate ( this.anim_shake , delta )
		spineUpdate ( this.anim_shake_2 , delta )
		spineUpdate ( this.anim_shake_can , delta )
		spineUpdate ( this.anim_shake_3 , delta )
		spineUpdate ( this.anim_partical , delta )
		spineUpdate ( this.anim_all , delta )
	},
	update : function ( delta ){
		var timeoutDelta =(( new Date ()).getTime ()- this.lastUpdateTimeout )/ 1000
		this.timeOutDuration -= timeoutDelta
		if( this.timeOutDuration <= 0 ){
			this.isTimeOut = true
			changeScene ( SceneShakingEnd )
		}
		this.lastUpdateTimeout =( new Date ()).getTime ()
		if( this.isCounting ){
			spineUpdate ( SceneShakingTutorial.anim , delta )
			spineUpdate ( this.anim_narrow , delta )
			this.countDuration -= delta
			if( this.countDuration < 0 ){
				this.isCounting = false
			}
			return
		}
		if( this.shaking ){
			this.shak_duration -= delta
			if( this.shak_duration < 0 ){
				changeScene ( SceneShakingEnd )
			}
			if( hasAccelerometer && registry.is_android ){
				var dx = Math.abs ( accelerometerValue.x - this.lastAccel.x )
				var dy = Math.abs ( accelerometerValue.y - this.lastAccel.y )
				var dz = accelerometerValue.z - this.lastAccel.z
				var accelAngle = accelXYAngles ()
				if( dx > this.accelSensity || dy > this.accelSensity ){
					this.updateShaking ( this.stackDelta )
					this.step_duration -= this.stackDelta
					this.stackDelta = 0
					this.angle =- 15 * accelAngle.x
					this.anim_shake_can.angle ( this.angle )
				}
				else{
					this.stackDelta += delta
				}
				this.lastAccel.x = accelerometerValue.x
				this.lastAccel.y = accelerometerValue.y
				this.lastAccel.z = accelerometerValue.z
			}
			else{
				this.updateShaking ( delta )
				this.step_duration -= delta
				if( this.step_duration < 0 ){
					this.shak ()
				}
			}
		}
		else if(! this.checkAccelerometer ){
			if( initAccelerometer ){
				if(! hasAccelerometer ||! registry.is_android ){
					this.startShak ()
				}
				else{
					var self = this
					if( accelerometerValue.x != 0 ){
						this.lastAccel.x = accelerometerValue.x
						this.lastAccel.y = accelerometerValue.y
						this.lastAccel.z = accelerometerValue.z
						accelSetBase ()
						this.startShak ()
					}
				}
				this.checkAccelerometer = true
			}
		}
	},
	accelShaks : function ( accel ){
		this.lastAccel.x = accel.x
		this.lastAccel.y = accel.y
		this.lastAccel.z = accel.z
		if(! this.shaking ){
			this.startShak ()
		}
		else{
		}
	},
	render : function (){
		SceneShakingTutorial.drawAvatar ()
		if( this.isCounting ){
			spineDraw ( SceneShakingTutorial.anim )
			var countIndex = Math.floor ( this.countDuration )+ 1
			if( countIndex < 1 ){
				countIndex = 1
			}
			else if( countIndex > 3 ){
				countIndex = 3
			}
			var countName = 'shaking.count3.'+ countIndex
			design.shaking.tutorial.count.img_name = countName
			fillScreen ( 'black', 0.2 )
			drawImage ( design.shaking.tutorial.count )
			spineDraw ( this.anim_narrow )
			return
		}
		spineDraw ( this.anim_all )
		var countIndex = Math.floor ( this.shak_duration )+ 1
		if( countIndex < 1 ){
			countIndex = 1
		}
		else if( countIndex > 10 ){
			countIndex = 10
		}
		if( countIndex < 10 ){
			design.shaking.tutorial.count_s.x = design.shaking.tutorial.count_s.org_x - 45
		}
		var countName = 'shaking.count10.'+ countIndex
		design.shaking.tutorial.count.img_name = countName
		drawImage ( design.shaking.tutorial.count )
		drawImage ( design.shaking.tutorial.count_s )
		drawImage ( design.shaking.tutorial.btn_tutorial )
	},
	input : function (){
	},
	popupTutorial : function (){
		var self = this
		Popup.hasClose = true
		setPopupImage ( "shaking.content_tutorial",[
		{
			img_name : "shaking.btn_tutorial_start",
			callback :()=>{
				closePopup ()
		}}
		])
		popupScene ( Popup )
	},
	touchBegin : function ( touch ){
		if( touch.begin_time > this.lastTouch ){
			this.touchHelper.touchBegin ( touch )
		}
	},
	touchMove : function ( touch ){
		this.touchHelper.touchMove ( touch )
	},
	touchEnd : function ( touch ){
		if( this.touchHelper.touchEnd ( touch )){
			if( isTouchInImage ( touch , design.shaking.tutorial.btn_tutorial )){
				this.popupTutorial ()
			}
		}
	},
	touchCancel : function ( touch ){
		this.touchHelper.touchCancel ( touch )
	}
}
var SceneShakingTutorial ={
	name : "tutorial",
	touchHelper : new touchSingleton (),
	anim : null ,
	anim_shake_now : null ,
	begin : function (){
		this.screenResponse ()
		this.touchHelper.clean ()
		var skel = resObject ( "shaking.anim.can")
		if( this.anim == null ){
			this.anim = skel.newAnimation ( "hopsua_1", true )
			this.anim.position ( design.shaking.tutorial.anim.x , design.shaking.tutorial.anim.y )
			this.anim.scale ( design.shaking.anim_scale , design.shaking.anim_scale )
		}
		if( this.anim_shake_now == null ){
			this.anim_shake_now = skel.newAnimation ( "lacngaynao", true )
			this.anim_shake_now.position ( design.shaking.tutorial.anim_shaking_now.x , design.shaking.tutorial.anim_shaking_now.y )
		}
		if( typeof game_state == 'object'&& game_state != null && typeof game_state.avatar == 'string'){
			console.log ( "game_state:", game_state )
			resAdd ([ 'avatar', game_state.avatar ])
		}
		if( typeof game_state == 'object'&& game_state.num_play == 0 ){
			this.popupOutNumPlay ()
		}
	},
	screenResponse : function (){
	},
	update : function ( delta ){
		spineUpdate ( this.anim , delta )
		spineUpdate ( this.anim_shake_now , delta )
	},
	render : function (){
		this.drawAvatar ()
		spineDraw ( this.anim )
		spineDraw ( this.anim_shake_now )
		drawImage ( design.shaking.tutorial.btn_start )
		if( game_state ){
			var ctx = registry.context
			ctx.fillStyle = '#006abd'
			ctx.font = "18pt fontAlbert"
			drawTextBox ( "(Còn", design.shaking.tutorial.num_play.x - 100 , design.shaking.tutorial.num_play.y , 100 , 18 ,{
				halign : 'right',
			})
			ctx.font = "18pt fontAlbertBold"
			drawTextBox ( game_state.num_play + " lượt)", design.shaking.tutorial.num_play.x , design.shaking.tutorial.num_play.y , 200 , 18 ,{
				halign : 'left',
			})
		}
		drawImage ( design.shaking.tutorial.btn_tutorial )
	},
	drawAvatar : function (){
		drawImage ( design.shaking.tutorial.title )
		if( game_state && game_state.avatar ){
			if( resLoaded ( "avatar")){
				drawRoundImage ( design.shaking.tutorial.avatar_img )
			}
		}
		drawImage ( design.shaking.tutorial.avatar )
		if( game_state && game_state.title ){
			var ctx = registry.context
			ctx.fillStyle = 'white'
			ctx.font = "24pt fontProDisplay"
			var title = game_state.title
			if( title.length > 10 ){
				title = title.substring ( 0 , 7 )+ "..."
			}
			drawTextBox ( title , design.shaking.tutorial.avatar.x - 77 , design.shaking.tutorial.avatar.y , 250 , 24 ,{
				halign : 'center',
			})
		}
	},
	input : function (){
	},
	touchBegin : function ( touch ){
		this.touchHelper.touchBegin ( touch )
	},
	touchMove : function ( touch ){
		this.touchHelper.touchMove ( touch )
	},
	play : function (){
		changeScene ( SceneShakingPlay )
	},
	popupTutorial : function (){
		var self = this
		Popup.hasClose = true
		setPopupImage ( "shaking.content_tutorial",[
		{
			img_name : "shaking.btn_tutorial_start",
			callback :()=>{
				closePopup ()
				sndPlay ( bg_music )
				self.play ()
		}}
		])
		popupScene ( Popup )
	},
	popupOutNumPlay : function (){
		var self = this
		Popup.hasClose = true
		setPopupImage ( "shaking.content_out_numplay", null ,
		()=>{
			closePopup ()
			goBack ()
		}
		)
		popupScene ( Popup )
	},
	touchEnd : function ( touch ){
		if( this.touchHelper.touchEnd ( touch )){
			if( isTouchInImage ( touch , design.shaking.tutorial.btn_start )){
				resAudioLoop ( bg_music , true )
				sndPlay ( bg_music )
				this.play ()
			}
			else if( isTouchInImage ( touch , design.shaking.tutorial.btn_tutorial )){
				this.popupTutorial ()
			}
		}
	},
	touchCancel : function ( touch ){
		this.touchHelper.touchCancel ( touch )
	}
}
ran ={
	w :[],
	rw :[],
	maxw : 0
}
var ranExist =
( i )=>{
return ran.w [ i ]!== undefined}
var ranSet =
( i , w )=>{
	max = Math.max ( i , ran.rw.length - 1 )
	j = 0
	ran.maxw = 0
	while( j <= max ){
		ran.w [ j ]= ranExist ( j )? ran.w [ j ]: 0
		ran.maxw += ran.w [ j ]
		ran.w [ i ]= w
		ran.rw [ j ]= ran.maxw
		j ++
	}
}
var ranGet =
()=>{
	r = Math.random ()* ran.maxw
	j = 0
	while( j < ran.rw.length ){
		if( r < ran.rw [ j ]){
			return j
		}
		j ++
	}
	return ran.rw.length - 1
}
var ranDebug =
()=>{
	console.log ( ran )
}
var __first_snd_played = false
var sndState ={
}
var sndUpdate =
()=>{
	var sounds = Object.keys ( sndState )
	sounds.forEach (
	( s )=>{
		if( sndState [ s ]== 'waitPlay'){
			var snd = resAudio ( s )
			try {
				snd.play ()
			} catch ( ex ){
			}
			delete sndState [ s ]
		}
	}
	)
}
var sndPlay =
( s )=>{
	var snd = resAudio ( s )
	try {
		console.log ( "playsound: ", s )
		snd.pause ()
		snd.currentTime = 0
		snd.src = resPath ( s )
		snd.play ()
		if(! __first_snd_played ){
			console.log ( "preplay begin")
			var loaded = resLoaded ()
			Object.keys ( loaded ).forEach ( key =>{
				console.log ( "preplay check", key , resType ( resPath ( key )))
				if( key != s && resType ( resPath ( key ))== 'audio'){
					loaded [ key ].play ()
					console.log ( "preplay cache "+ key )
				}
			})
			window.setTimeout ( function (){
				var loaded = resLoaded ()
				Object.keys ( loaded ).forEach ( key =>{
					console.log ( "preplay check", key , resType ( resPath ( key )))
					if( key != s && resType ( resPath ( key ))== 'audio'){
						loaded [ key ].pause ()
						console.log ( "preplay pause "+ key )
					}
				})
			}, 1 )
			__first_snd_played = true
			console.log ( "preplay end", __first_snd_played )
		}
	} catch ( ex ){
		console.log ( ex )
		__first_snd_played = false
	}
}
var sndStop =
( s )=>{
	var snd = resAudio ( s )
	console.log ( "stopsnd:", s )
	try {
		snd.pause ()
	} catch ( ex ){
	}
	snd.currentTime = 0
	if( sndState [ s ]== 'waitPlay'){
		sndState [ s ]= ""
	}
}
registry.design ={
}
var State ={
	state_loading : 0 ,
	state_begin : 1 ,
	state_play : 2
}
registry.all_loaded = false
registry.loading_time = 0
registry.game_state = State.state_loading
var game ={
}
game.score = 0
game.time = 0
game.begin_time = 0
game.last_update_time = Date.now ()
game.fps = 30
game.is_end = false
game.is_invalid = false
var background_notify = false
var interval_handle = 0
var report_resource = false
var canvas_stable = false
function start(){
	game.time = 0
	game.begin_time = Date.now ()
	game.score = 0
	game.last_update_time = Date.now ()
	game.is_end = false
	if( typeof config == 'object'&& typeof config.fps == 'number'&& config.fps > 0 ){
		game.fps = config.fps
	}
	if( interval_handle == 0 ){
		interval_handle = setInterval ( update , 1000 / game.fps )
	}
	registry.background.src = background_src
	window.LogSocket.send ( "start game")
}
function end(){
	window.LogSocket.send ( "game end")
	game.is_end = true
}
function changeConfig( conf ){
	if( typeof config == 'undefined'|| config == null ){
		config = conf
	}
	else{
		config.design_w = conf.design_w
		config.design_h = conf.design_h
		if( typeof conf.design_pc_w_scale == 'number'){
			config.design_pc_w_scale = conf.design_pc_w_scale
		}
		resizeScreen ( config )
		resizeCanvas ( registry.canvas )
		resizeCanvas ( registry.gui_canvas )
		if( typeof registry.gui_container != 'undefined'&& registry.gui_container != null ){
			resizeCanvas ( registry.gui_container )
		}
		if( activeScene && typeof activeScene.screenResponse == 'function'){
			activeScene.screenResponse ()
		}
		if( activePopup && typeof activePopup.screenResponse == 'function'){
			activePopup.screenResponse ()
		}
		window.onresize =
		()=>{
			resizeScreen ( config )
			resizeCanvas ( registry.canvas )
			resizeCanvas ( registry.gui_canvas )
			if( typeof registry.gui_container != 'undefined'&& registry.gui_container != null ){
				resizeCanvas ( registry.gui_container )
			}
			scaleBackground ( registry.background , resImage ( "background"), registry.background_options )
			if( activeScene && typeof activeScene.screenResponse == 'function'){
				activeScene.screenResponse ()
			}
			if( activePopup && typeof activePopup.screenResponse == 'function'){
				activePopup.screenResponse ()
			}
		}
	}
	console.log ( "config changed: w:", config.design_w , "h:", config.design_h )
}
function changeBackground( url ){
	console.log ( "background changed:", url )
	registry.background.src = url
	registry.background.onload =
	()=>{
		scaleBackground ( registry.background , null , registry.background_options )
	}
}
function changeFPS( fps ){
	if( typeof fps !== 'number'|| fps <= 0 ){
		return
	}
	game.fps = fps
	if( interval_handle > 0 ){
		clearInterval ( interval_handle )
	}
	interval_handle = setInterval ( update , 1000 / game.fps )
}
function update(){
	if( game.is_invalid ){
		return
	}
	if(! canvas_stable ){
		if( typeof registry.canvas !== 'undefined'){
			resizeScreen ( config )
			resizeCanvas ( registry.canvas )
			resizeCanvas ( registry.gui_canvas )
			if( typeof registry.gui_container != 'undefined'&& registry.gui_container != null ){
				resizeCanvas ( registry.gui_container )
			}
			if( activeScene && typeof activeScene.screenResponse == 'function'){
				activeScene.screenResponse ()
			}
			if( activePopup && typeof activePopup.screenResponse == 'function'){
				activePopup.screenResponse ()
			}
			window.onresize =
			()=>{
				resizeScreen ( config )
				resizeCanvas ( registry.canvas )
				resizeCanvas ( registry.gui_canvas )
				if( typeof registry.gui_container != 'undefined'&& registry.gui_container != null ){
					resizeCanvas ( registry.gui_container )
				}
				scaleBackground ( registry.background , resImage ( "background"), registry.background_options )
				if( activeScene && typeof activeScene.screenResponse == 'function'){
					activeScene.screenResponse ()
				}
				if( activePopup && typeof activePopup.screenResponse == 'function'){
					activePopup.screenResponse ()
				}
			}
			if( registry.canvas.width > 0 && registry.canvas.height > 0 ){
				registry.container.style.transform = registry.fix_container_transform
				registry.background.style.width = registry.screen_size.w + 'px'
				registry.background.style.height = registry.screen_size.h + 'px'
				registry.background.width = registry.screen_size.w
				registry.background.height = registry.screen_size.h
				canvas_stable = true
			}
		}
		return
	}
	if( game.is_end ){
		return
	}
	var delta =( Date.now ()- game.last_update_time )/ 1000.0
	if(! registry.all_loaded ){
		registry.loading_time += delta
		registry.all_loaded = true
		if( spine != undefined ){
			if(! spineLoadingComplete ()){
				game.last_update_time = Date.now ()
				registry.all_loaded = false
				return
			}
		}
		var loaded = resLoaded ()
		Object.keys ( loaded ).forEach (
		( img_name )=>{
			var img = loaded [ img_name ]
			if( img instanceof Image && img.width == 0 ){
				if(! report_resource && Date.now ()- game.last_update_time > 3000 ){
					window.LogSocket.send ( "warning load resource")
					report_resource = true
				}
				game.last_update_time = Date.now ()
				registry.all_loaded = false
			}
			else if(! background_notify && img_name == "background"){
				background_notify = true
				scaleBackground ( registry.background , resImage ( "background"), registry.background_options )
				console.log ( "background loaded")
			}
		}
		)
		game.last_update_time = Date.now ()
		return
	}
	if( registry.loading_time < 1.5 ){
		registry.loading_time += delta
		game.last_update_time = Date.now ()
		return
	}
	else if( registry.game_state == State.state_loading ){
		registry.loading.style.display = "none"
		let searchParams = new URL ( window.location.href ).searchParams
		var urlSession = searchParams.get ( "session")
		session = urlSession
		var urlSessionType = searchParams.get ( "type")
		session_type = urlSessionType
		login ( res =>{
			if( res == null || typeof res.profile == 'undefined'){
				popupError ()
				return
			}
			if( res.profile ){
				game_state.title = res.profile.full_name
				if( typeof res.profile.avatar == 'string'&& res.profile.avatar != game_state.avatar ){
					game_state.avatar = res.profile.avatar
					var resAvatar = new Image ()
					resAvatar.onload = function (){
						resPatch ( "avatar", resAvatar )
					}
					resAvatar.src = game_state.avatar
				}
				game_state.num_play = 0
				for( var i = 0;i < res.game_turns.length;i ++){
					if( res.game_turns [ i ].game_id == game_id ){
						game_state.num_play = res.game_turns [ i ].num_turn
						break
					}
				}
			}
			changeScene ( firstScene )
		})
		window.LogSocket.send ( "resource loaded")
		registry.game_state = State.state_begin
	}
	if( typeof sndUpdate === 'function'){
		sndUpdate ()
	}
	if( typeof animUpdate === 'function'){
		animUpdate ()
	}
	registry.context = registry.bg_context
	var ctx = registry.context
	var game_size = registry.game_size
	ctx.clearRect ( 0 , 0 , game_size.w , game_size.h )
	ctx.translate ( registry.x0 , registry.y0 )
	if( activePopup != null ){
		activePopup.input ()
		activePopup.update ( delta )
	}
	if( activeScene != null ){
		if( activePopup == null ){
			activeScene.input ()
			activeScene.update ( delta )
		}
		activeScene.render ()
	}
	ctx.translate (- registry.x0 ,- registry.y0 )
	registry.context = registry.gui_context
	var ctx2 = registry.context
	ctx2.clearRect ( 0 , 0 , game_size.w , game_size.h )
	if( activePopup != null ){
		ctx2.translate ( registry.x0 , registry.y0 )
		activePopup.render ()
		ctx2.translate (- registry.x0 ,- registry.y0 )
	}
	if( typeof __console == 'object'&& typeof __console.draw == 'function'){
		__console.draw ( registry.x0 , registry.y0 , registry.game_size.w )
	}
	game.last_update_time = Date.now ()
	game.time =( game.last_update_time - game.begin_time )/ 1000
}
function goBack(){
	game.is_invalid = true
	if( interval_handle > 0 ){
		clearInterval ( interval_handle )
		interval_handle = 0
	}
	resUnloadAll ()
	history.back ()
}
registry.canvas = document.getElementById ( "bg_canvas")
registry.gui_canvas = document.getElementById ( "gui_canvas")
registry.loading = document.getElementById ( "loading")
registry.container = document.getElementById ( "container")
registry.background = document.getElementById ( "background")
registry.bg_context = registry.canvas.getContext ( "2d")
registry.gui_context = registry.gui_canvas.getContext ( "2d")
registry.context = registry.bg_context
setCanvas ( registry.canvas )
start ()
registry.on ( "touchBegin",
( touch )=>{
	touch.x -= registry.x0
	touch.y -= registry.y0
	if( activePopup != null ){
		activePopup.touchBegin ( touch )
		return
	}
	if( activeScene != null ){
		activeScene.touchBegin ( touch )
	}
}
)
registry.on ( "touchMove",
( touch )=>{
	touch.x -= registry.x0
	touch.y -= registry.y0
	if( activePopup != null ){
		activePopup.touchMove ( touch )
		return
	}
	if( activeScene != null ){
		activeScene.touchMove ( touch )
	}
}
)
var isFixSound = false
registry.on ( "touchEnd",
( touch )=>{
	touch.x -= registry.x0
	touch.y -= registry.y0
	if(! isFixSound ){
		isFixSound = true
	}
	if( activePopup != null ){
		activePopup.touchEnd ( touch )
		return
	}
	if( activeScene != null ){
		activeScene.touchEnd ( touch )
	}
}
)
registry.on ( "touchCancel",
( touch )=>{
	touch.x -= registry.x0
	touch.y -= registry.y0
	if( activePopup != null && typeof activePopup.touchCancel == 'function'){
		activePopup.touchCancel ( touch )
		return
	}
	if( activeScene != null && typeof activeScene.touchCancel == 'function'){
		activeScene.touchCancel ( touch )
	}
}
)
function playGame( game , playContext , callback ){
	callback ({
	})
}
var firstScene = SceneShakingTutorial