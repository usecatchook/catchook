CREATE TYPE auth_type AS ENUM ('none', 'basic', 'bearer', 'apikey', 'signature');
CREATE TYPE protocol_type AS ENUM ('http', 'grpc', 'mqtt', 'websocket');
CREATE TYPE destination_type AS ENUM ('http', 'rabbitmq', 'database', 'file', 'queue', 'cli');
CREATE TYPE filter_mode AS ENUM ('nocode', 'code');
CREATE TYPE transformation_mode AS ENUM ('nocode', 'code');
CREATE TYPE webhook_status AS ENUM ('pending', 'filtered', 'transformed', 'delayed', 'delivered', 'failed');
CREATE TYPE delivery_status AS ENUM ('pending', 'success', 'failed', 'retrying');
CREATE TYPE filter_type AS ENUM ('condition', 'javascript', 'jsonpath', 'regex');
CREATE TYPE transformation_type AS ENUM ('header_add', 'header_remove', 'header_modify', 'body_add', 'body_remove', 'body_modify', 'format_json', 'format_xml', 'javascript', 'jsonpath');
CREATE TYPE step_type AS ENUM ('auth', 'filter', 'transformation', 'delivery');
CREATE TYPE step_status AS ENUM ('pending', 'success', 'failed', 'skipped');

-- Table sources
CREATE TABLE IF NOT EXISTS sources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    protocol protocol_type NOT NULL,
    auth_type auth_type NOT NULL DEFAULT 'none',
    auth_config JSONB DEFAULT '{}'::jsonb,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sources_user_id ON sources(user_id);
CREATE INDEX IF NOT EXISTS idx_sources_name ON sources(name);
CREATE INDEX IF NOT EXISTS idx_sources_is_active ON sources(is_active);
CREATE INDEX IF NOT EXISTS idx_sources_created_at ON sources(created_at);

