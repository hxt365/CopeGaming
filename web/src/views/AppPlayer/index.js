import React from "react";

import Display from "../../components/Display";

import "./style.scss";

export default function AppPlayer({ videoStream, inpChannel, onCloseApp }) {
  return (
    <div className="app-player">
      <button className="app-player__close" onClick={onCloseApp}>
        Exit
      </button>
      <Display streamSrc={videoStream} inpChannel={inpChannel} />
    </div>
  );
}
