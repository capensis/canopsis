import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('pbehavior');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      pbehaviorItems: 'items',
    }),
  },
  methods: {
    ...mapActions({
      createPbehavior: 'create',
      removePbehavior: 'remove',
      fetchPbehaviorsByEntityId: 'fetchListByEntityId',
    }),
  },
};
