import { createNamespacedHelpers } from 'vuex';

import viewMixin from './entities/view';

const { mapGetters, mapActions } = createNamespacedHelpers('view/widget');

/**
 * @mixin Helpers for widget
 * @see src/mixins/view.js
 */
export default {
  mixins: [
    viewMixin,
  ],
  computed: {
    ...mapGetters({
      getWidget: 'getItem',
      getWidgets: 'getItems',
    }),
  },
  methods: {
    ...mapActions({
      saveWidget: 'saveItem',
    }),
  },
};
