// @vitest-environment jsdom
import { afterEach, describe, expect, it, vi } from 'vitest';
import { toCSV, exportToPDF } from './export';

afterEach(() => vi.restoreAllMocks());

describe('CSV export', () => {
  it('neutralizes spreadsheet formulas while preserving quotes and line breaks', () => {
    expect(toCSV([
      { title: '=SUM(A1:A2)', details: 'a,"b"\nnext', count: -5 },
      { title: '  +cmd', details: '@formula', count: 0 },
      { title: '\tplain', details: '-1+1', count: null },
    ])).toBe('title,details,count\n\'=SUM(A1:A2),"a,""b""\nnext",-5\n\'  +cmd,\'@formula,0\n\'\tplain,\'-1+1,');
  });

  it('keeps the supplied filtered rows and every column in heterogeneous analytics data', () => {
    expect(toCSV([{ section: 'volume', date: '2026-10-22', count: 2 }, { section: 'resolution', priority: 'high', avg_hours: null }]))
      .toBe('section,date,count,priority,avg_hours\nvolume,2026-10-22,2,,\nresolution,,,high,');
  });
});

describe('print export', () => {
  it('prints text and tables without active HTML or an injected filename', () => {
    const printable = document.implementation.createHTMLDocument();
    const popup = { document: printable, focus: vi.fn(), print: vi.fn() };
    vi.spyOn(window, 'open').mockReturnValue(popup as unknown as Window);
    vi.useFakeTimers();
    try {
      exportToPDF('<table><tr><td onclick="alert(1)">Ticket &amp; details</td></tr></table><script>test()</script><img src="https://outside.test/track">', '</title><div id="filename-marker">');
      expect(printable.querySelector('#filename-marker')).toBeNull();
      expect(printable.querySelector('script, img, [onclick]')).toBeNull();
      expect(printable.querySelector('td')?.textContent).toBe('Ticket & details');
      vi.runAllTimers();
      expect(popup.print).toHaveBeenCalledTimes(1);
    } finally { vi.useRealTimers(); }
  });

  it('informs the user when the print popup is blocked', () => {
    vi.spyOn(window, 'open').mockReturnValue(null);
    const alert = vi.spyOn(window, 'alert').mockImplementation(() => {});
    exportToPDF('text', 'tickets');
    expect(alert).toHaveBeenCalledWith('Please allow popups to export as PDF.');
  });
});
