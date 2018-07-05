import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('context');

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
