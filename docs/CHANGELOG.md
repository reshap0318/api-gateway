# CHANGELOG

## v1.4.0 - 2026-07-09

- **Implementasi Phase 3B selesai (T-033 s/d T-036):** `05_ITL.md` diperbarui — seluruh task ditandai `[x]`. Semua fitur ini adalah delta pada modul User existing, reuse permission `user.edit` (tidak ada permission baru).
  - **Migration & Model (T-033):** `ALTER TABLE users` tambah `status` (`active`/`suspended`, default `active`), `failed_login_attempts` (default 0), `locked_until` (nullable). `UserDTO` diperluas dengan `status` dan `locked_until`.
  - **Update Status (T-034):** `PUT /api/users/:id/status` — independen dari mekanisme Lock. `AuthLogin` dicek status `suspended` sebelum validasi password → `401 "Akun tidak aktif, hubungi admin"`.
  - **Auto-Lock (T-035):** `registerFailedLogin` di `auth_service.go` — increment `failed_login_attempts` tiap password salah, lock otomatis (`locked_until = now + AUTH_LOCK_DURATION_MINUTES`) setelah mencapai `AUTH_MAX_FAILED_LOGIN_ATTEMPTS` (default 5/15 menit). Reset ke 0 saat login sukses. Cek lock dilakukan sebelum cek password di `AuthLogin`.
  - **Manual Unlock (T-036):** `POST /api/users/:id/unlock` — set `locked_until=NULL`, `failed_login_attempts=0`. Ditolak (`400`) jika akun memang tidak sedang terkunci.
  - **Frontend:** Panel "Status & Lock" baru di `pages/users/FormModal.vue` (muncul hanya saat edit) — Select Status dengan efek langsung (bukan bagian submit form utama), badge Lock Indicator (Locked/Unlocked), tombol "Unlock Now" (muncul hanya saat terkunci).
- **Diuji end-to-end:**
  - Playwright: ubah Status Viewer ke Suspend → toast sukses → (via curl) login user tsb ditolak `401 "Akun tidak aktif, hubungi admin"` → status dikembalikan ke Active.
  - curl: 5x login gagal berturut-turut → percobaan ke-6 dengan password **benar** tetap ditolak `401 "Akun terkunci sementara, coba lagi nanti"` (sesuai skenario ITL T-035).
  - Playwright: badge berubah jadi "Locked" + tombol "Unlock Now" muncul → klik → toast sukses, badge kembali "Unlocked" → (via curl) login dengan password benar berhasil (`200`, `status:"active"`, `locked_until:null`).

## v1.3.0 - 2026-07-09

- **Implementasi Phase 3 selesai (T-025, T-026, T-027, T-030, T-031, T-032):** `05_ITL.md` diperbarui — seluruh task ditandai `[x]` (T-028/T-029 tetap dihapus dari v1.1.0).
  - **Health Check (T-025/T-026/T-027):** `GatewayServiceHealthCheck` (manual, `POST /api/services/:id/health-check`) dan `GatewayServiceHealthCheckAll` (background, goroutine ticker in-process — bukan asynq scheduler, lihat catatan T-025) — keduanya GET `{base_url}/health` timeout 5s, update `health_status`/`health_checked_at`, tidak membuat audit log/notifikasi (murni observasi). FE: tombol Health Check di tabel Service (badge langsung ter-update).
  - **Audit Trail (T-030/T-031/T-032):** hook `recordAuditLog` dipasang di seluruh CUD Service & Route (create → snapshot penuh, update → `{before, after}`, delete → snapshot sebelum hapus). Permission baru `audit.index` (auto masuk Super Admin). Endpoint `GET /api/audit-logs` (+filter `entity_type`/`entity`/`actor`/`from`/`to`) dan `GET /api/audit-logs/:id`. FE: halaman `/audit-logs` dengan filter Entity Type/Actor/Date Range dan expand-row untuk detail JSON (pakai `UiTable` prop `expand` yang sudah ada).
- **Diuji end-to-end via Playwright** (accessibility snapshot, tanpa screenshot):
  - Health Check: klik tombol → status `unknown` → `down` (upstream tidak jalan, sesuai ekspektasi) → filter Health Status=Down menampilkan hanya service tsb.
  - Audit Log: create Service baru → entry `service/create` langsung muncul di `/audit-logs` dengan actor benar → expand row menampilkan JSON snapshot lengkap → filter Entity Type=Route menampilkan kosong (belum ada log route, sesuai ekspektasi).

