<template lang="pug">
  v-layout(column)
    state-setting-condition-field(
      v-for="key in sortedConditionsKeys",
      v-field="conditions[key]",
      :key="key",
      :label="$t(`stateSetting.states.${key}`)",
      :name="`${name}.${key}`",
      @input="errors.remove(name)"
    )
    c-alert(
      :value="errors.has(name)",
      type="error"
    ) {{ $t('stateSetting.conditionsError') }}
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';
import { validationAttachRequiredMixin } from '@/mixins/form/validation-attach-required';

import StateSettingConditionField from '../fields/state-setting-condition-field.vue';

export default {
  inject: ['$validator'],
  components: { StateSettingConditionField },
  mixins: [formValidationHeaderMixin, validationAttachRequiredMixin],
  model: {
    prop: 'conditions',
    event: 'input',
  },
  props: {
    conditions: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'conditions',
    },
  },
  computed: {
    sortedConditionsKeys() {
      return Object.keys(this.conditions).reverse();
    },
  },
  created() {
    this.attachRequiredRule(this.requiredRuleGetter);
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
  methods: {
    requiredRuleGetter() {
      return Object.values(this.conditions).some(({ enabled }) => enabled);
    },
  },
};
</script>
