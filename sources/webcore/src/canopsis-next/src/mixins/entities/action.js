import { createNamespacedHelpers } from 'vuex';

import popupMixin from '@/mixins/popup';

const { mapActions, mapGetters } = createNamespacedHelpers('action');

export default {
  mixins: [popupMixin],
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

