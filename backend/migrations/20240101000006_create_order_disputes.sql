-- Create order_disputes table
CREATE TABLE IF NOT EXISTS order_disputes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    order_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    reason VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    resolution TEXT,
    evidence JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    resolved_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_order_disputes_tenant_id ON order_disputes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_order_disputes_order_id ON order_disputes(order_id);
CREATE INDEX IF NOT EXISTS idx_order_disputes_customer_id ON order_disputes(customer_id);
CREATE INDEX IF NOT EXISTS idx_order_disputes_status ON order_disputes(status);
CREATE INDEX IF NOT EXISTS idx_order_disputes_created_at ON order_disputes(created_at);

-- Add foreign key constraints
ALTER TABLE order_disputes 
ADD CONSTRAINT fk_order_disputes_order_id 
FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE;

ALTER TABLE order_disputes 
ADD CONSTRAINT fk_order_disputes_customer_id 
FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add check constraints for valid status values
ALTER TABLE order_disputes 
ADD CONSTRAINT chk_order_disputes_status 
CHECK (status IN ('pending', 'escalated', 'resolved', 'closed'));

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_order_disputes_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_order_disputes_updated_at
    BEFORE UPDATE ON order_disputes
    FOR EACH ROW
    EXECUTE FUNCTION update_order_disputes_updated_at();