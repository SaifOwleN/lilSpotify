import { useEffect, useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import { ChangeState, CurrentSong, OpenApp } from "../wailsjs/go/main/App";

type Metadata = { [key: string]: string }

function App() {
    const [appStatus, setAppStatus] = useState<boolean | undefined>();
    const [status, setStatus] = useState<boolean | undefined>();
    const [metadata, setMetadata] = useState<Metadata>();
    const [posBar, setPosBar] = useState<number>()

    const openSpotify = async () => {
        await OpenApp();
    }

    const changeState = async () => {
        await ChangeState();
        setStatus(!status)
    }


    const fetchCurrentSong = async () => {
        try {
            const result = await CurrentSong();
            setMetadata(result)
            console.log(result)
            setStatus(result["status"] === "Playing")
            if (result["appStatus"] === "closed") {
                setAppStatus(false)
            } else {
                setAppStatus(true)
                setPosBar((parseFloat(result["position"]) / parseFloat(result["length"])) * 100)
            }
        } catch (error) {
            console.error("Error fetching current song:", error);
        }
    };

    useEffect(() => {
        const intervalId = setInterval(fetchCurrentSong, 1000); // Fetch every second

        return () => clearInterval(intervalId); // Cleanup interval on component unmount
    }, []);

    console.log(posBar)

    return appStatus === true && metadata ? (
        <div id="App" className='App'>
            <div className='artContainer'>

                <img src={metadata ? metadata["albumCover"] : logo} id="art" alt="album cover" />
                <div className='controls'>
                    <button className='playback' onClick={changeState}>{status ? "V" : "D"}</button>
                </div>
            </div>
            <div className='bot'>
                <div id="title" className="title">{metadata["title"]}</div>
                <div className='status'>
                    <div className='bar'>
                        <div className='innerBar' style={{ width: `${posBar}%` }}></div>
                    </div>
                    <div className='position'>
                        <span>{metadata["positionF"]}</span>
                        <span>{metadata["lengthF"]}</span>
                    </div>
                </div>
            </div>
        </div>
    ) : appStatus === false ? (<div className='App'><button className='openButton' onClick={() => openSpotify()}>Open Spotify</button></div>) : null
}

export default App
