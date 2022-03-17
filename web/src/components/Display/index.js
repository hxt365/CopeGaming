import React, { useRef, useEffect } from "react";

import loading from "../../assets/loading.svg";

import "./style.scss";

export default function Display({ streamSrc, inpChannel }) {
  const videoRef = useRef(null);

  useEffect(() => {
    if (!streamSrc) return;

    videoRef.current.srcObject = streamSrc;
  }, [streamSrc]);

  const sendMouseDown = (data) => {
    if (!inpChannel) return;

    inpChannel.send(
      JSON.stringify({
        type: "MOUSEDOWN",
        data: JSON.stringify(data),
      })
    );
  };

  const sendMouseUp = (data) => {
    if (!inpChannel) return;

    inpChannel.send(
      JSON.stringify({
        type: "MOUSEUP",
        data: JSON.stringify(data),
      })
    );
  };

  const sendMouseMove = (data) => {
    if (!inpChannel) return;

    inpChannel.send(
      JSON.stringify({
        type: "MOUSEMOVE",
        data: JSON.stringify(data),
      })
    );
  };

  const isLeftButton = (button) => {
    return button === 0 ? 1 : 0; // 1 is right button
  };

  const onMouseDown = (event) => {
    const boundRect = videoRef.current.getBoundingClientRect();
    sendMouseDown({
      isLeft: isLeftButton(event.button),
      x: event.clientX - boundRect.left,
      y: event.clientY - boundRect.top,
      width: boundRect.width,
      height: boundRect.height,
    });
  };

  const onMouseUp = (event) => {
    const boundRect = videoRef.current.getBoundingClientRect();
    sendMouseUp({
      isLeft: isLeftButton(event.button),
      x: event.clientX - boundRect.left,
      y: event.clientY - boundRect.top,
      width: boundRect.width,
      height: boundRect.height,
    });
  };

  const onMouseMove = (event) => {
    const boundRect = videoRef.current.getBoundingClientRect();
    sendMouseMove({
      isLeft: isLeftButton(event.button),
      x: event.clientX - boundRect.left,
      y: event.clientY - boundRect.top,
      width: boundRect.width,
      height: boundRect.height,
    });
  };

  const onContextMenu = (event) => {
    event.preventDefault();
    return false;
  };

  return (
    <div className="display">
      <video
        className="display__video"
        poster={loading}
        ref={videoRef}
        autoPlay
        onMouseDown={onMouseDown}
        onMouseUp={onMouseUp}
        onMouseMove={onMouseMove}
        onContextMenu={onContextMenu}
      />
    </div>
  );
}
