FROM python:3.11-slim AS app

WORKDIR /app

# Dev container
#ENV FIRESTORE_EMULATOR_HOST=host.docker.internal:8081

COPY requirements_lock.txt ./
RUN pip install --no-cache-dir uv
RUN python -m uv pip install --no-cache-dir -r requirements_lock.txt

COPY scripts/leaderboard ./scripts/leaderboard

EXPOSE 8080

CMD ["./scripts/leaderboard/__main__.py"]