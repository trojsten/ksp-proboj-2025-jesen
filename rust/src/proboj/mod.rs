#![allow(clippy::print_stdout)]

mod consts;
mod types;
mod vec2d;

use serde::Deserialize;
use std::collections::HashMap;
use std::io::{Write, stdout};

pub use consts::*;
pub use types::*;
pub use vec2d::Vec2D;

pub fn send_turns(turns: Vec<Turn>) {
    let turns = turns
        .into_iter()
        .map(|turn| match turn {
            Turn::BuyTurn(t) => serde_json::json!({"type": 0,  "data": t}),
            Turn::MoveTurn(t) => serde_json::json!({"type": 1, "data": t}),
            Turn::LoadTurn(t) => serde_json::json!({"type": 2, "data": t}),
            Turn::SiphonTurn(t) => serde_json::json!({"type": 3, "data": t}),
            Turn::ShootTurn(t) => serde_json::json!({"type": 4, "data": t}),
            Turn::RepairTurn(t) => serde_json::json!({"type": 5, "data": t}),
        })
        .collect::<Vec<_>>();

    let json = serde_json::to_string(&turns).unwrap();
    println!("{}", json);
    println!(".");
    stdout().flush().unwrap();
}

pub struct GameState {
    pub radius: f64,
    pub ships: HashMap<ShipId, Ship>,
    pub asteroids: HashMap<AsteroidId, Asteroid>,
    pub wormholes: HashMap<WormholeId, Wormhole>,
    pub players: HashMap<PlayerId, Player>,
    pub round: i64,
    pub my_id: PlayerId,
}

pub fn get_state() -> GameState {
    let mut input = String::new();
    let mut dot = String::new();

    std::io::stdin().read_line(&mut input).unwrap();
    std::io::stdin().read_line(&mut dot).unwrap();
    assert_eq!(dot.trim(), ".");

    #[derive(Deserialize)]
    struct StateMessage {
        map: GameMap,
        player_id: PlayerId,
    }

    let StateMessage { map, player_id } = serde_json::from_str(&input).unwrap();

    GameState {
        radius: map.radius,
        ships: map
            .ships
            .into_iter()
            .enumerate()
            .filter_map(|(i, ship)| ship.map(|s| (ShipId(i), s)))
            .collect(),
        asteroids: map
            .asteroids
            .into_iter()
            .enumerate()
            .filter_map(|(i, asteroid)| asteroid.map(|a| (AsteroidId(i), a)))
            .collect(),
        wormholes: map
            .wormholes
            .into_iter()
            .enumerate()
            .filter_map(|(i, wormhole)| wormhole.map(|w| (WormholeId(i), w)))
            .collect(),
        players: map
            .players
            .into_iter()
            .enumerate()
            .filter_map(|(i, player)| player.map(|p| (PlayerId(i), p)))
            .collect(),
        round: map.round,
        my_id: player_id,
    }
}
