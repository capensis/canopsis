import { findKey, isString, mapValues } from 'lodash';

export const SENSITIVE_HTTP_VALUE_MASK = '********';

const SENSITIVE_HTTP_KEY_PARTS = [
  'apikey',
  'authorization',
  'credential',
  'cookie',
  'password',
  'passwd',
  'privatekey',
  'pwd',
  'secret',
  'signature',
  'token',
];

/**
 * Checks whether an HTTP field name is commonly used for secret data.
 * Separators and casing are ignored so that names such as api_key,
 * client-secret, and accessToken are all handled consistently.
 *
 * @param {string} key
 * @returns {boolean}
 */
const isSensitiveHttpKey = (key) => {
  const normalizedKey = String(key).toLowerCase().replace(/[^a-z0-9]/g, '');

  return SENSITIVE_HTTP_KEY_PARTS.some(part => normalizedKey.includes(part));
};

/**
 * Returns an HTTP header value using a case-insensitive name comparison.
 *
 * @param {Object} headers
 * @param {string} name
 * @returns {string|string[]|undefined}
 */
const getHttpHeader = (headers, name) => {
  const normalizedName = name.toLowerCase();
  const headerKey = findKey(headers, (value, key) => key.toLowerCase() === normalizedName);

  return headerKey ? headers[headerKey] : undefined;
};

/**
 * Masks sensitive headers in a parsed HTTP request or response.
 *
 * @param {Object} data - Result returned by parseRawHttp
 * @returns {Object} Redacted copy of the parsed data
 */
export const maskSensitiveHttpHeaders = (data) => {
  const headers = data.headers && mapValues(
    data.headers,
    (value, key) => (isSensitiveHttpKey(key) ? SENSITIVE_HTTP_VALUE_MASK : value),
  );

  return {
    ...data,
    headers,
  };
};

/**
 * Decodes HTTP chunked transfer-encoding body.
 *
 * @param {string} body - Chunked-encoded body string
 * @returns {string} Decoded body content
 * @throws {Error} When chunk size line is not found or chunk size is invalid
 */
export const decodeChunkedBody = (body) => {
  let cursor = 0;
  let decoded = '';

  while (cursor < body.length) {
    const lineEnd = body.indexOf('\r\n', cursor);

    if (lineEnd === -1) {
      throw new Error('Invalid chunked body: chunk size line not found');
    }

    const chunkSizeHex = body.slice(cursor, lineEnd).trim();
    const chunkSize = parseInt(chunkSizeHex, 16);

    if (Number.isNaN(chunkSize)) {
      throw new Error(`Invalid chunk size: ${chunkSizeHex}`);
    }

    cursor = lineEnd + 2;

    if (chunkSize === 0) {
      break;
    }

    decoded += body.slice(cursor, cursor + chunkSize);
    cursor += chunkSize + 2;
  }

  return decoded;
};

/**
 * Parses raw HTTP request or response string into a structured object.
 *
 * @param {string} raw - Raw HTTP data string (supports escaped \\r\\n)
 * @returns {Object} Parsed result with startLine, headers, body, and type-specific fields
 * @returns {string} return.startLine - First line (request line or status line)
 * @returns {Object} return.headers - Headers object (values may be string or string[])
 * @returns {string} return.body - Decoded body content
 * @returns {'request'|'response'|'unknown'} return.type - Detected HTTP message type
 * @throws {TypeError} When raw is not a string
 * @throws {Error} When header/body separator is not found
 */
export const parseRawHttp = (raw) => {
  if (!isString(raw)) {
    throw new TypeError('Raw HTTP data must be a string');
  }

  const prepared = raw.includes('\\r\\n') ? raw.replace(/\\r\\n/g, '\r\n') : raw;
  const separator = '\r\n\r\n';
  const separatorIndex = prepared.indexOf(separator);

  if (separatorIndex === -1) {
    throw new Error('Cannot find header/body separator');
  }

  const head = prepared.slice(0, separatorIndex);
  let body = prepared.slice(separatorIndex + separator.length);
  const [startLine, ...headerLines] = head.split('\r\n');
  const headers = {};

  headerLines.forEach((line) => {
    const index = line.indexOf(':');

    if (index === -1) {
      return;
    }

    const name = line.slice(0, index).trim();
    const value = line.slice(index + 1).trim();

    if (headers[name] === undefined) {
      headers[name] = value;
    } else if (Array.isArray(headers[name])) {
      headers[name].push(value);
    } else {
      headers[name] = [headers[name], value];
    }
  });

  const transferEncoding = getHttpHeader(headers, 'transfer-encoding');
  const transferEncodings = Array.isArray(transferEncoding)
    ? transferEncoding
    : [transferEncoding];

  if (transferEncodings.some(value => isString(value) && value.toLowerCase().includes('chunked'))) {
    body = decodeChunkedBody(body);
  }

  const result = {
    startLine,
    headers,
    body,
  };
  const requestMatch = startLine.match(/^([A-Z]+)\s+(\S+)\s+(HTTP\/\d+(?:\.\d+)?)$/);

  if (requestMatch) {
    const [, method, target, httpVersion] = requestMatch;
    let protocol = null;
    let host = null;
    let path = '';
    let queryString = '';
    let query = {};

    try {
      const hostValues = getHttpHeader(headers, 'host');
      const hostHeader = Array.isArray(hostValues) ? hostValues[0] : hostValues;
      const base = hostHeader ? `http://${hostHeader}` : 'http://dummy-base';
      const url = new URL(target, base);

      protocol = url.protocol;
      host = url.host;
      path = url.pathname;
      queryString = url.search ? url.search.slice(1) : '';
      query = Object.fromEntries(url.searchParams.entries());
    } catch {
      path = target;
    }

    return {
      ...result,
      type: 'request',
      method,
      httpVersion,
      target,
      protocol,
      host,
      path,
      queryString,
      query,
    };
  }

  const responseMatch = startLine.match(/^(HTTP\/\d+(?:\.\d+)?)\s+(\d{3})\s*(.*)$/);

  if (responseMatch) {
    const [, httpVersion, statusCode, statusText] = responseMatch;

    return {
      ...result,
      type: 'response',
      httpVersion,
      statusCode: Number(statusCode),
      statusText,
    };
  }

  return {
    ...result,
    type: 'unknown',
  };
};
