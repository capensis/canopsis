import { get } from 'lodash';

export default {
  props: {
    config: {
      type: Object,
      default: () => ({}),
    },
  },
  provide() {
    return {
      $dayspanOptions: {
        options: this.config,
        getOptions(name) {
          return get(this.options, name, {});
        },
      },
    };
  },
};
