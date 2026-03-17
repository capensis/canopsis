<template>
  <c-collapse-panel>
    <template #header>
      <span class="white--text">
        <strong>{{ titlePrefix }} : </strong>
        {{ parsedRaw.startLine }}
      </span>
    </template>
    <template v-if="parsedRaw.headers">
      <c-information-block-row
        v-for="(value, key) in parsedRaw.headers"
        :key="key"
        :label="key"
        :value="value"
      />
    </template>
    <div v-if="parsedRaw.body" class="mt-2">
      <c-json-treeview
        v-if="bodyJsonObject"
        :json-object="bodyJsonObject"
      />
      <pre v-else class="text-break pre-wrap">{{ parsedRaw.body }}</pre>
    </div>
  </c-collapse-panel>
</template>

<script>
import { computed } from 'vue';

function decodeChunkedBody(body) {
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
    cursor += chunkSize + 2; // chunk + trailing \r\n
  }

  return decoded;
}

function parseRawHttp(raw) {
  if (typeof raw !== 'string') {
    throw new TypeError('Raw HTTP data must be a string');
  }

  // Если пришла строка с буквальными \r\n, превращаем их в реальные переводы строк
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
}

export default {
  props: {
    titlePrefix: {
      type: String,
      default: '',
    },
    raw: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const parsedRaw = computed(() => parseRawHttp(props.raw));

    const bodyJsonObject = computed(() => {
      try {
        return JSON.parse(parsedRaw.value.body);
      } catch (error) {
        return null;
      }
    });

    return {
      parsedRaw,
      bodyJsonObject,
    };
  },
};
</script>

<style lang="scss" scoped>
.ticket-status-jobs-http-data-collapse-panel__data {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 0.875rem;
}
</style>
