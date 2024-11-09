#!/usr/bin/env python

import datetime
import logging
import random
from collections import defaultdict
from dataclasses import asdict, dataclass
from enum import Enum
from typing import Any, Callable, Protocol

from evalica import Winner, elo, newman
from google.cloud import firestore
from google.cloud.firestore_v1 import FieldFilter

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class NoRatedChoices(Exception): ...


# TODO(rbtz): use protobuf enum.
class WinnerEnum(Enum):
    UNSPECIFIED = 0
    NONE = 1
    LEFT = 2
    RIGHT = 3
    BOTH = 4


@dataclass(frozen=True, repr=True)
class LeaderboardEntry:
    model: str
    votes: int
    elo_score: float
    elo_ci_lower: float
    elo_ci_upper: float
    newman_score: float
    newman_ci_lower: float
    newman_ci_upper: float


class RatingResultProtocol(Protocol):
    @property
    def scores(self) -> dict[str, float]: ...


RatingSystemProtocol = Callable[[list[str], list[str], list[Winner]], RatingResultProtocol]


@dataclass
class ConfidenceInterval:
    lower: float
    upper: float


def bootstrap_confidence_intervals(
    xs: list[str],
    ys: list[str],
    outcomes: list[Winner],
    rating_systems: dict[str, RatingSystemProtocol],
    n_bootstrap: int = 1000,
    confidence: float = 0.95,
) -> dict[str, dict[str, tuple[float, float]]]:
    """
    Perform bootstrap sampling of (xs, ys, outcomes) for multiple rating systems to estimate confidence intervals.
    Returns a dictionary keyed by rating system name, whose values are dictionaries keyed by model with
    tuples for (lower, upper) confidence intervals.
    """
    n = len(xs)
    if n != len(ys) or n != len(outcomes):
        raise ValueError("All lists (xs, ys, outcomes) must have the same length.")
    if not rating_systems:
        raise ValueError("At least one rating system must be provided.")

    distributions = {system_name: defaultdict(list) for system_name in rating_systems}

    for _ in range(n_bootstrap):
        indices = [random.randint(0, n - 1) for _ in range(n)]
        sample_xs = [xs[i] for i in indices]
        sample_ys = [ys[i] for i in indices]
        sample_outcomes = [outcomes[i] for i in indices]

        for system_name, rating_system in rating_systems.items():
            result = rating_system(sample_xs, sample_ys, sample_outcomes)
            for model_name, score in result.scores.items():
                distributions[system_name][model_name].append(score)

    confidence_intervals = {}
    for system_name, system_distributions in distributions.items():
        system_conf_intervals = {}
        for model_name, dist in system_distributions.items():
            dist_sorted = sorted(dist)
            dist_length = len(dist_sorted)
            if dist_length == 0:
                system_conf_intervals[model_name] = (0.0, 0.0)
                continue

            lower_index = int((1.0 - confidence) / 2 * dist_length)
            upper_index = int((1.0 + confidence) / 2 * dist_length) - 1

            if lower_index < 0:
                lower_index = 0
            if upper_index < 0:
                upper_index = 0
            if lower_index >= dist_length:
                lower_index = dist_length - 1
            if upper_index >= dist_length:
                upper_index = dist_length - 1

            system_conf_intervals[model_name] = (
                dist_sorted[lower_index],
                dist_sorted[upper_index],
            )
        confidence_intervals[system_name] = system_conf_intervals

    return confidence_intervals


