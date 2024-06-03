import more_itertools
from openai import AsyncOpenAI


class BatchModeration:
    def __init__(self, client: AsyncOpenAI, model: str = "text-moderation-latest"):
        self.client = client
        self.model = model

    async def classify(self, jokes: list[str]) -> list[tuple[bool, str]]:
        results = []
        for chunk in more_itertools.chunked(jokes, 32):
            moderation = await self.client.moderations.create(
                model=self.model,
                input=list(chunk),
            )
            for joke, result in zip(chunk, moderation.results):
                flagged = result.flagged
                if result.category_scores.harassment > 0.02:
                    flagged = True
                results.append((flagged, joke))
        return results
