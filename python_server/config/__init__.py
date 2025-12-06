"""
Config cho Python gRPC server
"""
import os
from dataclasses import dataclass


@dataclass
class ServerConfig:
    """Cấu hình server"""
    host: str = os.getenv("GRPC_SERVER_HOST", "0.0.0.0")
    port: int = int(os.getenv("GRPC_SERVER_PORT", "50051"))
    max_workers: int = int(os.getenv("GRPC_MAX_WORKERS", "10"))
    enable_tls: bool = os.getenv("GRPC_ENABLE_TLS", "false").lower() == "true"
    cert_file: str = os.getenv("GRPC_CERT_FILE", "")
    key_file: str = os.getenv("GRPC_KEY_FILE", "")
    
    def __str__(self):
        return f"ServerConfig(host={self.host}, port={self.port})"
