import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

import * as Sentry from "@sentry/react";
import { Integrations } from "@sentry/tracing";

// Sentry SDK initialization
Sentry.init({
  dsn: process.env.REACT_APP_SENTRY_DSN, // Retreiving the DSN from .env file
  integrations: [new Integrations.BrowserTracing()], // Enabling automatic browser tracing
  tracesSampleRate: 1.0, // Capturing all events
});

ReactDOM.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