## v1.2.0 - 2026-07-09

- **Implementasi Phase 2 selesai (T-021 s/d T-024):** `05_ITL.md` diperbarui — seluruh task Phase 2 ditandai `[x]`.
  - **Redis Pub/Sub multi-instance sync (T-021/T-022):** `RouteManager.RefreshAndPublish()` (dipakai on-save trigger, via adapter `routeCachePublisher` di `di/container.go`) vs `RouteManager.Refresh()` biasa (dipakai periodic ticker & Pub/Sub subscriber) — dipisah sengaja untuk mencegah infinite broadcast loop antar instance. `RedisCache` ditambah method `Publish`/`Subscribe`.
  - **Dynamic Rate Limit Resolution (T-023):** resolusi berjenjang Route → Service → Global sekarang dilakukan sekali saat `RouteManager.Refresh()` (bukan per-request) — hasilnya disimpan sebagai satu field `CachedRoute.RateLimit RateLimitConfig{Limit, WindowSecs}`, sesuai masukan user untuk menyederhanakan dari 2 field pointer terpisah. Enforcement pakai token bucket per key `route_id+client_IP` (`helpers.RateLimiter.AllowWithLimit`, refactor dari `Allow` yang lama) — limit 429 + header `Retry-After`/`X-RateLimit-*`.
  - **T-024** rate limit fields di form Service/Route sudah ada sejak Phase 1, diverifikasi ulang.
- **Diuji end-to-end:**
  - 2 instance backend (port 8080 & 8081) dijalankan bersamaan — route baru dibuat via instance A langsung muncul di `cache/status` instance B tanpa restart/request manual (Pub/Sub terbukti jalan).
  - Rate limit override 2 req/60s di satu Route: request 1-2 lolos, request ke-3 langsung `429` dengan `Retry-After`; instance B ikut menerapkan limit yang sama (config tersinkron), tapi counter tetap independen per instance (sesuai desain — hanya config yang di-broadcast, bukan counter).
  - Playwright: validasi form "Rate limit harus lebih dari 0" saat isi 0 pada Edit Route — sesuai skenario E2E T-024.

## v1.1.2 - 2026-07-09

- **Melengkapi FSD §3.2 & §3.4 (gap ditemukan user):** Service List/Table dan Route List/Table sebelumnya belum punya Search/Filter dan sejumlah kolom sesuai spesifikasi FSD.
  - `useCrud.ts`: `fetchAll()` diperluas menerima `extraParams` opsional (non-breaking) — dipersist lintas pagination.
  - Service List: tambah Search by Name (debounce 300ms), Filter Protocol/Is Active/Health Status, kolom Health Status (badge). Kolom `route_count` (di luar spek FSD) dihapus.
  - Route List: tambah Filter Service/Method/Is Active, kolom Rate Limit.
  - Diuji ulang dengan Playwright: search & filter Service, filter Service pada Route, semua mengirim query param yang benar dan hasil ter-filter sesuai.
- **Belum diimplementasikan (di luar scope Phase 1, sengaja ditunda):** tombol Health Check & View Routes di §3.2 (backend endpoint health-check adalah T-026/Phase 3), radio button literal untuk Permission Match Mode di §3.3 (saat ini pakai dropdown select, fungsinya sama), §3.6 Audit Log List (Phase 3) dan §3.7 User Status & Lock Panel (Phase 3B).

## v1.1.1 - 2026-07-09

- **Implementasi Phase 1 selesai (T-001 s/d T-020):** `05_ITL.md` diperbarui — seluruh task Phase 1 ditandai `[x]`.
  - Backend: migration 4 tabel gateway, model/repo/DTO/service/handler/route Service & Route Management, `RouteManager` (atomic-swap cache), Dynamic Proxy Engine (REST+WebSocket, permission any/all), cache refresh (on-save lokal + periodic + manual), 12 permission baru (`service.*`, `route.*`).
  - Frontend: store + halaman `gateway/services` dan `gateway/routes`, terintegrasi router & sidebar menu.
  - Diuji end-to-end dengan Playwright (accessibility snapshot, tanpa screenshot) + curl untuk skenario proxy (401/403/404/502).
