import argparse
import datetime
import json
import logging
import os
import re
from concurrent.futures import ThreadPoolExecutor, as_completed
from multiprocessing import cpu_count
from typing import Any

from litellm import completion
from tqdm import tqdm

SYSTEM = "system"
USER = "user"

ROLE = "role"
CONTENT = "content"

logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO)


def now() -> str:
    return datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")


def magic_query(
    prompt: list[dict[str, str]],
    *,
    model: str,
    temp: float = 0.01,
    max_tokens: int = 500,
) -> str:
    request = {
        "model": model,
        "temperature": temp,
        "max_tokens": max_tokens,
        "messages": prompt,
    }
    response = completion(**request)
    response_text = response["choices"][0]["message"]["content"]
    return response_text


def get_associations(theme: str, model: str) -> str:
    p: list[dict[str, str]] = []
    p.append(
        {
            ROLE: SYSTEM,
            CONTENT: """Given the theme, provide a numbered list of 20 free associations with it: different meanings, different contexts, stereotypes, puns, exaggerations, incongruity, cultural references, juxtaposition. Use short descriptions -- 1-5 words each.""",
        }
    )
    p.append({ROLE: USER, CONTENT: f"Theme: {theme}"})

    associations = magic_query(p, model=model, max_tokens=4000)
    return associations


def expand_associations(theme: str, associations: str, model: str) -> str:
    p: list[dict[str, str]] = []
    p.append(
        {
            ROLE: SYSTEM,
            CONTENT: """Read the theme and the list of associations with it.
Rewrite and expand each of the associations into longer and more detailed (10-15 words each). 
Avoid repetitions and usage the similar words across the sentences. Keep mentioning the original theme.
Provide a numbered list as a result.""",
        }
    )
    p.append({ROLE: USER, CONTENT: f"Theme: {theme}\n\nAssociations:\n{associations}"})
    res = magic_query(p, model=model, max_tokens=4000)
    return res


def refine_associations(theme: str, associations: str, model: str) -> str:
    p: list[dict[str, str]] = []
    p.append(
        {
            ROLE: SYSTEM,
            CONTENT: """Read the theme and analyze and review the list of items associated with it.
Write a shorter list of items using the following rules:
- Complementary items can be grouped into one item that combines the original descriptions.
- Weak, unsuccessful, or poorly related items to the original topic should be deleted.
- The final list should be no longer than 8 items, different items should have different meanings or contexts.
Provide a numbered list as a result.""",
        }
    )
    p.append({ROLE: USER, CONTENT: f"Theme: {theme}\n\nAssociations:\n{associations}"})
    res = magic_query(p, model=model, max_tokens=4000)
    return res


def make_jokes(theme: str, associations: str, model: str) -> str:
    p: list[dict[str, str]] = []
    p.append(
        {
            ROLE: SYSTEM,
            CONTENT: """Read the theme and analyze and review the list of items associated with it.
Based on that information, write a list of 7-10 jokes using the following rules:
- It should be "One-liner" -- a concise, self-contained joke, delivered in 1-2 sentences with a two-part structure where the first part (setup) establishes a scenario and the second part (punchline) delivers an unexpected twist or conclusion that subverts the setup.
- You may want to use one or several of following strategies:
  - Wordplay or pun
  - Misdirection
  - Exaggeration
  - Stereotyping
  - Satire
  - Absurdity or Surreal Humor
  - Dark Humor
  - Juxtaposition

Use these general principles (Anthropomorphism, Visual Imagery, Absurdity, Incongruity, Misdirection).
Provide a numbered list as a result.""",
        }
    )
    p.append({ROLE: USER, CONTENT: f"Theme: {theme}\n\nAssociations:\n{associations}"})
    res = magic_query(p, model=model, max_tokens=4000)
    return res


def make_jokes_ablated(theme: str, model: str) -> str:
    p: list[dict[str, str]] = []
    p.append(
        {
            ROLE: SYSTEM,
            CONTENT: """Read the theme and write a list of 7-10 jokes on that theme.
Provide a numbered list as a result.""",
        }
    )
    p.append({ROLE: USER, CONTENT: f"Theme: {theme}"})
    res = magic_query(p, model=model, max_tokens=4000)
    return res


parser = argparse.ArgumentParser(description="Process themes with humor generation.")
parser.add_argument("--model", type=str, default="openai/gpt-4o-2024-11-20", help="Model to use")
parser.add_argument(
    "--max-workers", type=int, default=cpu_count(), help="Number of workers for thread pool"
)
parser.add_argument(
    "--theme-set", type=str, default="theme_set_v2.txt", help="Path to theme set file"
)
args = parser.parse_args()

mmodel = args.model.split("/")[1]
os.makedirs("output", exist_ok=True)


def process_theme(theme: str, model: str) -> dict[str, Any]:
    output: dict[str, Any] = {}
    output["theme"] = theme
    output["ts0"] = now()
    logger.debug("Timestamp ts0: %s", output["ts0"])

    output["step1"] = get_associations(theme, model)
    logger.debug("Step1 output: %s", output["step1"])

    output["ts1"] = now()
    logger.debug("Timestamp ts1: %s", output["ts1"])

    output["step2"] = expand_associations(theme, output["step1"], model)
    logger.debug("Step2 output: %s", output["step2"])

    output["ts2"] = now()
    logger.debug("Timestamp ts2: %s", output["ts2"])

    output["step3"] = refine_associations(theme, output["step2"], model)
    logger.debug("Step3 output: %s", output["step3"])

    output["ts3"] = now()
    logger.debug("Timestamp ts3: %s", output["ts3"])

    output["step4"] = make_jokes(theme, output["step3"], model)
    logger.debug("Step4 output: %s", output["step4"])

    output["ts4"] = now()
    logger.debug("Timestamp ts4: %s", output["ts4"])

    output["ablated"] = make_jokes_ablated(theme, model)
    logger.debug("Ablated output: %s", output["ablated"])

    output["ts5"] = now()
    logger.debug("Timestamp ts5: %s", output["ts5"])

    outfn = "output/" + "_humorizer_v1_" + mmodel + "_" + re.sub(r"[^a-zA-Z]", "_", theme) + ".json"
    with open(outfn, "w", encoding="utf-8") as ofh:
        json.dump(output, ofh, indent=2)

    return output


with open(args.theme_set, "rt", encoding="utf-8") as fh:
    themes: list[str] = [line.strip() for line in fh if line.strip()]

with ThreadPoolExecutor(max_workers=args.max_workers) as executor:
    futures = {executor.submit(process_theme, theme, args.model): theme for theme in themes}
    for future in tqdm(as_completed(futures), total=len(futures), desc="Processing themes"):
        theme = futures[future]
        try:
            result = future.result()
            logger.debug("Completed processing theme: %s", theme)
        except Exception as e:
            logger.error("Error processing theme '%s': %s", theme, e)
