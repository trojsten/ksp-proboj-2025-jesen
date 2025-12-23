pub const RADIUS: f64 = 15000.0; // Game map radius

pub const MAX_ASTEROID_SIZE: f64 = 50.0; // Maximum size of generated asteroids
pub const MIN_ASTEROID_SIZE: f64 = MAX_ASTEROID_SIZE / 7.0 * 5.0; // Minimum size of generated asteroids

pub const ASTEROID_COUNT: f64 = 500.0; // Number of generated asteroids in the game
pub const WORMHOLE_COUNT: f64 = 25.0; // Number of generated wormhole pairs in the game

pub const WORMHOLE_RADIUS: f64 = 5.0; // Radius within which ships get teleported by wormholes
pub const WORMHOLE_TELEPORT_DISTANCE: f64 = WORMHOLE_RADIUS * 2.0; // Minimum distance from target wormhole (2x radius) to prevent teleport loops

pub const SHIP_MINING_DISTANCE: f64 = MAX_ASTEROID_SIZE; // Maximum distance for mining operations
pub const SHIP_MINING_AMOUNT: f64 = 10.0; // Units mined per tick
pub const SHIP_CONQUERING_DISTANCE: f64 = MAX_ASTEROID_SIZE; // Maximum distance for conquering operations
pub const SHIP_CONQUERING_RATE: f64 = 10.0; // Rate of conquering

pub const SHIP_MAX_HEALTH: f64 = 100.0; // Maximum health points for ships
pub const SHIP_START_FUEL: f64 = 100.0; // Starting fuel for new ships
pub const SHIP_MOVEMENT_FREE_SIZE: f64 = 1.0; // Movement delta per turn that is free of fuel cost
pub const SHIP_MOVEMENT_MAX_SIZE: f64 = 10000.0; // Maximum movement delta per turn - larger movements are scaled down
pub const SHIP_TRANSFER_DISTANCE: f64 = 20.0; // Maximum distance for resource transfer between ships
pub const SHIP_SHOOT_DISTANCE: f64 = 500.0; // Maximum shooting range for ships
pub const SHIP_SHOOT_DAMAGE: f64 = 25.0; // Damage dealt by ship weapons
pub const SHIP_REPAIR_DISTANCE: f64 = 50.0; // Maximum distance for ship repair operations
pub const SHIP_REPAIR_AMOUNT: f64 = 30.0; // Health points restored by repair
pub const SHIP_REPAIR_ROCK_COST: f64 = 15.0; // Rock cost per repair operation