- **Gap diketahui (dicatat di masing-masing task ITL, follow-up bukan blocker):** search-by-name Service (T-006), shortcut toggle Active di tabel (T-007), dialog cascade-delete khusus (T-008), filter Service/Method/Is Active di Route list (T-010, T-020), pengujian WebSocket dengan upstream sungguhan (T-015). Redis Pub/Sub publish saat on-save (bagian dari §2.18) sengaja ditunda ke T-021 (Phase 2) sesuai pembagian task asli.
- **Perbaikan saat implementasi:** validator `url` bawaan vuelidate menolak host tanpa TLD (mis. `http://localhost:9001`) — diganti validator custom berbasis `URL` constructor di `gatewayService.ts`.

## v1.1.0 - 2026-07-09

- **Ditambahkan — User Status & Lock (delta modul User existing):** kolom `status` (Active/Suspend, manual oleh admin) dipisah total dari `failed_login_attempts`/`locked_until` (Lock, otomatis oleh sistem setelah gagal login berulang). Reuse permission `user.edit` existing, tidak ada permission baru.
  - `01_PRD.md`: US-14, US-15; MoSCoW baru; NFR Security (Account Protection).
  - `02_FSD.md`: §2.25–§2.27, §3.7, §5.5, §6.5, Functional Hierarchy §6.
  - `03_Role_Matrix.md`: Catatan 4 (reuse `user.edit`, tidak ada permission baru), Data Ownership Rules.
  - `04_TDD.md`: delta ERD `users`, §4.7 API endpoints, env `AUTH_MAX_FAILED_LOGIN_ATTEMPTS`/`AUTH_LOCK_DURATION_MINUTES`.
  - `05_ITL.md`: Phase 3B — T-033 s/d T-036.
- **Dihapus — Route Testing Tool:** dibatalkan dari scope penuh (endpoint `POST /api/routes/:id/test`, permission `route.test`, UI panel). Nomor/ID yang sudah dipakai (US-10, FSD §2.12/§3.5/§6.4, TDD route.test row, ITL T-028/T-029) sengaja tidak dipakai ulang untuk menjaga jejak histori & validitas referensi silang lain.
- **Diubah — Manual Cache Refresh:** permission `gateway.cache-refresh` dihapus; endpoint `POST /api/gateway/cache/refresh` dan `GET /api/gateway/cache/status` sekarang cukup permission `route.create` **atau** `route.edit` (any-of, 1 endpoint 2 permission).
- **Diubah — Konvensi penamaan field request API:** suffix `_id`/`_ids` dihilangkan dari field **request** (bukan response): `service_id` → `service`, `permission_ids` → `permissions` (route endpoints); `entity_id` → `entity`, `actor_user_id` → `actor` (audit log query filter). Response DTO tetap boleh memakai suffix `_id`/`_ids` apa adanya.
- **Diubah — RouteManager refresh mechanism:** dipertegas sebagai **atomic swap** (build snapshot baru di struktur terpisah, baru ditukar setelah selesai) sehingga cache route yang sudah ada tidak pernah kosong/hilang sesaat selama proses refresh, baik sukses maupun gagal. Didokumentasikan di `01_PRD.md` (NFR Reliability), `02_FSD.md` (§2.18), dan `04_TDD.md` (§5 Konkurensi & Caching).

## v1.0.0 - 2026-07-08

- Initial generation of project plan for **API Gateway — Service & Route Management** (fitur delta dari base project boilerplate).
- **01_PRD.md:** Vision, target audience, 13 user stories, MoSCoW (Service CRUD, Route CRUD + Permission many-to-many, Dynamic Proxy Engine REST+WebSocket, refresh strategy, Redis Pub/Sub, dynamic rate limit, health check, route testing tool, audit trail sebagai Must/Should-have; gRPC & import-export sebagai Could-have/backlog).
- **02_FSD.md:** 24 functional requirements (§2.1–§2.24), 6 screen/form specs (§3.1–§3.6), use case diagram, 4 flow diagrams (§5.1–§5.4), 4 kelompok validasi (§6.1–§6.4).
- **03_Role_Matrix.md:** 3 role (Super Admin, Gateway Admin, Viewer), 12 permission baru (`service.*`, `route.*`, `gateway.cache-refresh`, `audit.index`), matrix lengkap, catatan data ownership & scoped access.
- **04_TDD.md:** Tech stack (reuse base project), system architecture diagram, ERD 4 tabel baru (`gateway_services`, `gateway_routes`, `gateway_route_permissions`, `gateway_audit_logs`), API contract §4.1–§4.6, infrastruktur & security (portable Docker/K8s/OCP).
- **05_ITL.md:** 32 task (T-001–T-032) terbagi 3 fase implementasi + catatan backlog (gRPC, import/export, circuit breaker, request transform header).
