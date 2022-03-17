import React from "react";
import animatedLogo from "../../assets/animated_logo.gif";

import "./style.scss";

export default function Welcoming() {
  return (
    <div className="welcoming">
      <img src={animatedLogo} alt="Welcoming" />
    </div>
  );
}
