<template lang="pug">
  v-layout(column)
    c-information-block(:title="$t('scenario.declareTicket')")
      v-layout
        v-flex(xs6)
          c-enabled-field(v-model="form.empty_response", :label="$t('scenario.emptyResponse')")
        v-flex(v-if="!form.empty_response", xs6)
          c-enabled-field(v-model="form.is_regexp", :label="$t('scenario.isRegexp')")
      c-text-pairs-field(
        v-if="!form.empty_response",
        v-field="form.mapping",
        :text-label="$t('scenario.key')",
        :name="name"
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
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'declare_ticket',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
