import React from "react";
import ProviderList from "../../components/ProviderList";

import "./style.scss";

export default function ProviderChoice({ onSelectProvider, onBack }) {
  return (
    <div className="provider-choice">
      <button className="provider-choice__close" onClick={onBack}>
        Back
      </button>
      <h1 className="provider-choice__title">Choose a game server</h1>
      <div className="provider-choice__provider-list">
        <ProviderList onSelectProvider={onSelectProvider} />
      </div>
    </div>
  );
}
