-- 数据库初始化SQL，需在 MySQL 8.0+ 执行
CREATE DATABASE IF NOT EXISTS collab_node_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE collab_node_platform;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    user_id CHAR(36) PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at DATETIME,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 任务表
CREATE TABLE IF NOT EXISTS tasks (
    task_id CHAR(36) PRIMARY KEY,
    task_name VARCHAR(128) NOT NULL,
    creator_id CHAR(36) NOT NULL,
    status TINYINT DEFAULT 0,
    created_at DATETIME,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users(user_id)
) ENGINE=InnoDB;

-- 任务成员表
CREATE TABLE IF NOT EXISTS task_members (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    task_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    role TINYINT NOT NULL,
    created_at DATETIME,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY task_user_unique (task_id, user_id),
    FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 节点表
CREATE TABLE IF NOT EXISTS nodes (
    node_id CHAR(36) PRIMARY KEY,
    task_id CHAR(36) NOT NULL,
    parent_node_id CHAR(36) DEFAULT NULL,
    node_name VARCHAR(128) NOT NULL,
    node_content TEXT,
    site VARCHAR(255) DEFAULT NULL,
    result TEXT DEFAULT NULL,
    next_step TEXT DEFAULT NULL,
    creator_id CHAR(36) NOT NULL,
    version INT DEFAULT 1,
    created_at DATETIME,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE,
    FOREIGN KEY (creator_id) REFERENCES users(user_id),
    FOREIGN KEY (parent_node_id) REFERENCES nodes(node_id) ON DELETE SET NULL
) ENGINE=InnoDB;

-- 操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    task_id CHAR(36) NOT NULL,
    node_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    operation_type ENUM('create','edit','delete') NOT NULL,
    operation_content TEXT,
    created_at DATETIME,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(task_id) ON DELETE CASCADE,
    FOREIGN KEY (node_id) REFERENCES nodes(node_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
) ENGINE=InnoDB;