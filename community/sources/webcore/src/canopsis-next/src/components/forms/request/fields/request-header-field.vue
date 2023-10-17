<template>
  <v-layout>
    <v-flex
      class="mr-3"
      xs6
    >
      <v-combobox
        v-field="header.text"
        v-validate="'required'"
        :label="$t('common.headerKey')"
        :error-messages="errors.collect(textFieldName)"
        :disabled="disabled"
        :items="textHints"
        :name="textFieldName"
      />
    </v-flex>
    <v-flex xs6>
      <v-combobox
        v-if="valueHints"
        v-field="header.value"
        v-validate="'required'"
        :label="$t('common.headerValue')"
        :error-messages="errors.collect(valueFieldName)"
        :disabled="disabled"
        :items="valueHints"
        :name="valueFieldName"
        :return-object="false"
      />
      <c-payload-text-field
        v-else
        v-field="header.value"
        :label="$t('common.headerValue')"
        :error-messages="errors.collect(valueFieldName)"
        :disabled="disabled"
        :variables="headersVariables"
        :name="valueFieldName"
        required
      />
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'header',
    event: 'input',
  },
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
