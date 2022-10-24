import {
  REQUEST_MESSAGES_TYPES,
  RESPONSE_MESSAGES_TYPES,
  MAX_RECONNECTS_COUNT,
  PING_INTERVAL,
  RECONNECT_INTERVAL,
  EVENTS_TYPES,
  ERROR_MESSAGES,
} from '../constants';

import SocketRoom from './socket-room';

class Socket {
  constructor() {
    this.rooms = {};
    this.token = '';
    this.url = '';
    this.protocols = '';
    this.messagesToSend = [];
    this.listeners = {};
    this.reconnecting = false;
    this.reconnectsCount = 0;
    this.lastPingedAt = 0;
    this.lastPongedAt = 0;
    this.authenticated = false;
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
    return this.connection?.readyState === WebSocket.OPEN;
  }

  /**
   * Getter for checking if connection is in connecting status
   *
   * @returns {*|boolean}
   */
  get isConnecting() {
    return this.connection?.readyState === WebSocket.CONNECTING;
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
    this.connection.addEventListener('open', this.baseOpenHandler);
    this.connection.addEventListener('close', this.baseCloseHandler);
    this.connection.addEventListener('error', this.baseErrorHandler);
    this.connection.addEventListener('message', this.baseMessageHandler);

    Object.entries(this.listeners).map(([event, listeners = []]) => (
      listeners.forEach(listener => this.connection.addEventListener(event, listener))
    ));

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

    if (this.connection) {
      this.connection.close();
      this.connection = undefined;
    }

    this.connect(this.url, this.protocols);
    this.joinToActiveRooms();

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
    this.listeners = [];

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

    if (!this.listeners[eventType]) {
      this.listeners[eventType] = [];
    }

    if (!this.listeners[eventType].includes(listener)) {
      this.listeners[eventType].push(listener);
    }

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

    if (this.listeners[eventType]) {
      this.listeners[eventType] = this.listeners[eventType].filter(item => item !== listener);
    }

    return this;
  }

  /**
   * Add new message to messagesToSend array
   *
   * @param {Object} message
   * @returns {Socket}
   */
  addMessageToSend(message) {
    const messageExists = this.messagesToSend
      .find(({ data, authNeeded }) => message.data === data && message.authNeeded === authNeeded);

    if (messageExists) {
      return this;
    }

    this.messagesToSend.push(message);

    return this;
  }

  /**
   * Send data to connection
   *
   * @param {Object} data
   * @param {boolean} [authNeeded = false]
   * @returns {Socket}
   */
  send(data, authNeeded = false) {
    if (!this.isConnectionOpen || (authNeeded && !this.authenticated)) {
      this.addMessageToSend({
        data,
        authNeeded,
      });

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
   * @param {boolean} [authNeeded = true]
   * @returns {SocketRoom}
   */
  join(room, authNeeded = true) {
    if (!this.rooms[room]) {
      this.rooms[room] = new SocketRoom(room, authNeeded);
    } else {
      this.rooms[room].increment();
    }

    this.send({
      room,
      type: REQUEST_MESSAGES_TYPES.join,
    }, authNeeded);

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
    if (this.authenticated && token === this.token) {
      return this;
    }

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

    this.pingTimeout = setTimeout(() => this.startPinging(), PING_INTERVAL);
  }

  /**
   * Stop custom ping mechanism
   */
  stopPinging() {
    clearTimeout(this.pingTimeout);
  }

  /**
   * Start reconnecting mechanism
   */
  startReconnecting() {
    if (this.isConnectionOpen || this.isConnecting) {
      return;
    }

    this.reconnecting = true;

    if (this.reconnectsCount >= MAX_RECONNECTS_COUNT) {
      const errorEvent = new ErrorEvent(EVENTS_TYPES.networkError, {
        message: `Amount of reconnecting hit the limit of ${MAX_RECONNECTS_COUNT}`,
      });

      this.connection.dispatchEvent(errorEvent);

      return;
    }

    this.reconnectingTimeout = setTimeout(() => {
      this.reconnect();
      this.startReconnecting();
    }, RECONNECT_INTERVAL);
  }

  /**
   * Joint to active rooms after reconnection
   */
  joinToActiveRooms() {
    Object.entries(this.rooms).forEach(([name, room]) => {
      room.decrement();

      this.join(name, room.authNeeded);
    });

    return this;
  }

  /**
   * Send messages which was added to messagesToSend array while connection was closed
   *
   * @returns {Socket}
   */
  sendMessagesToSend() {
    if (!this.messagesToSend.length) {
      return this;
    }

    const { toSend, toDelay } = this.messagesToSend.reduce((acc, message) => {
      if (message.data.type === REQUEST_MESSAGES_TYPES.leave) {
        return acc;
      }

      if (this.authenticated || !message.authNeeded) {
        acc.toSend.push(message);
      } else {
        acc.toDelay.push(message);
      }

      return acc;
    }, {
      toSend: [],
      toDelay: [],
    });

    toSend.forEach(({ data, authNeeded }) => this.send(data, authNeeded));

    this.messagesToSend = toDelay;

    return this;
  }

  /**
   * Stop reconnecting mechanism
   */
  stopReconnecting() {
    this.reconnecting = false;
    this.reconnectsCount = 0;

    clearTimeout(this.reconnectingTimeout);
  }

  /**
   * Base handler for 'open' event
   */
  baseOpenHandler() {
    this.stopReconnecting();

    if (this.token) {
      this.authenticate();
    }

    this.sendMessagesToSend();
    this.stopPinging();
    this.startPinging();
  }

  /**
   * Base handler for 'close' event
   */
  baseCloseHandler() {
    this.authenticated = false;

    if (this.reconnecting) {
      return;
    }

    this.stopReconnecting();
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

    this.authenticated = false;

    if (this.reconnecting) {
      return;
    }

    this.stopReconnecting();
    this.startReconnecting();
  }

  /**
   * Base handler for 'message' event
   *
   * @param {string} data
   */
  baseMessageHandler({ data }) {
    const { type, room, msg, error } = JSON.parse(data);

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
      case RESPONSE_MESSAGES_TYPES.authenticated:
        this.authenticated = true;
        this.sendMessagesToSend();
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
Socket.ERROR_MESSAGES = ERROR_MESSAGES;

export default Socket;
