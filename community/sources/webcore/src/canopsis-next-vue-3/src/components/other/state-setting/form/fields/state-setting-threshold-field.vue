<template>
  <c-layout column>
    <h5 class="text-subtitle-1 font-weight-bold">
      {{ t('stateSetting.targetEntityState') }} - {{ label }}
    </h5>

    <c-layout class="target-entity-state-row" align="center" gap="3">
      <div class="target-entity-state-row__toggle">
        <c-enabled-field
          :model-value="condition.enabled"
          no-margin
          hide-details
          hide-label
          @update:model-value="updateConditionField('enabled', $event)"
        />
      </div>

      <state-setting-threshold-method-field
        :model-value="condition.method"
        :name="methodName"
        :disabled="disabled"
        class="target-entity-state-row__method"
        standalone
        @update:model-value="updateConditionField('method', $event)"
      />

      <c-select
        :model-value="condition.state"
        :label="t('stateSetting.entitiesStates')"
        :items="states"
        :disabled="disabled"
        :required="!disabled"
        :name="stateName"
        :standalone="disabled"
        class="target-entity-state-row__field"
        item-title="text"
        item-value="value"
        @update:model-value="updateConditionField('state', $event)"
      />

      <c-select
        :model-value="condition.cond"
        :label="t('common.condition', 1)"
        :items="conditions"
        :disabled="disabled"
        :required="!disabled"
        :name="conditionName"
        :standalone="disabled"
        item-title="text"
        item-value="value"
        class="target-entity-state-row__field"
        @update:model-value="updateConditionField('cond', $event)"
      />

      <c-number-field
        :value="condition.value"
        :label="t('common.value')"
        :disabled="disabled"
        :required="!disabled"
        :name="valueName"
        :min="0"
        :max="maxValue"
        :standalone="disabled"
        class="target-entity-state-row__value"
        @update:value="updateConditionField('value', $event)"
      />
    </c-layout>

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

    const methodName = computed(() => `${props.name}.method`);
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
    const maxValue = computed(() => (isShareMethod.value ? 100 : undefined));

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
      const nextCondition = {
        ...props.condition,
        [field]: value,
      };

      if (
        field === 'method' &&
        value === STATE_SETTING_THRESHOLDS_METHODS.share &&
        Number(nextCondition.value) > 100
      ) {
        nextCondition.value = 100;
      }

      emit('update:condition', nextCondition);
    };

    return {
      t,
      methodName,
      stateName,
      conditionName,
      valueName,
      disabled,
      states,
      conditions,
      maxValue,
      summaryMessage,
      updateConditionField,
    };
  },
};
</script>

<style lang="scss" scoped>
.target-entity-state-row {
  &__toggle {
    flex: 0 0 auto;
  }

  &__method {
    flex: 3 1 0;
    min-width: 0;
  }

  &__field {
    flex: 3 1 0;
    min-width: 0;
  }

  &__value {
    flex: 2 1 0;
    min-width: 0;
  }
}
</style>
