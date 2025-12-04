from typing import IO
from pettingzoo import ParallelEnv
from pettingzoo.test import parallel_api_test
import subprocess
import numpy as np
from gymnasium import spaces
import functools
import json
import time

N_ENTITIES = 50


class ProbojEnv(ParallelEnv):
    _server: subprocess.Popen | None
    from_server: IO | None
    to_server: IO | None

    def read_input(self):
        print("input", self.round)
        data = self.from_server.readline()
        _ = self.from_server.readline()
        print("got")

        self.last_state = json.loads(data)

    def restart_server(self):
        if self._server:
            self._server.kill()
        if self.from_server:
            self.from_server.close()
        if self.to_server:
            self.to_server.close()

        self._server = subprocess.Popen(["./runner_linux"], cwd="..")
        time.sleep(1)
        self.from_server = open("/tmp/proboj_from", "r")
        self.to_server = open("/tmp/proboj_to", "w")
        self.read_input()

        self.to_server.write(
            #'[{"type": 0, "data": {"type": 1}}, {"type": 0, "data": {"type": 2}}, {"type": 0, "data": {"type": 5}}]\n.\n'
            '[{"type": 0, "data": {"type": 1}}, {"type": 0, "data": {"type": 2}}]\n.\n'
        )
        self.to_server.flush()
        self.read_input()

    metadata = {"render_modes": [], "name": "proboj"}

    def __init__(self, render_mode=None):
        # self.possible_agents = ["mothership", "sucker", "drill", "battle"]
        self.possible_agents = ["mothership", "sucker", "drill"]
        self.render_mode = render_mode
        self.round = 0
        self._server = None
        self.to_server = self.from_server = None
        self.restart_server()

    def render(self):
        pass

    agent_to_shiptype = {
        "mothership": 0,
        "sucker": 1,
        "drill": 2,
        "battle": 5,
    }

    def observe(self, agent):
        ship = self.get_agents_ship(agent)

        entities = []
        for s in self.last_state["map"]["ships"]:
            if not s:
                continue
            if s["player"] != self.last_state["player_id"]:
                type_ = 1
            else:
                if s["type"] == 0:
                    type_ = 2
                if s["type"] == 1:
                    type_ = 3
                if s["type"] == 2:
                    type_ = 4
                if s["type"] == 5:
                    type_ = 5

            entities.append(
                (
                    type_,
                    (s["position"]["x"] - ship["position"]["x"]) / 15000,
                    (s["position"]["y"] - ship["position"]["y"]) / 15000,
                    0,
                    (s["fuel"] + s["rock"]) / 25000,
                )
            )

        for a in self.last_state["map"]["asteroids"]:
            if not a:
                continue
            entities.append(
                (
                    6,
                    (a["position"]["x"] - ship["position"]["x"]) / 15000,
                    (a["position"]["y"] - ship["position"]["y"]) / 15000,
                    a["surface"]
                    / (a["size"] ** 2 * np.pi)
                    * (1 if a["owner_id"] == ship["player"] else -1),
                    (a["size"] ** 2 * np.pi) / 25000,
                )
            )

        entities.sort(key=lambda e: (e[1] ** 2 + e[2] ** 2))
        while len(entities) < N_ENTITIES:
            entities.append((0, 0, 0, 0, 0))

        return {
            "self_state": (
                ship["type"],
                ship["health"] / 100,
                ship["fuel"] / 25000,
                ship["rock"] / 25000,
                ship["position"]["x"] / 15000,
                ship["position"]["y"] / 15000,
                ship["vector"]["x"] / 15000,
                ship["vector"]["y"] / 15000,
            ),
            "closest_entities": entities[:N_ENTITIES],
        }

    def get_agents_ship(self, agent):
        ship = None

        for s in self.last_state["map"]["ships"]:
            if not s:
                continue
            if s["player"] != self.last_state["player_id"]:
                continue
            if s["type"] != self.agent_to_shiptype[agent]:
                continue
            ship = s

        if not ship:
            print(self.last_state["map"]["ships"])
            raise ValueError(f"no ship for {agent}")
        return ship

    def step(self, actions):
        turn = []
        self.round += 1

        for a in self.agents:
            act = actions[a]
            ship = self.get_agents_ship(a)
            if act["action_type"] == 0:
                continue
            elif act["action_type"] == 1:
                turn.append(
                    {
                        "type": 1,
                        "data": {
                            "ship_id": ship["id"],
                            "vector": {
                                "x": float(act["movement_vector"][0]),
                                "y": float(act["movement_vector"][1]),
                            },
                        },
                    }
                )
            elif act["action_type"] >= 2:
                turn.append(
                    {
                        "type": int(act["action_type"]),
                        "data": {
                            "source_id": ship["id"],
                            "destination_id": self.get_agents_ship(
                                {
                                    0: "mothership",
                                    1: "mothership",
                                    2: "sucker",
                                    3: "drill",
                                }[act["target_index"]]
                            )["id"],
                            "amount": 100,
                        },
                    }
                )

        if self.to_server:
            self.to_server.write(json.dumps(turn))
            self.to_server.write("\n.\n")
            self.to_server.flush()
            self.read_input()

        infos = {a: {} for a in self.agents}

        truncations = {a: self.round >= 1990 for a in self.agents}
        terminations = {}
        rewards = {}
        for a in self.agents:
            ship = self.get_agents_ship(a)
            terminations[a] = ship["health"] <= 0 or ship["is_destroyed"]
            if ship["type"] == 0:
                terminations[a] = False
            rewards[a] = (
                ship["health"]
                + ship["rock"]
                + ship["fuel"]
                + self.last_state["map"]["players"][self.last_state["player_id"]][
                    "score"
                ]
            )

        observations = {a: self.observe(a) for a in self.agents}

        if all(terminations.values()) or all(truncations.values()):
            self.agents = []

        return observations, rewards, terminations, truncations, infos

    @functools.lru_cache(maxsize=None)
    def action_space(self, agent):
        return spaces.Dict(
            {
                # 0: Wait/Idle
                # 1: Move
                # 2: Mine Fuel (for Sucker)
                # 3: Mine Rock (for Drill)
                # 4: shoot
                "action_type": spaces.Discrete(4),
                "movement_vector": spaces.Box(
                    low=np.array([-10, -10]),
                    high=np.array([10, 10]),
                    dtype=np.float32,
                ),
                # 3. Discrete Parameter (Used for Mine, Shoot, Siphon, Load commands)
                # Index 0 is always reserved for 'No Target' or can be used for 'Nearest Target'.
                # Index 1 is mothership, 2 sucker, 3 drill, 4 battle
                "target_index": spaces.Discrete(4),
            }
        )

    @functools.lru_cache(maxsize=None)
    def observation_space(self, agent) -> spaces.Space:
        return spaces.Dict(
            {
                "self_state": spaces.Box(low=-1, high=5, shape=(8,), dtype=np.float32),
                "closest_entities": spaces.Box(
                    low=-1.0, high=1.0, shape=(N_ENTITIES, 5), dtype=np.float32
                ),
            }
        )

    def reset(self, seed=None, options=None):
        self.agents = self.possible_agents[:]
        self.round = 0
        self.restart_server()

        return {a: self.observe(a) for a in self.agents}, {a: {} for a in self.agents}


parallel_api_test(ProbojEnv())
