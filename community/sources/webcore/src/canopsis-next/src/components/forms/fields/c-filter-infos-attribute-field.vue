<template lang="pug">
  v-layout
    v-flex(xs6)
      v-combobox(
        :value="valueParts.dictionary",
        v-validate="'required'",
        :items="dictionaries",
        :disabled="disabled",
        :label="label || $t('common.dictionary')",
        :return-object="false",
        :name="name",
        hide-details,
        @input="updateDictionary"
      )
    v-flex.ml-2(xs6)
      v-select(
        :value="valueParts.field",
        v-validate="'required'",
        :items="fieldItems",
        :disabled="disabled || !valueParts.dictionary",
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
      const [infos, dictionary, field] = this.value.split(this.divider);

      return {
        infos,
        dictionary,
        field,
      };
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
      this.updateModel([this.valueParts.infos, dictionary, field].join(this.divider));
    },

    updateDictionary(dictionary) {
      this.updateValue(dictionary, this.valueParts.field);
    },

    updateField(field) {
      this.updateValue(this.valueParts.dictionary, field);
    },
  },
};
</script>
