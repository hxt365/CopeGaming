import React, { useState, useEffect } from "react";
import { getProviderList } from "../../services/api/providers";

import "./style.scss";

export default function ProviderList({ onSelectProvider }) {
  const [providers, setProviders] = useState([]);

  useEffect(async () => {
    const resp = await getProviderList();
    if (resp.errorCode === undefined || resp.errorCode === 0) {
      if (resp.data?.providers !== null) {
        setProviders(resp.data.providers);
      }
    }
  }, []);

  return (
    <div className="provider-list">
      <ul className="provider-list__list">
        {providers.map((p) => {
          return (
            <li
              key={p.id}
              className="provider-list__list-item"
              onClick={() => {
                onSelectProvider(p.id);
              }}
            >
              {p.id}
            </li>
          );
        })}
      </ul>
    </div>
  );
}