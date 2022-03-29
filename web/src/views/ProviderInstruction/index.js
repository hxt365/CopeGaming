import React from "react";

import "./style.scss";

export default function ProviderChoice({ onBack }) {
  return (
    <div className="provider-instruction">
      <button className="provider-instruction__close" onClick={onBack}>
        Back
      </button>
      <h1>
        Follow below intructions to become a provider and make money while you
        sleep
      </h1>
      <p>
        Becoming a provider means that you offer your computer for other users
        to play games on it and earn money as &ldquo;renting fee&rdquo;.
        However, games are run in background and you can still use your computer
        as normal.
      </p>
      <br />
      <p> (Notice that we only support Linux computers for now)</p>
      <ol>
        <li>Install Golang (v1.16 or above)</li>
        <li>Install Docker and docker-compose</li>
        <li>
          Make sure that the current user on your computer has permissions to
          run Docker
        </li>
        <li>
          Clone <a href="https://github.com/hxt365/CopeGaming">the project</a>{" "}
          from Github
        </li>
        <li>Run provider/run.sh in the cloned folder</li>
        <li>You&apos;re now a provider</li>
      </ol>
    </div>
  );
}
