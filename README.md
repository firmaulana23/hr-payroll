# HR & Payroll System

Sistem mini untuk mengelola data karyawan, absensi, dan penggajian.

## 1. Ringkasan Teknis

### Backend
*   **Framework**: `Gin` (https://github.com/gin-gonic/gin)
*   **ORM**: `GORM` (https://gorm.io/)
*   **Database**: `PostgreSQL`
*   **Lain-lain**:
    *   `godotenv`: Untuk memuat konfigurasi dari file `.env`.
    *   `swaggo`: Untuk auto-generasi dokumentasi API Swagger dari komentar kode.

### Frontend
*   **Framework**: JavaScript (ES6)
*   **Styling**: CSS Murni (tanpa library)

## 2. Desain Database

Database menggunakan PostgreSQL. Skema tabel dibuat secara otomatis oleh GORM (`AutoMigrate`) berdasarkan struct domain di Go.

### Tabel: `employees`
Menyimpan data master karyawan.

| Nama Kolom   | Tipe Data        | Keterangan                  |
|--------------|------------------|-----------------------------|
| `id`         | `bigint`         | **Primary Key** (auto-increment) |
| `name`       | `text`           | Nama lengkap karyawan       |
| `base_salary`| `float8`         | Gaji pokok                  |
| `allowance`  | `float8`         | Tunjangan tetap             |
| `position`   | `text`           | Jabatan karyawan            |
| `created_at` | `timestamptz`    | Waktu pembuatan record      |
| `updated_at` | `timestamptz`    | Waktu pembaruan record      |

### Tabel: `attendances`
Mencatat kehadiran harian setiap karyawan.

| Nama Kolom   | Tipe Data        | Keterangan                  |
|--------------|------------------|-----------------------------|
| `id`         | `bigint`         | **Primary Key** (auto-increment) |
| `employee_id`| `bigint`         | **Foreign Key** ke `employees.id` |
| `date`       | `timestamptz`    | Tanggal absensi             |
| `status`     | `text`           | `PRESENT`, `ABSENT`, `LEAVE` |
| `check_in`   | `timestamptz`    | Waktu masuk (jika `PRESENT`) |
| `check_out`  | `timestamptz`    | Waktu pulang (jika `PRESENT`)|
| `created_at` | `timestamptz`    | Waktu pembuatan record      |

*Constraint Unik*: `(employee_id, date)` untuk memastikan satu karyawan hanya punya satu record absensi per hari.

### Tabel: `payrolls`
Menyimpan data slip gaji yang digenerate setiap bulan untuk setiap karyawan.

| Nama Kolom         | Tipe Data        | Keterangan                        |
|--------------------|------------------|-----------------------------------|
| `id`               | `bigint`         | **Primary Key** (auto-increment)     |
| `employee_id`      | `bigint`         | **Foreign Key** ke `employees.id`   |
| `period`           | `timestamptz`      | Periode gaji (misal: 2025-11-01)  |
| `base_salary`      | `float8`         | Gaji pokok saat digenerate        |
| `allowance`        | `float8`         | Tunjangan saat digenerate         |
| `total_absent`     | `bigint`         | Jumlah absen di periode tersebut  |
| `absence_deduction`| `float8`         | Total potongan karena absen       |
| `take_home_pay`    | `float8`         | Gaji bersih yang diterima         |
| `generated_at`     | `timestamptz`    | Waktu slip gaji dibuat            |

*Constraint Unik*: `(employee_id, period)` untuk memastikan satu karyawan hanya punya satu slip gaji per periode.

## 3. Flow Bisnis

1.  **Manajemen Karyawan**:
    *   Admin dapat **menambahkan** data karyawan baru (nama, posisi, gaji pokok, tunjangan).
    *   Admin dapat **melihat** daftar semua karyawan.
    *   Admin dapat **mengubah** data karyawan yang sudah ada.

2.  **Pencatatan Absensi**:
    *   Setiap hari, admin dapat mencatat status kehadiran karyawan:
        *   **Check-in**: Menandai karyawan hadir dan mencatat waktu masuk.
        *   **Check-out**: Memperbarui record kehadiran hari itu dengan waktu pulang.
        *   **Mark Absent**: Menandai karyawan tidak hadir.
        *   **Mark on Leave**: Menandai karyawan cuti.
    *   Admin dapat melihat riwayat absensi seorang karyawan dalam rentang tanggal tertentu.

3.  **Penggajian**:
    *   Pada akhir bulan, admin dapat men-**generate slip gaji** untuk seorang karyawan pada periode tertentu.
    *   Sistem akan menghitung gaji dengan rumus:
        *   Mencari jumlah hari absen (`ABSENT`) dari tabel `attendances` selama periode berjalan.
        *   Menghitung potongan absen: `Potongan = (Gaji Pokok / 22) * Jumlah Absen`. (Asumsi 22 hari kerja sebulan).
        *   Menghitung gaji bersih: `Gaji Bersih = Gaji Pokok + Tunjangan - Potongan`.
    *   Hasil perhitungan disimpan di tabel `payrolls`.
    *   Admin dapat melihat daftar semua slip gaji yang pernah dibuat.

## 4. Struktur Aplikasi (Backend)

Aplikasi backend menggunakan arsitektur **Hexagonal (Ports & Adapters)** untuk memisahkan logika bisnis dari detail teknis (seperti database atau framework HTTP).

*   `cmd/main.go`: Entry point aplikasi. Bertugas melakukan dependency injection (DI) dan menjalankan server.
*   `config/`: Memuat konfigurasi dari environment variables atau file `.env`.
*   `database/`: Inisialisasi koneksi ke database.
*   `internal/domain/`: **Inti aplikasi (Core Domain)**.
    *   Berisi definisi *entity* (struct Go seperti `Employee`, `Payroll`) dan *interface* untuk *repository* dan *service*. Bagian ini tidak bergantung pada library eksternal manapun.
*   `internal/repository/`: **Adapter Database (Secondary Adapter)**.
    *   Implementasi dari *repository interface* yang didefinisikan di domain. Bertanggung jawab untuk berkomunikasi dengan database (GORM).
*   `internal/service/`: **Use Case / Logika Bisnis**.
    *   Implementasi dari *service interface*. Di sinilah semua logika bisnis utama berada (misalnya, cara menghitung gaji).
*   `internal/delivery/`: **Adapter Input (Primary Adapter)**.
    *   `handler/`: Menerima request HTTP, memvalidasi input, memanggil service yang sesuai, dan mengembalikan response JSON.
    *   `http/`: Mengatur routing URL (misal: `/employees` ke `EmployeeHandler`).
*   `docs/`: File dokumentasi Swagger yang digenerate otomatis.

## 5. Cara Menjalankan

### Prasyarat
*   Go (versi 1.20+)
*   PostgreSQL (berjalan di local/Docker)
*   `make`
*   `swag` (untuk regenerasi docs): `go install github.com/swaggo/swag/cmd/swag@latest`

### Setup Backend
1.  **Masuk ke direktori backend**:
    ```bash
    cd backend
    ```

2.  **Buat file `.env`**:
    Salin dari file sampel dan sesuaikan dengan konfigurasi database Anda.
    ```bash
    cp .sample.env .env
    ```
    Contoh isi `.env`:
    ```
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_NAME=hr_payroll_db
    DB_PORT=5432
    ```

3.  **Install dependencies**:
    ```bash
    go mod tidy
    ```

4.  **Jalankan Server Backend**:
    Server akan berjalan di `http://localhost:8080`.
    ```bash
    make run
    ```
    GORM akan otomatis membuat tabel (`AutoMigrate`) saat server pertama kali dijalankan.

### Setup Frontend
1.  **Masuk ke direktori frontend**:
    ```bash
    cd frontend
    ```

2.  **Jalankan server sederhana**:
    Frontend akan disajikan di `http://localhost:3000`.
    ```bash
    # Menggunakan Python 3
    python3 -m http.server 3000
    ```

### Menjalankan dengan `make` (Root Direktori)
Cara termudah adalah menggunakan `Makefile` di root project.

1.  **Jalankan Backend & Frontend (background process)**:
    ```bash
    make dev
    ```
    *   Backend: `http://localhost:8080`
    *   Frontend: `http://localhost:3000`

2.  **Buka frontend di browser**:
    ```bash
    make open
    ```

3.  **Lihat log**:
    Log disimpan di folder `.tmp/`.

4.  **Hentikan server**:
    ```bash
    make dev-stop
    ```

### Dokumentasi API
Dokumentasi Swagger tersedia di:
**http://localhost:8080/swagger/index.html**

## 6. Catatan Tambahan

*   **Auto Migration**: Fitur `AutoMigrate` dari GORM digunakan untuk kemudahan development. Untuk production, disarankan menggunakan sistem migrasi yang lebih robust seperti `golang-migrate`.
*   **Error Handling**: Error handling masih dasar. Bisa ditingkatkan dengan response error yang lebih terstruktur.
*   **Frontend**: Frontend dibuat sangat sederhana untuk mendemonstrasikan fungsionalitas backend. Belum ada handling untuk semua edge case (misal, state saat loading).
*   **Testing**: Belum ada unit test atau integration test. Ini adalah langkah penting selanjutnya yang perlu ditambahkan.
