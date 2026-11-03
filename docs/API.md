# Native API

Base: `/voice/api/` when managed, `/api/` standalone. Every endpoint except
`bootstrap` and `/healthz` requires the app session cookie. Mutations require the
exact origin and `X-CSRF-Token` from `GET me`; send `X-Workspace-Expected-Account`
with its user ID to reject a split browser session. MyAccounts controls managed
users; standalone local users require the apphost owner's explicit action.

| Endpoint | Purpose |
|---|---|
| GET /healthz | Database/process health, app ID/version; outside API base |
| GET models | Available/installed engines and licence/runtime metadata |
| POST models/{id}/install | Owner-only `{license:true}`; queued download status in models |
| DELETE models/{id} | Owner-only model removal; retains voices/recordings |
| GET voices | Presets and voices owned by or granted to this user |
| POST voices | Multipart audio, name, tags, model, consent=true; validates and encodes |
| PATCH / DELETE voices/{id} | Owner rename/tags/favourite or permanent deletion |
| POST voices/{id}/favorite | `{favorite:true}`; account-specific preset favourite |
| GET voices/{id}/reference | Owner-only decoded reference WAV |
| POST voices/{id}/preview | Queue a permission-checked speech sample |
| POST voices/{id}/share | Owner grant by `{user:"OneAccount email"}`; user must have opened app |
| DELETE voices/{id}/share/{userID} | Owner revocation |
| GET / POST projects | List/save owner-scoped script projects |
| GET / POST dictionaries | Account exact text replacements `[{from,to}]` |
| POST tts/jobs | Queue speech; see below |
| GET tts/jobs/{id} | Owner-scoped status, progress, completed segments, duration, waveform |
| POST tts/jobs/{id}/cancel | Idempotent stop request |
| POST tts/jobs/{id}/retry | Resume failed/cancelled job with preserved finished segments |
| GET history?q= | Search up to 200 recent jobs |
| DELETE history/{id} | Owner permanent removal; running/queued jobs must stop first |
| GET outputs/{id} | Owner audio playback |
| GET outputs/{id}/download?format=wav\|mp3 | Owner attachment export |
| GET outputs/{id}/segments/{segmentID} | Owner PCM segment |
| GET exports/zip?ids=ID,ID | Up to 20 complete owner recordings |
| GET / PATCH settings | Performance/storage details; owner-only thread update |

Speech request:

```json
{"requestId":"unique-retry-id","name":"Welcome","text":"Hello from My Voice.","voice":"alba","model":"pocket18","language":"en","speed":1,"pitch":0,"volume":1,"pause":0.3,"dictionary":[{"from":"OWS","to":"One Work space"}],"segments":[]}
```

For dialogue, supply `{id,text,voice}` segments using voices compatible with the
selected engine. To regenerate one segment, preserve model/audio controls,
dictionary, segment order and all other segments, and add `parent` (complete job
ID) plus `replace` (selected segment ID). A new job owns the rebuilt exports.
`requestId` is optional; repeating it for the same user returns the original job.

Other Workspace apps can use the existing curated, authenticated My AI native
connector (voice list, generation, history, job status). This adds no blanket LAN
bearer token or arbitrary HTTP proxy. General external token issuance is not
included. Model installation/settings remain owner-controlled.
