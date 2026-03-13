import { normalizeHandlebarsNbsp } from '@/helpers/handlebars/normalize';

describe('normalizeHandlebarsNbsp', () => {
  it('Returns empty string when no argument provided', () => {
    expect(normalizeHandlebarsNbsp()).toBe('');
  });

  it('Returns empty string when empty string provided', () => {
    expect(normalizeHandlebarsNbsp('')).toBe('');
  });

  it('Replaces &nbsp; in simple handlebars block', () => {
    const template = '{{value&nbsp;test}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('{{value test}}');
  });

  it('Replaces \\u00A0 in simple handlebars block', () => {
    const template = '{{value\u00A0test}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('{{value test}}');
  });

  it('Replaces &nbsp; in triple handlebars block', () => {
    const template = '{{{value&nbsp;test}}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('{{{value test}}}');
  });

  it('Replaces \\u00A0 in triple handlebars block', () => {
    const template = '{{{value\u00A0test}}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('{{{value test}}}');
  });

  it('Replaces multiple &nbsp; in handlebars block', () => {
    const template = '{{value&nbsp;test&nbsp;more}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('{{value test more}}');
  });

  it('Replaces multiple \\u00A0 in handlebars block', () => {
    const template = '{{value\u00A0test\u00A0more}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('{{value test more}}');
  });

  it('Replaces mixed &nbsp; and \\u00A0 in handlebars block', () => {
    const template = '{{value&nbsp;test\u00A0more}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('{{value test more}}');
  });

  it('Replaces NBSP in multiple handlebars blocks', () => {
    const template = '{{first&nbsp;value}} and {{second&nbsp;value}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('{{first value}} and {{second value}}');
  });

  it('Does not replace &nbsp; outside handlebars blocks', () => {
    const template = '<p>text&nbsp;here</p>{{value}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('<p>text&nbsp;here</p>{{value}}');
  });

  it('Does not replace \\u00A0 outside handlebars blocks', () => {
    const template = '<p>text\u00A0here</p>{{value}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('<p>text\u00A0here</p>{{value}}');
  });

  it('Replaces NBSP only inside handlebars blocks', () => {
    const template = '<p>text&nbsp;outside</p>{{inside&nbsp;value}}<span>more&nbsp;text</span>';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('<p>text&nbsp;outside</p>{{inside value}}<span>more&nbsp;text</span>');
  });

  it('Handles complex template with multiple block types', () => {
    const template = '{{#if condition&nbsp;test}}{{value&nbsp;here}}{{{html&nbsp;content}}}{{/if}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('{{#if condition test}}{{value here}}{{{html content}}}{{/if}}');
  });

  it('Handles template with nested handlebars blocks', () => {
    const template = '{{#each items&nbsp;list}}{{item&nbsp;value}}{{/each}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('{{#each items list}}{{item value}}{{/each}}');
  });

  it('Handles template without any handlebars blocks', () => {
    const template = '<div>Some&nbsp;text&nbsp;here</div>';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('<div>Some&nbsp;text&nbsp;here</div>');
  });

  it('Handles template with only regular spaces in handlebars blocks', () => {
    const template = '{{value test more}}';
    const result = normalizeHandlebarsNbsp(template);

    expect(result).toBe('{{value test more}}');
  });
});
