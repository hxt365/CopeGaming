import React from "react";

import "./style.scss";

export default function Poster({ src, width = 300, height = 350 }) {
  return (
    <div className="poster">
      <img src={src} style={{ height: height, width: width }} />
    </div>
  );
}