def run_once(firestore_client: firestore.Client) -> None:
    choices_ref = firestore_client.collection("choices")
    choices_docs = choices_ref.stream()
    choices: list[dict[str, Any]] = []

    doc: firestore.DocumentSnapshot
    total_choices = 0
    for doc in choices_docs:
        total_choices += 1
        choice_data = doc.to_dict()
        if choice_data is None:
            continue
        winner = choice_data.get("winner", WinnerEnum.UNSPECIFIED.value)
        if winner is None or winner not in {
            WinnerEnum.LEFT.value,
            WinnerEnum.RIGHT.value,
            WinnerEnum.BOTH.value,
            WinnerEnum.NONE.value,
        }:
            continue
        choices.append(choice_data)

    if not choices:
        raise NoRatedChoices(f"No rated choices found. Total choices: {total_choices}")

    joke_map: dict[str, Any] = {}
    jokes_docs = firestore_client.collection("jokes").stream()
    for doc in jokes_docs:
        jokes_dict = doc.to_dict()
        if jokes_dict is None:
            continue
        joke_map[doc.id] = jokes_dict

    model_votes: defaultdict[str, int] = defaultdict(int)
    for choice in choices:
        for joke_id in (choice.get("left_joke_id"), choice.get("right_joke_id")):
            if joke_id is None:
                continue
            joke = joke_map.get(joke_id)
            if joke is None:
                continue
            model = joke.get("model")
            if model is None:
                continue
            model_votes[model] += 1

    xs: list[str] = []
    ys: list[str] = []
    outcomes: list[Winner] = []

    skip_count = 0
    for choice in choices:
        left_joke_id = choice.get("left_joke_id")
        right_joke_id = choice.get("right_joke_id")
        w: WinnerEnum = WinnerEnum(choice.get("winner", 0))
        if left_joke_id is None or right_joke_id is None or w is None:
            logger.debug("Skipping invalid choice: %s", choice)
            skip_count += 1
            continue
        if w not in {WinnerEnum.LEFT, WinnerEnum.RIGHT, WinnerEnum.BOTH, WinnerEnum.NONE}:
            if w == WinnerEnum.UNSPECIFIED:
                continue
            logger.debug("Skipping invalid winner: %d", w.value)
            skip_count += 1
            continue

        left_joke = joke_map.get(left_joke_id)
        right_joke = joke_map.get(right_joke_id)

        if left_joke is None or right_joke is None:
            logger.debug("Skipping invalid joke: %s, %s", left_joke_id, right_joke_id)
            skip_count += 1
            continue

        left_model = left_joke.get("model")
        right_model = right_joke.get("model")

        if left_model is None or right_model is None:
            logger.debug("Skipping invalid model: %s, %s", left_model, right_model)
            skip_count += 1
            continue
        if left_model == right_model:
            logger.debug("Skipping same model: %s", left_model)
            continue

        xs.append(left_model)
        ys.append(right_model)
        match w:
            case WinnerEnum.LEFT:
                outcomes.append(Winner.X)
            case WinnerEnum.RIGHT:
                outcomes.append(Winner.Y)
            case WinnerEnum.BOTH | WinnerEnum.NONE:
                outcomes.append(Winner.Draw)
    logger.info(
        f"Choices processed successfully: {len(xs)=}, {len(ys)=}, {len(outcomes)=}, {skip_count=}"
    )
    if not xs or not ys or not outcomes:
        raise NoRatedChoices("No valid comparisons found.")

    rating_systems_dict = {
        "elo": elo,
        "newman": lambda x_list, y_list, out_list: newman(
            x_list, y_list, out_list, tolerance=1e-6, limit=1000
        ),
    }

    n_bootstrap_samples = 1000
    confidence_level = 0.95
    confidence_intervals = bootstrap_confidence_intervals(
        xs,
        ys,
        outcomes,
        rating_systems=rating_systems_dict,
        n_bootstrap=n_bootstrap_samples,
        confidence=confidence_level,
    )
    logger.info(f"Bootstrap confidence intervals computed successfully: {confidence_intervals=}")

    elo_result = elo(xs, ys, outcomes)
    newman_result = newman(xs, ys, outcomes, tolerance=1e-6, limit=1000)

    leaderboard: list[LeaderboardEntry] = []
    for model_name in model_votes.keys():
        votes_count = model_votes.get(model_name, 0)
        elo_val = elo_result.scores.get(model_name, 0.0)
        elo_ci = confidence_intervals.get("elo", {}).get(model_name, (elo_val, elo_val))
        newman_val = newman_result.scores.get(model_name, 0.0)
        newman_ci = confidence_intervals.get("newman", {}).get(model_name, (newman_val, newman_val))
        logger.info(
            f"Leaderboard entry: {model_name=}, {votes_count=}, {elo_val=}, {elo_ci=}, {newman_val=}, {newman_ci=}"
        )

        elo_ci_lower_diff = float(elo_val - elo_ci[0])  # difference to lower CI boundary
        elo_ci_upper_diff = float(elo_ci[1] - elo_val)  # difference to upper CI boundary
        newman_ci_lower_diff = float(newman_val - newman_ci[0])
        newman_ci_upper_diff = float(newman_ci[1] - newman_val)

        leaderboard.append(
            LeaderboardEntry(
                model=model_name,
                votes=votes_count,
                elo_score=float(elo_val),
                elo_ci_lower=elo_ci_lower_diff,
                elo_ci_upper=elo_ci_upper_diff,
                newman_score=float(newman_val),
                newman_ci_lower=newman_ci_lower_diff,
                newman_ci_upper=newman_ci_upper_diff,
            )
        )
    logger.info(f"Leaderboard computed successfully: {leaderboard=}")

    leaderboard_doc = {
        "leaderboard": [asdict(entry) for entry in leaderboard],
        "created_at": firestore.SERVER_TIMESTAMP,
    }
    leaderboard_ref = firestore_client.collection("leaderboard").document()
    leaderboard_ref.set(leaderboard_doc)
    logger.info("Leaderboard saved successfully")

    # remove expired choices
    time_threshold = datetime.datetime.now(tz=datetime.timezone.utc) - datetime.timedelta(hours=1)
    non_voted = (
        firestore_client.collection("choices")
        .where(filter=FieldFilter("winner", "==", WinnerEnum.UNSPECIFIED.value))
        .stream()
    )
    with firestore_client.batch() as batch:
        for doc in non_voted:
            if doc.get("created_at") > time_threshold:
                continue
            batch.delete(doc.reference)


run_once(
    firestore_client=firestore.Client("humor-arena"),
)
