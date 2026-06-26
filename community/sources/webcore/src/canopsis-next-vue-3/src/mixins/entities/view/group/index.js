import { mapActions, mapState, mapStores } from 'pinia';

import { MAX_LIMIT } from '@/constants';

import { useAuthStore } from '@/stores/modules/auth';
import { useViewStore } from '@/stores/modules/view';

/**
 * @mixin Helpers for the view group entity
 */
export const entitiesViewGroupMixin = {
  computed: {
    ...mapStores(useViewStore),
    ...mapState(useViewStore, {
      groupsPending: 'pending',
      groups: 'items',
      getGroupById: 'getGroupById',
      getViewById: 'getViewById',
      getViewTabById: 'getViewTabById',
    }),
  },

  methods: {
    ...mapActions(useViewStore, {
      fetchGroupsList: 'fetchList',
      fetchGroupsListWithoutStore: 'fetchListWithoutStore',
      createGroup: 'create',
      createPrivateGroup: 'createPrivateGroup',
      updateGroup: 'update',
      updatePrivateGroup: 'updatePrivateGroup',
      removeGroup: 'remove',
      removePrivateGroup: 'removePrivateGroup',
    }),

    ...mapActions(useAuthStore, ['fetchCurrentUser']),

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
