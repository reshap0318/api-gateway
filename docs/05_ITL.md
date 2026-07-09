# 05 - Implementation Task List (ITL)

## API Gateway — Service & Route Management

---

## Implementation Phases

| Phase | Fokus | Tasks |
|---|---|---|
| **Phase 1** | Foundation — DB schema, Service & Route CRUD, Dynamic Proxy Engine (REST + WebSocket), refresh dasar | T-001 s/d T-020 |
| **Phase 2** | Consistency & Scale — Redis Pub/Sub multi-instance, Rate limiting dinamis | T-021 s/d T-024 |
| **Phase 3** | Extended Features — Health check, Audit trail | T-025 s/d T-032 (T-028, T-029 dihapus — lihat catatan) |
| **Phase 3B** | User Status & Lock (Delta Modul User Existing) | T-033 s/d T-036 |
| **Phase 4** | Backlog (belum dijadwalkan versi ini) | gRPC proxy support, Import/Export config — lihat catatan di akhir dokumen |

---

## Phase 1 — Foundation

### T-001: Database Migration — Gateway Tables

- **Feature/Module:** Foundation
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §3 (skema data implisit dari seluruh Detailed Functional Requirements §2.1–§2.24)
- **TDD Ref:**
  - ERD §3 — `gateway_services`, `gateway_routes`, `gateway_route_permissions`, `gateway_audit_logs`

### T-002: Service Model, Repository & DTO

- **Feature/Module:** Service Management
- **Priority:** P0
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.1 Create Service, §2.2 Edit Service, §2.3 Delete Service, §2.4 List/View Service
- **TDD Ref:**
  - ERD §3 — `gateway_services`

### T-003: Route Model, Repository, DTO & Route-Permission Join

- **Feature/Module:** Route Management
- **Priority:** P0
- **Estimated Effort:** 4h
- **Status:** [x]
- **FSD Ref:**
  - §2.7 Create Route, §2.8 Edit Route, §2.9 Delete Route, §2.10 List/View Route, §2.11 Assign Permissions to Route
- **TDD Ref:**
  - ERD §3 — `gateway_routes`, `gateway_route_permissions`

### T-004: RouteManager — In-Memory Cache & Startup Load

- **Feature/Module:** Dynamic Proxy Engine
- **Priority:** P0
- **Estimated Effort:** 5h
- **Status:** [x]
- **FSD Ref:**
  - §2.13 Request Matching, §2.18 Cache Refresh On-Save Trigger (bagian lokal), §2.19 Cache Refresh Periodic Fallback
- **TDD Ref:**
  - §5 Infrastructure & Security — Konkurensi & Caching (RouteManager)

### T-005: Create Service API

- **Feature/Module:** Service Management
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.1 Create Service
  - §3.1 Create/Edit Service Form
  - §6.1 Create/Edit Service Validation
- **TDD Ref:**
  - POST /api/services
- **E2E Scenario (Given-When-Then):**
  - **Given** Gateway Admin berada di halaman Service Management dan klik "Tambah Service"
  - **When** mengisi Name, Base URL, Protocol, lalu submit
  - **Then** service baru muncul di tabel dengan status Active, dan toast sukses ditampilkan
  - **Given** Admin mengisi Name yang sudah dipakai service lain
  - **When** submit
  - **Then** error inline "Nama service sudah digunakan" muncul pada field Name

### T-006: List/View Service API

- **Feature/Module:** Service Management
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.4 List/View Service
  - §3.2 Service List/Table
- **TDD Ref:**
  - GET /api/services, GET /api/services/:id
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin berada di halaman Service Management
  - **When** halaman selesai loading
  - **Then** tabel Service tampil dengan kolom Name, Base URL, Protocol, Rate Limit, Health Status, Is Active, Actions
  - **Given** Admin mengetik keyword di Search by Name
  - **When** hasil pencarian dimuat
  - **Then** tabel hanya menampilkan Service yang cocok dengan keyword

### T-007: Edit Service API (termasuk Toggle Active)

- **Feature/Module:** Service Management
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [x] (edit via form selesai; shortcut toggle switch langsung di baris tabel belum ada, saat ini toggle Is Active hanya lewat form Edit — follow-up)
- **FSD Ref:**
  - §2.2 Edit Service, §2.5 Toggle Active/Inactive
  - §3.1 Create/Edit Service Form
  - §6.1 Create/Edit Service Validation
