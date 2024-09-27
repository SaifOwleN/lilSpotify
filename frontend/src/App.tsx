import { useEffect, useState } from 'react';
import { Init } from '../wailsjs/go/main/App';
import './App.css';
import Player from './player';


function App() {
	const [App, setApp] = useState<"dbus" | "api" | undefined>();


	function startDbus() {
		console.log("xd")
		setApp("dbus")
	}


	useEffect(() => {
		if (App == "api") {
			(async () => { await Init("api") })()
		}

	}, [App])



	return (
		<div className='home'>
			<button onClick={() => startDbus()}>Dbus (Linux Only)</button>
			<button onClick={() => setApp("api")}>Spotify Api</button>
		</div>
	)

}

export default App
