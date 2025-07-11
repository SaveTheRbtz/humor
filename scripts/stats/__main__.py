#!/usr/bin/env python
import sys
import json
from google.cloud import firestore
from collections import Counter

# Initialize the Firestore client
client = firestore.Client('humor-arena')

# Function to download a Firestore collection (kind) to a local JSON file
def download_collection_to_json(client, kind, output_file):
    # Query all entities from the specified collection (kind)
    query = client.collection(kind).stream()
    entities = list(query)

    # Prepare a list to store all documents as dictionaries
    documents = []

    # Loop through each entity and convert it to a dictionary
    for entity in entities:
        entity_dict = entity.to_dict()
        entity_dict['id'] = entity.id
        documents.append(entity_dict)

    # Write the list of documents to a JSON file
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(documents, f, ensure_ascii=False, indent=4, default=str)

    print(f"Downloaded {len(entities)} documents from collection '{kind}' to '{output_file}'.")
    return documents

tag = sys.argv[1]

d1 = download_collection_to_json(client, 'jokes', f'{tag}.jokes.json')
d2 = download_collection_to_json(client, 'choices', f'{tag}.choices.json')

report = []
report.append( f'{len(d1)} jokes' )
report.append( f'{len(d2)} choices' )
report.append( f"{len(set([it['session_id'] for it in d2]))} sessions" )
report.append( str(Counter([it['winner'] for it in d2])) + ' winners, where NONE = 1; LEFT = 2; RIGHT = 3; BOTH = 4;' )

print("\n".join(report))

with open(f'{tag}.report.txt', 'w') as f:
    f.write("\n".join(report))
