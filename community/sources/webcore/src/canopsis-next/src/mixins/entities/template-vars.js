import { mapActions, mapState, mapStores } from 'pinia';

import { useTemplateVarsStore } from '@/stores/modules/template-vars';

export const entitiesTemplateVarsMixin = {
  computed: {
    ...mapStores(useTemplateVarsStore),
    ...mapState(useTemplateVarsStore, {
      templateVarsPending: 'pending',
      templateVars: 'items',
    }),
  },

  methods: {
    ...mapActions(useTemplateVarsStore, {
      fetchTemplateVars: 'fetchList',
    }),
  },
};

export default entitiesTemplateVarsMixin;
