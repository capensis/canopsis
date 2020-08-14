import Vue from 'vue';

export default {
  provide() {
    return {
      $system: this.system,
    };
  },
  data() {
    return {
      system: {
        timezone: this.timezone,
      },
    };
  },
  methods: {
    setSystemData(options) {
      Object.entries(options).forEach(([key, value]) => {
        Vue.set(this.system, key, value);
      });
    },
  },
};
