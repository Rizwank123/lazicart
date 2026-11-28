const path = require("path");
const { getDefaultConfig } = require("expo/metro-config");

const config = getDefaultConfig(__dirname);

/**
 * Zustand v5 ships an ESM build that reads `import.meta.env.MODE`.
 *
 * The web bundle is served as a classic script, not as an ES module, so
 * `import.meta` is a syntax error there. The whole bundle then fails to
 * parse and the app only ever shows the splash screen.
 *
 * The package also ships a CommonJS build without `import.meta`. The native
 * build already picks that one up through the "react-native" export
 * condition, but the web build resolves the "import" condition and ends up
 * with the ESM file instead.
 *
 * This resolver points the web build at the CommonJS file, so every platform
 * runs the same code and no build ships an `import.meta` expression.
 */
const zustandRoot = path.dirname(require.resolve("zustand/package.json"));
const zustandCommonJsMiddleware = path.join(zustandRoot, "middleware.js");

const previousResolveRequest = config.resolver.resolveRequest;

config.resolver.resolveRequest = (context, moduleName, platform) => {
  if (platform === "web" && moduleName === "zustand/middleware") {
    return { type: "sourceFile", filePath: zustandCommonJsMiddleware };
  }
  if (previousResolveRequest) {
    return previousResolveRequest(context, moduleName, platform);
  }
  return context.resolveRequest(context, moduleName, platform);
};

module.exports = config;
