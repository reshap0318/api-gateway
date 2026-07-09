# 03 - Role & Permissions Matrix

## API Gateway — Service & Route Management

> Dokumen ini hanya mencakup **modul baru** (Service, Route, Gateway runtime ops, Audit). Role/permission untuk modul existing (User, Role, Permission, Notification) mengikuti Role Matrix yang sudah berjalan di base project dan tidak diulang di sini.

---

## 1. User Roles Definition

| Role Name | Description | Access Level |
|---|---|---|
| **Super Admin** | Pemilik penuh sistem, termasuk modul existing (User/Role/Permission) dan modul Gateway baru ini. | Full access — semua permission |
| **Gateway Admin** | Bertanggung jawab operasional Gateway sehari-hari: mendaftarkan service, mengatur routing & permission per endpoint, menguji route, memantau kesehatan upstream. | Full CRUD pada modul Service & Route, plus operational actions (health check, cache refresh, test route). Tidak memiliki akses ke modul User/Role/Permission Management (itu domain Super Admin). |
| **Viewer** | Auditor/support/observer yang perlu visibilitas tanpa kewenangan mengubah konfigurasi. | Read-only pada Service, Route, dan Audit Log. |

---

## 2. Permissions List

| Permission | Module | Action | Description |
|---|---|---|---|
| `service.index` | service | index | Melihat daftar & detail Service |
| `service.create` | service | create | Mendaftarkan Service baru |
| `service.edit` | service | edit | Mengubah Service (termasuk toggle active) |
| `service.delete` | service | delete | Menghapus (soft-delete) Service |
| `service.health-check` | service | health-check | Memicu manual health check pada Service |
| `route.index` | route | index | Melihat daftar & detail Route |
| `route.create` | route | create | Membuat Route baru |
| `route.edit` | route | edit | Mengubah Route (termasuk assign permission) |
| `route.delete` | route | delete | Menghapus (soft-delete) Route |
| `audit.index` | audit | index | Melihat Audit Log |

> **Catatan 1:** Traffic yang diproxy oleh Dynamic Proxy Engine (§FSD 2.13–2.17) **tidak** dikontrol oleh permission di atas — otorisasinya dinamis, bergantung pada permission yang di-assign ke masing-masing Route (lihat `route_has_permissions` di TDD). Permission-permission tersebut (misal `user.show`, `finance.approve`, dsb.) dibuat & dikelola lewat modul Permission Management yang sudah ada, bukan bagian dari daftar di atas.
>
> **Catatan 2:** Tidak ada permission `gateway.cache-refresh` terpisah. Manual cache refresh (`POST /api/gateway/cache/refresh`) cukup mensyaratkan **salah satu** dari `route.create` **atau** `route.edit` (1 endpoint, 2 kemungkinan permission) — siapapun yang boleh membuat/mengubah Route otomatis boleh memicu refresh.
>
> **Catatan 3:** Permission `route.test` dihapus bersama fitur Route Testing Tool yang dibatalkan dari scope (v1.1.0).
>
> **Catatan 4:** Fitur User Status (Active/Suspend) dan Account Lock/Unlock (§FSD 2.25–2.27) **tidak menambah permission baru** — keduanya direuse lewat permission `user.edit` yang sudah ada di modul User existing, karena secara konsep adalah bagian dari "mengubah data user".

---

## 3. The Matrix Table

| Permission | Super Admin | Gateway Admin | Viewer |
|---|---|---|---|
| `service.index` | ✅ | ✅ | ✅ |
| `service.create` | ✅ | ✅ | ❌ |
| `service.edit` | ✅ | ✅ | ❌ |
| `service.delete` | ✅ | ✅ | ❌ |
| `service.health-check` | ✅ | ✅ | ❌ |
| `route.index` | ✅ | ✅ | ✅ |
| `route.create` | ✅ | ✅ | ❌ |
| `route.edit` | ✅ | ✅ | ❌ |
| `route.delete` | ✅ | ✅ | ❌ |
| `audit.index` | ✅ | ✅ | ✅ |

---

## 4. Data Ownership Rules

- **Service** dan **Route** tidak memiliki konsep "kepemilikan individual" (bukan data per-user) — semua record adalah shared organizational config. Siapapun dengan permission yang sesuai (`service.edit`/`route.edit`) dapat mengubah record milik siapapun; tidak ada scoping "hanya boleh edit yang dibuat sendiri".
- **Audit Log** bersifat immutable — tidak ada permission `edit`/`delete` untuk entity ini karena memang tidak boleh diubah oleh siapapun (termasuk Super Admin) lewat aplikasi.
- **Route Permission Assignment** (`route_has_permissions`) tunduk pada permission `route.edit` — siapa yang boleh mengubah Route juga otomatis boleh mengubah permission apa saja yang di-assign ke Route tersebut. Tidak ada permission terpisah untuk "assign permission" vs "edit route lainnya", karena keduanya adalah satu operasi form yang sama (§FSD 2.11).
- **User Status & Lock** — `status` (Active/Suspend) diubah manual lewat permission `user.edit`. `locked_until`/`failed_login_attempts` (Lock) di-set otomatis oleh sistem saat login gagal berulang (bukan lewat permission — ini alur autentikasi, bukan operasi CRUD admin), tapi Manual Unlock tetap tunduk pada `user.edit` yang sama.

## 5. Scoped Access Rules

- Tidak ada scoped access (satker/unit-based) untuk modul ini di v1 — seluruh Service dan Route bersifat global/organization-wide, karena Gateway adalah infrastruktur bersama, bukan data per-unit. Jika kebutuhan multi-tenant/scoped muncul di masa depan, dapat ditambahkan kolom `owner_unit_id` dan permission varian `route.index-satker` mengikuti pola scoped permission standar (`{module}.{action}-satker`) — dicatat sebagai potential future extension, bukan bagian dari v1.
