import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('view/group');

/**
 * @mixin Helpers for the view group entity
 */
export default {
  computed: {
    ...mapGetters({
      groups: 'items',
    }),

    groupsOrdered() {
      return this.groups.sort((a = {}, b = {}) => (a.position || Infinity) - (b.position || Infinity));
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
