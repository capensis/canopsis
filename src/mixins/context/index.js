import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('entity');

/**
 * @mixin Helpers' for context store
 */
export default {
  computed: {
    ...mapGetters({
      pending: 'pending',
      contextEntities: 'items',
      contextEntitiesMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchContextEntities: 'fetchList',
      remove: 'remove',
      updateContextEntity: 'update',
    }),
  },
};
