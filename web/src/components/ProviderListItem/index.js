import React from "react";

import "./style.scss";

const RED = "#e50914";
const YELLOW = "#f6d70b";
const WHITE = "#f5f5f1";

export default function ProviderListItem({ provider }) {
  const getColor = (cpu, mem) => {
    if (cpu >= 75 || mem >= 75) return RED;
    if (cpu >= 50 || mem >= 50) return YELLOW;
    return WHITE;
  };

  return (
    <div
      className="provider-list-item"
      style={{ color: getColor(provider.cpuPercent, provider.memPercent) }}
    >
      <div>
        <span className="provider-list-item__host-name">
          {provider.hostName}
        </span>
        <br />
        <span className="provider-list-item__info">{`(${provider.platform}, ${provider.cpuName}, ${provider.memSize}GB RAM)`}</span>
      </div>
      <div className="provider-list-item__stats">
        <span>CPU: {Math.ceil(provider.cpuPercent)}%</span>
        <br />
        <span>Mem: {Math.ceil(provider.memPercent)}%</span>
      </div>
    </div>
  );
}
