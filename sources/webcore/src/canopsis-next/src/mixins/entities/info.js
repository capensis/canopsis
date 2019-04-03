import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('info');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      version: 'version',
    }),
  },
  methods: {
    ...mapActions({
      fetchLoginInfos: 'fetchLoginInfos',
    }),
  },
};
