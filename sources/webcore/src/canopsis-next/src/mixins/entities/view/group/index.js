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
  },
  methods: {
    ...mapActions({
      fetchGroupsList: 'fetchList',
      createGroup: 'create',
    }),
  },
};
