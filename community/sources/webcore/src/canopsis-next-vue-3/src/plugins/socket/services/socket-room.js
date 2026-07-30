import { REQUEST_MESSAGES_TYPES } from '@/plugins/socket/constants';

class SocketRoom {
  constructor(name, payload, authNeeded, send) {
    this.name = name;
    this.payload = payload;
    this.authNeeded = authNeeded;
    this.count = 1;
    this.listeners = [];
    this.parentSend = send;
    this.joined = false;
    this.messagesToSend = [];
  }

  addListener(listener) {
    this.listeners.push(listener);

    return this;
  }

  removeListener(listener) {
    this.listeners = this.listeners.filter(item => item !== listener);

    return this;
  }

  increment() {
    this.count += 1;

    return this;
  }

  decrement() {
    this.count -= 1;

    return this;
  }

  isEmpty() {
    return this.count <= 0;
  }

  call(context, ...restArgs) {
    this.listeners.forEach(listener => listener.call(context, ...restArgs));

    return this;
  }

  send(payload) {
    const message = {
      type: REQUEST_MESSAGES_TYPES.send,
      room: this.name,
      payload,
    };

    if (!this.joined) {
      this.messagesToSend.push(message);

      return this;
    }

    this.parentSend?.(message);

    return this;
  }

  setJoined(joined) {
    this.joined = joined;

    if (joined) {
      this.messagesToSend.forEach(message => this.parentSend?.(message));
      this.messagesToSend = [];
    }
  }

  updatePayloadField(field, value) {
    if (!this.payload) {
      this.payload = {};
    }

    this.payload[field] = value;
  }
}

export default SocketRoom;
