<template lang="pug">
  v-layout(column)
    h5.subheading.font-weight-bold {{ $t('stateSetting.targetEntityState') }} - {{ label }}
    v-layout.target-entity-state-row(row, align-center)
      v-flex(xs3)
        state-setting-calculation-method-field(v-field="value.type")
      v-flex(xs3)
        c-select-field(
          v-field="value.state",
          :label="$t('stateSetting.entitiesStates')",
          :items="states",
          name="state",
          required
        )
      v-flex(xs3)
        c-select-field(
          v-field="value.condition",
          :label="$tc('common.condition', 1)",
          :items="conditions",
          name="condition",
          required
        )
      v-flex(xs3)
        c-number-field(
          v-field="value.value",
          :label="$t('common.value')",
          required
        )
</template>

<script>
import { ENTITIES_STATES, STATE_SETTING_CONDITIONS } from '@/constants';

import StateSettingCalculationMethodField from './state-setting-calculation-method-field.vue';

export default {
  inject: ['$validator'],
  components: {
    StateSettingCalculationMethodField,
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
    states() {
      return Object.values(ENTITIES_STATES)
        .map(state => ({
          text: this.$t(`common.stateTypes.${state}`),
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
