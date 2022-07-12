import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

const { mapGetters, mapActions } = createNamespacedHelpers('view/group');
const { mapActions: mapAuthActions } = createNamespacedHelpers('auth');

/**
 * @mixin Helpers for the view group entity
 */
export const entitiesViewGroupMixin = {
  computed: {
    ...mapGetters({
      groupsPending: 'pending',
      groups: 'items',
      getGroupById: 'getItemById',
    }),
  },
  methods: {
    ...mapActions({
      fetchGroupsList: 'fetchList',
      fetchGroupsListWithoutStore: 'fetchListWithoutStore',
      createGroup: 'create',
      updateGroup: 'update',
      removeGroup: 'remove',
    }),

    ...mapAuthActions(['fetchCurrentUser']),

    fetchAllGroupsListWithWidgets() {
      return this.fetchGroupsList({
        params: {
          limit: MAX_LIMIT,
          page: 1,
          with_views: true,
          with_tabs: true,
          with_widgets: true,
          with_flags: true,
        },
      });
    },

    fetchAllGroupsListWithWidgetsWithCurrentUser() {
      return Promise.all([
        this.fetchAllGroupsListWithWidgets(),
        this.fetchCurrentUser(),
      ]);
    },
  },
};
