"""
JSON handler - Xử lý logic JSON từ gRPC request
"""
import json
import logging

logger = logging.getLogger(__name__)


class JsonHandler:
    """Xử lý JSON request từ gRPC"""
    
    @staticmethod
    def process_json(json_data: str) -> dict:
        """
        Xử lý JSON nhận được
        
        Args:
            json_data: JSON string
            
        Returns:
            dict: Kết quả xử lý
        """
        try:
            logger.info("Processing JSON request")
            
            # Parse JSON
            data = json.loads(json_data)
            
            # Logic xử lý (ví dụ: validation, transformation)
            result = {
                "original_data": data,
                "processed": True,
                "item_count": len(data) if isinstance(data, (list, dict)) else 1,
                "status": "processed"
            }
            
            logger.info("JSON processed successfully")
            return result
            
        except json.JSONDecodeError as e:
            logger.error(f"JSON decode error: {str(e)}")
            return {
                "error": f"Invalid JSON: {str(e)}",
                "status": "error"
            }
        except Exception as e:
            logger.error(f"Error processing JSON: {str(e)}")
            return {
                "error": str(e),
                "status": "error"
            }
