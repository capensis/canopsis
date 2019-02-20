<template lang="pug">
  div
    h2 Request
    v-layout(justify-space-between, align-center)
      v-flex(xs6).pa-1
        v-select(
        :value="request.method",
        :items="availableMethods",
        label="Method",
        v-validate="'required'",
        name="request.method",
        :error-messages="errors.collect('request.method')",
        @input="updateField('method', $event)"
        )
      v-flex(xs6).pa-1
        v-text-field(
        v-model="request.url",
        label="URL",
        v-validate="'required'",
        name="request.url",
        :error-messages="errors.collect('request.url')",
        @input="updateField('url', $event)"
        )
    text-pairs(
    :items="request.headers",
    :title="'Headers'",
    :textLabel="'Header key'",
    :valueLabel="'Header value'",
    @input="updateField('headers', $event)"
    )
    v-layout
      v-flex
        v-textarea(
        v-model="request.payload",
        label="Payload",
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
