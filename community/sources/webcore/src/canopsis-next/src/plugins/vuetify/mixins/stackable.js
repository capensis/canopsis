import { getMaxZIndex } from '@/helpers/vuetify';

export default {
  methods: {
    getMaxZIndex(exclude = []) {
      return getMaxZIndex(this.$el, 300, exclude);
    },
  },
};
