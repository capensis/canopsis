import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('authProtocol');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      ldapConfigPending: 'ldapConfigPending',
      ldapConfig: 'ldapConfig',
    }),
  },
  methods: {
    ...mapActions({
      fetchLDAPConfigWithoutStore: 'fetchLDAPConfigWithoutStore',
    }),
  },
};
