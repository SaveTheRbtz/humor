import asyncio
import csv
import glob
import json
import logging
import os
import random as insecure_random
import re
from dataclasses import asdict
from typing import Iterable

import aiofiles
from openai import AsyncOpenAI
from tqdm.asyncio import tqdm_asyncio

from concurrency import limited_concurrency
from main import JokeBuilder
from models import Joke

MIN_JOKE_LENGTH = 20
MAX_JOKE_LENGTH = 256

OUTPUT_PATH = "output/processed.jsonl"

logging.basicConfig(level=logging.WARN)


class FineTuningDataset:
    def __init__(self, client: AsyncOpenAI, model: str = "gpt-4o"):
        self.client = client
        self.model = model

    async def _get_joke_topic(self, joke: str) -> str:
        completion = await self.client.chat.completions.create(
            messages=[
                {
                    "role": "system",
                    "content": """You are a joke classification model. Read the joke and determine the topic of the joke. Provide a short one-two word topic of the joke (not too specific, not too generic).""",
                },
                {
                    "role": "user",
                    "content": joke,
                },
            ],
            model=self.model,
        )
        return completion.choices[0].message.content.strip().lower()

    async def reverse_gen(self, joke: str) -> Joke:
        topic = await self._get_joke_topic(joke)
        builder = JokeBuilder(self.client)
        assocv1 = await builder.get_topic_associations(topic)
        assocv2 = await builder.expand_associations(assocv1)
        assocv3 = await builder.refine_associations(assocv2)
        return Joke(
            topic=topic, assocv1=assocv1, assocv2=assocv2, assocv3=assocv3, joke=joke
        )



async def get_cleaned_jokes(client: AsyncOpenAI, path: str) -> list[str]:
    _cache_path = "output/not_flagged.jsonl"
    if os.path.exists(_cache_path):
        async with aiofiles.open(_cache_path) as f:
            return [json.loads(line)["joke"] async for line in f]

    jokes: list[str] = []
    for fn in glob.glob(path):
        with open(fn) as f:
            reader = csv.reader(f)
            next(reader)
            jokes.extend([row[1] for row in reader])

    jokes = filter_jokes(cleanup_jokes(jokes))
    insecure_random.shuffle(jokes)

    moderation = BatchModeration(client)
    results = await moderation.classify(jokes)

    not_flagged = [joke for flagged, joke in results if not flagged]
    async with aiofiles.open(_cache_path, "w") as f:
        for joke in not_flagged:
            await f.write(json.dumps({"joke": joke}) + "\n")
    return not_flagged


def filter_jokes(jokes: list[str]) -> list[str]:
    jokes = [joke for joke in jokes if len(joke) >= MIN_JOKE_LENGTH]
    jokes = [joke for joke in jokes if len(joke) <= MAX_JOKE_LENGTH]
    jokes = [joke for joke in jokes if "read more" not in joke.lower()]
    return jokes


def cleanup_jokes(jokes: Iterable[str]) -> list[str]:
    jokes = [joke.strip() for joke in jokes]
    jokes = [joke[3:] if joke.startswith("Q: ") else joke for joke in jokes]
    jokes = [joke.replace(" A: ", "") for joke in jokes]
    jokes = [re.sub(r"\.+", ".", joke) for joke in jokes]
    jokes = [re.sub(r"\s+", " ", joke) for joke in jokes]
    return jokes


async def main():
    os.makedirs("output", exist_ok=True)

    async with AsyncOpenAI(api_key=os.environ.get("OPENAI_API_KEY")) as client:
        not_flagged = await get_cleaned_jokes(
            client, "data/short-jokes-dataset/data/*.csv"
        )

        if os.path.exists(OUTPUT_PATH):
            async with aiofiles.open(OUTPUT_PATH, "r") as f:
                processed = [json.loads(line) async for line in f]
                processed = {j["joke"] for j in processed}
                processed = filter_jokes(cleanup_jokes(processed))
                not_flagged = filter_jokes(cleanup_jokes(not_flagged))
                not_flagged = [joke for joke in not_flagged if joke not in processed]

        @limited_concurrency(limit=10)
        async def process_joke(joke, builder, f):
            try:
                result = await builder.reverse_gen(joke)
                json_record = asdict(result)
                await f.write(json.dumps(json_record) + "\n")
            except Exception as e:
                print(f"Error processing joke: {joke}")
                print(e)
                await asyncio.sleep(1)

        builder = FineTuningDataset(client)
        async with aiofiles.open(OUTPUT_PATH, "a") as f:
            # cleanup again, in case rules have changed.
            logging.warning("Cleaning up jokes: %d", len(not_flagged))
            jokes = cleanup_jokes(not_flagged)
            logging.warning("Filtering jokes: %d", len(jokes))
            jokes = filter_jokes(jokes)
            logging.warning("Processing jokes: %d", len(jokes))

            tasks = []
            for joke in jokes:
                tasks.append(process_joke(joke, builder, f))

            await tqdm_asyncio.gather(*tasks)


if __name__ == "__main__":
    asyncio.run(main())
