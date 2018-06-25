import { createNamespacedHelpers } from 'vuex';

import viewMixin from './view';

const { mapGetters, mapActions } = createNamespacedHelpers('view/widget');

export default {
  mixins: [
    viewMixin,
  ],
  computed: {
    ...mapGetters({
      getWidget: 'getItem',
    }),
  },
  methods: {
    ...mapActions({
      saveWidget: 'saveItem',
    }),
  },
};
