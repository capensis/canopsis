import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('event');

/**
 * @mixin
 */
export default {
  methods: {
    ...mapActions({
      createEventAction: 'create',
    }),
  },
};
