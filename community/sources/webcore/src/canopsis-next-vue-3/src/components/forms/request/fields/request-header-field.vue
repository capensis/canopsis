<template>
  <v-row>
    <v-col
      class="mr-3"
      cols="6"
    >
      <c-combobox
        :model-value="header.text"
        :label="$t('common.headerKey')"
        :disabled="disabled"
        :items="textHints"
        :name="textFieldName"
        required
        @update:model-value="updateField('text', $event)"
      />
    </v-col>
    <v-col cols="6">
      <c-combobox
        v-if="valueHints"
        :model-value="header.value"
        :label="$t('common.headerValue')"
        :disabled="disabled"
        :items="valueHints"
        :name="valueFieldName"
        :return-object="false"
        required
        @update:model-value="updateField('value', $event)"
      />
      <c-payload-text-field
        v-else
        :value="header.value"
        :label="$t('common.headerValue')"
        :disabled="disabled"
        :variables="headersVariables"
        :name="valueFieldName"
        required
        @update:value="updateField('value', $event)"
      />
    </v-col>
  </v-row>
</template>

<script>
import { formMixin } from '@/mixins/form';

export default {
  emits: ['update:header'],
  inject: ['$validator'],
  mixins: [formMixin],
  props: {
    header: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'header',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    headersHints: {
      type: Array,
      default: () => [],
    },
    headersVariables: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    textHints() {
      return this.headersHints.map(({ text }) => text);
    },

    valueHints() {
      return this.headersHints.find(({ text }) => text === this.header.text)?.value;
    },

    textFieldName() {
      return `${this.name}.text`;
    },

    valueFieldName() {
      return `${this.name}.value`;
    },
  },
};
</script>
