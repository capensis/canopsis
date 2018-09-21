import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('view/row');

/**
 * @mixin Helpers for the view row entity
 */
export default {
  methods: {
    ...mapActions({
      createRowInStore: 'createInStore',
      createRow: 'create',
      updateRow: 'update',
    }),
  },
};
