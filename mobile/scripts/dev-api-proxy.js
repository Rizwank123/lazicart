#!/usr/bin/env node
/**
 * Development API proxy for the web build.
 *
 * A browser may only read a response from another origin when the server
 * answers the CORS preflight that the browser sends before a request with a
 * body or with an authorization header. The API gateway does not answer that
 * preflight, so every sign in, every add to cart and every other request
 * that sends a body is blocked by the browser before it reaches the gateway.
 *
 * This proxy sits between the browser and the gateway. It answers the
 * preflight and adds the CORS headers to every response, so the web build can
 * use the whole API without the gateway having to know about browsers.
 *
 * It is a development helper. A production deployment serves the web build
 * and the API from the same origin instead, which needs no proxy at all.
 *
 * Environment:
 *   API_PROXY_PORT    the port to listen on, 7071 by default
 *   API_PROXY_TARGET  the gateway to forward to, http://localhost:7070
 */
const http = require("http");

const PORT = Number(process.env.API_PROXY_PORT ?? 7071);
const TARGET = new URL(process.env.API_PROXY_TARGET ?? "http://localhost:7070");

// The headers the browser asks about in the preflight. Idempotency-Key is
// included because the checkout sends it.
const CORS_HEADERS = {
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Methods": "GET,POST,PUT,PATCH,DELETE,OPTIONS",
  "Access-Control-Allow-Headers": "Origin,Content-Type,Accept,Authorization,X-User-Id,Idempotency-Key",
  "Access-Control-Expose-Headers": "Content-Type",
  "Access-Control-Max-Age": "86400",
};

function forward(req, res) {
  // The browser asks for /api/... and the gateway serves /api/... as well,
  // so the path is passed through unchanged.
  const headers = { ...req.headers, host: TARGET.host };
  delete headers.connection;

  const upstream = http.request(
    {
      protocol: TARGET.protocol,
      hostname: TARGET.hostname,
      port: TARGET.port || (TARGET.protocol === "https:" ? 443 : 80),
      method: req.method,
      path: req.url,
      headers,
    },
    (upstreamRes) => {
      const responseHeaders = { ...upstreamRes.headers };
      // The gateway copies the CORS headers of the service it called, so a
      // response can already carry them. They are removed first, because the
      // browser rejects a response whose Access-Control-Allow-Origin header
      // holds more than one value.
      for (const name of Object.keys(responseHeaders)) {
        if (name.toLowerCase().startsWith("access-control-")) {
          delete responseHeaders[name];
        }
      }
      Object.assign(responseHeaders, CORS_HEADERS);
      // The body is streamed, so the length is not known up front.
      delete responseHeaders["transfer-encoding"];
      res.writeHead(upstreamRes.statusCode ?? 502, responseHeaders);
      upstreamRes.pipe(res);
    },
  );

  upstream.on("error", (error) => {
    if (res.headersSent) {
      res.destroy();
      return;
    }
    res.writeHead(502, { ...CORS_HEADERS, "Content-Type": "application/json" });
    res.end(JSON.stringify({ message: `The API gateway is unreachable: ${error.message}` }));
  });

  req.pipe(upstream);
}

const server = http.createServer((req, res) => {
  if (req.method === "OPTIONS") {
    res.writeHead(204, CORS_HEADERS);
    res.end();
    return;
  }
  forward(req, res);
});

server.listen(PORT, () => {
  console.log(`API proxy listening on http://localhost:${PORT} -> ${TARGET.origin}`);
});

for (const signal of ["SIGINT", "SIGTERM"]) {
  process.on(signal, () => {
    server.close(() => process.exit(0));
  });
}
