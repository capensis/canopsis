import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('alarm/links');

/**
 * @mixin
 */
export const entitiesAlarmLinksMixin = {
  methods: {
    ...mapActions({
      fetchAlarmLinkWithoutStore: 'fetchItemWithoutStore',
    }),
  },
};
