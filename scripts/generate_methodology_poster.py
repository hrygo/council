import markdown
import os
import re
import subprocess
# from weasyprint import HTML, CSS (Disabled due to missing system libs)


def generate_poster():
    # Paths
    project_root = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    input_file = os.path.join(project_root, 'docs/methodology/AI_DRIVEN_PRODUCT_PROCESS.md')
    output_pdf = os.path.join(project_root, 'docs/methodology/AI_DRIVEN_PRODUCT_PROCESS_Poster.pdf')
    resource_base = os.path.join(project_root, 'docs/methodology')

    print(f"Reading markdown from: {input_file}")
    with open(input_file, 'r', encoding='utf-8') as f:
        md_content = f.read()

    # Pre-processing
    # 1. Transform the specific image tag into a distinct div for specific internal styling
    md_content = re.sub(
        r'!\[(.*?)\]\((.*?)\)', 
        r'<figure class="main-visual"><img src="\2" alt="\1"><figcaption>\1</figcaption></figure>', 
        md_content
    )

    # 2. Handle Mermaid blocks for Rendering
    #    We wrap them in a div class="mermaid" so the library can find and render them.
    md_content = re.sub(
        r'```mermaid([\s\S]*?)```', 
        r'<div class="mermaid">\1</div>', 
        md_content
    )

    # 3. Fix List Rendering (Markdown requires a blank line before lists)
    #    Look for lines ending with a colon (English or Chinese) followed immediately by a list item
    md_content = re.sub(r'(:|：)\n(-|\*|\d+\.) ', r'\1\n\n\2 ', md_content)

    # Convert to HTML
    html_body = markdown.markdown(md_content, extensions=['tables', 'fenced_code', 'toc', 'sane_lists'])

    # Post-Processing: Wrap H2 sections in Glass Cards
    # We look for <h2>... (anything not <h2>) ... until next <h2> or end
    # Note: Regex allows for attributes in h2 tag just in case
    html_body = re.sub(
        r'(<h2.*?>[\s\S]*?)(?=<h2|\Z)', 
        r'<section class="glass-card">\1</section>', 
        html_body
    )

    # Ultra-Performance "Flat Glass" CSS (Mobile Optimized)
    css = """
    body {
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
        font-size: 11pt;
        line-height: 1.5;
        color: #1e293b; 
        margin: 0;
        padding: 15mm;
        width: 100%;
        box-sizing: border-box;
        max-width: 100%;
        margin-left: auto;
        margin-right: auto;
        
        /* Solid Flat Background - No Gradients = Fast Rendering */
        background-color: #f4f6f9; 
    }

    /* Flat Card Style - Simulates Glass with Solid Colors */
    .glass-card {
        background: #ffffff; /* Solid White = No Alpha Blending calc */
        border: 1px solid #e2e8f0;
        border-radius: 12px;
        padding: 25px;
        margin-bottom: 30px;
        /* Simple optimized shadow */
        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05); 
        break-inside: avoid;
    }

    /* Headings */
    h1 {
        font-size: 26pt; 
        font-weight: 800;
        color: #0073e5; 
        text-align: center;
        margin-bottom: 40px;
        text-transform: uppercase;
        letter-spacing: -0.5px;
        /* text-shadow: removed for performance */
        
        border-bottom: 3px solid #ff531a;
        display: inline-block;
        padding-bottom: 8px;
        position: relative;
        left: 50%;
        transform: translateX(-50%);
    }

    blockquote {
        background: #f0f9ff; /* Solid light cyan */
        border-left: 4px solid #ff531a; 
        margin: 20px 0;
        padding: 15px 20px;
        color: #334155;
        font-size: 10.5pt;
        font-style: italic;
        border-radius: 0 8px 8px 0;
    }

    h2 {
        font-size: 18pt; 
        font-weight: 700;
        color: #0073e5; 
        margin-top: 0;
        padding-bottom: 10px;
        border-bottom: 2px solid #e0f2fe; /* Solid line instead of shadow/gradient */
        margin-bottom: 20px;
    }

    h3 {
        font-size: 14pt; 
        font-weight: 600;
        color: #00bfff; 
        margin-top: 25px;
        margin-bottom: 15px;
    }
    
    strong {
        color: #0073e5;
        font-weight: 700;
    }
    
    /* Horizontal Rule (Divider) */
    hr {
        border: none;
        height: 1px;
        background-color: #cbd5e1; /* Solid color */
        margin: 40px 0;
    }
    
    /* Main Visual Image */
    figure.main-visual {
        margin: 0 0 40px 0;
        text-align: center;
        width: 100%;
    }
    
    figure.main-visual img {
        background: white;
        padding: 8px;
        border: 1px solid #e2e8f0;
        border-radius: 12px;
        box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.05);
        max-width: 100%;
        width: auto;
    }
    
    figcaption {
        color: #64748b;
        font-size: 9pt;
        margin-top: 10px;
        font-weight: 500;
    }

    /* Standard Images */
    img {
        max-width: 100%;
        height: auto;
        border-radius: 6px;
    }

    /* Tables */
    table {
        width: 100%;
        border-collapse: separate; 
        border-spacing: 0;
        margin: 30px 0;
        font-size: 10pt;
        background: white;
        border-radius: 8px;
        overflow: hidden;
        border: 1px solid #e2e8f0;
    }
    
    th, td {
        padding: 10px 15px;
        border-bottom: 1px solid #f1f5f9;
        color: #334155;
    }
    
    th {
        background-color: #f8fafc; /* Solid very light grey */
        color: #0073e5;
        font-weight: 700;
        text-transform: uppercase;
        font-size: 9pt;
        letter-spacing: 0.5px;
    }
    
    th:nth-child(1), td:nth-child(1),
    th:nth-child(2), td:nth-child(2) {
        white-space: nowrap;
        width: 1%; /* Shrink to fit content */
        font-weight: 700;
    }

    tr:last-child td {
        border-bottom: none;
    }

    /* Code blocks */
    pre {
        background: #f8fafc; /* Very light grey */
        border: 1px solid #e2e8f0;
        border-radius: 10px;
        padding: 25px;
        font-family: 'JetBrains Mono', monospace;
        font-size: 13pt;
        color: #0f172a;
        margin: 30px 0;
        overflow-x: auto;
    }
    
    code {
        font-family: 'JetBrains Mono', monospace;
        background: rgba(0, 115, 229, 0.1);
        color: #0073e5;
        padding: 2px 6px;
        border-radius: 4px;
        font-size: 0.9em;
    }
    
    pre code {
        background: none;
        color: inherit;
        padding: 0;
    }

    /* Mermaid Container */
    .mermaid {
        display: flex;
        justify-content: center;
        width: 100%;
        background: white;
        padding: 40px;
        border-radius: 16px;
        margin: 40px 0;
        box-sizing: border-box;
        border: 1px solid rgba(0, 115, 229, 0.1);
        box-shadow: 0 4px 20px rgba(0,0,0,0.05);
    }

    .mermaid svg {
        max-width: 100% !important;
        height: auto !important;
    }
    
    /* Lists */
    ul, ol {
        margin: 1em 0;
        padding-left: 1.5em;
    }
    
    li {
        margin-bottom: 0.6em;
        color: #334155;
    }
    
    /* Custom Scrollbar for overflow elements (just in case) */
    ::-webkit-scrollbar {
        width: 8px;
        height: 8px;
    }
    ::-webkit-scrollbar-thumb {
        background: #cbd5e1;
        border-radius: 4px;
    }
    """

    html_content = f"""
    <!DOCTYPE html>
    <html lang="zh-CN">
    <head>
        <meta charset="UTF-8">
        <title>Methodology Poster</title>
        <style>
        {css}
        </style>
        <!-- Inject Mermaid - White Theme -->
        <script type="module">
            import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs';
            mermaid.initialize({{ 
                startOnLoad: true, 
                theme: 'base',
                themeVariables: {{
                    'primaryColor': '#e0f2fe',
                    'primaryTextColor': '#0073e5',
                    'primaryBorderColor': '#00bfff',
                    'lineColor': '#0073e5',
                    'secondaryColor': '#ffedd5',
                    'tertiaryColor': '#fff'
                }},
                securityLevel: 'loose'
            }});
        </script>
    </head>
    <body class="glass-theme">
        <div id="poster-content">
            {html_body}
        </div>
    </body>
    </html>
    """

    # Output HTML
    html_output_file = os.path.join(project_root, 'docs/methodology/AI_DRIVEN_PRODUCT_PROCESS_Poster.html')
    with open(html_output_file, 'w', encoding='utf-8') as f:
        f.write(html_content)
    
    print(f"HTML generated at: {html_output_file}")
    
    # Automate PDF Generation via Puppeteer
    print("🚀 Triggering PDF conversion (Puppeteer)...")
    node_script = os.path.join(project_root, 'scripts/pdf_tool/convert.js')
    
    try:
        # Check if node is available
        subprocess.run(['node', '-v'], check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        
        # Run conversion
        result = subprocess.run(['node', node_script], cwd=project_root, text=True)
        if result.returncode == 0:
            print("✅ PDF Generation Complete!")
        else:
            print("❌ PDF Generation Failed. Please check Node.js output above.")
            
    except FileNotFoundError:
        print("❌ Error: 'node' executable not found. Please ensure Node.js is installed to generate PDF.")
    except Exception as e:
        print(f"❌ Unexpected error during PDF generation: {e}")

if __name__ == "__main__":
    generate_poster()
