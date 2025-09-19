-- Create order_disputes table
CREATE TABLE IF NOT EXISTS order_disputes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'escalated', 'resolved', 'closed')),
    resolution TEXT,
    evidence JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_order_disputes_tenant_id ON order_disputes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_order_disputes_order_id ON order_disputes(order_id);
CREATE INDEX IF NOT EXISTS idx_order_disputes_customer_id ON order_disputes(customer_id);
CREATE INDEX IF NOT EXISTS idx_order_disputes_status ON order_disputes(status);
CREATE INDEX IF NOT EXISTS idx_order_disputes_created_at ON order_disputes(created_at);

-- Add trigger to update updated_at timestamp
CREATE TRIGGER update_order_disputes_updated_at
    BEFORE UPDATE ON order_disputes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();