import React, { useEffect } from "react";
import useTimer from "easytimer-react-hook";

export default function TimeCounter() {
  const [timer] = useTimer({});

  useEffect(() => {
    timer.start({});

    return () => timer.stop();
  }, []);

  return <div>{timer.getTimeValues().toString()}</div>;
}
