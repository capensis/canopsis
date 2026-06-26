import {
  ERROR_MESSAGES,
  EVENTS_TYPES,
  MAX_RECONNECTS_COUNT,
  PING_INTERVAL,
  RECONNECT_INTERVAL,
  REQUEST_MESSAGES_TYPES,
  RESPONSE_MESSAGES_TYPES,
} from '../constants'

import SocketRoom from './socket-room'

class Socket {
  constructor() {
    this.rooms = {}
    this.token = ''
    this.url = ''
    this.protocols = ''
    this.messagesToSend = []
    this.listeners = {}
    this.reconnecting = false
    this.reconnectsCount = 0
    this.lastPingedAt = 0
    this.lastPongedAt = 0
    this.authenticated = false
    this.baseOpenHandler = this.baseOpenHandler.bind(this)
    this.baseCloseHandler = this.baseCloseHandler.bind(this)
    this.baseErrorHandler = this.baseErrorHandler.bind(this)
    this.baseMessageHandler = this.baseMessageHandler.bind(this)
  }

  get isConnectionOpen() {
    return this.connection?.readyState === WebSocket.OPEN
  }

  get isConnecting() {
    return this.connection?.readyState === WebSocket.CONNECTING
  }

  connect(url, protocols) {
    if (this.connection && (this.isConnectionOpen || this.isConnecting)) {
      return this
    }

    this.url = url
    this.protocols = protocols
    this.connection = new WebSocket(url, protocols)
    this.connection.addEventListener('open', this.baseOpenHandler)
    this.connection.addEventListener('close', this.baseCloseHandler)
    this.connection.addEventListener('error', this.baseErrorHandler)
    this.connection.addEventListener('message', this.baseMessageHandler)

    Object.entries(this.listeners).forEach(([event, listeners = []]) => {
      listeners.forEach(listener => this.connection.addEventListener(event, listener))
    })

    return this
  }

  reconnect() {
    this.reconnectsCount += 1
    this.lastPingedAt = 0
    this.lastPongedAt = 0

    if (this.connection) {
      this.connection.close()
      this.connection = undefined
    }

    this.connect(this.url, this.protocols)
    this.joinToActiveRooms()

    return this
  }

  disconnect(code, reason) {
    this.listeners = {}

    if (this.connection) {
      this.connection.close(code, reason)
      this.connection = undefined
    }

    return this
  }

  on(eventType, listener) {
    this.connection?.addEventListener(eventType, listener)

    if (!this.listeners[eventType]) {
      this.listeners[eventType] = []
    }

    if (!this.listeners[eventType].includes(listener)) {
      this.listeners[eventType].push(listener)
    }

    return this
  }

  off(eventType, listener) {
    this.connection?.removeEventListener(eventType, listener)

    if (this.listeners[eventType]) {
      this.listeners[eventType] = this.listeners[eventType].filter(item => item !== listener)
    }

    return this
  }

  addMessageToSend(message) {
    const messageExists = this.messagesToSend.find(({ payload, authNeeded }) => (
      message.payload === payload && message.authNeeded === authNeeded
    ))

    if (!messageExists) {
      this.messagesToSend.push(message)
    }

    return this
  }

  send(payload, authNeeded = false) {
    if (!this.isConnectionOpen || (authNeeded && !this.authenticated)) {
      this.addMessageToSend({ payload, authNeeded })

      return this
    }

    try {
      this.connection.send(JSON.stringify(payload))
    } catch (err) {
      console.error(err)
    }

    return this
  }

  sendJoin(room, payload = {}, authNeeded = true) {
    this.send({
      room,
      payload,
      type: REQUEST_MESSAGES_TYPES.join,
    }, authNeeded)
  }

  join(room, payload = {}, authNeeded = true) {
    if (!this.rooms[room]) {
      this.rooms[room] = new SocketRoom(room, payload, authNeeded, this.send.bind(this))
      this.sendJoin(room, payload, authNeeded)
    } else {
      this.rooms[room].increment()
    }

    return this.rooms[room]
  }

  leave(room) {
    const socketRoom = this.rooms[room]

    if (socketRoom) {
      socketRoom.decrement()

      if (!socketRoom.count) {
        delete this.rooms[room]

        this.send({
          room,
          type: REQUEST_MESSAGES_TYPES.leave,
        })
      }
    }

    return socketRoom ?? new SocketRoom(room)
  }

