<template lang="pug">
  v-layout(justify-space-between, align-center)
    v-flex.pr-2(xs6)
      v-select(
        v-field="request.method",
        v-validate="'required'",
        :items="availableMethods",
        :label="methodLabel || $t('common.method')",
        :error-messages="errors.collect(methodFieldName)",
        :name="methodFieldName"
      )
    v-flex.pl-2(xs6)
      v-text-field(
        v-field="request.url",
        v-validate="'required|url'",
        :label="urlLabel || $t('common.url')",
        :error-messages="errors.collect(urlFieldName)",
        :name="urlFieldName"
      )
        v-tooltip(v-if="helpText", slot="append", left)
          v-icon(slot="activator") help
          div(v-html="helpText")
</template>

<script>
import { REQUEST_METHODS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'request',
    event: 'input',
  },
  props: {
    request: {
      type: Object,
      required: true,
    },
    methodLabel: {
      type: String,
      required: false,
    },
    urlLabel: {
      type: String,
      required: false,
    },
    helpText: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'request',
    },
  },
  computed: {
    availableMethods() {
      return Object.values(REQUEST_METHODS);
    },

    methodFieldName() {
      return `${this.name}.method`;
    },

    urlFieldName() {
      return `${this.name}.url`;
    },
  },
};
</script>
