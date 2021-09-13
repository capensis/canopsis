import { find } from 'lodash';

import { SOCKET_EVENTS_TYPES, MAX_RECONNECTS_COUNT } from '../constants';

import SocketRoom from './socket-room';

class Socket {
  constructor() {
    this.rooms = {};
    this.url = '';
    this.protocols = '';
    this.sendQueue = [];
    this.reconnectsCount = 0;
    this.baseOpenHandler = this.baseOpenHandler.bind(this);
    this.baseErrorHandler = this.baseErrorHandler.bind(this);
    this.baseMessageHandler = this.baseMessageHandler.bind(this);
  }

  /**
   * Getter for checking if connection is open
   *
   * @returns {*|boolean}
   */
  get isConnectionOpen() {
    return this.connection && this.connection.readyState === WebSocket.OPEN;
  }

  /**
   * Create websocket connection and add base event listeners
   *
   * @param {string} url
   * @param {string | string[]} [protocols]
   * @returns {Socket}
   */
  connect(url, protocols) {
    this.url = url;
    this.protocols = protocols;
    this.connection = new WebSocket(url, protocols);

    this.on('open', this.baseOpenHandler);
    this.on('error', this.baseErrorHandler);
    this.on('message', this.baseMessageHandler);

    return this;
  }

  /**
   * Reconnect and update reconnectsCount
   *
   * @returns {Socket}
   */
  reconnect() {
    this.reconnectsCount += 1;
    this.disconnect();
    this.connect(this.url, this.protocols);
    this.rooms.forEach(({ name }) => this.join(name));

    return this;
  }

  /**
   * Disconnect websocket connection
   *
   * @param {number} [code]
   * @param {string} [reason]
   * @returns {Socket}
   */
  disconnect(code, reason) {
    if (this.connection) {
      this.connection.close(code, reason);
      this.connection = undefined;
    }

    return this;
  }

  /**
   * Add event listener for connection
   *
   * @param {string} eventType
   * @param {Function} listener
   * @returns {Socket}
   */
  on(eventType, listener) {
    this.connection.addEventListener(eventType, listener);

    return this;
  }

  /**
   * Remove event listener from connection
   *
   * @param {string} eventType
   * @param {Function} listener
   * @returns {Socket}
   */
  off(eventType, listener) {
    this.connection.removeEventListener(eventType, listener);

    return this;
  }

  /**
   * Send data to connection
   *
   * @param {Object} data
   * @returns {Socket}
   */
  send(data) {
    if (!this.isConnectionOpen && !find(this.sendQueue, data)) {
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

  /**
   * Join to a room
   *
   * @param {string} room
   * @returns {Socket}
   */
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

    return this;
  }

  /**
   * Leave a room
   *
   * @param {string} room
   * @returns {Socket}
   */
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

  /**
   * Base handler for 'open' event
   */
  baseOpenHandler() {
    this.reconnectsCount = 0;

    if (this.sendQueue.length) {
      this.sendQueue.forEach(data => this.send(data));
    }

    this.sendQueue = [];
  }

  /**
   * Base handler for 'error' event
   *
   * @param {string} err
   */
  baseErrorHandler(err) {
    if (this.isConnectionOpen) {
      console.error(err);

      return;
    }

    if (this.reconnectsCount >= MAX_RECONNECTS_COUNT) {
      throw new Error(err);
    }

    this.reconnect();
  }

  /**
   * Base handler for 'message' event
   *
   * @param {string} data
   */
  baseMessageHandler({ data }) {
    try {
      const { room, msg, error } = JSON.parse(data);

      if (error) {
        throw error;
      }

      if (this.rooms[room] && msg) {
        this.rooms[room].call(null, msg);
      }
    } catch (err) {
      console.error(err);
    }
  }

  /**
   * Get socket room by name
   *
   * @param {string} name
   * @returns {SocketRoom}
   */
  getRoom(name) {
    const room = this.rooms[name];

    if (!room) {
      throw new Error('Unknown room');
    }

    return room;
  }
}

export default Socket;
