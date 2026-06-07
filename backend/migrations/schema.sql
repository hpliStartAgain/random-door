-- schema.sql — 15 张表建表 DDL（对应 database-detailed-design.md）
-- 枚举用 VARCHAR + 应用层校验；MVP 不建物理外键，*_id 字段建普通索引。
-- 字符集统一 utf8mb4。

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  anonymous_id VARCHAR(128) NOT NULL,
  nickname VARCHAR(64),
  avatar_url VARCHAR(512),
  current_city_id BIGINT,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  UNIQUE KEY uk_anonymous_id (anonymous_id),
  KEY idx_user_current_city (current_city_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS cities (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  province VARCHAR(64) NOT NULL,
  lat DOUBLE NOT NULL,
  lng DOUBLE NOT NULL,
  intro TEXT,
  cover_image_url VARCHAR(512),
  dialect_sample VARCHAR(255),
  dialect_explanation TEXT,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  UNIQUE KEY uk_city_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS city_tags (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  city_id BIGINT NOT NULL,
  tag VARCHAR(64) NOT NULL,
  created_at DATETIME NOT NULL,
  UNIQUE KEY uk_ct_city_tag (city_id, tag),
  KEY idx_ct_city (city_id),
  KEY idx_ct_tag (tag)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS landmarks (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  city_id BIGINT NOT NULL,
  name VARCHAR(128) NOT NULL,
  image_url VARCHAR(512),
  description TEXT,
  created_at DATETIME NOT NULL,
  UNIQUE KEY uk_lm_city_name (city_id, name),
  KEY idx_lm_city (city_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS foods (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  city_id BIGINT NOT NULL,
  name VARCHAR(128) NOT NULL,
  image_url VARCHAR(512),
  description TEXT,
  created_at DATETIME NOT NULL,
  UNIQUE KEY uk_food_city_name (city_id, name),
  KEY idx_food_city (city_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS characters (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  city_id BIGINT NOT NULL,
  name VARCHAR(128) NOT NULL,
  character_type VARCHAR(32) NOT NULL,
  avatar_url VARCHAR(512),
  persona TEXT NOT NULL,
  dialect_style TEXT,
  role_title VARCHAR(128),
  life_span VARCHAR(64),
  intro_quote VARCHAR(255),
  prompt TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  UNIQUE KEY uk_char_city_name (city_id, name),
  KEY idx_char_city (city_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS city_visits (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  city_id BIGINT NOT NULL,
  visit_mode VARCHAR(32) NOT NULL,
  source VARCHAR(64),
  from_city_id BIGINT,
  dice_roll_id BIGINT,
  created_at DATETIME NOT NULL,
  KEY idx_cv_user_city (user_id, city_id),
  KEY idx_cv_user_time (user_id, created_at),
  KEY idx_cv_city (city_id),
  KEY idx_cv_from_city (from_city_id),
  KEY idx_cv_dice_roll (dice_roll_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS dice_rolls (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  from_city_id BIGINT,
  to_city_id BIGINT NOT NULL,
  direction VARCHAR(32) NOT NULL,
  distance_km INT NOT NULL,
  target_lat DOUBLE,
  target_lng DOUBLE,
  created_at DATETIME NOT NULL,
  KEY idx_dr_user_time (user_id, created_at),
  KEY idx_dr_from_city (from_city_id),
  KEY idx_dr_to_city (to_city_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS checkins (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  city_id BIGINT NOT NULL,
  landmark_id BIGINT,
  visit_id BIGINT,
  generated_image_url VARCHAR(512),
  checkin_mode VARCHAR(32),
  created_at DATETIME NOT NULL,
  KEY idx_ck_user_time (user_id, created_at),
  KEY idx_ck_user_city (user_id, city_id),
  KEY idx_ck_city (city_id),
  KEY idx_ck_landmark (landmark_id),
  KEY idx_ck_visit (visit_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS achievements (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(128) NOT NULL,
  description TEXT,
  rule_type VARCHAR(64) NOT NULL,
  rule_value VARCHAR(255) NOT NULL,
  badge_url VARCHAR(512),
  created_at DATETIME NOT NULL,
  UNIQUE KEY uk_ach_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS user_achievements (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  achievement_id BIGINT NOT NULL,
  unlocked_at DATETIME NOT NULL,
  UNIQUE KEY uk_ua_user_ach (user_id, achievement_id),
  KEY idx_ua_achievement (achievement_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS chat_messages (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  city_id BIGINT NOT NULL,
  character_id BIGINT NOT NULL,
  role VARCHAR(16) NOT NULL,
  content TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  KEY idx_cm_user_char_time (user_id, character_id, created_at),
  KEY idx_cm_city (city_id),
  KEY idx_cm_character (character_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS comments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  target_type VARCHAR(32) NOT NULL,
  target_id BIGINT NOT NULL,
  user_id BIGINT,
  nickname VARCHAR(64) NOT NULL,
  content VARCHAR(500) NOT NULL,
  created_at DATETIME NOT NULL,
  KEY idx_comments_target_time (target_type, target_id, created_at),
  KEY idx_comments_user_time (user_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ai_tasks (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  type VARCHAR(32) NOT NULL,
  status VARCHAR(32) NOT NULL,
  input_json JSON NOT NULL,
  result_url VARCHAR(512),
  error TEXT,
  attempts INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  KEY idx_ai_tasks_user_time (user_id, created_at),
  KEY idx_ai_tasks_status_time (status, updated_at),
  KEY idx_ai_tasks_type_status (type, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS ai_usage_logs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  usage_type VARCHAR(32) NOT NULL,
  usage_date DATE NOT NULL,
  count INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  UNIQUE KEY uk_ai_usage_user_type_date (user_id, usage_type, usage_date),
  KEY idx_ai_usage_date (usage_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
