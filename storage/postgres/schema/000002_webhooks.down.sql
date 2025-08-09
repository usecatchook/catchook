-- Drop triggers
DROP TRIGGER IF EXISTS update_sources_updated_at ON sources;
DROP TRIGGER IF EXISTS update_webhook_events_updated_at ON webhook_events;
DROP TRIGGER IF EXISTS update_rules_updated_at ON rules;
DROP TRIGGER IF EXISTS update_deliveries_updated_at ON deliveries;
DROP TRIGGER IF EXISTS update_destinations_updated_at ON destinations;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();


-- Drop indexes on sources
DROP INDEX IF EXISTS idx_sources_user_id;
DROP INDEX IF EXISTS idx_sources_is_active;
DROP INDEX IF EXISTS idx_sources_created_at;

-- Drop indexes on webhook_events
DROP INDEX IF EXISTS idx_webhook_events_source_id;
DROP INDEX IF EXISTS idx_webhook_events_status;
DROP INDEX IF EXISTS idx_webhook_events_created_at;

-- Drop indexes on rules
DROP INDEX IF EXISTS idx_rules_user_id;
DROP INDEX IF EXISTS idx_rules_is_active;
DROP INDEX IF EXISTS idx_rules_source_id;
DROP INDEX IF EXISTS idx_rules_destination_id;
DROP INDEX IF EXISTS idx_rules_created_at;

-- Drop indexes on deliveries
DROP INDEX IF EXISTS idx_deliveries_webhook_event_id;
DROP INDEX IF EXISTS idx_deliveries_destination_id;
DROP INDEX IF EXISTS idx_deliveries_status;
DROP INDEX IF EXISTS idx_deliveries_created_at;

-- Drop indexes on destinations
DROP INDEX IF EXISTS idx_destinations_user_id;
DROP INDEX IF EXISTS idx_destinations_is_active;
DROP INDEX IF EXISTS idx_destinations_created_at;

-- Drop tables (ordre important pour respecter les FK)
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS rules;
DROP TABLE IF EXISTS webhook_events;
DROP TABLE IF EXISTS destinations;
DROP TABLE IF EXISTS sources;

-- Drop enum types
DROP TYPE IF EXISTS auth_type;
DROP TYPE IF EXISTS protocol_type;
DROP TYPE IF EXISTS destination_type;
DROP TYPE IF EXISTS rule_mode;
DROP TYPE IF EXISTS webhook_status;
DROP TYPE IF EXISTS delivery_status;
