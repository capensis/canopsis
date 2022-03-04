<template lang="pug">
  v-layout
    v-flex(xs6)
      v-text-field(
        :value="valueParts.field",
        v-validate="'required'",
        :disabled="disabled",
        :label="label || $t('common.field')",
        :error-messages="errors.collect(name)",
        :name="name",
        hide-details,
        @input="updateValue"
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
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'extra_infos',
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
      const [extraInfos] = this.value.split(this.divider, 1);
      const field = this.value.slice(extraInfos.length + 1);

      return {
        extraInfos,
        field,
      };
    },
  },
  methods: {
    updateValue(field) {
      this.updateModel([this.valueParts.extraInfos, field].join(this.divider));
    },
  },
};
</script>
