-- EVE Navigation & Intelligence System - SQL Views
-- These views provide foundation for route planning and analysis

-- =============================================================================
-- v_stargate_graph: Bidirectional stargate connectivity graph
-- Extracts gate connections from JSON destination field
-- Used for pathfinding and route calculation
-- =============================================================================
CREATE VIEW IF NOT EXISTS v_stargate_graph AS
SELECT 
    s.solarSystemID as from_system_id,
    CAST(json_extract(s.destination, '$.solarSystemID') AS INTEGER) as to_system_id,
    s._key as gate_id,
    s.typeID as gate_type_id
FROM mapStargates s
WHERE json_extract(s.destination, '$.solarSystemID') IS NOT NULL

UNION ALL

-- Reverse direction (bidirectional edges)
SELECT 
    CAST(json_extract(s.destination, '$.solarSystemID') AS INTEGER) as from_system_id,
    s.solarSystemID as to_system_id,
    s._key as gate_id,
    s.typeID as gate_type_id
FROM mapStargates s
WHERE json_extract(s.destination, '$.solarSystemID') IS NOT NULL;

-- =============================================================================
-- v_system_info: Enhanced system information with parsed names and security zones
-- Provides human-readable system data for routing and analysis
-- =============================================================================
CREATE VIEW IF NOT EXISTS v_system_info AS
SELECT 
    sys._key as system_id,
    sys._key as solar_system_id,  -- _key IS the solar system ID
    COALESCE(json_extract(sys.name, '$.en'), json_extract(sys.name, '$.de')) as system_name,
    sys.securityStatus as security_status,
    CASE 
        WHEN sys.securityStatus >= 0.45 THEN 'High-Sec'
        WHEN sys.securityStatus > 0.0 THEN 'Low-Sec'
        WHEN sys.securityStatus <= 0.0 AND sys.wormholeClassID IS NULL THEN 'Null-Sec'
        WHEN sys.wormholeClassID IS NOT NULL THEN 'Wormhole'
        ELSE 'Unknown'
    END as security_zone,
    sys.constellationID as constellation_id,
    sys.regionID as region_id,
    COALESCE(json_extract(r.name, '$.en'), json_extract(r.name, '$.de')) as region_name,
    COALESCE(json_extract(c.name, '$.en'), json_extract(c.name, '$.de')) as constellation_name,
    sys.border,
    sys.corridor,
    sys.hub,
    sys.wormholeClassID as wormhole_class_id
FROM mapSolarSystems sys
LEFT JOIN mapRegions r ON sys.regionID = r._key
LEFT JOIN mapConstellations c ON sys.constellationID = c._key;

-- =============================================================================
-- v_system_security_zones: Security zone statistics by region/constellation
-- Useful for risk assessment and region analysis
-- =============================================================================
CREATE VIEW IF NOT EXISTS v_system_security_zones AS
SELECT 
    region_id,
    region_name,
    security_zone,
    COUNT(*) as system_count,
    ROUND(AVG(security_status), 3) as avg_security
FROM v_system_info
GROUP BY region_id, region_name, security_zone;

-- =============================================================================
-- v_region_stats: Comprehensive region statistics
-- Total systems, average security, and border system counts
-- =============================================================================
CREATE VIEW IF NOT EXISTS v_region_stats AS
SELECT 
    region_id,
    region_name,
    COUNT(*) as total_systems,
    ROUND(AVG(security_status), 3) as avg_security,
    SUM(CASE WHEN border = 1 THEN 1 ELSE 0 END) as border_systems,
    SUM(CASE WHEN security_zone = 'High-Sec' THEN 1 ELSE 0 END) as high_sec_count,
    SUM(CASE WHEN security_zone = 'Low-Sec' THEN 1 ELSE 0 END) as low_sec_count,
    SUM(CASE WHEN security_zone = 'Null-Sec' THEN 1 ELSE 0 END) as null_sec_count,
    SUM(CASE WHEN security_zone = 'Wormhole' THEN 1 ELSE 0 END) as wormhole_count
FROM v_system_info
GROUP BY region_id, region_name;

-- =============================================================================
-- v_trade_hubs: Pre-calculated major trade hub information
-- Known trade hub system IDs (from EVE data):
-- Jita (The Forge): 30000142
-- Amarr (Domain): 30002187
-- Dodixie (Sinq Laison): 30002659
-- Rens (Heimatar): 30002510
-- Hek (Metropolis): 30002053
-- =============================================================================
CREATE VIEW IF NOT EXISTS v_trade_hubs AS
SELECT 
    system_id,
    system_name,
    region_name,
    security_status,
    security_zone,
    CASE system_id
        WHEN 30000142 THEN 'Jita'
        WHEN 30002187 THEN 'Amarr'
        WHEN 30002659 THEN 'Dodixie'
        WHEN 30002510 THEN 'Rens'
        WHEN 30002053 THEN 'Hek'
        ELSE NULL
    END as hub_name
FROM v_system_info
WHERE system_id IN (30000142, 30002187, 30002659, 30002510, 30002053);
