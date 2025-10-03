from typing import List
from .proboj import (
    Client,
    BuyTurn,
    MoveTurn,
    LoadTurn,
    SiphonTurn,
    ShootTurn,
    RepairTurn,
    ShipType,
    Position,
    Turn,
)


class MyClient(Client):
    def turn(self) -> List[Turn]:
        my_ships = self.get_my_ships()
        if not my_ships:
            return []

        turns: List[Turn] = []

        # Example: Move first ship if it exists
        if my_ships:
            first_ship = my_ships[0]
            # Move the ship by a small amount
            turns.append(MoveTurn(first_ship.id, Position(10, 5)))

        # Example: Buy a battleship if we have enough resources
        my_player = self.get_my_player()
        if my_player and my_player.rock >= 100:
            turns.append(BuyTurn(ShipType.BATTLE_SHIP))

        # Example: Shoot at enemy ships if we have battleships
        for ship in my_ships:
            if ship.type == ShipType.BATTLE_SHIP:
                # Find enemy ships
                enemy_ships = []
                if self.game_map:
                    for other_ship in self.game_map.ships:
                        if other_ship and other_ship.player_id != self.my_player_id:
                            enemy_ships.append(other_ship)

                # Shoot at nearest enemy if in range
                if enemy_ships:
                    nearest_enemy = min(
                        enemy_ships, key=lambda e: ship.position.distance(e.position)
                    )
                    if ship.position.distance(nearest_enemy.position) <= 100:
                        turns.append(ShootTurn(ship.id, nearest_enemy.id))

        return turns


if __name__ == "__main__":
    client = MyClient()
    client.run()

