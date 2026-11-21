FROM python:3.12

WORKDIR /app

COPY alembic/requirements.txt ./requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

COPY alembic/ ./

ENTRYPOINT ["alembic"]
CMD ["upgrade", "head"]
