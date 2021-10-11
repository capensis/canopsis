import Socket from './services/socket';

export default {
  install(Vue) {
    const socket = new Socket();

    Object.defineProperty(Vue.prototype, '$socket', {
      get() {
        return socket;
      },
    });
  },
};
