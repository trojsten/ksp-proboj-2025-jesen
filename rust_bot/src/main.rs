use rust_bot::proboj::*;

fn main() {
    loop {
        let state = get_state();
        let myself = &state.players[&state.my_id];

        let mut turns = vec![Turn::buy_turn(ShipType::DrillShip)];

        let my_ships = state
            .ships
            .values()
            .filter(|ship| ship.player == state.my_id && ship.id != myself.mothership.id);

        for ship in my_ships {
            let target = Vec2D { x: 100.0, y: 100.0 };

            let acceleration = target - ship.position;

            turns.push(Turn::move_turn(
                ship.id,
                acceleration.clamp_magnitude(SHIP_MOVEMENT_FREE_SIZE),
            ));
        }

        send_turns(turns);
    }
}
