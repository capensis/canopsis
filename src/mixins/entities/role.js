import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('role');

export default {
  computed: {
    ...mapGetters({
      roles: 'items',
      pending: 'pending',
    }),
  },
  methods: {
    ...mapActions({
      fetchRolesList: 'fetchList',
    }),
  },
};
