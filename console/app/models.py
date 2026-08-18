from datetime import datetime
from pydantic import BaseModel

class AgentHeartbeat(BaseModel):
    agent_id: str
    hostname: str
    os: str
    architecture: str
    agent_version: str
    sent_at: datetime
    
class AgentStatus(BaseModel):
    agent_id: str
    hostname: str
    os: str
    architecture: str
    agent_version: str
    last_seen: datetime