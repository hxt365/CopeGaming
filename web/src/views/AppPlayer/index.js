import React from "react";
import Display from "../../components/Display";
import TimeCounter from "../../components/TimeCounter";

import "./style.scss";

export default function AppPlayer({ videoStream, inpChannel, onCloseApp }) {
  return (
    <div className="app-player">
      <div className="app-player__timer">
        <TimeCounter />
      </div>
      <button className="app-player__close" onClick={onCloseApp}>
        Exit
      </button>
      <div className="app-player__display">
        <Display streamSrc={videoStream} inpChannel={inpChannel} />
      </div>
    </div>
  );
}
