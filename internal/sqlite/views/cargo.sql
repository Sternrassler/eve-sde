-- EVE Cargo & Hauling System - SQL Views
-- These views provide foundation for cargo calculations and hauling optimization

-- =============================================================================
-- v_item_volumes: Item volume data for transport calculations
-- Provides volume, capacity, and value density information for all published items
-- =============================================================================
CREATE VIEW IF NOT EXISTS v_item_volumes AS
SELECT 
    t._key as type_id,
    COALESCE(json_extract(t.name, '$.en'), json_extract(t.name, '$.de')) as item_name,
    t.volume,
    t.capacity,
    t.packagedVolume,
    t.basePrice,
    g.categoryID as category_id,
    COALESCE(json_extract(g.name, '$.en'), json_extract(g.name, '$.de')) as category_name,
    mg.marketGroupID,
    -- ISK/m³ Ratio for Value-Density calculations
    CASE WHEN t.volume > 0 THEN CAST(t.basePrice AS REAL) / t.volume ELSE 0 END as isk_per_m3
FROM types t
LEFT JOIN groups g ON t.groupID = g._key
LEFT JOIN marketGroups mg ON t.marketGroupID = mg._key
WHERE t.published = 1;

-- =============================================================================
-- v_ship_cargo_capacities: Ship cargo capacity information
-- Provides base cargo capacities for all published ships (without skill bonuses)
-- Skill bonuses must be applied in application code
-- =============================================================================
CREATE VIEW IF NOT EXISTS v_ship_cargo_capacities AS
SELECT
    t._key as ship_type_id,
    COALESCE(json_extract(t.name, '$.en'), json_extract(t.name, '$.de')) as ship_name,
    t.volume as ship_volume,
    t.capacity as base_cargo_capacity,
    -- Special cargo holds via dogma attributes (base values)
    -- 1556 = Fleet Hangar Capacity
    (SELECT value FROM typeDogma WHERE typeID = t._key AND attributeID = 1556) as base_fleet_hangar_capacity,
    -- 38 = Capacity (used for ore holds on mining ships)
    (SELECT value FROM typeDogma WHERE typeID = t._key AND attributeID = 38) as base_ore_hold_capacity,
    -- Ship classification
    g._key as group_id,
    COALESCE(json_extract(g.name, '$.en'), json_extract(g.name, '$.de')) as group_name,
    c._key as category_id
FROM types t
JOIN groups g ON t.groupID = g._key
JOIN categories c ON g.categoryID = c._key
WHERE c._key = 6  -- Ships category
AND t.published = 1;

-- =============================================================================
-- v_route_security_analysis: Route security analysis for hauling
-- Provides security classification and risk indicators for all systems
-- =============================================================================
CREATE VIEW IF NOT EXISTS v_route_security_analysis AS
SELECT
    sys._key as system_id,
    COALESCE(json_extract(sys.name, '$.en'), json_extract(sys.name, '$.de')) as system_name,
    sys.securityStatus as security_status,
    CASE 
        WHEN sys.securityStatus >= 0.5 THEN 'High-Sec'
        WHEN sys.securityStatus > 0.0 THEN 'Low-Sec'
        ELSE 'Null-Sec'
    END as security_class,
    -- Chokepoint detection (fewer gates = higher risk)
    (SELECT COUNT(*) FROM mapStargates WHERE solarSystemID = sys._key) as gate_count,
    sys.border as is_border_system,
    sys.corridor as is_corridor_system,
    r.regionID,
    COALESCE(json_extract(r.name, '$.en'), json_extract(r.name, '$.de')) as region_name
FROM mapSolarSystems sys
LEFT JOIN mapRegions r ON sys.regionID = r._key;
