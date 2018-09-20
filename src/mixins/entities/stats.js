import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('stats');

export default {
  computed: {
    ...mapGetters({
      statsList: 'getList',
      statItem: 'getItem',
    }),
  },
  methods: {
    ...mapActions({
      fetchStats: 'fetchStats',
    }),
  },
};
