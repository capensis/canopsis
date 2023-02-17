<template lang="pug">
  v-layout(row)
    v-flex.pr-2(xs3)
      v-select(
        v-field="condition.type",
        :items="conditionTypes",
        :label="$t('common.type')",
        :disabled="disabled"
      )
    v-flex.px-2(xs4)
      v-text-field(
        v-field="condition.attribute",
        v-validate="'required'",
        :label="$t('common.attribute')",
        :name="conditionFieldName",
        :error-messages="errors.collect(conditionFieldName)",
        :disabled="disabled"
      )
    v-flex.pl-2(xs5)
      v-layout(row, align-center)
        c-payload-text-field(
          v-field="condition.value",
          :label="$t('common.value')",
          :disabled="disabled",
          :variables="variables",
          clearable
        )
        v-btn(
          v-if="!disabled",
          :disabled="disabledRemove",
          icon,
          small,
          @click="removeCondition"
        )
          v-icon(color="error", small) delete
</template>

<script>
import {
  EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES,
  EVENT_FILTER_EXTERNAL_DATA_CONDITION_VALUES,
} from '@/constants';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'condition',
    event: 'input',
  },
  props: {
    condition: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledRemove: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    variables() {
      return Object.values(EVENT_FILTER_EXTERNAL_DATA_CONDITION_VALUES).map(({ value, text }) => ({
        value,
        text: this.$t(`eventFilter.externalDataValues.${text}`),
      }));
    },

    conditionTypes() {
      return Object.values(EVENT_FILTER_EXTERNAL_DATA_CONDITION_TYPES)
        .map(type => ({ text: this.$t(`eventFilter.externalDataConditionTypes.${type}`), value: type }));
    },

    conditionFieldName() {
      return `${this.name}.condition`;
    },
  },
  methods: {
    removeCondition() {
      this.$emit('remove', this.condition);
    },
  },
};
</script>
