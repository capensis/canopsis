import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('role');

export default {
  computed: {
    ...mapGetters({
      roles: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchRolesList: 'fetchList',
    }),
  },
};
