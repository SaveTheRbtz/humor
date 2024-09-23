#!/usr/bin/env python

"""
gcloud init
gcloud auth application-default login
gcloud config set project humor-arena
"""


import json
import random

import click
from google.cloud import firestore
from tqdm import tqdm


@click.command()
@click.option("--themes-file", default="themes.jsonl", help="Path to themes.jsonl")
@click.option("--jokes-file", default="jokes.jsonl", help="Path to jokes.jsonl")
@click.option("--project", default="humor-arena", help="Firestore project ID")
def import_data(themes_file, jokes_file, project):
    db = firestore.Client(project=project)
    theme_id_map = {}
    with open(themes_file, "r") as f:
        for line in tqdm(f):
            theme_data = json.loads(line)
            theme_name = theme_data["name"]
            existing_themes = (
                db.collection("themes")
                .where("text", "==", theme_name)
                .limit(1)
                .stream()
            )
            existing_theme = next(existing_themes, None)
            if existing_theme:
                theme_doc_id = existing_theme.id
            else:
                theme_doc = {
                    "text": theme_name,
                    "random": random.random(),
                }
                doc_ref = db.collection("themes").document()
                doc_ref.set(theme_doc)
                theme_doc_id = doc_ref.id
            theme_id_map[theme_data["_id"]] = theme_doc_id
    with open(jokes_file, "r") as f:
        for line in tqdm(f):
            joke_data = json.loads(line)
            joke_text = joke_data["text"]
            joke_topic = joke_data["topic"]
            theme_doc_id = theme_id_map.get(joke_topic)
            if not theme_doc_id:
                continue
            existing_jokes = db.collection('jokes').where('text', '==', joke_text).limit(1).stream()
            existing_joke = next(existing_jokes, None)
            if existing_joke:
                continue
            joke_doc = {
                'text': joke_text,
                'theme_id': theme_doc_id,
                'random': random.random(),
            }
            doc_ref = db.collection('jokes').document()
            doc_ref.set(joke_doc)

import_data()
