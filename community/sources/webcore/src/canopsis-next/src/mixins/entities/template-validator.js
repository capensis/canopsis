import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('templateValidator');

export const entitiesTemplateValidatorMixin = {
  methods: {
    ...mapActions({
      validateDeclareTicketRulesVariables: 'validateDeclareTicketRulesVariables',
      validateScenariosVariables: 'validateScenariosVariables',
    }),
  },
};
