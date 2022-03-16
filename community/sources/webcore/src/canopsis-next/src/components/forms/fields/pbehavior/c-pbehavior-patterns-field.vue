<template lang="pug">
  c-pattern-groups-field(
    v-field="groups",
    :disabled="disabled",
    :name="name",
    :attributes="pbehaviorAttributes"
  )
</template>

<script>
import { PATTERN_OPERATORS, PBEHAVIOR_PATTERN_FIELDS } from '@/constants';

export default {
  model: {
    prop: 'groups',
    event: 'input',
  },
  props: {
    groups: {
      type: Array,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      required: false,
    },
  },
  computed: {
    nameOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-pbehavior-field',
        },
      };
    },

    reasonOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-pbehavior-reason-field',
        },
      };
    },

    typeOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-pbehavior-type-field',
        },
      };
    },

    pbehaviorAttributes() {
      return [
        {
          text: this.$t('common.name'),
          value: PBEHAVIOR_PATTERN_FIELDS.name,
          options: this.nameOptions,
        },
        {
          text: this.$t('common.reason'),
          value: PBEHAVIOR_PATTERN_FIELDS.reason,
          options: this.reasonOptions,
        },
        {
          text: this.$t('common.type'),
          value: PBEHAVIOR_PATTERN_FIELDS.type,
          options: this.typeOptions,
        },
      ];
    },
  },
};
</script>
