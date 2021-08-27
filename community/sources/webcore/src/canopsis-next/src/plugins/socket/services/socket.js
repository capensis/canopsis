import { SOCKET_EVENTS_TYPES } from '../constants';

import SocketRoom from './socket-room';

const MAX_RECONNECTS_COUNT = 10;

class Socket {
  constructor() {
    this.rooms = {};
    this.url = '';
    this.protocols = '';
    this.sendQueue = [];
    this.reconnectsCount = 0;
    this.openHandler = this.openHandler.bind(this);
    this.errorHandler = this.errorHandler.bind(this);
    this.messageHandler = this.messageHandler.bind(this);
  }

  get isConnectionOpen() {
    return this.connection && this.connection.readyState === WebSocket.OPEN;
  }

  connect(url, protocols) {
    this.url = url;
    this.protocols = protocols;
    this.connection = new WebSocket(url, protocols);

    this.on('open', this.openHandler);
    this.on('error', this.errorHandler);
    this.on('message', this.messageHandler);

    return this;
  }

  reconnect() {
    this.reconnectsCount += 1;
    this.disconnect();
    this.connect(this.url, this.protocols);
    this.rooms.forEach(({ name }) => this.join(name));
  }

  disconnect(code, reason) {
    if (this.connection) {
      this.connection.close(code, reason);
      this.connection = undefined;
    }

    return this;
  }

  openHandler() {
    this.reconnectsCount = 0;

    if (this.sendQueue.length) {
      this.sendQueue.forEach(data => this.send(data));
    }

    this.sendQueue = [];
  }

  errorHandler(err) {
    if (this.reconnectsCount >= MAX_RECONNECTS_COUNT) {
      throw new Error(err);
    }

    this.reconnect();
  }

  messageHandler(data) {
    try {
      const { room, msg } = JSON.parse(data);

      if (this.rooms[room]) {
        this.rooms[room].call(null, msg);
      }
    } catch (err) {
      console.error(err);
    }
  }

  join(room) {
    this.send({
      room,
      type: SOCKET_EVENTS_TYPES.join,
    });

    if (!this.rooms[room]) {
      this.rooms[room] = new SocketRoom(room);
    } else {
      this.rooms[room].increment();
    }

    return this.rooms[room];
  }

  leave(room) {
    this.send({
      room,
      type: SOCKET_EVENTS_TYPES.leave,
    });

    if (this.rooms[room]) {
      this.rooms[room].decrement();

      if (!this.rooms[room].count) {
        delete this.rooms[room];
      }
    }

    return this;
  }

  on(eventType, listener) {
    this.connection.addEventListener(eventType, listener);

    return this;
  }

  off(eventType, listener) {
    this.connection.removeEventListener(eventType, listener);

    return this;
  }

  room(name) {
    return this.rooms[name];
  }

  send(data) {
    if (!this.isConnectionOpen) {
      this.sendQueue.push(data);

      return this;
    }

    try {
      this.connection.send(JSON.stringify(data));
    } catch (err) {
      console.error(err);
    }

    return this;
  }
}

export default Socket;
