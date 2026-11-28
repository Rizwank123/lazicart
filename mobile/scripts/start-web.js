#!/usr/bin/env node
/**
 * Start the development API proxy and the Expo web server together.
 *
 * The web build needs the proxy, because the gateway does not answer the
 * CORS preflight that the browser sends. See dev-api-proxy.js.
 *
 * Both processes are stopped together, so no proxy is left behind.
 */
const { spawn } = require("child_process");
const path = require("path");

const proxy = spawn(process.execPath, [path.join(__dirname, "dev-api-proxy.js")], {
  stdio: "inherit",
});

const expo = spawn("npx", ["expo", "start", "--web"], {
  stdio: "inherit",
  shell: process.platform === "win32",
});

let stopping = false;

function stop(code) {
  if (stopping) return;
  stopping = true;
  proxy.kill("SIGTERM");
  expo.kill("SIGTERM");
  process.exit(code ?? 0);
}

expo.on("exit", (code) => stop(code));
proxy.on("exit", (code) => {
  if (code) stop(code);
});

process.on("SIGINT", () => stop(0));
process.on("SIGTERM", () => stop(0));
