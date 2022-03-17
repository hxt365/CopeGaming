import React, { useState, useEffect } from "react";
import AppChoice from "./views/AppChoice";
import AppPlayer from "./views/AppPlayer";

import { decodeBase64, encodeBase64 } from "./utils";
import { addRemoteSdp, addIceCandidate } from "./services/webrtc";
import { getDevice } from "./services/api/apps";

import "./App.scss";
import Welcoming from "./views/Welcoming";

function App() {
  const [welcoming, setWelcoming] = useState(true);
  const [ws, setWs] = useState(null);
  const [pc, setPc] = useState(null);
  const [inpChannel, setInpChannel] = useState(null);
  const [videoStream, setVideoStream] = useState(null);
  const [selectedApp, setSelectedApp] = useState("");

  useEffect(() => {
    setTimeout(() => {
      setWelcoming(false);
    }, 2500);
  });

  useEffect(() => {
    const ws = new WebSocket(process.env.REACT_APP_WS_ENDPOINT);

    ws.onopen = () => {
      setWs(ws);
    };

    ws.onerror = () => {
      throw Error("Failed to connect to the server");
    };

    return () => ws.close();
  }, []);

  useEffect(() => {
    if (pc === null) return;

    const msg = {
      type: "start",
      data: JSON.stringify({
        appID: selectedApp,
        device: getDevice(),
      }),
    };
    ws.send(JSON.stringify(msg));

    ws.onmessage = async (event) => {
      const msg = JSON.parse(event.data);
      if (msg.type === "sdp") {
        const offer = JSON.parse(decodeBase64(msg.data));
        const answer = await addRemoteSdp(pc, offer);
        ws.send(
          JSON.stringify({
            type: "sdp",
            data: encodeBase64(JSON.stringify(answer)),
          })
        );
      } else if (msg.type === "ice-candidate") {
        const ice = JSON.parse(decodeBase64(msg.data));
        addIceCandidate(pc, ice);
      }
    };
  }, [pc]);

  useEffect(() => {
    if (inpChannel === null) return;

    const onKeyDown = (event) => {
      if (inpChannel.readyState !== "open") return;

      inpChannel.send(
        JSON.stringify({
          type: "KEYDOWN",
          data: JSON.stringify({
            keyCode: event.keyCode,
          }),
        })
      );
    };

    const onKeyUp = (event) => {
      if (inpChannel.readyState !== "open") return;

      inpChannel.send(
        JSON.stringify({
          type: "KEYUP",
          data: JSON.stringify({
            keyCode: event.keyCode,
          }),
        })
      );
    };

    document.addEventListener("keydown", onKeyDown);
    document.addEventListener("keyup", onKeyUp);

    return () => {
      document.removeEventListener("keydown", onKeyDown);
      document.removeEventListener("keyup", onKeyUp);
    };
  }, [inpChannel]);

  const startApp = async (appId) => {
    console.log("Start playing", appId);

    const newPc = new RTCPeerConnection({
      iceServers: [
        {
          urls: "stun:stun.l.google.com:19302",
        },
      ],
    });

    newPc.ondatachannel = (event) => {
      const channel = event.channel;
      if (channel.label === "app-input") {
        channel.onopen = () => {
          console.log("got input datachannel");
          setInpChannel(channel);
        };
      } else if (channel.label === "health-check") {
        let healthCheckIntId;

        channel.onopen = () => {
          console.log("got health-check datachannel");
          healthCheckIntId = setInterval(() => {
            channel.send({});
          }, 2000);
        };

        channel.onclose = () => {
          clearInterval(healthCheckIntId);
        };
      }
    };

    newPc.ontrack = (event) => {
      console.log("got track", event);
      if (event.streams && event.streams[0]) {
        setVideoStream(event.streams[0]);
      } else {
        let inboundStream = null;
        if (!videoStream) {
          inboundStream = new MediaStream();
          inboundStream.addTrack(event.track);
        } else {
          inboundStream = { ...videoStream };
          inboundStream.addTrack(event.track);
        }
        setVideoStream(inboundStream);
      }
    };

    newPc.onicecandidate = (event) => {
      const iceCandidate = event.candidate;

      if (iceCandidate) {
        ws.send(
          JSON.stringify({
            type: "ice-candidate",
            data: encodeBase64(JSON.stringify(iceCandidate)),
          })
        );
      }
    };

    newPc.oniceconnectionstatechange = (event) => {
      console.log(event.target.iceConnectionState);
    };

    setPc(newPc);
  };

  const selectApp = (appId) => {
    setSelectedApp(appId);
    startApp(appId);
  };

  const closeApp = () => {
    if (pc !== null) {
      pc.close();
    }

    setPc(null);
    setVideoStream(null);
    setInpChannel(null);
    setSelectedApp("");
  };

  return (
    <div className="App">
      {welcoming && <Welcoming />}
      {selectedApp !== "" ? (
        <AppPlayer
          videoStream={videoStream}
          inpChannel={inpChannel}
          onCloseApp={closeApp}
        />
      ) : (
        <AppChoice onSelectApp={selectApp} />
      )}
    </div>
  );
}

export default App;
