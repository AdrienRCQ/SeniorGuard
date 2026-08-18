from datetime import datetime, timezone
from fastapi import FastAPI
from app.models import AgentHeartbeat, AgentStatus

app = FastAPI(
    title= "SeniorGuard Console API",
    version= "0.1.1"
)

# definition de la liste des agents (vide au début)
agents : dict[str, AgentStatus] = {}

@app.get("/health")
def health() -> dict[str,str]:
    return {"status": "healthy"}
    
@app.post("/api/v1/agents/heartbeat")
def receive_heartbeat(heartbeat:AgentHeartbeat) -> dict[str,str]:
    now = datetime.now(timezone.utc)
    agents[heartbeat.agent_id] = AgentStatus(
        agent_id =heartbeat.agent_id,
        hostname =heartbeat.hostname,
        os = heartbeat.os,
        architecture = heartbeat.architecture,
        agent_version = heartbeat.agent_version,
        last_seen = now
    )
    return {"status": "accepted"}

@app.get("/api/v1/agents")
def get_agents() -> list[AgentStatus]:
    return list(agents.values())