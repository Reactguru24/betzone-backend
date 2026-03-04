-- Migration: Add provider_player_id column to users table
ALTER TABLE users ADD COLUMN provider_player_id VARCHAR(50) UNIQUE;
