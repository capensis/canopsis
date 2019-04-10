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
      isLDAPAuthentEnabled: 'isLDAPAuthentEnabled',
      isCASAuthentEnabled: 'isCASAuthentEnabled',
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
