import asyncio
import os

from convex import ConvexClient, ConvexInt64
from dotenv import load_dotenv
from openai import AsyncOpenAI

from main import JokeBuilder
from moderation import BatchModeration
from models import Joke, Policy




"""
Humor Mechanics: Advancing Humor Generation with Multistep Reasoning

https://arxiv.org/html/2405.07280v1
"""

DEFAULT_MODEL = "gpt-4o"
POLICY1 = """
- Wordplay or pun: The use of words with multiple meanings or similar sounds to create humorous ambiguity or surprise.
- Misdirection: Leading the audience to expect a certain outcome or narrative, only to reveal a different, often contradictory or absurd, conclusion.
- Exaggeration: Amplifying a characteristic, situation, or behavior to absurd levels to highlight its comedic potential.
- Stereotyping: Utilizing exaggerated and oversimplified characterizations of groups or individuals for comedic effect.
- Satire: Using humor, irony, or exaggeration to critique or mock people, institutions, societal norms, or other targets.
- Absurdity or Surreal Humor: Creating humor through scenarios or statements that are illogical, bizarre, or defy common sense.
- Dark Humor: Making light of subjects that are generally considered serious, taboo, or morbid.
- Juxtaposition: Placing two contrasting ideas or scenarios side by side for comedic effect, highlighting their differences.
"""


async def main(
    *,
    db: ConvexClient,
):
    default_policy = Policy(
        name="default",
        text=POLICY1,
        model=DEFAULT_MODEL,
    )

    async with AsyncOpenAI(api_key=os.environ.get("OPENAI_API_KEY")) as client:
        builder = JokeBuilder(client)
        moderator = BatchModeration(client)

        model = db.mutation("mutations/models:upsertModel", {"name": DEFAULT_MODEL, "temperature": 1.0})
        policy = db.mutation("mutations/policies:upsertPolicy", {
            "name": default_policy.name,
            "text": default_policy.text,
            "model": model,
        })
        with open("data/english-nouns.txt") as lines:
            for noun in lines:
                noun = noun.strip()

                topic = db.mutation("mutations/topics:upsertTopic", {"name": noun})
                try:
                    jokes: list[Joke] = await builder.jokes(default_policy.text, noun)
                    joke0 = jokes[0]
                except Exception as e:
                    print(f"Error generating jokes for noun: {noun}: {e}")
                    continue

                assocv1 = db.mutation(
                    "mutations/assocs:upsertAssocv1",
                    {
                        "text": joke0.assocv1,
                        "topic": topic,
                        "model": model,
                    },
                )
                assocv2 = db.mutation(
                    "mutations/assocs:upsertAssocv2",
                    {
                        "text": joke0.assocv2,
                        "topic": topic,
                        "model": model,
                    },
                )
                assocv3 = db.mutation(
                    "mutations/assocs:upsertAssocv3",
                    {
                        "text": joke0.assocv3,
                        "topic": topic,
                        "model": model,
                    },
                )

                for i, (flagged, _) in enumerate(await moderator.classify(
                    [j.joke for j in jokes]
                )):
                    if flagged:
                        print(f"Flagged joke {i}: {jokes[i].joke}")
                        continue
                    db.mutation(
                        "mutations/jokes:insertJoke",
                        {
                            "policy": policy,
                            "topic": topic,
                            "model": model,
                            "assocv1": assocv1,
                            "assocv2": assocv2,
                            "assocv3": assocv3,
                            "text": jokes[i].joke,
                            "score": ConvexInt64(0),
                            "views": ConvexInt64(0),
                        },
                    )

if __name__ == "__main__":
    load_dotenv(".env.local")
    load_dotenv()

    # XXX async
    db = ConvexClient(os.getenv("CONVEX_URL"))
    asyncio.run(main(db=db))
