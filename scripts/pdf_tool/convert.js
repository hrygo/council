const puppeteer = require('puppeteer');
const path = require('path');

(async () => {
    try {
        // Adjust paths relative to this script: ../../docs/methodology/...
        const htmlPath = path.resolve(__dirname, '../../docs/methodology/AI_DRIVEN_PRODUCT_PROCESS_Poster.html');
        const pdfPath = path.resolve(__dirname, '../../docs/methodology/AI_DRIVEN_PRODUCT_PROCESS_Poster.pdf');

        console.log(`Reading HTML from: ${htmlPath}`);
        console.log(`Writing PDF to: ${pdfPath}`);

        const browser = await puppeteer.launch({
            headless: "new",
            args: ['--no-sandbox', '--disable-setuid-sandbox']
        });
        const page = await browser.newPage();

        // Load the HTML file
        // Note: Make sure the base URL for local resources (images) works.
        // Puppeteer loads file://, so relative paths in HTML (./images/...) are resolved relative to HTML file location.
        // This is perfectly aligned with our Python script output.
        await page.goto(`file://${htmlPath}`, {
            waitUntil: 'networkidle0'
        });

        // Wait for Mermaid to Render
        try {
            await page.waitForSelector('.mermaid svg', { timeout: 10000 });
            console.log('Mermaid diagrams rendered.');
        } catch (e) {
            console.log('No Mermaid diagrams found or timed out. Proceeding...');
        }

        // Calculate Height
        const bodyHeight = await page.evaluate(() => document.body.scrollHeight + 50); // Small Buffer
        const customWidth = '210mm'; // A4 Width

        console.log(`Generating Long PDF with Dimensions: ${customWidth} x ${bodyHeight}px`);

        await page.pdf({
            path: pdfPath,
            width: customWidth,
            height: `${bodyHeight}px`,
            printBackground: true,
            displayHeaderFooter: false,
            margin: { top: '0', right: '0', bottom: '0', left: '0' }
        });

        await browser.close();
        console.log('PDF generated successfully!');
    } catch (e) {
        console.error('Error generating PDF:', e);
        process.exit(1);
    }
})();
