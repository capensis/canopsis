import {
  maskSensitiveHttpHeaders,
  parseRawHttp,
  SENSITIVE_HTTP_VALUE_MASK,
} from '@/helpers/request';

describe('request helpers', () => {
  describe('maskSensitiveHttpHeaders', () => {
    it('masks sensitive headers without changing the request line or body', () => {
      const body = JSON.stringify({
        password: 'body-secret',
        nested: {
          accessToken: 'body-token',
        },
      });
      const parsed = parseRawHttp([
        'POST /hooks?access_token=query-secret&visible=value HTTP/1.1',
        'Host: example.test',
        'Authorization: Bearer header-secret',
        'X-API-Key: api-key-secret',
        'Content-Type: application/json',
        '',
        body,
      ].join('\r\n'));

      const result = maskSensitiveHttpHeaders(parsed);

      expect(result.startLine).toBe(parsed.startLine);
      expect(result.headers).toEqual({
        Host: 'example.test',
        Authorization: SENSITIVE_HTTP_VALUE_MASK,
        'X-API-Key': SENSITIVE_HTTP_VALUE_MASK,
        'Content-Type': 'application/json',
      });
      expect(result.body).toBe(body);
    });

    it('does not mutate the parsed HTTP data', () => {
      const parsed = {
        startLine: 'HTTP/1.1 200 OK',
        headers: { Authorization: 'secret' },
        body: '{"token":"secret"}',
      };

      maskSensitiveHttpHeaders(parsed);

      expect(parsed).toEqual({
        startLine: 'HTTP/1.1 200 OK',
        headers: { Authorization: 'secret' },
        body: '{"token":"secret"}',
      });
    });
  });

  describe('parseRawHttp', () => {
    it('parses header names case-insensitively', () => {
      const result = parseRawHttp([
        'GET /health HTTP/1.1',
        'host: example.test',
        'transfer-encoding: chunked',
        '',
        '2',
        'OK',
        '0',
        '',
      ].join('\r\n'));

      expect(result.host).toBe('example.test');
      expect(result.body).toBe('OK');
    });

    it('detects chunked encoding in repeated headers', () => {
      const result = parseRawHttp([
        'HTTP/1.1 200 OK',
        'Transfer-Encoding: gzip',
        'Transfer-Encoding: chunked',
        '',
        '2',
        'OK',
        '0',
        '',
      ].join('\r\n'));

      expect(result.body).toBe('OK');
    });
  });
});
