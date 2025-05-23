-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Waktu pembuatan: 22 Bulan Mei 2025 pada 17.59
-- Versi server: 10.4.28-MariaDB
-- Versi PHP: 8.2.4

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `test_pt_xyz`
--

-- --------------------------------------------------------

--
-- Struktur dari tabel `customers`
--

CREATE TABLE `customers` (
  `id` varchar(100) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` text NOT NULL,
  `nik` varchar(16) NOT NULL,
  `full_name` varchar(100) DEFAULT NULL,
  `legal_name` varchar(100) DEFAULT NULL,
  `place_of_birth` varchar(100) DEFAULT NULL,
  `date_of_birth` datetime DEFAULT NULL,
  `salary` decimal(18,2) DEFAULT NULL,
  `ktp_photo_url` text DEFAULT NULL,
  `selfie_photo_url` text DEFAULT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `customers`
--

INSERT INTO `customers` (`id`, `email`, `password`, `nik`, `full_name`, `legal_name`, `place_of_birth`, `date_of_birth`, `salary`, `ktp_photo_url`, `selfie_photo_url`, `created_at`, `updated_at`, `deleted_at`) VALUES
('2e4745e6-77ad-41d0-8a2a-368a301ab752', 'annisa@xyz.com', '$2a$12$Rvslxj25D4OU7w3Ercz/IucMiDkEp1dOCSwq902oWpy0mqcUx2GAq', '3175091401920002', 'Annisa Rahim', 'Annisa Rahim', 'Jakarta', '2022-01-01 00:00:00', 8500000.00, 'http://localhost:8000/storage/59c4a222-06a5-487e-a9ba-4a99b0362bc7.png', 'http://localhost:8000/storage/7410f8ef-2b34-4571-b4ae-421528678ff2.png', '2025-05-20 06:06:22', '2025-05-22 22:58:29', NULL),
('7158e6f8-9cb5-4059-8556-840bf2facd5b', 'budi@xyz.com', '$2a$12$Rvslxj25D4OU7w3Ercz/IucMiDkEp1dOCSwq902oWpy0mqcUx2GAq', '3175091401924356', 'Budi Doremi', 'Budi Doremi', 'Bali', '2022-01-01 00:00:00', 8500000.00, 'http://localhost:8000/storage/a82ff124-5a1b-4879-8dba-71016fef0a80.png', 'http://localhost:8000/storage/87bec52b-2e8a-470d-bfdf-bc80b7cea84d.png', '2025-05-21 07:41:14', '2025-05-22 22:58:33', NULL);

-- --------------------------------------------------------

--
-- Struktur dari tabel `limit`
--

CREATE TABLE `limit` (
  `id` char(36) NOT NULL,
  `customer_id` char(36) DEFAULT NULL,
  `tenor_months` int(20) DEFAULT NULL,
  `limit_amount` decimal(18,2) DEFAULT NULL,
  `status` enum('available','booked') DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `limit`
--

INSERT INTO `limit` (`id`, `customer_id`, `tenor_months`, `limit_amount`, `status`, `created_at`, `updated_at`, `deleted_at`) VALUES
('12c4d394-8c32-4bf2-9b28-b83e3fc19bbb', '7158e6f8-9cb5-4059-8556-840bf2facd5b', 1, 100000.00, 'available', '2025-05-22 09:16:47', '2025-05-22 09:16:47', NULL),
('bdc34a06-1fa9-4f94-a7ef-4cd74d41774a', '7158e6f8-9cb5-4059-8556-840bf2facd5b', 3, 300000.00, 'available', '2025-05-22 09:17:12', '2025-05-22 09:17:12', NULL),
('c7e398e5-fd09-44a9-83da-57f99ccc6aef', '2e4745e6-77ad-41d0-8a2a-368a301ab752', 6, 2000000.00, 'available', '2025-05-21 10:13:27', '2025-05-22 14:42:35', NULL),
('cec87fd0-2a35-4b80-9d1c-7cc38817f68d', '2e4745e6-77ad-41d0-8a2a-368a301ab752', 3, 1500000.00, 'available', '2025-05-22 08:08:22', '2025-05-22 14:42:56', NULL),
('d7a82828-c65d-4e33-a04b-4727132daa3b', '2e4745e6-77ad-41d0-8a2a-368a301ab752', 2, 1200000.00, 'available', '2025-05-20 10:30:58', '2025-05-22 14:43:02', NULL),
('da888ba6-355b-11f0-a1ec-088fc37ce4dd', '2e4745e6-77ad-41d0-8a2a-368a301ab752', 1, 1000000.00, 'available', '2025-05-22 16:18:34', '2025-05-22 14:43:08', NULL),
('e341c797-0bf1-4036-849d-3d4b7c0acd19', '7158e6f8-9cb5-4059-8556-840bf2facd5b', 6, 700000.00, 'available', '2025-05-22 09:17:22', '2025-05-22 09:17:22', NULL),
('f477c286-2a4f-4134-bdd2-da8c29c90a48', '7158e6f8-9cb5-4059-8556-840bf2facd5b', 2, 200000.00, 'available', '2025-05-22 09:17:04', '2025-05-22 09:17:04', NULL);

-- --------------------------------------------------------

--
-- Struktur dari tabel `transactions`
--

CREATE TABLE `transactions` (
  `id` char(36) NOT NULL,
  `customer_id` char(36) DEFAULT NULL,
  `contract_number` varchar(50) NOT NULL,
  `channel` varchar(50) DEFAULT NULL,
  `otr_amount` decimal(18,2) DEFAULT NULL,
  `admin_fee` decimal(18,2) DEFAULT NULL,
  `installment` decimal(18,2) DEFAULT NULL,
  `interest` decimal(18,2) DEFAULT NULL,
  `asset_name` varchar(100) DEFAULT NULL,
  `tenor_months` int(11) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `users`
--

CREATE TABLE `users` (
  `id` char(36) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` text NOT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `users`
--

INSERT INTO `users` (`id`, `email`, `password`, `created_at`, `updated_at`, `deleted_at`) VALUES
('4246fb58-ff45-4d2a-8946-93e541fc39fd', 'admin@admin.com', '$2a$12$Rvslxj25D4OU7w3Ercz/IucMiDkEp1dOCSwq902oWpy0mqcUx2GAq', '2025-05-21 16:52:12', '2025-05-21 16:52:12', NULL);

--
-- Indexes for dumped tables
--

--
-- Indeks untuk tabel `customers`
--
ALTER TABLE `customers`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `nik` (`nik`),
  ADD KEY `email` (`email`);

--
-- Indeks untuk tabel `limit`
--
ALTER TABLE `limit`
  ADD PRIMARY KEY (`id`),
  ADD KEY `customer_id` (`customer_id`);

--
-- Indeks untuk tabel `transactions`
--
ALTER TABLE `transactions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `contract_number` (`contract_number`),
  ADD KEY `customer_id` (`customer_id`);

--
-- Indeks untuk tabel `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `email` (`email`);

--
-- Ketidakleluasaan untuk tabel pelimpahan (Dumped Tables)
--

--
-- Ketidakleluasaan untuk tabel `limit`
--
ALTER TABLE `limit`
  ADD CONSTRAINT `limit_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`);

--
-- Ketidakleluasaan untuk tabel `transactions`
--
ALTER TABLE `transactions`
  ADD CONSTRAINT `transactions_ibfk_1` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
