-- Migration: Add display columns to categories table
-- Created: 2024-12-19

-- Add show_in_menu and is_featured columns to categories table
ALTER TABLE categories 
ADD COLUMN IF NOT EXISTS show_in_menu BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS is_featured BOOLEAN NOT NULL DEFAULT false;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_categories_show_in_menu ON categories(show_in_menu);
CREATE INDEX IF NOT EXISTS idx_categories_is_featured ON categories(is_featured);

-- Update existing categories to have reasonable defaults
-- Set active categories to show in menu by default
UPDATE categories SET show_in_menu = true WHERE is_active = true;

-- Comment: This migration adds the missing display columns that are used
-- in the category filtering functionality but were missing from the original schema.