import asyncio
from typing import Any, Callable, Coroutine, TypeVar, cast

F = TypeVar("F", bound=Callable[..., Coroutine[Any, Any, Any]])


async def semaphore_gather(*tasks, limit: int):
    semaphore = asyncio.Semaphore(limit)

    async def sem_task(task):
        async with semaphore:
            return await task

    return await asyncio.gather(*(sem_task(task) for task in tasks))


def limited_concurrency(limit: int) -> Callable[[F], F]:
    """Limit the number of concurrent calls to the decorated async function."""

    def decorator(func: F) -> F:
        semaphore = asyncio.Semaphore(limit)

        async def wrapper(*args: object, **kwargs: object) -> Any:
            async with semaphore:
                return await func(*args, **kwargs)

        return cast(F, wrapper)

    return decorator
