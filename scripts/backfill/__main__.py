#!/usr/bin/env python

"""
gcloud init
gcloud auth application-default login
gcloud config set project humor-arena
"""

import csv
import logging
import random

import click
from google.cloud import firestore
from tqdm import tqdm

logging.basicConfig(level=logging.INFO)

MODEL_CODES = {
    "o1-preview_ablated": "foxtrot-palegreen",
    "claude-3-5-haiku-20241022": "hotel-honeydew",
    "claude-3-opus-20240229": "oscar-olive",
    "claude-3-5-sonnet-20241022": "sierra-salmon",
}


@click.command()
@click.option("--jokes-file", default="jokes.tsv", help="Path to jokes.tsv")
@click.option("--project", default="humor-arena", help="Firestore project ID")
@click.option("--dry-run", is_flag=True, help="Perform a dry run without updating the database")
def import_data(jokes_file, project, dry_run):
    db = firestore.Client(project=project)
    existing_themes_docs = db.collection("themes").get()
    theme_text_map = {doc.get("text"): doc for doc in existing_themes_docs}

    with open(jokes_file, "r", encoding="utf-8") as f:
        for line in tqdm(csv.DictReader(f, delimiter="\t", fieldnames=["model", "theme", "text"])):
            model, theme, text = line["model"].strip(), line["theme"].strip(), line["text"].strip()
            if not model or not theme or not text:
                continue
            if len(text) < 20:
                logging.warning(f"Joke too short: {text}")
                continue
            if text.endswith(":"):
                logging.warning(f"Joke ends with colon: {text}")
                continue

            existing_theme = theme_text_map.get(theme)
            if not existing_theme:
                theme_doc = {
                    "text": theme,
                    "random": random.random(),
                }
                if dry_run:
                    logging.info(f"Would create theme: {theme_doc}")
                else:
                    doc_ref = db.collection("themes").document()
                    doc_ref.set(theme_doc)
                    doc_snapshot = doc_ref.get()
                    theme_text_map[theme] = doc_snapshot

            existing_jokes = db.collection("jokes").where("text", "==", text).limit(1).get()
            if existing_jokes:
                logging.info(f"Joke already exists: {text}")
                continue

            joke_doc = {
                "text": text,
                "model": model,
                "model_code": MODEL_CODES.get(model),
                "code": MODEL_CODES.get(model),
                "policy": "v2",
                "theme": theme,
                "theme_set": "v2",
                "theme_id": theme_text_map[theme].id,
                "random": random.random(),
                "active": True,
            }
            if dry_run:
                logging.info(f"Would create joke: {joke_doc}")
            else:
                doc_ref = db.collection("jokes").document()
                doc_ref.set(joke_doc)


import_data()
