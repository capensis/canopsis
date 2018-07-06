import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('context');

/**
 * @mixin Helpers' for context store
 */
export default {
  computed: {
    ...mapGetters(['pending']),
    ...mapGetters({
      contextEntities: 'items',
      contextEntitiesMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchContextEntities: 'fetchList',
    }),
  },
};
