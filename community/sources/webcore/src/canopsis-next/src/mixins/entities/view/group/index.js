import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

const { mapGetters, mapActions } = createNamespacedHelpers('view');
const { mapActions: mapAuthActions } = createNamespacedHelpers('auth');

/**
 * @mixin Helpers for the view group entity
 */
export const entitiesViewGroupMixin = {
  computed: {
    ...mapGetters({
      groupsPending: 'pending',
      groups: 'items',
      getGroupById: 'getGroupById',
      getViewById: 'getViewById',
      getViewTabById: 'getViewTabById',
    }),
  },
  methods: {
    ...mapActions({
      fetchGroupsList: 'fetchList',
      fetchGroupsListWithoutStore: 'fetchListWithoutStore',
      createGroup: 'create',
      createPrivateGroup: 'createPrivateGroup',
      updateGroup: 'update',
      updatePrivateGroup: 'updatePrivateGroup',
      removeGroup: 'remove',
      removePrivateGroup: 'removePrivateGroup',
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
          with_private: true,
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
