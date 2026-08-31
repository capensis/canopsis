<template>
  <c-layout column>
    <h5 class="text-subtitle-1 font-weight-bold">
      {{ t('stateSetting.targetEntityState') }} - {{ label }}
    </h5>

    <v-row class="target-entity-state-row align-center" dense>
      <v-col cols="1">
        <c-enabled-field
          :model-value="condition.enabled"
          no-margin
          hide-details
          @update:model-value="updateConditionField('enabled', $event)"
        />
      </v-col>

      <v-col cols="3">
        <state-setting-threshold-method-field
          :model-value="condition.method"
          :disabled="disabled"
          @update:model-value="updateConditionField('method', $event)"
        />
      </v-col>

      <v-col cols="3">
        <c-select
          :model-value="condition.state"
          :label="t('stateSetting.entitiesStates')"
          :items="states"
          :disabled="disabled"
          :required="!disabled"
          :name="stateName"
          item-title="text"
          item-value="value"
          @update:model-value="updateConditionField('state', $event)"
        />
      </v-col>

      <v-col cols="3">
        <c-select
          :model-value="condition.cond"
          :label="t('common.condition', 1)"
          :items="conditions"
          :disabled="disabled"
          :required="!disabled"
          :name="conditionName"
          item-title="text"
          item-value="value"
          @update:model-value="updateConditionField('cond', $event)"
        />
      </v-col>

      <v-col cols="2">
        <c-number-field
          :value="condition.value"
          :label="t('common.value')"
          :disabled="disabled"
          :required="!disabled"
          :name="valueName"
          :min="0"
          @update:value="updateConditionField('value', $event)"
        />
      </v-col>
    </v-row>

    <v-expand-transition>
      <span v-if="summaryMessage">
        <strong>{{ t('common.summary') }}:</strong>
        <span class="ml-2">
          {{ summaryMessage }}
        </span>
      </span>
    </v-expand-transition>
  </c-layout>
</template>

<script>
import { pick } from 'lodash';
import { computed } from 'vue';

import {
  ALARM_STATES,
  STATE_SETTING_THRESHOLDS_CONDITIONS,
  STATE_SETTING_THRESHOLDS_METHODS,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';

import StateSettingThresholdMethodField from './state-setting-threshold-method-field.vue';

export default {
  components: { StateSettingThresholdMethodField },

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

  emits: ['update:condition'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const stateName = computed(() => `${props.name}.state`);
    const conditionName = computed(() => `${props.name}.condition`);
    const valueName = computed(() => `${props.name}.value`);
    const disabled = computed(() => !props.condition.enabled);

    const states = computed(() => Object.entries(ALARM_STATES).map(([key, value]) => ({
        text: t(`common.stateTypes.${value}`),
        value: key,
      })),
    );

    const conditions = computed(() => Object.values(STATE_SETTING_THRESHOLDS_CONDITIONS).map(condition => ({
        value: condition,
        text: t(`stateSetting.thresholdConditions.${condition}`),
      })),
    );

    const isShareMethod = computed(
      () => props.condition.method === STATE_SETTING_THRESHOLDS_METHODS.share,
    );

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

    const updateConditionField = (field, value) => {
      emit('update:condition', {
        ...props.condition,
        [field]: value,
      });
    };

    return {
      t,
      stateName,
      conditionName,
      valueName,
      disabled,
      states,
      conditions,
      summaryMessage,
      updateConditionField,
    };
  },
};
</script>

<style lang="scss">
.target-entity-state-row {
  gap: 12px;
}
</style>
