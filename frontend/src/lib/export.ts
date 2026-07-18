// Lightweight export helpers — no external dependencies.
// CSV is generated in-memory and downloaded via a Blob URL.
// PDF uses the browser's native print dialog (window.print) on a styled HTML view.

function escapeCSV(value: unknown): string {
  if (value === null || value === undefined) return '';
  const str = String(value);
  // Wrap in quotes and escape embedded quotes.
  if (/[",\n\r]/.test(str)) {
    return `"${str.replace(/"/g, '""')}"`;
  }
  return str;
}

/**
 * Convert an array of plain objects into CSV text.
 * Column headers are derived from the keys of the first object.
 */
export function toCSV(data: Record<string, any>[]): string {
  if (!data.length) return '';
  const headers = Object.keys(data[0]);
  const lines = [headers.join(',')];
  for (const row of data) {
    lines.push(headers.map((h) => escapeCSV(row[h])).join(','));
  }
  return lines.join('\n');
}

/**
 * Trigger a browser download of `data` (array of objects) as a CSV file.
 */
export function exportToCSV(data: Record<string, any>[], filename: string): void {
  const csv = toCSV(data);
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename.endsWith('.csv') ? filename : `${filename}.csv`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
}

/**
 * Export arbitrary HTML content to PDF using the browser print dialog.
 * Opens the HTML in a new window and invokes print, which lets the user
 * save as PDF. Falls back gracefully when popups are blocked.
 */
export function exportToPDF(html: string, filename: string): void {
  const win = window.open('', '_blank', 'width=900,height=700');
  if (!win) {
    alert('Please allow popups to export as PDF.');
    return;
  }
  win.document.write(`<!doctype html><html><head><meta charset="utf-8"><title>${filename}</title>
    <style>
      body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; color: #4d5f7a; padding: 2rem; }
      h1 { font-size: 1.5rem; }
      table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
      th, td { text-align: left; padding: 0.5rem; border-bottom: 1px solid #e0e5ec; }
      th { text-transform: uppercase; font-size: 0.75rem; color: #7d8da8; }
      @media print { body { padding: 0; } }
    </style></head><body>${html}</body></html>`);
  win.document.close();
  win.focus();
  // Give the new document a tick to render before printing.
  setTimeout(() => {
    win.print();
    // Some browsers close the print dialog automatically; leave the window
    // open so the user can re-print or save.
  }, 300);
}
