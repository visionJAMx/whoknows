from sqlalchemy import Column, Integer, String, Text, TIMESTAMP, CheckConstraint
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()

class User(Base):
    __tablename__ = 'users'
    id = Column(Integer, primary_key=True, autoincrement=True)
    username = Column(String, unique=True, nullable=False)
    email = Column(String, unique=True, nullable=False)
    password_hash = Column(String, nullable=False)

class Page(Base):
    __tablename__ = 'pages'
    title = Column(String, primary_key=True)
    url = Column(String, unique=True, nullable=False)
    language = Column(String, nullable=False, default='en')
    last_updated = Column(TIMESTAMP, nullable=False)
    content = Column(Text, nullable=False)
    __table_args__ = (
        CheckConstraint("language IN ('en', 'da')", name='check_language'),
    )