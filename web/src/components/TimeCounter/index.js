import React from "react";
import useTimer from "easytimer-react-hook";

export default function TimeCounter() {
  const [timer] = useTimer({});

  timer.start({});

  return <div>{timer.getTimeValues().toString()}</div>;
}
