"""
Python gRPC Server - Entry point
"""
import logging
import sys
import os
from concurrent import futures
import grpc

# Add generated folder to path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), 'generated'))

# Config logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Import config
from config import ServerConfig

# Import proto-generated files
import data_transfer_pb2
import data_transfer_pb2_grpc
import common_pb2

# Import handlers
from handlers.file_handler import FileHandler
from handlers.json_handler import JsonHandler


class DataTransferServicer(data_transfer_pb2_grpc.DataTransferServiceServicer):
    """Implement DataTransferService"""
    
    def ProcessFile(self, request_iterator, context):
        """
        Xử lý file stream
        
        Args:
            request_iterator: Iterator của FileChunk messages
            context: gRPC context
            
        Returns:
            ProcessFileResponse
        """
        try:
            logger.info("ProcessFile called")
            file_data = b''
            filename = ""
            
            # Collect all chunks
            for chunk in request_iterator:
                filename = chunk.filename
                file_data += chunk.data
                logger.info(
                    f"Received chunk {chunk.chunk_index + 1}/{chunk.total_chunks} "
                    f"for {filename}"
                )
            
            logger.info(f"All chunks received for {filename}")
            
            # Process file
            result = FileHandler.process_file(file_data, filename)
            
            # Build response
            response = data_transfer_pb2.ProcessFileResponse()
            response.status.status = common_pb2.ProcessingStatus.SUCCESS
            response.status.message = "File processed successfully"
            response.status.code = 0
            response.result = str(result)
            
            for key, value in result.items():
                response.data[key] = str(value)
            
            return response
            
        except Exception as e:
            logger.error(f"Error in ProcessFile: {str(e)}")
            response = data_transfer_pb2.ProcessFileResponse()
            response.status.status = common_pb2.ProcessingStatus.ERROR
            response.status.message = str(e)
            response.status.code = 1
            return response
    
    def ProcessJson(self, request, context):
        """
        Xử lý JSON request
        
        Args:
            request: JsonRequest message
            context: gRPC context
            
        Returns:
            JsonResponse
        """
        try:
            logger.info("ProcessJson called")
            
            # Process JSON
            result = JsonHandler.process_json(request.json_data)
            
            # Build response
            response = data_transfer_pb2.JsonResponse()
            response.status.status = common_pb2.ProcessingStatus.SUCCESS
            response.status.message = "JSON processed successfully"
            response.status.code = 0
            response.result = str(result)
            
            for key, value in result.items():
                response.data[key] = str(value)
            
            return response
            
        except Exception as e:
            logger.error(f"Error in ProcessJson: {str(e)}")
            response = data_transfer_pb2.JsonResponse()
            response.status.status = common_pb2.ProcessingStatus.ERROR
            response.status.message = str(e)
            response.status.code = 1
            return response


def serve():
    """Start gRPC server"""
    config = ServerConfig()
    logger.info(f"Starting server: {config}")
    
    # Create server
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=config.max_workers))
    
    # Register servicer
    data_transfer_pb2_grpc.add_DataTransferServiceServicer_to_server(
        DataTransferServicer(), server
    )
    
    # Add port
    port_str = f"{config.host}:{config.port}"
    server.add_insecure_port(port_str)
    
    logger.info(f"Server listening on {port_str}")
    
    # Start server
    server.start()
    
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Shutting down server...")
        server.stop(0)


if __name__ == "__main__":
    serve()
