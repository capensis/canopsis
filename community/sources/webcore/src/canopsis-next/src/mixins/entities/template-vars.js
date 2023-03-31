import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('templateVars');

export const entitiesTemplateVarsMixin = {
  computed: {
    ...mapGetters({
      templateVars: 'items',
      templateVarsPending: 'pending',
    }),
  },
  methods: {
    ...mapActions({
      fetchTemplateVars: 'fetchList',
    }),
  },
};
