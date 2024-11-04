import datetime
import logging
import os
import random
import time
from collections import defaultdict
from dataclasses import asdict, dataclass
from enum import Enum
from typing import Any

from evalica import Winner, elo, newman
from google.cloud import firestore

logging.basicConfig(level=logging.DEBUG)
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


def bootstrap_confidence_intervals(
    xs: list[str],
    ys: list[str],
    outcomes: list[Winner],
    n_bootstrap: int = 1000,
    confidence: float = 0.95,
) -> tuple[dict, dict]:
    """
    Perform bootstrap sampling of (xs, ys, outcomes) to estimate confidence intervals.
    Returns dictionaries for ELO and Newman, keyed by model, with tuples for (lower, upper) intervals.
    """
    n = len(xs)
    if n != len(ys) or n != len(outcomes):
        raise ValueError("All lists (xs, ys, outcomes) must have the same length.")

    elo_distributions = defaultdict(list)
    newman_distributions = defaultdict(list)

    for _ in range(n_bootstrap):
        indices = [random.randint(0, n - 1) for _ in range(n)]
        sample_xs = [xs[i] for i in indices]
        sample_ys = [ys[i] for i in indices]
        sample_outcomes = [outcomes[i] for i in indices]

        elo_result_sample = elo(sample_xs, sample_ys, sample_outcomes)
        newman_result_sample = newman(
            sample_xs, sample_ys, sample_outcomes, tolerance=1e-6, limit=1000
        )

        for model_name, score in elo_result_sample.scores.items():
            elo_distributions[model_name].append(score)
        for model_name, score in newman_result_sample.scores.items():
            newman_distributions[model_name].append(score)

    elo_confidence_intervals = {}
    newman_confidence_intervals = {}
    for model_name, dist in elo_distributions.items():
        dist_sorted = sorted(dist)
        lower_index = int((1.0 - confidence) / 2 * n_bootstrap)
        upper_index = int((1.0 + confidence) / 2 * n_bootstrap) - 1
        elo_confidence_intervals[model_name] = (dist_sorted[lower_index], dist_sorted[upper_index])
    for model_name, dist in newman_distributions.items():
        dist_sorted = sorted(dist)
        lower_index = int((1.0 - confidence) / 2 * n_bootstrap)
        upper_index = int((1.0 + confidence) / 2 * n_bootstrap) - 1
        newman_confidence_intervals[model_name] = (
            dist_sorted[lower_index],
            dist_sorted[upper_index],
        )
    return elo_confidence_intervals, newman_confidence_intervals


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
        winner = choice_data.get("winner")
        if winner is None:
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
        joke_ids: list[str | None] = []
        match WinnerEnum(choice.get("winner")):
            case WinnerEnum.LEFT:
                joke_ids = [choice.get("left_joke_id")]
            case WinnerEnum.RIGHT:
                joke_ids = [choice.get("right_joke_id")]
            case WinnerEnum.BOTH | WinnerEnum.NONE:
                joke_ids = [choice.get("left_joke_id"), choice.get("right_joke_id")]
            case _:
                continue

        for joke_id in joke_ids:
            if joke_id is None:
                continue
            model = joke_map[joke_id].get("model")
            if model is None:
                continue
            model_votes[model] += 1

    leaderboard: list[LeaderboardEntry] = []

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
            logger.debug("Skipping invalid winner: %d", w)
            skip_count += 1
            continue

        left_joke = joke_map.get(left_joke_id)
        right_joke = joke_map.get(right_joke_id)

        if left_joke is None or right_joke is None:
            logger.debug("Skipping invalid joke: %s", left_joke_id, right_joke_id)
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

        match w:
            case WinnerEnum.LEFT:
                xs.append(left_model)
                ys.append(right_model)
                outcomes.append(Winner.X)
            case WinnerEnum.RIGHT:
                xs.append(left_model)
                ys.append(right_model)
                outcomes.append(Winner.Y)
            case WinnerEnum.BOTH | WinnerEnum.NONE:
                xs.append(left_model)
                ys.append(right_model)
                outcomes.append(Winner.Draw)

    if not xs or not ys or not outcomes:
        raise NoRatedChoices("No valid comparisons found.")

    elo_result = elo(xs, ys, outcomes)
    newman_result = newman(xs, ys, outcomes, tolerance=1e-6, limit=1000)

    n_bootstrap_samples = 1000
    confidence_level = 0.95
    elo_cis, newman_cis = bootstrap_confidence_intervals(
        xs, ys, outcomes, n_bootstrap=n_bootstrap_samples, confidence=confidence_level
    )
    logger.info(f"Bootstrap confidence intervals computed successfully: {elo_cis=}, {newman_cis=}")

    all_models = (
        set(model_votes.keys()) | set(elo_result.scores.index) | set(newman_result.scores.index)
    )
    for model_name in all_models:
        votes_count = model_votes.get(model_name, 0)
        elo_val = elo_result.scores.get(model_name, 0.0)
        elo_ci = elo_cis.get(model_name, (elo_val, elo_val))
        newman_val = newman_result.scores.get(model_name, 0.0)
        newman_ci = newman_cis.get(model_name, (newman_val, newman_val))
        logger.info(
            f"Leaderboard entry: {model_name=}, {votes_count=}, {elo_val=}, {elo_ci=}, {newman_val=}, {newman_ci=}"
        )
        leaderboard.append(
            LeaderboardEntry(
                model=model_name,
                votes=votes_count,
                elo_score=float(elo_val),
                elo_ci_lower=float(elo_val - elo_ci[0]),
                elo_ci_upper=float(elo_ci[1] - elo_val),
                newman_score=float(newman_val),
                newman_ci_lower=float(newman_val - newman_ci[0]),
                newman_ci_upper=float(newman_ci[1] - newman_val),
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
    expired_docs = (
        firestore_client.collection("choices")
        .where(filter=firestore.FieldFilter("winner", "==", value=None))
        .where(filter=firestore.FieldFilter("created_at", "<", time_threshold))
        .stream()
    )
    with firestore_client.batch() as batch:
        for doc in expired_docs:
            batch.delete(doc.reference)


while True:
    try:
        logger.info("Running leaderboard update")
        run_once(
            firestore_client=firestore.Client("humor-arena"),
        )
    except Exception as e:
        logger.error(f"Error: {e}")
    finally:
        time.sleep(600)
