-- 分布式微服务架构平台数据库初始化脚本

-- 创建用户服务数据库
CREATE DATABASE microservice_user;

-- 创建订单服务数据库
CREATE DATABASE microservice_order;

-- 创建支付服务数据库
CREATE DATABASE microservice_payment;

-- 切换到用户服务数据库
\c microservice_user;

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    avatar VARCHAR(255),
    status INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 创建用户表索引
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- 插入测试用户数据
INSERT INTO users (username, email, password_hash, phone, status) VALUES
('admin', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1VrVvmXGHxVYpkmLCS2kstNfQUb.inu', '13800138000', 1),
('user1', 'user1@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1VrVvmXGHxVYpkmLCS2kstNfQUb.inu', '13800138001', 1),
('user2', 'user2@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMye1VrVvmXGHxVYpkmLCS2kstNfQUb.inu', '13800138002', 1)
ON CONFLICT (username) DO NOTHING;

-- 切换到订单服务数据库
\c microservice_order;

-- 创建订单表
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    order_no VARCHAR(50) UNIQUE NOT NULL,
    user_id INTEGER NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    status INTEGER DEFAULT 1,
    payment_status INTEGER DEFAULT 0,
    shipping_address TEXT,
    remark TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 创建订单商品表
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    product_id INTEGER NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_price DECIMAL(10,2) NOT NULL,
    quantity INTEGER NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建订单表索引
CREATE INDEX IF NOT EXISTS idx_orders_order_no ON orders(order_no);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_payment_status ON orders(payment_status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);

-- 插入测试订单数据
INSERT INTO orders (order_no, user_id, total_amount, status, payment_status, shipping_address) VALUES
('ORD202401010001', 1, 299.99, 1, 0, '北京市朝阳区xxx街道xxx号'),
('ORD202401010002', 2, 199.99, 1, 1, '上海市浦东新区xxx路xxx号'),
('ORD202401010003', 3, 399.99, 2, 1, '广州市天河区xxx大道xxx号')
ON CONFLICT (order_no) DO NOTHING;

-- 插入测试订单商品数据
INSERT INTO order_items (order_id, product_id, product_name, product_price, quantity, total_price) VALUES
(1, 1001, '商品A', 99.99, 3, 299.97),
(2, 1002, '商品B', 199.99, 1, 199.99),
(3, 1003, '商品C', 199.99, 2, 399.98);

-- 切换到支付服务数据库
\c microservice_payment;

-- 创建支付记录表
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    payment_no VARCHAR(50) UNIQUE NOT NULL,
    order_no VARCHAR(50) NOT NULL,
    user_id INTEGER NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    payment_method INTEGER NOT NULL, -- 1:支付宝 2:微信 3:银行卡
    payment_status INTEGER DEFAULT 0, -- 0:待支付 1:已支付 2:支付失败 3:已退款
    transaction_id VARCHAR(100),
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建退款记录表
CREATE TABLE IF NOT EXISTS refunds (
    id SERIAL PRIMARY KEY,
    refund_no VARCHAR(50) UNIQUE NOT NULL,
    payment_id INTEGER NOT NULL REFERENCES payments(id),
    order_no VARCHAR(50) NOT NULL,
    user_id INTEGER NOT NULL,
    refund_amount DECIMAL(10,2) NOT NULL,
    refund_reason TEXT,
    refund_status INTEGER DEFAULT 0, -- 0:申请中 1:已退款 2:退款失败
    refunded_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建支付表索引
CREATE INDEX IF NOT EXISTS idx_payments_payment_no ON payments(payment_no);
CREATE INDEX IF NOT EXISTS idx_payments_order_no ON payments(order_no);
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(payment_status);
CREATE INDEX IF NOT EXISTS idx_payments_method ON payments(payment_method);
CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at);
CREATE INDEX IF NOT EXISTS idx_refunds_refund_no ON refunds(refund_no);
CREATE INDEX IF NOT EXISTS idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX IF NOT EXISTS idx_refunds_order_no ON refunds(order_no);
CREATE INDEX IF NOT EXISTS idx_refunds_user_id ON refunds(user_id);
CREATE INDEX IF NOT EXISTS idx_refunds_status ON refunds(refund_status);

-- 插入测试支付数据
INSERT INTO payments (payment_no, order_no, user_id, amount, payment_method, payment_status, transaction_id, paid_at) VALUES
('PAY202401010001', 'ORD202401010001', 1, 299.99, 1, 0, NULL, NULL),
('PAY202401010002', 'ORD202401010002', 2, 199.99, 2, 1, 'WX20240101001', '2024-01-01 10:30:00+08'),
('PAY202401010003', 'ORD202401010003', 3, 399.99, 1, 1, 'ALI20240101001', '2024-01-01 11:15:00+08')
ON CONFLICT (payment_no) DO NOTHING;

-- 创建数据库函数：更新updated_at字段
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 为各表创建更新时间触发器
\c microservice_user;
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

\c microservice_order;
CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

\c microservice_payment;
CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_refunds_updated_at BEFORE UPDATE ON refunds FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();