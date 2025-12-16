#!/usr/bin/env python3
import os
import sys
import subprocess
import re
import time
import argparse

# --- Configuration ---
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.dirname(SCRIPT_DIR)
DOCS_DIR = os.path.join(PROJECT_ROOT, "docs")
HISTORY_DIR = os.path.join(DOCS_DIR, "reports")
DIALECTA_BIN = "dialecta"

# --- Helpers ---

def read_file(path):
    try:
        with open(path, 'r', encoding='utf-8') as f:
            return f.read()
    except Exception as e:
        print(f"Error reading {path}: {e}")
        # Return empty string if history file is missing/empty
        return ""

def write_file(path, content):
    with open(path, 'w', encoding='utf-8') as f:
        f.write(content)

def clean_ansi(text):
    ansi_escape = re.compile(r'\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])')
    return ansi_escape.sub('', text)

def run_dialecta(input_text, context_prompt=""):
    """Invokes dialecta CLI via stdin, streaming output to console in real-time with throttling."""
    full_prompt = f"{context_prompt}\n\nMATERIAL TO ANALYZE:\n{input_text}"
    
    try:
        process = subprocess.Popen(
            [DIALECTA_BIN, "-"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,  # Merge stderr into stdout for unified streaming
            text=True,
            bufsize=1,  # Line buffering
            cwd=DOCS_DIR
        )
        
        # Write input and close stdin
        try:
            process.stdin.write(full_prompt)
            process.stdin.close()
        except BrokenPipeError:
            pass

        captured_output = []
        last_status_time = 0
        STATUS_THROTTLE = 2.0  # Seconds
        
        # Stream output line by line
        for line in process.stdout:
            is_thinking = "Thinking..." in line
            
            if is_thinking:
                current_time = time.time()
                if current_time - last_status_time < STATUS_THROTTLE:
                    # Skip printing and skip saving to keep report clean
                    continue
                last_status_time = current_time
                sys.stdout.write(line)
                sys.stdout.flush()
                # Do NOT append transient 'Checking/Thinking' status lines to the final report file
                continue
            
            # Print and capture all other lines (content, headers, etc.)
            sys.stdout.write(line)
            sys.stdout.flush()
            captured_output.append(line)
            
        return_code = process.wait()
        full_output = "".join(captured_output)
        
        if return_code != 0:
            print(f"Dialecta Error (Exit Code {return_code})")
            return None
            
        return clean_ansi(full_output)

    except FileNotFoundError:
        print("Error: 'dialecta' binary not found in PATH.")
        sys.exit(1)
    except Exception as e:
        print(f"Unexpected error running dialecta: {e}")
        return None

# --- Core Task ---

def convene_council(draft_path, prd_path, history_path):
    print(f" >>> 🗣️  Convening the Council...")
    print(f"      Target: {os.path.basename(draft_path)}")
    print(f"      Context: {os.path.basename(prd_path)}")
    if history_path and os.path.exists(history_path):
        print(f"      History: {os.path.basename(history_path)}")
        history_content = read_file(history_path)
    else:
        print(f"      History: (None provided or empty)")
        history_content = "No prior history available."
    
    draft_content = read_file(draft_path)
    prd_content = read_file(prd_path)
    
    combined_input = f"--- PRD (CONTEXT) ---\n{prd_content}\n\n--- DESIGN DRAFT (TARGET) ---\n{draft_content}"
    
    instruction = """
    INSTRUCTIONS FOR ADJUDICATOR:
    1. Critically review the 'DESIGN DRAFT' against the 'PRD' and 'CONTEXT HISTORY'.
    2. Focus on: Architecture, User Experience, and Feasibility.
    3. Provide actionable recommendations.
    4. DO NOT Generate Code Blocks for the whole file. Focus on the 'Verdict' and 'Analysis'.
    """
    
    report = run_dialecta(combined_input, context_prompt=f"{instruction}\n\nCONTEXT HISTORY:\n{history_content}")
    
    if not report:
        return None
        
    # Save Report
    timestamp = int(time.time())
    report_filename = f"debate_{timestamp}.md"
    report_path = os.path.join(HISTORY_DIR, report_filename)
    write_file(report_path, report)
    print(f"      ✅ Report Generated: {report_path}")
    
    return report_path

# --- Main CLI ---

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Run a Council Session (Dialecta Debate).")
    parser.add_argument("draft", help="Path to the Design Draft file")
    parser.add_argument("prd", help="Path to the PRD file")
    parser.add_argument("--history", help="Path to the History Summary file", required=False)
    
    args = parser.parse_args()
    
    # Ensure reports dir
    os.makedirs(HISTORY_DIR, exist_ok=True)
    
    # 1. Run Debate
    report_path = convene_council(args.draft, args.prd, args.history)
    
    if not report_path:
        sys.exit(1)
        
    print(f"\nOUTPUT_REPORT={report_path}")
    sys.exit(0)
