<template lang="pug">
  v-layout(column)
    request-form(v-field="webhook", :name="name")
    h4.ml-1 {{ $t('scenario.declareTicket') }}
    c-enabled-field(v-model="webhook.empty_response", :label="$t('declareTicket.emptyResponse')")
    c-enabled-field(v-model="webhook.is_regexp", :label="$t('declareTicket.isRegexp')")
    c-text-pairs-field(
      v-field="webhook.declare_ticket",
      :text-label="$t('scenario.key')",
      :name="name",
      text-required
    )
</template>

<script>
import { formMixin } from '@/mixins/form';

import RequestForm from '@/components/forms/request/request-form.vue';

export default {
  inject: ['$validator'],
  components: { RequestForm },
  mixins: [formMixin],
  model: {
    prop: 'webhook',
    event: 'input',
  },
  props: {
    webhook: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
  },
};
</script>
