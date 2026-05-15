<template>
  <v-layout class="gap-5" column>
    <v-layout class="gap-4">
      <c-name-field
        v-field="form.name"
        autofocus
        required
      />
      <c-name-field
        v-field="form.system_name"
        :label="$t('common.systemName')"
        name="system_name"
      />
    </v-layout>

    <c-enabled-field
      v-field="form.emit_trigger"
      :label="$t('common.emitTrigger')"
      hide-details
      no-margin
    />

    <declare-ticket-rule-webhooks-field v-field="form.webhooks" :template-vars="templateVars" />
  </v-layout>
</template>

<script>
import DeclareTicketRuleWebhooksField from './fields/declare-ticket-rule-webhooks-field.vue';

export default {
  inject: ['$validator'],
  components: { DeclareTicketRuleWebhooksField },
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
