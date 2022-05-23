import { React, useState } from "react";
import { DebounceInput } from "react-debounce-input";

import ProviderList from "../../components/ProviderList";

import "./style.scss";

export default function ServerList({ onBack }) {
  const [onwerID, setOwnerID] = useState("");

  const updateOwnerID = (e) => {
    setOwnerID(e.target.value);
  };

  return (
    <div className="server-list">
      <button className="server-list__close" onClick={onBack}>
        Back
      </button>
      <h1>
        List of servers of&nbsp;
        <DebounceInput
          className="server-list__owner-id"
          placeholder="Your ID"
          debounceTimeout={300}
          onChange={updateOwnerID}
        />
      </h1>
      <div className="server-list__provider-list">
        <ProviderList ownerID={onwerID} />
      </div>
    </div>
  );
}
