<template lang="pug">
  v-layout
    v-flex(xs6)
      v-select(
        :value="dictionary",
        v-validate="'required'",
        :items="dictionaries",
        :disabled="disabled",
        :label="label || $t('common.dictionary')",
        :name="name",
        hide-details,
        @input="updateDictionary"
      )
    v-flex.ml-2(xs6)
      v-select(
        :value="field",
        v-validate="'required'",
        :items="fieldItems",
        :disabled="disabled || !dictionary",
        :label="label || $t('common.field')",
        :name="name",
        hide-details,
        @input="updateField"
      )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
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
    valueParts() {
      return this.value.split(this.divider);
    },

    dictionary() {
      return this.valueParts[1];
    },

    field() {
      return this.valueParts[2];
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
  methods: {
    updateValue(dictionary, field) {
      const [infos] = this.valueParts;

      this.updateModel([infos, dictionary, field].join(this.divider));
    },

    updateDictionary(dictionary) {
      this.updateValue(dictionary, this.field);
    },

    updateField(field) {
      this.updateValue(this.dictionary, field);
    },
  },
};
</script>