- **TDD Ref:**
  - PUT /api/services/:id
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin membuka form Edit pada salah satu Service
  - **When** mengubah Base URL dan submit
  - **Then** tabel menampilkan Base URL baru dan toast sukses muncul
  - **Given** Admin menekan toggle Is Active pada baris tabel
  - **When** toggle di-klik
  - **Then** status berubah langsung tanpa membuka form penuh

### T-008: Delete Service API

- **Feature/Module:** Service Management
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [x] (backend cascade-confirmation logic — 409 + jumlah route aktif — selesai; dialog FE khusus untuk alur cascade belum ada, baru delete biasa — follow-up)
- **FSD Ref:**
  - §2.3 Delete Service
- **TDD Ref:**
  - DELETE /api/services/:id
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin klik Delete pada Service yang tidak punya Route aktif
  - **When** konfirmasi delete
  - **Then** Service hilang dari tabel dan toast sukses muncul
  - **Given** Admin klik Delete pada Service yang masih punya Route aktif
  - **When** dialog konfirmasi muncul
  - **Then** ditampilkan peringatan jumlah Route aktif terkait sebelum admin melanjutkan

### T-009: Create Route API

- **Feature/Module:** Route Management
- **Priority:** P1
- **Estimated Effort:** 4h
- **Status:** [x]
- **FSD Ref:**
  - §2.7 Create Route, §2.11 Assign Permissions to Route
  - §3.3 Create/Edit Route Form
  - §6.2 Create/Edit Route Validation
- **TDD Ref:**
  - POST /api/routes
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin berada di halaman Route Management dan klik "Tambah Route"
  - **When** memilih Service, mengisi Method dan Path Pattern, memilih Permission + Match Mode, lalu submit
  - **Then** Route baru muncul di tabel dengan permission dan match mode yang dipilih
  - **Given** Admin mengisi kombinasi Service+Method+Path Pattern yang sudah ada
  - **When** submit
  - **Then** error inline "Route dengan method dan path ini sudah terdaftar untuk service ini" muncul

### T-010: List/View Route API

- **Feature/Module:** Route Management
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.10 List/View Route
  - §3.4 Route List/Table
- **TDD Ref:**
  - GET /api/routes, GET /api/routes/:id
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin berada di halaman Route Management dan memfilter berdasarkan Service tertentu
  - **When** filter diterapkan
  - **Then** tabel hanya menampilkan Route milik Service tersebut, lengkap dengan chip Permission dan Match Mode

### T-011: Edit Route API

- **Feature/Module:** Route Management
- **Priority:** P1
- **Estimated Effort:** 4h
- **Status:** [x]
- **FSD Ref:**
  - §2.8 Edit Route
  - §3.3 Create/Edit Route Form
  - §6.2 Create/Edit Route Validation
- **TDD Ref:**
  - PUT /api/routes/:id
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin membuka form Edit pada suatu Route
  - **When** mengubah daftar Permission yang di-assign dan submit
  - **Then** tabel menampilkan chip Permission baru sesuai perubahan

### T-012: Delete Route API

- **Feature/Module:** Route Management
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.9 Delete Route
- **TDD Ref:**
  - DELETE /api/routes/:id
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin klik Delete pada suatu Route
  - **When** konfirmasi delete
  - **Then** Route hilang dari tabel dan traffic ke path tersebut langsung mendapat 404 setelah refresh cache

### T-013: Dynamic Proxy Engine — Request Matching & REST Proxying

- **Feature/Module:** Dynamic Proxy Engine
- **Priority:** P0
- **Estimated Effort:** 6h
- **Status:** [x]
- **FSD Ref:**
  - §2.13 Request Matching, §2.16 REST Proxying
  - §5.2 Dynamic Proxy Request Flow
  - §6.3 Dynamic Proxy Runtime Errors
- **TDD Ref:**
  - §4.6 Dynamic Proxied Requests (Catch-all)

### T-014: Permission Resolution (any/all) di Proxy

- **Feature/Module:** Dynamic Proxy Engine
- **Priority:** P0
- **Estimated Effort:** 4h
- **Status:** [x]
- **FSD Ref:**
  - §2.14 Permission Resolution
  - §6.3 Dynamic Proxy Runtime Errors
- **TDD Ref:**
  - §4.6 Dynamic Proxied Requests (Catch-all)

