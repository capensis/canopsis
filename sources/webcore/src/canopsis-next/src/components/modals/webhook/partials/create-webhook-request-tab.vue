<template lang="pug">
  div
    h2 {{ $t('webhook.tabs.request.title') }}
    v-layout(justify-space-between, align-center)
      v-flex(xs6).pa-1
        v-select(
        :value="request.method",
        :items="availableMethods",
        :label="$t('webhook.tabs.request.fields.method')",
        v-validate="'required'",
        name="request.method",
        :error-messages="errors.collect('request.method')",
        @input="updateField('method', $event)"
        )
      v-flex(xs6).pa-1
        v-text-field(
        v-model="request.url",
        :label="$t('webhook.tabs.request.fields.url')",
        v-validate="'required'",
        name="request.url",
        :error-messages="errors.collect('request.url')",
        @input="updateField('url', $event)"
        )
    text-pairs(
    :items="request.headers",
    :title="$t('webhook.tabs.request.fields.headers')",
    :textLabel="$t('webhook.tabs.request.fields.headerKey')",
    :valueLabel="$t('webhook.tabs.request.fields.headerValue')",
    @input="updateField('headers', $event)"
    )
    v-layout
      v-flex
        v-textarea(
        v-model="request.payload",
        :label="$t('webhook.tabs.request.fields.payload')",
        v-validate="'required'",
        name="request.payload",
        :error-messages="errors.collect('request.payload')",
        @input="updateField('payload', $event)"
        )
</template>

<script>
import formMixin from '@/mixins/form';

import TextPairs from '@/components/forms/fields/text-pairs.vue';

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
  },
  data() {
    return {
      headers: [],
      availableMethods: [
        'POST', 'GET', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'CONNECT', 'OPTIONS', 'TRACE',
      ],
    };
  },
};
</script>
