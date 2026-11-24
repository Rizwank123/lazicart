# LahziCart Mobile

Production-oriented Expo React Native app integrated with the LahziCart Go microservices backend.

## Stack

- Expo SDK 55
- Expo Router
- TypeScript strict mode
- TanStack Query
- Zustand
- Axios
- React Hook Form + Zod
- SecureStore
- AsyncStorage
- Expo Image
- FlashList
- Reanimated + Gesture Handler

## Setup

1. Install dependencies:

```bash
npm install
```

2. Create `.env` from `.env.example`:

```bash
EXPO_PUBLIC_API_BASE_URL=http://localhost:7070
EXPO_PUBLIC_DEFAULT_CURRENCY=INR
```

For Android emulator, use your machine IP or `http://10.0.2.2:7070` if the gateway runs on the host.

3. Start the app:

```bash
npm start
```

## Quality Checks

```bash
npm run typecheck
npm run lint
```

## Backend Notes

See `docs/BACKEND_ANALYSIS.md` for the endpoint map, business rules, and missing backend APIs discovered before implementation.
