import { createNamespacedHelpers } from 'vuex';
import { sortBy } from 'lodash';

const { mapGetters, mapActions } = createNamespacedHelpers('view/group');

/**
 * @mixin Helpers for the view group entity
 */
export default {
  computed: {
    ...mapGetters({
      groupsPending: 'pending',
      groups: 'items',
      getGroupById: 'getItemById',
    }),

    groupsOrdered() {
      return sortBy(this.groups, ['position']);
    },
  },
  methods: {
    ...mapActions({
      fetchGroupsList: 'fetchList',
      createGroup: 'create',
      updateGroup: 'update',
      removeGroup: 'remove',
    }),
  },
};
