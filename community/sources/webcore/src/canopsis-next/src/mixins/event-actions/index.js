import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('event');

/**
 * @mixin
 */
export const eventActionsMixin = {
  methods: {
    ...mapActions({
      createEventAction: 'create',
    }),
  },
};
