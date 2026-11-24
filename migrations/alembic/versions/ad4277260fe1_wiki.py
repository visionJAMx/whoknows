"""wiki

Revision ID: ad4277260fe1
Revises: 8b0e1920971c
Create Date: 2026-11-19 16:04:19.630162

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = 'ad4277260fe1'
down_revision: Union[str, Sequence[str], None] = '8b0e1920971c'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    """Upgrade schema."""
    pass


def downgrade() -> None:
    """Downgrade schema."""
    pass
