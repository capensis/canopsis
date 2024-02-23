<template>
  <v-layout column>
    <v-layout>
      <v-flex xs6>
        <c-enabled-field v-field="form.enabled" />
      </v-flex>
      <v-flex xs6>
        <c-enabled-field
          v-field="form.emit_trigger"
          :label="$t('common.emitTrigger')"
          name="emit_trigger"
        />
      </v-flex>
    </v-layout>
    <v-layout class="mb-2">
      <v-flex
        class="mr-3"
        xs8
      >
        <c-name-field
          v-field="form.name"
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
    <declare-ticket-rule-webhooks-field v-field="form.webhooks" />
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
  },
};
</script>
