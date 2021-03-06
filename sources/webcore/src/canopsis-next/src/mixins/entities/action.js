import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('action');

export default {
  computed: {
    ...mapGetters({
      actionsPending: 'pending',
      actions: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchActionsList: 'fetchList',
      refreshActionsList: 'fetchListWithPreviousParams',
      createAction: 'create',
      removeAction: 'remove',
      updateAction: 'update',
    }),
  },
};

