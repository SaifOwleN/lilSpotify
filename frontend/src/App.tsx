import { useEffect, useState } from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import { ChangeState, CurrentSong, NextSong, OpenApp, PrevSong, Seek } from "../wailsjs/go/main/App";
import { IoPauseOutline, IoPlayOutline, IoPlaySkipBackOutline, IoPlaySkipForwardOutline } from "react-icons/io5";

type appStatus = "opened" | "closed"
type playbackStatus = "Paused" | "Playing"

export type Metadata = {
  "appStatus": appStatus
  "status": playbackStatus
  "artist": string
  "albumCover": string
  "title": string
  "album": string
  "length": string
  "lengthF": string
  "position": string
  "positionF": string

}

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

  const nextSong = async () => {
    await NextSong();
  }
  const prevSong = async () => {
    await PrevSong();
  }

  const seek = async (event: React.MouseEvent<HTMLDivElement>) => {
    const bar = (event.target as HTMLElement).getBoundingClientRect();

    const x = (event.clientX - bar.left) / bar.width
    console.log(x * parseInt(metadata["lengthR"]))
    await Seek(Math.floor(x * parseInt(metadata["lengthR"])), metadata["trackId"])

  }


  const fetchCurrentSong = async () => {
    try {
      const result = await CurrentSong();
      console.log(result)
      setMetadata(result)
      console.log(result["status"] === "Paused")
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
    const intervalId = setInterval(fetchCurrentSong, 1000);

    return () => clearInterval(intervalId);
  }, []);

  useEffect(() => {
    const handlePauseListner = (event: KeyboardEvent) => {
      if (event.code === "Space") changeState()
    }
    window.addEventListener('keydown', handlePauseListner)
    return () => window.removeEventListener('keydown', handlePauseListner)
  })


  return appStatus === true && metadata ? (
    <div id="App" className='App'>
      <div className='artContainer'>

        <img src={metadata ? metadata["albumCover"] : logo} id="art" alt="album cover" />
        <div className='controls'>
          <button className='playback' onClick={prevSong}><IoPlaySkipBackOutline /></button>
          <button className='playback' onClick={changeState}>{status ? <IoPauseOutline /> : <IoPlayOutline />}</button>
          <button className='playback' onClick={nextSong}><IoPlaySkipForwardOutline /></button>
        </div>
      </div>
      <div className='bot'>
        <div id="title" className="title">{metadata["title"]}</div>
        <div className='status'>
          <div onClick={seek} className='bar'>
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
