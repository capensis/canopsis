import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('view/widget');

/**
 * @mixin Helpers for the view widget entity
 */
export default {
  methods: {
    ...mapActions({
      createWidget: 'create',
      updateWidget: 'update',
    }),
  },
};
