import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('authProtocol');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      configPending: 'configPending',
    }),
  },
  methods: {
    ...mapActions({
      fetchLDAPConfigWithoutStore: 'fetchLDAPConfigWithoutStore',
      updateLDAPConfig: 'updateLDAPConfig',
      fetchCASConfigWithoutStore: 'fetchCASConfigWithoutStore',
      updateCASConfig: 'updateCASConfig',
    }),
  },
};
