import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('info');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      version: 'version',
      logo: 'logo',
      appTitle: 'appTitle',
      footer: 'footer',
      edition: 'edition',
      stack: 'stack',
      isLDAPAuthEnabled: 'isLDAPAuthEnabled',
      isCASAuthEnabled: 'isCASAuthEnabled',
      casConfig: 'casConfig',
    }),
  },
  methods: {
    ...mapActions({
      fetchLoginInfos: 'fetchLoginInfos',
      fetchAppInfos: 'fetchAppInfos',
    }),
  },
};
