import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('linkRule');

export const entitiesLinkRuleMixin = {
  computed: {
    ...mapGetters({
      linkRulesMeta: 'meta',
      linkRulesPending: 'pending',
      linkRules: 'items',
    }),
  },
  methods: {
    ...mapActions({
      fetchLinkRulesList: 'fetchList',
      createLinkRule: 'create',
      updateLinkRule: 'update',
      removeLinkRule: 'remove',
      bulkRemoveLinkRules: 'bulkRemove',
    }),
  },
};
