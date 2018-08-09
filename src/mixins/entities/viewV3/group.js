import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('viewV3/group');

/**
 * @mixin Helpers for the view entity
 */
export default {
  methods: {
    ...mapActions({
      fetchGroupList: 'fetchList',
      createGroup: 'create',
    }),
  },
  computed: {
    ...mapGetters({
      groupList: 'items',
    }),
  },
};
