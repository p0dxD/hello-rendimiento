# hello-rendimiento

A tiny Go web app for testing [rendimiento.ai](https://rendimiento.joserod.space) end to end.

- `GET /`: greeting page (set `GREETING` to change it)
- `GET /healthz`: health check
- `GET /api/info`: JSON

It has no Dockerfile on purpose: rendimiento detects Go and adds one in the onboarding PR.