### T-015: WebSocket Proxying Support

- **Feature/Module:** Dynamic Proxy Engine
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [x] (satu code path dengan REST via `httputil.ReverseProxy`; belum diuji dengan upstream WebSocket sungguhan — hanya diverifikasi secara struktural)
- **FSD Ref:**
  - §2.17 WebSocket Proxying
- **TDD Ref:**
  - §4.6 Dynamic Proxied Requests (Catch-all)

### T-016: Cache Refresh — On-Save Trigger

- **Feature/Module:** Dynamic Proxy Engine
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x] (refresh lokal on-save selesai; publish Redis Pub/Sub adalah T-021, Phase 2)
- **FSD Ref:**
  - §2.18 Cache Refresh On-Save Trigger
  - §5.1 Create Route Flow
- **TDD Ref:**
  - POST /api/services, PUT /api/services/:id, DELETE /api/services/:id, POST /api/routes, PUT /api/routes/:id, DELETE /api/routes/:id

### T-017: Cache Refresh — Periodic Fallback

- **Feature/Module:** Dynamic Proxy Engine
- **Priority:** P0
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.19 Cache Refresh Periodic Fallback
- **TDD Ref:**
  - §5 Infrastructure & Security — Observability (`GATEWAY_CACHE_REFRESH_INTERVAL_SECONDS`)

### T-018: Manual Cache Refresh Endpoint

- **Feature/Module:** Dynamic Proxy Engine
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.20 Cache Refresh Manual Trigger
- **TDD Ref:**
  - POST /api/gateway/cache/refresh, GET /api/gateway/cache/status
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin berada di halaman Gateway Status dan klik "Refresh Cache Now"
  - **When** proses refresh selesai
  - **Then** timestamp "Last refreshed at" ter-update dan toast sukses muncul

### T-019: Frontend — Service Management UI

- **Feature/Module:** Service Management
- **Priority:** P1
- **Estimated Effort:** 6h
- **Status:** [x] (rute aktual `/gateway/services`, bukan `/services` seperti draf awal skenario)
- **FSD Ref:**
  - §3.1 Create/Edit Service Form, §3.2 Service List/Table
- **TDD Ref:**
  - §4.2 Service Management Endpoints
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin login sebagai Gateway Admin dan navigasi ke `/services`
  - **When** halaman dimuat
  - **Then** tabel Service tampil sesuai §3.2, tombol "Tambah Service" terlihat dan dapat membuka form §3.1

### T-020: Frontend — Route Management UI

- **Feature/Module:** Route Management
- **Priority:** P1
- **Estimated Effort:** 7h
- **Status:** [x] (rute aktual `/gateway/routes` — daftar flat semua Route dengan filter Service/Method/Is Active, bukan nested `/services/:id/routes` seperti draf awal skenario)
- **FSD Ref:**
  - §3.3 Create/Edit Route Form, §3.4 Route List/Table
- **TDD Ref:**
  - §4.3 Route Management Endpoints
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin navigasi ke `/services/:id/routes`
  - **When** halaman dimuat
  - **Then** tabel Route tampil sesuai §3.4, filter Service/Method/Is Active berfungsi

---

## Phase 2 — Consistency & Scale

### T-021: Redis Pub/Sub Publisher on Save

- **Feature/Module:** Dynamic Proxy Engine
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.21 Multi-instance Sync (Redis Pub/Sub)
  - §5.3 Multi-instance Cache Sync Flow
- **TDD Ref:**
  - §5 Infrastructure & Security — Observability (`GATEWAY_REFRESH_CHANNEL`)

### T-022: Redis Pub/Sub Subscriber per Instance

- **Feature/Module:** Dynamic Proxy Engine
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.21 Multi-instance Sync (Redis Pub/Sub)
  - §5.3 Multi-instance Cache Sync Flow
- **TDD Ref:**
  - §2 System Architecture (Redis Subscriber goroutine)

### T-023: Dynamic Rate Limit Resolution (Route → Service → Global)

- **Feature/Module:** Dynamic Proxy Engine
- **Priority:** P1
- **Estimated Effort:** 4h
- **Status:** [x] (resolusi chain dilakukan sekali saat RouteManager.Refresh(), bukan per-request — hasil resolved disimpan sebagai CachedRoute.RateLimit)
- **FSD Ref:**
  - §2.15 Rate Limit Enforcement (Dynamic)
  - §6.3 Dynamic Proxy Runtime Errors
