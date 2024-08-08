import { useEffect, useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import { CurrentSong } from "../wailsjs/go/main/App";

type Metadata = { [key: string]: string }

function App() {
    const [metadata, setMetadata] = useState<Metadata>();
    const [posBar, setPosBar] = useState<number>()
    const updateResultText = (result: Metadata) => setMetadata(result);


    const fetchCurrentSong = async () => {
        try {
            const result = await CurrentSong();
            updateResultText(result);
            setPosBar((parseFloat(result["position"]) / parseFloat(result["length"])) * 100)
        } catch (error) {
            console.error("Error fetching current song:", error);
        }
    };

    useEffect(() => {
        const intervalId = setInterval(fetchCurrentSong, 1000); // Fetch every second

        return () => clearInterval(intervalId); // Cleanup interval on component unmount
    }, []);

    console.log(posBar)

    return metadata ? (
        <div id="App" className='App'>
            <img src={metadata ? metadata["albumCover"] : logo} id="logo" alt="logo" />
            <div id="title" className="title">{metadata["title"]}</div>
            <div className='status'><span>{metadata["positionF"]}</span><div className='bar'><div className='innerBar' style={{ width: `${posBar}%` }}></div></div><span>{metadata["lengthF"]}</span></div>
        </div>
    ) : null
}

export default App
