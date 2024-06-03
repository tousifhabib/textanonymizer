from typing import List
from flask import Flask, request, jsonify
import spacy

app = Flask(__name__)

nlp = spacy.load('en_core_web_sm')


def extract_entities(text: str, label: str) -> List[str]:
    doc = nlp(text)
    return list(set(ent.text for ent in doc.ents if ent.label_ == label))


def get_person_names(text: str) -> List[str]:
    return extract_entities(text, 'PERSON')


def anonymize_text(text: str) -> str:
    doc = nlp(text)
    anonymized_tokens = []
    for token in doc:
        if token.ent_type_ == 'PERSON':
            anonymized_tokens.append('REDACTED')
        else:
            anonymized_tokens.append(token.text)
    return ' '.join(anonymized_tokens)


@app.route('/anonymize', methods=['POST'])
def anonymize():
    data = request.get_json()

    print(data)
    text = data.get('text', '')

    anonymized_text = anonymize_text(text)

    print("output")
    print(anonymized_text)

    response = {
        'choices': [
            {
                'message': {
                    'content': anonymized_text
                }
            }
        ]
    }

    return jsonify(response)


if __name__ == '__main__':
    app.run(port=5000)