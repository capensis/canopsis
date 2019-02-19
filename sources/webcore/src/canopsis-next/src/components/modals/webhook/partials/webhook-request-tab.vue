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
    v-layout(
    v-for="(header, index) in request.headers",
    :key="header.id",
    justify-space-between,
    align-center
    )
      v-flex.pa-1
        v-text-field(
        :value="header.key",
        label="Header key",
        :name="`headers[${index}].key`",
        :error-messages="errors.collect(`headers[${index}].key`)"
        v-validate="'required'",
        @input="updateField(`headers[${index}].key`, $event)"
        )
      v-flex.pa-1
        v-text-field(
        :value="header.value",
        label="Header value",
        :name="`headers[${index}].value`",
        :error-messages="errors.collect(`headers[${index}].value`)"
        v-validate="'required'",
        @input="updateField(`headers[${index}].value`, $event)"
        )
      v-btn(icon, @click="removeItemFromArray('headers', index)")
        v-icon close
    v-layout
      v-btn(color="primary", @click="addHeader") Add header
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
import uid from '@/helpers/uid';

import formDeepMixin from '@/mixins/form/deep';

export default {
  inject: ['$validator'],
  mixins: [formDeepMixin],
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
  methods: {
    addHeader() {
      this.addItemIntoArray('headers', { id: uid(), key: '', value: '' });
    },
  },
};
</script>
