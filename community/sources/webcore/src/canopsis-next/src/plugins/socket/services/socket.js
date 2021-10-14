import { find } from 'lodash';

import {
  REQUEST_MESSAGES_TYPES,
  RESPONSE_MESSAGES_TYPES,
  MAX_RECONNECTS_COUNT,
  PING_INTERVAL,
  RECONNECT_INTERVAL,
  EVENTS_TYPES,
} from '../constants';

import SocketRoom from './socket-room';

class Socket {
  constructor() {
    this.rooms = {};
    this.token = '';
    this.url = '';
    this.protocols = '';
    this.sendQueue = [];
    this.reconnectsCount = 0;
    this.lastPingedAt = 0;
    this.lastPongedAt = 0;
    this.isReconnecting = true;
    this.baseOpenHandler = this.baseOpenHandler.bind(this);
    this.baseCloseHandler = this.baseCloseHandler.bind(this);
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
    this.on('close', this.baseCloseHandler);
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
    this.lastPingedAt = 0;
    this.lastPongedAt = 0;
    this.disconnect();
    this.connect(this.url, this.protocols);

    if (this.token) {
      this.authenticate();
    }

    Object.keys(this.rooms).forEach(name => this.join(name));

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
   * @returns {SocketRoom}
   */
  join(room) {
    if (!this.rooms[room]) {
      this.rooms[room] = new SocketRoom(room);
    } else {
      this.rooms[room].increment();
    }

    this.send({
      room,
      type: REQUEST_MESSAGES_TYPES.join,
    });

    return this.rooms[room];
  }

  /**
   * Leave a room
   *
   * @param {string} room
   * @returns {SocketRoom}
   */
  leave(room) {
    const socketRoom = this.rooms[room];

    if (socketRoom) {
      socketRoom.decrement();

      if (!socketRoom.count) {
        delete this.rooms[room];
      }
    }

    this.send({
      room,
      type: REQUEST_MESSAGES_TYPES.leave,
    });

    return socketRoom ?? new SocketRoom(room);
  }

  /**
   * Authenticate by token
   *
   * @param {string} [token = this.token]
   * @return {Socket}
   */
  authenticate(token = this.token) {
    this.token = token;

    this.send({
      token,
      type: REQUEST_MESSAGES_TYPES.authenticate,
    });

    return this;
  }

  /**
   * Send ping message
   *
   * @return {Socket}
   */
  ping() {
    this.send({
      type: REQUEST_MESSAGES_TYPES.ping,
    });

    this.lastPingedAt = Date.now();

    return this;
  }

  /**
   * Start custom ping mechanism
   */
  startPinging() {
    if (!this.isConnectionOpen) {
      return;
    }

    if (this.lastPingedAt > this.lastPongedAt) {
      this.connection.close();

      /**
       * We need to use this code block to avoid problem with a long waiting for connection closing
       * without internet on Google Chrome on Linux
       */
      const closeEvent = new CloseEvent(EVENTS_TYPES.customClose, { code: 1006, reason: '', wasClean: false });

      this.connection.dispatchEvent(closeEvent);

      return;
    }

    this.ping();

    setTimeout(() => this.startPinging(), PING_INTERVAL);
  }

  /**
   * Start reconnecting mechanism
   */
  startReconnecting() {
    if (this.isConnectionOpen) {
      return;
    }

    if (this.reconnectsCount >= MAX_RECONNECTS_COUNT) {
      throw new Error('Network problem');
    }

    this.reconnect();

    setTimeout(() => this.startReconnecting(), RECONNECT_INTERVAL);
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
    this.startPinging();
  }

  /**
   * Base handler for 'close' event
   */
  baseCloseHandler() {
    if (this.reconnectsCount) {
      return;
    }

    this.startReconnecting();
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

    if (this.reconnectsCount) {
      return;
    }

    this.startReconnecting();
  }

  /**
   * Base handler for 'message' event
   *
   * @param {string} data
   */
  baseMessageHandler({ data }) {
    const { type, room, msg, error } = JSON.parse(data);

    if (type === RESPONSE_MESSAGES_TYPES.error) {
      const event = new ErrorEvent('error', { message: error });

      this.connection.dispatchEvent(event);
      return;
    }

    if (type === RESPONSE_MESSAGES_TYPES.pong) {
      this.lastPongedAt = Date.now();
      return;
    }

    switch (type) {
      case RESPONSE_MESSAGES_TYPES.pong:
        this.lastPongedAt = Date.now();
        break;
      case RESPONSE_MESSAGES_TYPES.ok:
        // eslint-disable-next-line no-unused-expressions
        this.rooms[room]?.call(null, msg);
        break;
      case RESPONSE_MESSAGES_TYPES.error:
        this.connection.dispatchEvent(
          new ErrorEvent('error', { message: error }),
        );
        break;
      case RESPONSE_MESSAGES_TYPES.close:
        this.connection.dispatchEvent(
          new Event(EVENTS_TYPES.closeRoom),
        );
        break;
      default:
        this.connection.dispatchEvent(
          new ErrorEvent('error', { message: 'Unknown message type' }),
        );
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

Socket.EVENTS_TYPES = EVENTS_TYPES;

export default Socket;