- **TDD Ref:**
  - §3 Catatan Skema (resolusi berjenjang Route → Service → Global)

### T-024: Frontend — Rate Limit Fields di Service/Route Form

- **Feature/Module:** Service Management, Route Management
- **Priority:** P2
- **Estimated Effort:** 2h
- **Status:** [x] (sudah dibangun di Phase 1 T-019/T-020; diverifikasi ulang di Phase 2)
- **FSD Ref:**
  - §3.1 Create/Edit Service Form (Rate Limit per Minute), §3.3 Create/Edit Route Form (Rate Limit per Minute override)
- **TDD Ref:**
  - PUT /api/services/:id, PUT /api/routes/:id
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin mengisi Rate Limit per Minute pada form Edit Route dengan angka 0
  - **When** submit
  - **Then** error inline "Rate limit harus lebih dari 0" muncul

---

## Phase 3 — Extended Features

### T-025: Health Check Background Job

- **Feature/Module:** Health Check Upstream
- **Priority:** P2
- **Estimated Effort:** 4h
- **Status:** [x] (implementasi pakai goroutine ticker in-process — pola sama dengan RouteManager periodic refresh — bukan asynq scheduler, supaya tidak perlu deploy `cmd/worker` terpisah hanya untuk fitur ini; interval detik presisi via `GATEWAY_HEALTHCHECK_INTERVAL_SECONDS` tetap terpenuhi)
- **FSD Ref:**
  - §2.22 Health Check Upstream (Background Job)
  - §5.4 Health Check Flow
- **TDD Ref:**
  - §5 Infrastructure & Security — Observability (`GATEWAY_HEALTHCHECK_INTERVAL_SECONDS`)

### T-026: Manual Health Check Endpoint

- **Feature/Module:** Health Check Upstream
- **Priority:** P2
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.6 Manual Health Check Trigger
- **TDD Ref:**
  - POST /api/services/:id/health-check
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin klik "Health Check" pada salah satu baris Service
  - **When** proses selesai
  - **Then** badge Health Status pada baris tersebut ter-update (Up/Down) beserta waktu pengecekan terakhir

### T-027: Frontend — Health Status Badge & Manual Trigger

- **Feature/Module:** Health Check Upstream
- **Priority:** P2
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §3.2 Service List/Table (Filter Health Status, kolom Health Status)
- **TDD Ref:**
  - GET /api/services, POST /api/services/:id/health-check
- **E2E Scenario (Given-When-Then):**
  - **Given** Admin memfilter tabel Service dengan Health Status = Down
  - **When** filter diterapkan
  - **Then** hanya Service berstatus Down yang tampil

> **T-028 dan T-029 dihapus.** Route Testing Tool dibatalkan dari scope (v1.1.0) — lihat `CHANGELOG.md`. ID task sengaja tidak dipakai ulang.

### T-030: Audit Log — Record Change (Hook di Service/Route CUD)

- **Feature/Module:** Audit Trail
- **Priority:** P2
- **Estimated Effort:** 4h
- **Status:** [x]
- **FSD Ref:**
  - §2.23 Record Change (Audit Trail)
- **TDD Ref:**
  - ERD §3 — `gateway_audit_logs`; hook pada seluruh endpoint §4.2 dan §4.3 yang mengubah data

### T-031: Audit Log — View API

- **Feature/Module:** Audit Trail
- **Priority:** P2
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.24 View Audit Log
  - §3.6 Audit Log List
- **TDD Ref:**
  - GET /api/audit-logs, GET /api/audit-logs/:id

### T-032: Frontend — Audit Log List UI

- **Feature/Module:** Audit Trail
- **Priority:** P2
- **Estimated Effort:** 3h
- **Status:** [x] (kolom "Entity Name/ID" disederhanakan jadi "Entity ID" saja — tidak join ke nama Service/Route asli untuk hindari query tambahan; Date Range pakai 2 native date input, bukan komponen picker khusus karena belum ada di `components/utils`)
- **FSD Ref:**
  - §3.6 Audit Log List
- **TDD Ref:**
  - GET /api/audit-logs
