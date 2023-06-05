<template lang="pug">
  v-layout
    v-flex.pa-2.mr-2
      c-pattern-operator-information {{ $t('common.and') }}
    v-flex(xs11)
      v-layout(column)
        c-pattern-rules-field(
          :rules="group.rules",
          :attributes="attributes",
          :disabled="disabled",
          :readonly="readonly",
          @input="updateRules"
        )
</template>

<script>
import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
  model: {
    prop: 'group',
    event: 'input',
  },
  props: {
    group: {
      type: Object,
      required: true,
    },
    attributes: {
      type: Array,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    updateRules(rules) {
      if (rules.length) {
        this.updateField('rules', rules);
      } else {
        this.$emit('remove');
      }
    },
  },
};
</script>
