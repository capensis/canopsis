import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('action');

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
      fetchActionsListWithoutStore: 'fetchListWithoutStore',
    }),
  },
};
