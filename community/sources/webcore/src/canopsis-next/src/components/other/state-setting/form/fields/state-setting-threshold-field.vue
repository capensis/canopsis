<template>
  <v-layout column>
    <h5 class="subheading font-weight-bold">
      {{ $t('stateSetting.targetEntityState') }} - {{ label }}
    </h5>
    <v-layout
      class="target-entity-state-row"
      align-center
    >
      <v-flex xs1>
        <v-switch
          v-field="condition.enabled"
          color="primary"
        />
      </v-flex>
      <v-flex xs3>
        <state-setting-threshold-method-field
          v-field="condition.method"
          :disabled="disabled"
        />
      </v-flex>
      <v-flex xs3>
        <c-select-field
          v-field="condition.state"
          :label="$t('stateSetting.entitiesStates')"
          :items="states"
          :disabled="disabled"
          :required="!disabled"
          :name="stateName"
        />
      </v-flex>
      <v-flex xs3>
        <c-select-field
          v-field="condition.cond"
          :label="$tc('common.condition', 1)"
          :items="conditions"
          :disabled="disabled"
          :required="!disabled"
          :name="conditionName"
        />
      </v-flex>
      <v-flex xs2>
        <c-number-field
          v-field="condition.value"
          :label="$t('common.value')"
          :disabled="disabled"
          :required="!disabled"
          :name="valueName"
          :min="0"
        />
      </v-flex>
    </v-layout>
    <v-expand-transition>
      <span v-if="summaryMessage">
        <strong>{{ $t('common.summary') }}:</strong>
        <span class="ml-2">{{ summaryMessage }}</span>
      </span>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import { pick } from 'lodash';

import {
  ENTITIES_STATES_KEYS,
  STATE_SETTING_THRESHOLDS_CONDITIONS,
  STATE_SETTING_THRESHOLDS_METHODS,
} from '@/constants';

import StateSettingThresholdMethodField from './state-setting-threshold-method-field.vue';

export default {
  inject: ['$validator'],
  components: { StateSettingThresholdMethodField },
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
    state: {
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
      return Object.values(STATE_SETTING_THRESHOLDS_CONDITIONS)
        .map(condition => ({
          value: condition,
          text: this.$t(`stateSetting.thresholdConditions.${condition}`),
        }));
    },

    isShareMethod() {
      return this.condition.method === STATE_SETTING_THRESHOLDS_METHODS.share;
    },

    summaryMessage() {
      const fieldsForSummary = pick(this.condition, ['cond', 'state', 'value']);
      const hasFieldsForSummary = Object.values(fieldsForSummary).every(value => !!String(value));

      return hasFieldsForSummary
        ? this.$t('stateSetting.targetEntityThresholdSummary', {
          state: this.state,
          method: this.condition.method,
          condition: this.$t(`stateSetting.thresholdConditions.${this.condition.cond}`).toLowerCase(),
          impactingEntitiesState: this.condition.state,
          value: `${this.condition.value}${this.isShareMethod ? '%' : ''}`,
        })
        : '';
    },
  },
};
</script>

<style lang="scss">
.target-entity-state-row {
  gap: 12px;
}
</style>