-- Table webhook_events
CREATE TABLE IF NOT EXISTS webhook_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    pipeline_id UUID REFERENCES pipelines(id),
    payload JSONB NOT NULL,
    original_payload JSONB,
    metadata JSONB DEFAULT '{}'::jsonb,
    filter_results JSONB DEFAULT '{}'::jsonb,
    transformation_results JSONB DEFAULT '{}'::jsonb,
    status webhook_status NOT NULL DEFAULT 'pending',
    error_message TEXT,
    scheduled_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_webhook_events_source_id ON webhook_events(source_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_pipeline_id ON webhook_events(pipeline_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_status ON webhook_events(status);
CREATE INDEX IF NOT EXISTS idx_webhook_events_created_at ON webhook_events(created_at);
CREATE INDEX IF NOT EXISTS idx_webhook_events_processed_at ON webhook_events(processed_at);
CREATE INDEX IF NOT EXISTS idx_webhook_events_scheduled_at ON webhook_events(scheduled_at);

-- Table destinations
CREATE TABLE destinations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    destination_type destination_type NOT NULL DEFAULT 'http',
    config JSONB DEFAULT '{}'::jsonb,
    is_active BOOLEAN DEFAULT TRUE,
    delay_seconds INTEGER DEFAULT 0,
    retry_attempts INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_destinations_user_id ON destinations(user_id);
CREATE INDEX IF NOT EXISTS idx_destinations_is_active ON destinations(is_active);
CREATE INDEX IF NOT EXISTS idx_destinations_created_at ON destinations(created_at);

-- Table pipelines
CREATE TABLE IF NOT EXISTS pipelines (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    destination_id UUID NOT NULL REFERENCES destinations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    execution_order INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, source_id, destination_id, name)
);

CREATE INDEX IF NOT EXISTS idx_pipelines_user_id ON pipelines(user_id);
CREATE INDEX IF NOT EXISTS idx_pipelines_source_id ON pipelines(source_id);
CREATE INDEX IF NOT EXISTS idx_pipelines_destination_id ON pipelines(destination_id);
CREATE INDEX IF NOT EXISTS idx_pipelines_is_active ON pipelines(is_active);
CREATE INDEX IF NOT EXISTS idx_pipelines_execution_order ON pipelines(execution_order);

-- Table filters
CREATE TABLE IF NOT EXISTS filters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pipeline_id UUID NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '',
    filter_type filter_type NOT NULL,
    mode filter_mode NOT NULL DEFAULT 'nocode',
    config JSONB DEFAULT '{}'::jsonb,
    code TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    execution_order INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_filters_pipeline_id ON filters(pipeline_id);
CREATE INDEX IF NOT EXISTS idx_filters_is_active ON filters(is_active);
CREATE INDEX IF NOT EXISTS idx_filters_execution_order ON filters(execution_order);
CREATE INDEX IF NOT EXISTS idx_filters_filter_type ON filters(filter_type);

-- Table transformations
CREATE TABLE IF NOT EXISTS transformations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pipeline_id UUID NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT DEFAULT '',
    transformation_type transformation_type NOT NULL,
    mode transformation_mode NOT NULL DEFAULT 'nocode',
    config JSONB DEFAULT '{}'::jsonb,
    code TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    execution_order INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transformations_pipeline_id ON transformations(pipeline_id);
CREATE INDEX IF NOT EXISTS idx_transformations_is_active ON transformations(is_active);
CREATE INDEX IF NOT EXISTS idx_transformations_execution_order ON transformations(execution_order);
CREATE INDEX IF NOT EXISTS idx_transformations_type ON transformations(transformation_type);

-- Table deliveries
CREATE TABLE IF NOT EXISTS deliveries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    webhook_event_id UUID NOT NULL REFERENCES webhook_events(id),
    destination_id UUID NOT NULL REFERENCES destinations(id),
    status delivery_status NOT NULL DEFAULT 'pending',
    response_code INTEGER,
    attempt INTEGER DEFAULT 0,
    last_error TEXT,
    scheduled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_deliveries_webhook_event_id ON deliveries(webhook_event_id);
CREATE INDEX IF NOT EXISTS idx_deliveries_destination_id ON deliveries(destination_id);
CREATE INDEX IF NOT EXISTS idx_deliveries_status ON deliveries(status);
CREATE INDEX IF NOT EXISTS idx_deliveries_created_at ON deliveries(created_at);

-- Table webhook_steps
CREATE TABLE IF NOT EXISTS webhook_steps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    webhook_event_id UUID NOT NULL REFERENCES webhook_events(id) ON DELETE CASCADE,
    pipeline_id UUID REFERENCES pipelines(id),
    step_type step_type NOT NULL,
    step_name VARCHAR(100) NOT NULL,
    step_id UUID,
    execution_order INTEGER NOT NULL,
    status step_status NOT NULL DEFAULT 'pending',
    input_data JSONB DEFAULT '{}'::jsonb,
    output_data JSONB DEFAULT '{}'::jsonb,
    error_message TEXT,
    duration_ms INTEGER,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_webhook_steps_webhook_event_id ON webhook_steps(webhook_event_id);
CREATE INDEX IF NOT EXISTS idx_webhook_steps_pipeline_id ON webhook_steps(pipeline_id);
CREATE INDEX IF NOT EXISTS idx_webhook_steps_step_type ON webhook_steps(step_type);
CREATE INDEX IF NOT EXISTS idx_webhook_steps_status ON webhook_steps(status);
CREATE INDEX IF NOT EXISTS idx_webhook_steps_execution_order ON webhook_steps(execution_order);
CREATE INDEX IF NOT EXISTS idx_webhook_steps_started_at ON webhook_steps(started_at);
CREATE INDEX IF NOT EXISTS idx_webhook_steps_duration_ms ON webhook_steps(duration_ms);



CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_sources_updated_at
    BEFORE UPDATE ON sources
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_webhook_events_updated_at
    BEFORE UPDATE ON webhook_events
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_pipelines_updated_at
    BEFORE UPDATE ON pipelines
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_filters_updated_at
    BEFORE UPDATE ON filters
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transformations_updated_at
    BEFORE UPDATE ON transformations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_deliveries_updated_at
    BEFORE UPDATE ON deliveries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_destinations_updated_at
    BEFORE UPDATE ON destinations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_webhook_steps_updated_at
    BEFORE UPDATE ON webhook_steps
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();