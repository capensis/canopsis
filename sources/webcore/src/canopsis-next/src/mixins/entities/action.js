import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('action');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      actions: 'items',
      actionsPending: 'pending',
    }),
  },
  methods: {
    ...mapActions({
      fetchActionsList: 'fetchList',
      removeAction: 'remove',
    }),
  },
};
