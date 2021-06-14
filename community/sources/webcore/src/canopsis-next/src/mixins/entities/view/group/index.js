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
      bulkCreateGroupsWithoutStore: 'bulkCreateWithoutStore',
    }),

    ...mapAuthActions(['fetchCurrentUser']),

    fetchAllGroupsListWithViews() {
      return this.fetchGroupsList({
        params: {
          limit: MAX_LIMIT,
          page: 1,
          with_views: true,
          with_flags: true,
        },
      });
    },

    fetchAllGroupsListWithViewsWithCurrentUser() {
      return Promise.all([
        this.fetchGroupsList({
          params: {
            limit: MAX_LIMIT,
            page: 1,
            with_views: true,
            with_flags: true,
          },
        }),
        this.fetchCurrentUser(),
      ]);
    },
  },
};
