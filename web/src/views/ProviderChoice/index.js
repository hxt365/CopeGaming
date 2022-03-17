import React from "react";
import ProviderList from "../../components/ProviderList";

import "./style.scss";

export default function ProviderChoice({ onSelectProvider }) {
  return (
    <div className="provider-choice">
      <h1 className="provider-choice__title">Choose a provider</h1>
      <ProviderList onSelectProvider={onSelectProvider} />
    </div>
  );
}