  authenticate(token = this.token) {
    if (this.authenticated && token === this.token) {
      return this
    }

    this.token = token

    this.send({
      token,
      type: REQUEST_MESSAGES_TYPES.authenticate,
    })

    return this
  }

  ping() {
    this.send({
      type: REQUEST_MESSAGES_TYPES.ping,
    })

    this.lastPingedAt = Date.now()

    return this
  }

  startPinging() {
    if (!this.isConnectionOpen) {
      return
    }

    if (this.lastPingedAt > this.lastPongedAt) {
      this.connection.close()

      return
    }

    this.ping()
    this.pingTimeout = setTimeout(() => this.startPinging(), PING_INTERVAL)
  }

  stopPinging() {
    clearTimeout(this.pingTimeout)
  }

  startReconnecting() {
    if (this.isConnectionOpen || this.isConnecting) {
      return
    }

    this.reconnecting = true

    if (this.reconnectsCount >= MAX_RECONNECTS_COUNT) {
      this.connection?.dispatchEvent(new ErrorEvent(EVENTS_TYPES.networkError, {
        message: `Amount of reconnecting hit the limit of ${MAX_RECONNECTS_COUNT}`,
      }))

      return
    }

    this.reconnectingTimeout = setTimeout(() => {
      this.reconnect()
    }, RECONNECT_INTERVAL)
  }

  joinToActiveRooms() {
    Object.entries(this.rooms).forEach(([name, room]) => {
      room.decrement()
      this.sendJoin(name, room.payload, room.authNeeded)
    })

    return this
  }

  sendMessagesToSend() {
    if (!this.messagesToSend.length) {
      return this
    }

    const { toSend, toDelay } = this.messagesToSend.reduce((acc, message) => {
      if (message.payload.type === REQUEST_MESSAGES_TYPES.leave) {
        return acc
      }

      if (this.authenticated || !message.authNeeded) {
        acc.toSend.push(message)
      } else {
        acc.toDelay.push(message)
      }

      return acc
    }, {
      toSend: [],
      toDelay: [],
    })

    toSend.forEach(({ payload, authNeeded }) => this.send(payload, authNeeded))
    this.messagesToSend = toDelay

    return this
  }

  stopReconnecting() {
    this.reconnecting = false
    clearTimeout(this.reconnectingTimeout)
  }

  baseOpenHandler() {
    this.reconnectsCount = 0
    this.stopReconnecting()

    if (this.token) {
      this.authenticate()
    }

    this.sendMessagesToSend()
    this.stopPinging()
    this.startPinging()
  }

  baseCloseHandler() {
    this.authenticated = false

    if (this.reconnecting) {
      return
    }

    this.stopReconnecting()
    this.startReconnecting()
  }

  baseErrorHandler(err) {
    if (this.isConnectionOpen) {
      console.error(err)

      return
    }

    this.authenticated = false
    this.stopReconnecting()
    this.startReconnecting()
  }

  baseMessageHandler({ data }) {
    const { type, room, payload, error } = JSON.parse(data)

    switch (type) {
      case RESPONSE_MESSAGES_TYPES.pong:
        this.lastPongedAt = Date.now()
        break
      case RESPONSE_MESSAGES_TYPES.ok:
        this.rooms[room]?.call(null, payload)
        break
      case RESPONSE_MESSAGES_TYPES.error:
        this.connection.dispatchEvent(
          new ErrorEvent('error', { message: error, error: { code: error, room, payload } }),
        )
        break
      case RESPONSE_MESSAGES_TYPES.close:
        this.connection.dispatchEvent(new Event(EVENTS_TYPES.closeRoom))
        break
      case RESPONSE_MESSAGES_TYPES.authenticated:
        this.authenticated = true
        this.sendMessagesToSend()
        break
      case RESPONSE_MESSAGES_TYPES.joined:
        this.rooms[room]?.setJoined?.(true)
        break
      case RESPONSE_MESSAGES_TYPES.left:
        this.rooms[room]?.setJoined?.(false)
        break
      default:
        this.connection.dispatchEvent(new ErrorEvent('error', { message: 'Unknown message type' }))
    }
  }

  getRoom(name) {
    const room = this.rooms[name]

    if (!room) {
      throw new Error('Unknown room')
    }

    return room
  }
}

Socket.EVENTS_TYPES = EVENTS_TYPES
Socket.ERROR_MESSAGES = ERROR_MESSAGES

export default Socket
