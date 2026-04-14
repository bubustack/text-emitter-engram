# 📝 Text Emitter Engram

A lightweight Engram that emits configured text (greetings, banners, or status
messages) without calling external services. Useful for demos, placeholders, or
control flows within a Story.

## 🌟 Highlights

- Supports batch and streaming deployments with the same config.
- Emits deterministic text driven entirely by `spec.with`.
- Can introduce controlled delays or disable emission dynamically.

## 🚀 Quick Start

```bash
go test ./...
```

Apply `Engram.yaml`, reference the template in a Story step, and supply the
message you want emitted.

## ⚙️ Configuration (`Engram.spec.with`)

| Field | Default | Description |
| --- | --- | --- |
| `message` | `"Hello! How can I help you today?"` | Text emitted when the step runs. |
| `enabled` | `true` | Controls whether the message is emitted. |
| `delayMs` | `0` | Delay before emitting the message (ms). |

## 📥 Inputs

- Accepts any object payload; unused by default but available for overrides.

## 📤 Outputs

- Returns an object with `type` (defaults to `text.emitted`) and `text`.

## 🧪 Local Development

- `go test ./...` – Run unit tests.
- `go vet ./...` – Static analysis pass before release.

## 🤝 Community & Support

- [Contributing](./CONTRIBUTING.md)
- [Support](./SUPPORT.md)
- [Security Policy](./SECURITY.md)
- [Code of Conduct](./CODE_OF_CONDUCT.md)
- [Discord](https://discord.gg/dysrB7D8H6)

## 📄 License

Copyright 2025 BubuStack.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
