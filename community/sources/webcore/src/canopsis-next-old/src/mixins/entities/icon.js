import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('icon');

export const entitiesIconMixin = {
  computed: {
    ...mapGetters({
      icons: 'items',
      iconsPending: 'pending',
      iconsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchIconsList: 'fetchList',
      fetchIconsListWithPreviousParams: 'fetchListWithPreviousParams',
      createIcon: 'create',
      updateIcon: 'update',
      removeIcon: 'remove',
    }),
  },
};
