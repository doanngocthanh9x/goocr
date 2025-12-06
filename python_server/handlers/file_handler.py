"""
File handler - Xử lý logic file từ gRPC request
"""
import json
import logging

logger = logging.getLogger(__name__)


class FileHandler:
    """Xử lý file stream từ gRPC"""
    
    @staticmethod
    def process_file(file_data: bytes, filename: str) -> dict:
        """
        Xử lý file nhận được
        
        Args:
            file_data: Nội dung file
            filename: Tên file
            
        Returns:
            dict: Kết quả xử lý
        """
        try:
            logger.info(f"Processing file: {filename}")
            
            # Logic xử lý file (ví dụ: đọc, phân tích)
            file_size = len(file_data)
            
            result = {
                "filename": filename,
                "size": file_size,
                "status": "processed",
                "content_preview": file_data[:100].decode('utf-8', errors='ignore')
            }
            
            logger.info(f"File processed successfully: {filename}")
            return result
            
        except Exception as e:
            logger.error(f"Error processing file: {str(e)}")
            return {
                "error": str(e),
                "status": "error"
            }
