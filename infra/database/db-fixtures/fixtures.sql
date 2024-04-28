INSERT INTO `user` (`id`, `name`, `user_name`, `email`, `password`, `cpf`, `cellphone_number`, `status`, `is_active`, `password_reset`, `created_at`, `updated_at`, `deleted_at`) VALUES
('8f1949d9-7c7f-4945-b86e-acb6bf27e53a', 'Get The Kill', 'GTK-MK2', 'davi@gmail.com', '$2a$13$C1lDXle1G3yk5bQkoq6tFOz4Oh78G7dAHY6C6sVvCdJhx647Y12QK', '36228024051', '31983685543', 'logged', 1, 0, '2024-03-16 12:46:45.000', NULL, NULL);

INSERT INTO `document` (`id`, `type`, `path`, `mime_type`, `post_id`, `owner_id`, `deleted_reason`, `created_at`, `updated_at`, `deleted_at`) VALUES
('ebb692fd-73a8-45ef-994e-06c1170b04d6', 'profile_picture', './tmp/ca60a8ec-d5b9-4d64-9791-b272d1180656_.jpg', 'image/jpeg', NULL, '8f1949d9-7c7f-4945-b86e-acb6bf27e53a', NULL, NULL, NULL, NULL);