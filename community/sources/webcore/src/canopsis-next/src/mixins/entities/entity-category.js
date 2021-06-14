import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('entityCategory');

export default {
  computed: {
    ...mapGetters({
      entityCategoriesPending: 'pending',
      entityCategories: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchEntityCategoriesList: 'fetchList',
      updateEntityCategory: 'update',
      createEntityCategory: 'create',
      removeEntityCategory: 'remove',
    }),
  },
};
