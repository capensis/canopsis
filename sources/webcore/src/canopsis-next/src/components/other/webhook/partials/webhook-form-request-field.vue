<template lang="pug">
  v-card.pa-4.mt-2
    v-layout(justify-space-between, align-center)
      v-flex.pa-1(xs6)
        v-select(
          v-field="request.method",
          v-validate="'required'",
          :items="availableMethods",
          :disabled="disabled",
          :label="$t('webhook.tabs.request.fields.method')",
          :error-messages="errors.collect(getFieldName('method'))",
          :name="getFieldName('method')"
        )
      v-flex.pa-1(xs6)
        v-text-field(
          v-field="request.url",
          v-validate="'required'",
          :disabled="disabled",
          :label="$t('webhook.tabs.request.fields.url')",
          :error-messages="errors.collect(getFieldName('url'))",
          :name="getFieldName('url')"
        )
    v-layout(row, wrap)
      v-flex(xs12)
        v-switch(
          v-model="request.withAuth",
          :label="$t('webhook.tabs.request.fields.authSwitch')",
          :disabled="disabled",
          color="primary"
        )
      template(v-if="request.withAuth")
        v-flex(xs12)
          h4.ml-1 {{ $t('webhook.tabs.request.fields.auth') }}
        v-flex.pa-1(xs6)
          v-text-field(
            v-field="request.auth.username",
            :disabled="disabled",
            :label="$t('webhook.tabs.request.fields.username')"
          )
        v-flex.pa-1(xs6)
          v-text-field(
            v-field="request.auth.password",
            :disabled="disabled",
            :label="$t('webhook.tabs.request.fields.password')"
          )
    text-pairs(
      v-field="request.headers",
      :disabled="disabled",
      :name="getFieldName('headers')",
      :title="$t('webhook.tabs.request.fields.headers')",
      :textLabel="$t('webhook.tabs.request.fields.headerKey')",
      :valueLabel="$t('webhook.tabs.request.fields.headerValue')",
      valueValidationRules="required"
    )
    v-layout
      v-flex
        v-textarea(
          v-field="request.payload",
          v-validate="'required'",
          :disabled="disabled",
          :read-only="disabled",
          :label="$t('webhook.tabs.request.fields.payload')",
          :error-messages="errors.collect(getFieldName('payload'))",
          :name="getFieldName('payload')"
        )
</template>

<script>
import formMixin from '@/mixins/form';

import TextPairs from '@/components/forms/fields/text-pairs.vue';
import { AVAILABLE_REQUEST_METHODS } from '@/constants';

export default {
  inject: ['$validator'],
  components: { TextPairs },
  mixins: [formMixin],
  model: {
    prop: 'request',
    event: 'input',
  },
  props: {
    request: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'request',
    },
  },
  data() {
    return {
      headers: [],
    };
  },
  computed: {
    availableMethods() {
      return Object.values(AVAILABLE_REQUEST_METHODS);
    },
  },
  methods: {
    getFieldName(name) {
      return `${this.name}.${name}`;
    },
  },
};
</script>
