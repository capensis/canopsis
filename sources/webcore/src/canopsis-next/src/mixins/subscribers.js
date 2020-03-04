import { isFunction } from 'lodash';

export default {
  data() {
    return {
      subscribers: [],
    };
  },
  methods: {
    subscribe(callback) {
      if (isFunction(callback)) {
        this.subscribers.push(callback);
      }
    },

    unsubscribe(callback) {
      if (isFunction(callback)) {
        this.subscribers.filter(subscriber => callback !== subscriber);
      }
    },

    async notify() {
      await Promise.all(this.subscribers.map(subscriber => subscriber()));
    },
  },
};
