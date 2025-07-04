#!/usr/bin/env python3
"""
Extract JavaScript code from GWT-style array format
"""

import re
import json
import ast

def extract_js_from_gwt_array(file_path):
    """
    Extract JavaScript code from a GWT-style array format file
    """
    print(f"Reading from: {file_path}")
    
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        print(f"File size: {len(content)} characters")
        
        # Check if file starts with array
        if not content.strip().startswith('['):
            print("File doesn't appear to be in the expected format (should start with [")
            return False
        
        # Remove any leading/trailing whitespace
        content = content.strip()
        
        # The content is a JavaScript array with a single large string
        # We need to extract the string content and unescape it
        print("Parsing JavaScript array...")
        
        # Find the opening and closing brackets
        start_idx = content.find('[')
        end_idx = content.rfind(']')
        
        if start_idx == -1 or end_idx == -1:
            print("Could not find array brackets")
            return False
        
        # Extract the content between brackets
        array_content = content[start_idx + 1:end_idx]
        
        # The array contains a single large string with escaped characters
        # We need to parse it as a JavaScript string literal
        print("Extracting and unescaping JavaScript string...")
        
        # Remove the outer quotes and unescape the string
        if array_content.startswith("'") and array_content.endswith("'"):
            # Single quoted string
            js_string = array_content[1:-1]
        elif array_content.startswith('"') and array_content.endswith('"'):
            # Double quoted string
            js_string = array_content[1:-1]
        else:
            print("Could not find string delimiters")
            return False
        
        # Unescape JavaScript string literals
        js_code = unescape_js_string(js_string)
        
        # Write the extracted code to a new file
        output_file = file_path.replace('.js', '_extracted.js')
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write(js_code)
        
        print(f"Successfully extracted JavaScript code to: {output_file}")
        print(f"Extracted code size: {len(js_code)} characters")
        
        # Show a preview of the extracted code
        preview = js_code[:200] + "..." if len(js_code) > 200 else js_code
        print(f"\nPreview of extracted code:\n{preview}")
        
        print("\nExtraction completed successfully!")
        return True
        
    except Exception as e:
        print(f"Error during extraction: {e}")
        return False

def unescape_js_string(js_string):
    """
    Unescape a JavaScript string literal
    """
    # Replace common JavaScript escape sequences
    unescaped = js_string
    
    # Handle newlines
    unescaped = unescaped.replace('\\n', '\n')
    unescaped = unescaped.replace('\\r', '\r')
    unescaped = unescaped.replace('\\t', '\t')
    
    # Handle quotes
    unescaped = unescaped.replace('\\"', '"')
    unescaped = unescaped.replace("\\'", "'")
    
    # Handle backslashes
    unescaped = unescaped.replace('\\\\', '\\')
    
    # Handle other escape sequences
    unescaped = unescaped.replace('\\b', '\b')
    unescaped = unescaped.replace('\\f', '\f')
    unescaped = unescaped.replace('\\v', '\v')
    
    return unescaped

if __name__ == "__main__":
    input_file = "./fromCache.js"
    success = extract_js_from_gwt_array(input_file)
    
    if not success:
        print("Extraction failed!") 