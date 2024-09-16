import asyncio
import os

from models import Joke, Policy

from openai import AsyncOpenAI

"""
Humor Mechanics: Advancing Humor Generation with Multistep Reasoning

https://arxiv.org/html/2405.07280v1
"""

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


async def building_blocks(client: AsyncOpenAI, joke: str):
    completion = await client.chat.completions.create(
        messages=[
            {
                "role": "system",
                "content": """You are an expert in humor theory, with extensive knowledge in various forms of comedy ranging from slapstick to satire and dry humor.
Your expertise includes understanding the structural elements of jokes, cultural influences on humor, and the psychology behind what makes things funny.
Additionally, you are a practicing stand-up comedian with years of experience in writing and performing jokes, giving you a unique perspective on how humor resonates with different audiences.
Your task is to analyze the provided joke, identifying its humorous elements, potential audience reception, and categorizing it based on humor types.
Discuss the elements that contribute to its humor, such as the structure of the joke, the play on words, the unexpected twist, or cultural references.
""",
            },
            {
                "role": "user",
                "content": joke,
            },
        ],
        model="gpt-4o",
    )
    return completion.choices[0].message.content.strip()


async def summarize_into_policy(client: AsyncOpenAI, blocks: list[str]):
    completion = await client.chat.completions.create(
        messages=[
            {
                "role": "system",
                "content": """Read several texts about different jokes above. Formulate a list of typical (repeated) elements, used for construction of jokes. Avoid providing semantic doubles, provide only unique elements and their formal abstract descriptions (without details of any particular joke).""",
            },
            {
                "role": "user",
                "content": "\n\n".join(blocks),
            },
        ],
        model="gpt-4o",
    )
    return completion.choices[0].message.content.strip()


async def build_policy(client: AsyncOpenAI, jokes: list[str]):
    blocks = await asyncio.gather(*[building_blocks(client, joke) for joke in jokes])
    policy = await summarize_into_policy(client, blocks)
    print(policy)


class JokeBuilder:
    def __init__(self, client: AsyncOpenAI, model: str = "gpt-4o"):
        self.client = client
        self.model = model

    async def expand_associations(self, associations: str):
        completion = await self.client.chat.completions.create(
            messages=[
                {
                    "role": "system",
                    "content": """Read the theme and the list of associations with it.
Rewrite and expand each of the associations into longer and more detailed (10-15 words each).
Avoid repetitions and usage the similar words across the sentences. Keep mentioning the original theme.
Provide a numbered list as a result.""",
                },
                {
                    "role": "user",
                    "content": associations,
                },
            ],
            model=self.model,
        )
        return completion.choices[0].message.content.strip()

    async def refine_associations(self, associations: str):
        completion = await self.client.chat.completions.create(
            messages=[
                {
                    "role": "system",
                    "content": """Read the theme and analyze and review the list of items associated with it.
Write a shorter list of items using the following rules:
- Complementary items can be grouped into one item that combines the original descriptions.
- Weak, unsuccessful, or poorly related items to the original topic should be deleted.
- The final list should be no longer than 6 items, different items should have different meanings or contexts.
Provide a numbered list as a result.""",
                },
                {
                    "role": "user",
                    "content": associations,
                },
            ],
            model=self.model,
        )
        return completion.choices[0].message.content.strip()

    async def generate_jokes(self, subject: str, policy: str, associations: str) -> str:
        completion = await self.client.chat.completions.create(
            messages=[
                {
                    "role": "system",
                    "content": f"""
Read the theme and analyze and review the list of items associated with it.
Based on that information, write a list of 7-10 jokes using the following rules:
- It should be "One-liner" -- a concise, self-contained joke, delivered in 1-2 sentences with a two-part structure where the first part (setup) establishes a scenario and the second part (punchline) delivers an unexpected twist or conclusion that subverts the setup.
- You may want to use one or several of following strategies:
{policy}
One joke per line, no numbering is needed.""",
                },
                {
                    "role": "user",
                    "content": f"Theme is '{subject}'",
                },
                {
                    "role": "user",
                    "content": associations,
                },
            ],
            model=self.model,
        )
        return completion.choices[0].message.content.strip()

    async def get_topic_associations(self, subject: str):
        completion = await self.client.chat.completions.create(
            messages=[
                {
                    "role": "system",
                    "content": """Given the theme, provide a numbered list of 20 free associations with it: different meanings, different contexts, stereotypes, puns, exaggerations, incongruity, cultural references, juxtaposition. Use short descriptions -- 1-5 words each.""",
                },
                {
                    "role": "user",
                    "content": subject,
                },
            ],
            model=self.model,
        )
        return completion.choices[0].message.content.strip()

    async def jokes(self, policy: str, subject: str) -> list[Joke]:
        assocv1 = await self.get_topic_associations(subject)
        assocv2 = await self.expand_associations(assocv1)
        assocv3 = await self.refine_associations(assocv2)
        joke_str = await self.generate_jokes(subject, policy, assocv3)
        jokes = [line for line in joke_str.split("\n") if len(line.strip()) > 25]
        return [
            Joke(
                topic=subject,
                assocv1=assocv1,
                assocv2=assocv2,
                assocv3=assocv3,
                joke=joke,
            )
            for joke in jokes
        ]


async def main():
    async with AsyncOpenAI(api_key=os.environ.get("OPENAI_API_KEY")) as client:
        builder = JokeBuilder(client)
        with open("data/english-nouns.txt") as lines:
            for noun in lines:
                jokes = await builder.jokes(POLICY1, noun)
                print(jokes)


if __name__ == "__main__":
    asyncio.run(main())