- **E2E Scenario (Given-When-Then):**
  - **Given** Auditor login dan navigasi ke `/audit-logs`
  - **When** memfilter berdasarkan Entity Type = Route dan rentang tanggal tertentu
  - **Then** tabel menampilkan hanya entry Route dalam rentang tanggal tsb, dan expand detail menampilkan JSON before/after

---

## Phase 3B — User Status & Lock (Delta Modul User Existing)

### T-033: Migration & Model Delta — `users` (status, failed_login_attempts, locked_until)

- **Feature/Module:** User Status & Lock
- **Priority:** P1
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.25 Update User Status, §2.26 Automatic Account Lock, §2.27 Manual Unlock
- **TDD Ref:**
  - ERD §3 — delta tabel `users` (`status`, `failed_login_attempts`, `locked_until`)

### T-034: Update User Status API + Frontend Toggle

- **Feature/Module:** User Status & Lock
- **Priority:** P1
- **Estimated Effort:** 3h
- **Status:** [x]
- **FSD Ref:**
  - §2.25 Update User Status
  - §3.7 User Status & Lock Panel
  - §6.5 User Status & Lock Validation
- **TDD Ref:**
  - PUT /api/users/:id/status
- **E2E Scenario (Given-When-Then):**
  - **Given** Super Admin membuka detail seorang User dan mengubah Status ke "Suspend"
  - **When** submit
  - **Then** badge Status berubah menjadi "Suspend" dan user tersebut tidak bisa login lagi
  - **Given** User dengan status Suspend mencoba login
  - **When** submit credential (walau benar)
  - **Then** error "Akun tidak aktif, hubungi admin" ditampilkan

### T-035: Auto-Lock Logic pada Login Flow

- **Feature/Module:** User Status & Lock
- **Priority:** P2
- **Estimated Effort:** 4h
- **Status:** [x]
- **FSD Ref:**
  - §2.26 Automatic Account Lock (Failed Login Threshold)
  - §5.5 Login with Status & Lock Check Flow
  - §6.5 User Status & Lock Validation
- **TDD Ref:**
  - §5 Infrastructure & Security — Observability (`AUTH_MAX_FAILED_LOGIN_ATTEMPTS`, `AUTH_LOCK_DURATION_MINUTES`); §4.7 (catatan business logic internal login)
- **E2E Scenario (Given-When-Then):**
  - **Given** User memasukkan password salah 5 kali berturut-turut (default threshold)
  - **When** percobaan ke-5 gagal
  - **Then** akun terkunci dan percobaan login ke-6 (walau password benar) menampilkan "Akun terkunci sementara, coba lagi nanti"
  - **Given** User login dengan password benar sebelum mencapai threshold
  - **When** login sukses
  - **Then** `failed_login_attempts` ter-reset ke 0

### T-036: Manual Unlock API + Frontend Button

- **Feature/Module:** User Status & Lock
- **Priority:** P2
- **Estimated Effort:** 2h
- **Status:** [x]
- **FSD Ref:**
  - §2.27 Manual Unlock
  - §3.7 User Status & Lock Panel
  - §6.5 User Status & Lock Validation
- **TDD Ref:**
  - POST /api/users/:id/unlock
- **E2E Scenario (Given-When-Then):**
  - **Given** Super Admin membuka detail User yang sedang berstatus Locked
  - **When** klik "Unlock Now"
  - **Then** badge Lock hilang dan user tersebut bisa login kembali dengan password yang benar

---

## Backlog (Belum Dijadwalkan — Phase 4+)

Item berikut disepakati sebagai **Could-have/Won't (v1)** di PRD §4 — tidak dibuatkan task block penuh karena belum masuk scope implementasi, dicatat sebagai arah pengembangan berikutnya:

- **gRPC Proxy Support** — butuh transport HTTP/2 terpisah dari `httputil.ReverseProxy` yang dipakai REST/WebSocket; akan jadi modul tambahan pada `protocol` enum di `gateway_services` (`grpc`) beserta handler proxy khusus.
- **Import/Export Config (JSON/YAML)** — endpoint tambahan untuk backup/restore `gateway_services` + `gateway_routes` + relasinya, mempermudah migrasi antar environment.
- **Circuit Breaker per Service** — mekanisme trip otomatis saat upstream error-rate tinggi.
- **Request/Response Transform Header** (`X-Forwarded-User`, `X-Request-ID`) — auto-injection header tambahan ke upstream untuk tracing lintas service.
