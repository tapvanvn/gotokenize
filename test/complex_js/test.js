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