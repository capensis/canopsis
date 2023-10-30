<template lang="pug">
  v-layout(column)
    h5.subheading.font-weight-bold {{ $t('stateSetting.targetEntityState') }} - {{ label }}
    v-layout.target-entity-state-row(row, align-center)
      v-flex(xs1)
        v-switch(
          v-field="condition.enabled",
          color="primary"
        )
      v-flex(xs3)
        state-setting-condition-method-field(
          v-field="condition.type",
          :disabled="disabled"
        )
      v-flex(xs3)
        c-select-field(
          v-field="condition.state",
          :label="$t('stateSetting.entitiesStates')",
          :items="states",
          :disabled="disabled",
          :required="!disabled",
          :name="stateName"
        )
      v-flex(xs3)
        c-select-field(
          v-field="condition.cond",
          :label="$tc('common.condition', 1)",
          :items="conditions",
          :disabled="disabled",
          :required="!disabled",
          :name="conditionName"
        )
      v-flex(xs2)
        c-number-field(
          v-field="condition.value",
          :label="$t('common.value')",
          :disabled="disabled",
          :required="!disabled",
          :name="valueName"
        )
</template>

<script>
import { ENTITIES_STATES_KEYS, STATE_SETTING_CONDITIONS } from '@/constants';

import StateSettingConditionMethodField from './state-setting-condition-method-field.vue';

export default {
  components: { StateSettingConditionMethodField },
  model: {
    prop: 'condition',
    event: 'input',
  },
  props: {
    condition: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: '',
    },
  },
  computed: {
    stateName() {
      return `${this.name}.state`;
    },

    conditionName() {
      return `${this.name}.condition`;
    },

    valueName() {
      return `${this.name}.value`;
    },

    disabled() {
      return !this.condition.enabled;
    },

    states() {
      return Object.values(ENTITIES_STATES_KEYS)
        .map(state => ({
          text: this.$t(`stateSetting.states.${state}`),
          value: state,
        }));
    },

    conditions() {
      return Object.values(STATE_SETTING_CONDITIONS)
        .map(condition => ({
          value: condition,
          text: this.$t(`stateSetting.conditions.${condition}`),
        }));
    },
  },
};
</script>

<style lang="scss">
.target-entity-state-row {
  gap: 12px;
}
</style>
