from dataclasses import dataclass

@dataclass(frozen=True, order=True)
class Joke:
    topic: str
    assocv1: str
    assocv2: str
    assocv3: str
    joke: str


@dataclass(frozen=True)
class Policy:
    name: str
    text: str
    model: str