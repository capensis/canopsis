import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('template/vars');

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
