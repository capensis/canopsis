import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('eventFilter');

export const entitiesEventFilterMixin = {
  computed: {
    ...mapGetters({
      eventFiltersPending: 'pending',
      eventFilters: 'items',
      eventFiltersMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchEventFiltersList: 'fetchList',
      refreshEventFiltersList: 'fetchListWithPreviousParams',
      fetchEventFilterErrorsListWithoutStore: 'fetchEventFilterErrorsListWithoutStore',
      markNewEventFilterErrorsAsRead: 'markNewEventFilterErrorsAsRead',
      createEventFilter: 'create',
      updateEventFilter: 'update',
      removeEventFilter: 'remove',
    }),
  },
};
