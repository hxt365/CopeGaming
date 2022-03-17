import React from "react";
import ProviderList from "../../components/ProviderList";

import "./style.scss";

export default function ProviderChoice({ onSelectProvider, onBack }) {
  return (
    <div className="provider-choice">
      <button className="provider-choice__close" onClick={onBack}>
        Back
      </button>
      <h1 className="provider-choice__title">Choose a provider</h1>
      <ProviderList onSelectProvider={onSelectProvider} />
    </div>
  );
}
