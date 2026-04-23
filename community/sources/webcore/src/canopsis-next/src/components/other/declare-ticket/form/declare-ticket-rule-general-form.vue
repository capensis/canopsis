<template>
  <v-layout class="gap-2" column>
    <v-layout class="gap-3">
      <v-flex xs8>
        <c-name-field
          v-field="form.name"
          autofocus
          required
        />
      </v-flex>
      <v-flex xs4>
        <c-name-field
          v-field="form.system_name"
          :label="$t('common.systemName')"
          name="system_name"
        />
      </v-flex>
    </v-layout>
    <c-enabled-field
      v-field="form.emit_trigger"
      :label="$t('common.emitTrigger')"
      name="emit_trigger"
    />
    <declare-ticket-rule-webhooks-field v-field="form.webhooks" :template-vars="templateVars" />
  </v-layout>
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

import DeclareTicketRuleWebhooksField from './fields/declare-ticket-rule-webhooks-field.vue';

export default {
  inject: ['$validator'],
  components: { DeclareTicketRuleWebhooksField },
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
