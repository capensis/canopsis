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
  if (typeof raw !== 'string') {
    throw new TypeError('Raw HTTP data must be a string');
  }

  const prepared = raw.includes('\\r\\n')
    ? raw.replace(/\\r\\n/g, '\r\n')
    : raw;

  const separator = '\r\n\r\n';
  const separatorIndex = prepared.indexOf(separator);

  if (separatorIndex === -1) {
    throw new Error('Cannot find header/body separator');
  }

  const head = prepared.slice(0, separatorIndex);
  let body = prepared.slice(separatorIndex + separator.length);

  const lines = head.split('\r\n');
  const startLine = lines[0];
  const headerLines = lines.slice(1);

  const headers = {};

  for (const line of headerLines) {
    const index = line.indexOf(':');
    if (index === -1) continue;

    const name = line.slice(0, index).trim();
    const value = line.slice(index + 1).trim();

    if (headers[name] === undefined) {
      headers[name] = value;
    } else if (Array.isArray(headers[name])) {
      headers[name].push(value);
    } else {
      headers[name] = [headers[name], value];
    }
  }

  if (
    typeof headers['Transfer-Encoding'] === 'string'
    && headers['Transfer-Encoding'].toLowerCase().includes('chunked')
  ) {
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
      const hostHeader = Array.isArray(headers.Host) ? headers.Host[0] : headers.Host;
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

    result.type = 'request';
    result.method = method;
    result.httpVersion = httpVersion;
    result.target = target;
    result.protocol = protocol;
    result.host = host;
    result.path = path;
    result.queryString = queryString;
    result.query = query;

    return result;
  }

  const responseMatch = startLine.match(/^(HTTP\/\d+(?:\.\d+)?)\s+(\d{3})\s*(.*)$/);

  if (responseMatch) {
    const [, httpVersion, statusCode, statusText] = responseMatch;

    result.type = 'response';
    result.httpVersion = httpVersion;
    result.statusCode = Number(statusCode);
    result.statusText = statusText;

    return result;
  }

  result.type = 'unknown';
  return result;
};

/**
 * Convert object to FormData object for requests
 *
 * @param {Object} obj
 * @returns {FormData}
 */
export const convertObjectToFormData = (obj = {}) => Object.entries(obj).reduce((acc, [key, value]) => {
  acc.append(key, value);
  return acc;
}, new FormData());
