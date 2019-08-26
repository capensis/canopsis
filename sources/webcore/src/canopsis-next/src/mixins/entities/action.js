import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('action');

/**
 * @mixin
 */
export default {
  computed: {},
  methods: {
    ...mapActions({
      fetchActionsList: 'fetchList',
    }),
  },
};
