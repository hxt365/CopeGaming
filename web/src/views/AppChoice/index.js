import React from "react";

import AppList from "../../components/AppList";

import "./style.scss";

export default function AppChoice({ onSelectApp, children }) {
  return (
    <div className="app-choice">
      <h1 className="app-choice__title">Choose a game</h1>
      <AppList onSelectApp={onSelectApp} />
      {children}
    </div>
  );
}
