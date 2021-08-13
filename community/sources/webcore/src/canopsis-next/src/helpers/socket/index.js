import qs from 'qs';
import { merge } from 'lodash';

export class Socket {
  constructor(url, { query, protocols } = {}) {
    const queryString = query ? `?${qs.stringify(query)}` : '';

    this.socket = new WebSocket(`${url}${queryString}`, protocols);

    return this;
  }

  on(event, handler) {
    this.socket.addEventListener(event, handler);

    return this;
  }

  off(event, handler) {
    this.socket.removeEventListener(event, handler);

    return this;
  }

  close(code, reason) {
    this.socket.close(code, reason);

    return this;
  }

  send(data) {
    this.socket.send(data);

    return this;
  }

  static create({ baseUrl, ...restOptions }) {
    class SocketInstance extends Socket {
      constructor(url, options) {
        super(`${baseUrl}${url}`, merge(restOptions, options));
      }
    }

    return SocketInstance;
  }
}
