# ¿Who Knows? - Flask application

A search engine from 2009, upgraded to run with Python 3 and Flask 3.

**Note**: This application is intentionally full of problems and vulnerabilities. Do not run it in a production environment. 

## Installation

Create and activate a virtual environment:

```bash
python3 -m venv .venv
source .venv/bin/activate
```

Install the dependencies:

```bash
python3 -m pip install -r backend/requirements.txt
```

To initialize a new database:

```bash
make init
```

Note: Windows does not natively support Make. 


## Running the application

Start a development server on port `8080`:

```bash
make run
```
Or:

```bash
python3 backend/app.py
```

## Test the application

To run the tests:

```bash
make test
```

Or:

```bash
PYTHONPATH=backend python3 backend/app_tests.py
```
