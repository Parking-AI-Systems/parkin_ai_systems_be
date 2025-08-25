-- Initialize database tables for parkin-ai-system
-- This file will be executed when PostgreSQL container starts

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    role VARCHAR(50) DEFAULT 'user',
    avatar_url VARCHAR(255),
    wallet_balance DECIMAL(10,2) DEFAULT 0.00,
    gender VARCHAR(10),
    birth_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create api_tokens table for JWT refresh tokens
CREATE TABLE IF NOT EXISTS api_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create parking_lots table
CREATE TABLE IF NOT EXISTS parking_lots (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    total_slots INTEGER DEFAULT 0,
    available_slots INTEGER DEFAULT 0,
    hourly_rate DECIMAL(8,2),
    description TEXT,
    amenities TEXT,
    operating_hours VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create parking_slots table
CREATE TABLE IF NOT EXISTS parking_slots (
    id BIGSERIAL PRIMARY KEY,
    parking_lot_id BIGINT NOT NULL REFERENCES parking_lots(id) ON DELETE CASCADE,
    slot_number VARCHAR(50) NOT NULL,
    slot_type VARCHAR(50) DEFAULT 'standard',
    is_available BOOLEAN DEFAULT true,
    is_reserved BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create vehicles table
CREATE TABLE IF NOT EXISTS vehicles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    license_plate VARCHAR(20) UNIQUE NOT NULL,
    vehicle_type VARCHAR(50),
    brand VARCHAR(100),
    model VARCHAR(100),
    color VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create parking_orders table
CREATE TABLE IF NOT EXISTS parking_orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parking_lot_id BIGINT NOT NULL REFERENCES parking_lots(id) ON DELETE CASCADE,
    parking_slot_id BIGINT REFERENCES parking_slots(id) ON DELETE SET NULL,
    vehicle_id BIGINT REFERENCES vehicles(id) ON DELETE SET NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    duration_hours DECIMAL(4,2),
    total_amount DECIMAL(10,2),
    status VARCHAR(50) DEFAULT 'active',
    payment_status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_api_tokens_user_id ON api_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_api_tokens_token ON api_tokens(token);
CREATE INDEX IF NOT EXISTS idx_parking_slots_lot_id ON parking_slots(parking_lot_id);
CREATE INDEX IF NOT EXISTS idx_vehicles_user_id ON vehicles(user_id);
CREATE INDEX IF NOT EXISTS idx_parking_orders_user_id ON parking_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_parking_orders_status ON parking_orders(status);

-- Insert sample data
INSERT INTO users (username, password_hash, full_name, email, phone, role) VALUES 
('admin', '$2a$14$example_hash_here', 'System Administrator', 'admin@parkin-ai.com', '+1234567890', 'admin'),
('testuser', '$2a$14$example_hash_here', 'Test User', 'test@example.com', '+1234567891', 'user')
ON CONFLICT (email) DO NOTHING;

INSERT INTO parking_lots (name, address, latitude, longitude, total_slots, available_slots, hourly_rate, description) VALUES 
('Downtown Plaza', '123 Main Street, Downtown', 1.2966, 103.8520, 100, 95, 5.00, 'Premium parking in the heart of downtown'),
('Airport Terminal', 'Changi Airport Terminal 1', 1.3644, 103.9915, 500, 480, 3.50, 'Long-term and short-term parking near airport')
ON CONFLICT DO NOTHING;
