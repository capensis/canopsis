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
          :max="valueMax"
        />
      </v-flex>
    </v-layout>
    <v-expand-transition>
      <span v-if="summaryMessage">
        <strong>{{ $t('common.summary') }}:</strong>
        <span class="ml-2">
          {{ summaryMessage }}
        </span>
      </span>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import { computed } from 'vue';
import { pick } from 'lodash';

import { ALARM_STATES, STATE_SETTING_THRESHOLDS_CONDITIONS, STATE_SETTING_THRESHOLDS_METHODS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

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
  setup(props) {
    const { t } = useI18n();

    const stateName = computed(() => `${props.name}.state`);
    const conditionName = computed(() => `${props.name}.condition`);
    const valueName = computed(() => `${props.name}.value`);
    const disabled = computed(() => !props.condition.enabled);
    const isShareMethod = computed(() => props.condition.method === STATE_SETTING_THRESHOLDS_METHODS.share);
    const valueMax = computed(() => (isShareMethod.value ? 99 : undefined));

    const states = computed(() => Object.entries(ALARM_STATES)
      .map(([key, value]) => ({
        text: t(`common.stateTypes.${value}`),
        value: key,
      })));

    const conditions = computed(() => Object.values(STATE_SETTING_THRESHOLDS_CONDITIONS)
      .map(condition => ({
        value: condition,
        text: t(`stateSetting.thresholdConditions.${condition}`),
      })));

    const summaryMessage = computed(() => {
      const fieldsForSummary = pick(props.condition, ['cond', 'state', 'value']);
      const hasFieldsForSummary = Object.values(fieldsForSummary).every(value => !!String(value));

      return hasFieldsForSummary
        ? t('stateSetting.targetEntityThresholdSummary', {
          state: props.state,
          method: props.condition.method,
          condition: t(`stateSetting.thresholdConditions.${props.condition.cond}`).toLowerCase(),
          dependenciesEntitiesState: props.condition.state,
          value: `${props.condition.value}${isShareMethod.value ? '%' : ''}`,
        })
        : '';
    });

    return {
      stateName,
      conditionName,
      valueName,
      disabled,
      valueMax,
      states,
      conditions,
      summaryMessage,
    };
  },
};
</script>

<style lang="scss">
.target-entity-state-row {
  gap: 12px;
}
</style>
