import React, { useState, useEffect } from "react";
import { getProviderList } from "../../services/api/providers";
import ProviderListItem from "../ProviderListItem";

import "./style.scss";

export default function ProviderList({ onSelectProvider, ownerID }) {
  const [providers, setProviders] = useState([]);

  useEffect(async () => {
    const resp = await getProviderList(ownerID);
    if (resp.errorCode === undefined || resp.errorCode === 0) {
      if (resp.data?.providers !== null) {
        setProviders(resp.data.providers);
      }
    }
  }, [ownerID]);

  return (
    <div className="provider-list">
      <ul className="provider-list__list">
        {providers.map((p) => {
          return (
            <li
              key={p.id}
              className="provider-list__list-item"
              onClick={() => {
                if (onSelectProvider !== undefined) onSelectProvider(p.id);
              }}
            >
              <ProviderListItem provider={p} />
            </li>
          );
        })}
      </ul>
    </div>
  );
}
