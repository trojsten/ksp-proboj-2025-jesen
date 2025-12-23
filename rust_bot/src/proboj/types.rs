use super::Vec2D;
use serde::{Deserialize, Serialize};
use serde_repr::{Deserialize_repr, Serialize_repr};

#[repr(u8)]
#[derive(Clone, Debug, Serialize_repr, Deserialize_repr)]
pub enum ShipType {
    MotherShip,
    SuckerShip,
    DrillShip,
    TankerShip,
    TruckShip,
    BattleShip,
}

#[derive(Clone, Copy, Debug, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub struct ShipId(pub(super) usize);

#[derive(Clone, Debug, Deserialize)]
pub struct Ship {
    pub id: ShipId,
    pub player: PlayerId,
    pub position: Vec2D,
    #[serde(rename = "vector")]
    pub velocity: Vec2D,
    pub health: i64,
    pub fuel: f64,
    #[serde(rename = "type")]
    pub ship_type: ShipType,
    pub rock: i64,
    pub is_destroyed: bool,
}

#[repr(u8)]
#[derive(Clone, Debug, Deserialize_repr, PartialEq, Eq)]
pub enum AsteroidType {
    RockAsteroid,
    FuelAsteroid,
}

#[derive(Clone, Copy, Debug, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub struct AsteroidId(pub(super) usize);

#[derive(Clone, Debug, Deserialize)]
pub struct Asteroid {
    pub id: AsteroidId,
    pub position: Vec2D,
    #[serde(rename = "type")]
    pub asteroid_type: AsteroidType,
    pub size: f64,
    pub owner_id: i64,
    pub surface: f64,
}

#[derive(Clone, Copy, Debug, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub struct WormholeId(pub(super) usize);

#[derive(Clone, Debug, Deserialize)]
pub struct Wormhole {
    pub id: WormholeId,
    pub target_id: WormholeId,
    pub position: Vec2D,
}

#[derive(Clone, Copy, Debug, Serialize, Deserialize, PartialEq, Eq, Hash)]
pub struct PlayerId(pub(super) usize);

#[derive(Clone, Debug, Deserialize)]
pub struct Player {
    pub id: PlayerId,
    pub name: String,
    pub color: String,
    pub mothership: Ship,
    pub alive: bool,
    pub score: i64,
}

#[derive(Clone, Debug, Deserialize)]
pub(super) struct GameMap {
    pub radius: f64,
    pub ships: Vec<Option<Ship>>,
    pub asteroids: Vec<Option<Asteroid>>,
    pub wormholes: Vec<Option<Wormhole>>,
    pub players: Vec<Option<Player>>,
    pub round: i64,
}

#[derive(Clone, Debug, Serialize)]
pub struct BuyTurn {
    #[serde(rename = "type")]
    pub ship_type: ShipType,
}

#[derive(Clone, Debug, Serialize)]
pub struct MoveTurn {
    pub ship_id: ShipId,
    #[serde(rename = "vector")]
    pub acceleration: Vec2D,
}

#[derive(Clone, Debug, Serialize)]
pub struct LoadTurn {
    pub source_id: ShipId,
    pub destination_id: ShipId,
    pub amount: i64,
}

#[derive(Clone, Debug, Serialize)]
pub struct SiphonTurn {
    pub source_id: ShipId,
    pub destination_id: ShipId,
    pub amount: i64,
}

#[derive(Clone, Debug, Serialize)]
pub struct ShootTurn {
    pub source_id: ShipId,
    pub destination_id: ShipId,
}

#[derive(Clone, Debug, Serialize)]
pub struct RepairTurn {
    pub ship_id: ShipId,
}

#[derive(Clone, Debug)]
pub enum Turn {
    BuyTurn(BuyTurn),
    MoveTurn(MoveTurn),
    LoadTurn(LoadTurn),
    SiphonTurn(SiphonTurn),
    ShootTurn(ShootTurn),
    RepairTurn(RepairTurn),
}

impl Turn {
    pub fn buy_turn(ship_type: ShipType) -> Turn {
        Turn::BuyTurn(BuyTurn { ship_type })
    }

    pub fn move_turn(ship_id: ShipId, acceleration: Vec2D) -> Turn {
        Turn::MoveTurn(MoveTurn {
            ship_id,
            acceleration,
        })
    }

    pub fn load_turn(source_id: ShipId, destination_id: ShipId, amount: i64) -> Turn {
        Turn::LoadTurn(LoadTurn {
            source_id,
            destination_id,
            amount,
        })
    }

    pub fn siphon_turn(source_id: ShipId, destination_id: ShipId, amount: i64) -> Turn {
        Turn::SiphonTurn(SiphonTurn {
            source_id,
            destination_id,
            amount,
        })
    }

    pub fn shoot_turn(source_id: ShipId, destination_id: ShipId) -> Turn {
        Turn::ShootTurn(ShootTurn {
            source_id,
            destination_id,
        })
    }

    pub fn repair_turn(ship_id: ShipId) -> Turn {
        Turn::RepairTurn(RepairTurn { ship_id })
    }
}
