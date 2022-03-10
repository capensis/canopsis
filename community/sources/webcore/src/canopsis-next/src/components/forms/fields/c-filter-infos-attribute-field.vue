<template lang="pug">
  v-layout
    v-flex(xs6)
      v-combobox(
        v-field="value.dictionary",
        v-validate="'required'",
        :items="dictionaries",
        :disabled="disabled",
        :label="label || $t('common.dictionary')",
        :return-object="false",
        :name="dictionaryName",
        :error-messages="errors.collect(dictionaryName)",
        hide-details
      )
    v-flex.ml-2(xs6)
      v-select(
        v-field="value.field",
        v-validate="'required'",
        :items="fieldItems",
        :disabled="disabled || !value.dictionary",
        :label="label || $t('common.field')",
        :name="fieldName",
        :error-messages="errors.collect(fieldName)",
        hide-details
      )
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'infos',
    },
    divider: {
      type: String,
      default: '.',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    dictionaryName() {
      return `${this.name}.dictionary`;
    },

    fieldName() {
      return `${this.name}.field`;
    },

    fieldItems() {
      return [
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.value'),
          value: 'value',
        },
      ];
    },

    dictionaries() {
      return this.items;
    },
  },
};
</script>
