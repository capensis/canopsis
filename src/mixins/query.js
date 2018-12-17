import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('query');

export default {
  computed: {
    ...mapGetters(['getQueryById', 'getQueryNonceById']),
  },
  methods: {
    ...mapActions({
      updateQuery: 'update',
      mergeQuery: 'merge',
      removeQuery: 'remove',

      forceUpdateQuery: 'forceUpdate',
    }),
  },
};
