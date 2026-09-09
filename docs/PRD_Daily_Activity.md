# PRD — Daily Activity Management (Personal ERP)

**Module:** Daily Activity
**Status:** Draft v1.0
**Stack:** Go (backend) + MySQL + Next.js (frontend)

## 1. Ringkasan

Modul Daily Activity memungkinkan pengguna untuk mencatat, mengorganisasi, dan meninjau tugas serta aktivitas harian mereka di dalam sistem Personal ERP single-user. Modul ini menyediakan pencatatan waktu yang ringan dilengkapi dengan kategori, tingkat prioritas, status, dan statistik harian, sehingga pengguna dapat merencanakan hari mereka dan meninjau bagaimana aktivitas tersebut berjalan.

## 2. Tujuan

- Mencatat aktivitas untuk hari apa pun (direncanakan, sedang berlangsung, selesai, dibatalkan).
- Mengorganisasi aktivitas dengan kategori, prioritas, dan catatan.
- Melihat ringkasan harian yang jelas: total aktivitas, jumlah yang selesai, dan waktu yang dihabiskan.
- Memfilter dan mencari aktivitas dalam rentang tanggal.
- Menjaga model data tetap sederhana dan dapat diperluas untuk modul ERP di masa mendatang.

## 3. Non-Goals (v1)

- Multi User/Auth.
- Tidak ada aktivitas berulang atau sinkronisasi kalender.
- Tidak ada lampiran.
- Tidak ada timer pelacakan waktu aktivitas (mulai/berhenti) — hanya entri waktu manual.

## 4. Model Data

### 4.1 `categories` (Kategori)

| Column     | Type         | Notes                            |
| ---------- | ------------ | -------------------------------- |
| id         | BIGINT PK AI |                                  |
| name       | VARCHAR(100) | required, unique                 |
| color      | VARCHAR(7)   | hex color for UI, e.g. `#2563eb` |
| icon       | VARCHAR(32)  | optional emoji/icon label        |
| created_at | TIMESTAMP    | default now                      |

Seed categories on first run: Work, Personal, Health, Learning, Errands.

### 4.2 `activities`

| Column        | Type                                    | Notes                                          |
| ------------- | --------------------------------------- | ---------------------------------------------- |
| id            | BIGINT PK AI                            |                                                |
| title         | VARCHAR(255)                            | required                                       |
| description   | TEXT NULL                               |                                                |
| category_id   | BIGINT NULL                             | FK → categories.id (ON DELETE SET NULL)        |
| activity_date | DATE                                    | the day this activity belongs to               |
| start_time    | TIME NULL                               | optional start                                 |
| end_time      | TIME NULL                               | optional end                                   |
| duration_min  | INT NULL                                | minutes; auto-calc if times given, else manual |
| status        | ENUM planned/in_progress/done/cancelled | default planned                                |
| priority      | ENUM low/medium/high                    | default medium                                 |
| notes         | TEXT NULL                               |                                                |
| created_at    | TIMESTAMP                               | default now                                    |
| updated_at    | TIMESTAMP                               | auto-update on change                          |

## 5. API Design (REST, base `/api`)

| Method | Endpoint                | Description                                                                     |
| ------ | ----------------------- | ------------------------------------------------------------------------------- |
| GET    | `/api/health`           | liveness + db check                                                             |
| GET    | `/api/categories`       | list categories                                                                 |
| POST   | `/api/categories`       | create category                                                                 |
| PUT    | `/api/categories/{id}`  | update category                                                                 |
| DELETE | `/api/categories/{id}`  | delete category                                                                 |
| GET    | `/api/activities`       | list; filters: `date`, `start`, `end`, `category_id`, `status`, `priority`, `q` |
| GET    | `/api/activities/{id}`  | get one                                                                         |
| POST   | `/api/activities`       | create                                                                          |
| PUT    | `/api/activities/{id}`  | update                                                                          |
| DELETE | `/api/activities/{id}`  | delete                                                                          |
| GET    | `/api/activities/stats` | `date`, `start`, `end` → totals                                                 |

**Response envelope:** `{ "data": ... }` for success, `{ "error": "message" }` for errors.

**Stats response:**

```json
{
  "total": 12,
  "done": 7,
  "cancelled": 1,
  "in_progress": 2,
  "planned": 2,
  "total_minutes": 480,
  "by_category": [{ "category": "Work", "count": 5, "minutes": 300 }]
}
```

## 6. Validation Rules

- `title` required, max 255 chars.
- `activity_date` required, valid ISO date.
- `status` ∈ {planned, in_progress, done, cancelled}; `priority` ∈ {low, medium, high}.
- If both `start_time` and `end_time` present, `end_time` must be after `start_time`; `duration_min` auto-computed.
- Unknown `category_id` → 400 (or null allowed).

## 7. Frontend Pages

| Route                   | Purpose                                                               |
| ----------------------- | --------------------------------------------------------------------- |
| `/`                     | Dashboard: today's summary + top activities                           |
| `/activities`           | Full activity log with filters (date range, status, category, search) |
| `/activities/new`       | Create activity (also inline modal)                                   |
| `/activities/[id]/edit` | Edit activity                                                         |

- Category management UI on the activities page (add/edit/delete).

## 8. Acceptance Criteria

1. User can create, view, edit, and delete activities.
2. Activities can be filtered by date range, status, category, and free-text search.
3. Dashboard shows today's stats (counts by status, total minutes, by-category breakdown).
4. Duration auto-computes from start/end times and is editable.
5. Categories are CRUD-able and seeded on first run.
6. All endpoints return consistent JSON; errors are meaningful.
7. Backend builds with `go build`, frontend passes `eslint` and `next build`.

## 9. Future Ideas

- Recurring activities, daily goals/templates.
- Weekly/monthly reports and charts.
- Auth + multi-user.
- Calendar view / integration.

## 10. fitur selanjutnya

- buat fitur dashboard besok malam.
