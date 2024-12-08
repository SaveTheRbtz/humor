import argparse
import csv
import json
import re
import sys


def clean(s: str) -> str:
    s = re.sub(r"\*\*Semi.*?\*\* +", "", s)
    s = re.sub(r"\*\*.*?\*\*: ", "", s)
    s = re.sub(r"\*\*.*?\*\* +- ", "", s)
    s = re.sub(r"\*\*.*?:\*\* ", "", s)
    s = s.replace("*", "")
    s = re.sub(r"^\d+\. ", "", s)
    s = s.strip()
    if s.startswith("- "):
        s = s[2:].strip()

    if s and s[0] == '"' and s[-1] == '"':
        s = s[1:-1]
    s = re.sub(r" +", " ", s)

    return s.strip()


parser = argparse.ArgumentParser(description="Process themes with humor generation.")
parser.add_argument("--model", type=str, required=True, help="Model to use")
parser.add_argument("input_file", type=str, help="Path to the input file")
args = parser.parse_args()

with open(args.input_file, "rt", encoding="utf-8") as f:
    data = json.load(f)

if data["step4"].count("\n\n") > 2:
    result_items = [
        {"text": re.sub(r"^\d+\. *", "", it)} for it in data["step4"].strip().split("\n\n") if it
    ]
else:
    result_items = [
        {"text": re.sub(r"^\d+\. *", "", it)} for it in data["step4"].strip().split("\n") if it
    ]

for j in result_items[1:-1]:
    j = re.sub(r"\(.*\)", "", j["text"])
    j = clean(j)
    if j.endswith(":"):
        continue
    csv.writer(sys.stdout, delimiter="\t").writerow([args.model, data["theme"], j])
