<template>
  <v-layout
    class="state-setting-thresholds-step"
    column
  >
    <state-setting-threshold-field
      v-field="thresholds[key]"
      v-for="key in sortedThresholdsKeys"
      :key="key"
      :label="getLabel(key)"
      :name="`${name}.${key}`"
      :state="key"
      @input="errors.remove(name)"
    />
    <c-alert
      :value="errors.has(name)"
      type="error"
    >
      {{ $t('stateSetting.conditionsError') }}
    </c-alert>
  </v-layout>
</template>

<script>
import { ALARM_STATES } from '@/constants';

import { formValidationHeaderMixin } from '@/mixins/form';
import { validationAttachRequiredMixin } from '@/mixins/form/validation-attach-required';

import StateSettingThresholdField from '../fields/state-setting-threshold-field.vue';

export default {
  inject: ['$validator'],
  components: { StateSettingThresholdField },
  mixins: [formValidationHeaderMixin, validationAttachRequiredMixin],
  model: {
    prop: 'thresholds',
    event: 'input',
  },
  props: {
    thresholds: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'thresholds',
    },
  },
  computed: {
    sortedThresholdsKeys() {
      return Object.keys(this.thresholds).reverse();
    },
  },
  created() {
    this.attachRequiredRule(this.requiredRuleGetter);
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
  methods: {
    getLabel(key) {
      return this.$t(`common.stateTypes.${ALARM_STATES[key]}`);
    },

    requiredRuleGetter() {
      return Object.values(this.thresholds).some(({ enabled }) => enabled);
    },
  },
};
</script>

<style lang="scss" scoped>
.state-setting-thresholds-step {
  gap: 16px;
}
</style>
