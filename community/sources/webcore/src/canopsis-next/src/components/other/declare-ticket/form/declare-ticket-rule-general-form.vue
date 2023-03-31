<template lang="pug">
  v-layout(column)
    v-layout(row)
      v-flex(xs6)
        c-enabled-field(v-field="form.enabled")
      v-flex(xs6)
        c-enabled-field(v-field="form.emit_trigger", :label="$t('common.emitTrigger')", name="emit_trigger")
    v-layout.mb-2(row)
      v-flex.mr-3(xs8)
        c-name-field(v-field="form.name", required)
      v-flex(xs4)
        c-name-field(v-field="form.system_name", :label="$t('common.systemName')", name="system_name")
    declare-ticket-rule-webhooks-field(v-field="form.webhooks")
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
  },
};
</script>
