<template>
  <c-layout class="text-subtitle-1" column>
    <span>
      {{ t('stateSetting.computeMethod') }}:
      <v-progress-circular
        v-if="pending"
        class="ml-1"
        color="primary"
        size="20"
        width="3"
        indeterminate
      />
      <b v-else>
        {{ stateMethodName }}
      </b>
    </span>
    <v-expand-transition>
      <div v-if="!pending && stateMethodSummaryText">
        <p v-if="isInheritedMethod">
          <i18n-t keypath="stateSetting.stateIsInheritFrom" tag="span">
            <template #name>
              <b>{{ entity.name }}</b>
            </template>
          </i18n-t>
          <v-btn
            :aria-label="t('stateSetting.seeFilterPattern')"
            class="ml-2"
            color="primary"
            variant="outlined"
            size="small"
            @click="showStateSettingsPatterns"
          >
            {{ t('stateSetting.seeFilterPattern') }}
          </v-btn>
        </p>
        <c-layout v-else-if="isDependenciesMethod" column>
          <span class="mb-2">
            {{
              t('stateSetting.entityThresholdSummary', {
                name: entity.name,
                state: entityStateString,
                method: currentCondition.method,
                condition: conditionMethodSummary,
                dependenciesEntitiesState: currentCondition.state,
                value: conditionValue,
              })
            }}
          </span>
          <c-layout
            v-for="{ message, count } in counts"
            :key="message"
            class="text-body-2 font-weight-regular"
            align="center"
          >
            <div class="state-settings-summary__count-label">
              {{ message }}:
            </div>
            <b>{{ count }}</b>
          </c-layout>
        </c-layout>
      </div>
    </v-expand-transition>
  </c-layout>
</template>

<script>
import isUndefined from 'lodash/isUndefined';
import { computed } from 'vue';

import { ALARM_STATES, MODALS, STATE_SETTING_METHODS, STATE_SETTING_THRESHOLDS_METHODS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

export default {
  props: {
    entity: {
      type: Object,
      required: true,
    },
    pending: {
      type: Boolean,
      required: true,
    },
    stateSetting: {
      type: Object,
      required: false,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();

    const isInheritedMethod = computed(() => props.stateSetting?.method === STATE_SETTING_METHODS.inherited);

    const isDependenciesMethod = computed(
      () => props.stateSetting?.method === STATE_SETTING_METHODS.dependencies,
    );

    const entityState = computed(() => props.entity.state);

    const entityStateString = computed(() => t(`common.stateTypes.${entityState.value}`));

    const currentCondition = computed(() => {
      const { ok, minor, major, critical } = props.stateSetting?.state_thresholds ?? {};

      return {
        [ALARM_STATES.ok]: ok,
        [ALARM_STATES.minor]: minor,
        [ALARM_STATES.major]: major,
        [ALARM_STATES.critical]: critical,
      }[entityState.value];
    });

    const isConditionShareMethod = computed(
      () => currentCondition.value?.method === STATE_SETTING_THRESHOLDS_METHODS.share,
    );

    const conditionMethodSummary = computed(() =>
      t(`stateSetting.thresholdConditions.${currentCondition.value.cond}`).toLowerCase(),
    );

    const conditionValue = computed(
      () => `${currentCondition.value.value}${isConditionShareMethod.value ? '%' : ''}`,
    );

    const stateMethodName = computed(() => props.stateSetting?.title || t('common.event', 2));

    const stateMethodSummaryText = computed(() => {
      if (!props.stateSetting) {
        return '';
      }

      if (isInheritedMethod.value) {
        return t('stateSetting.stateIsInheritFrom');
      }

      if (!currentCondition.value) {
        return '';
      }

      return t('stateSetting.entityThresholdSummary', {
        state: entityStateString.value,
        method: currentCondition.value.method,
        condition: conditionMethodSummary.value,
        dependenciesEntitiesState: currentCondition.value.state,
        value: conditionValue.value,
      });
    });

    const counts = computed(() => {
      const {
        depends_count: dependsCount,
        threshold_state_depends_count: thresholdStateDependsCount,
        threshold_state: thresholdState,
      } = props.stateSetting ?? {};

      return isUndefined(dependsCount)
        ? []
        : [
            {
              message: t('stateSetting.dependsCount'),
              count: dependsCount,
            },
            {
              message: t('stateSetting.stateDependsCount', { state: thresholdState }),
              count: thresholdStateDependsCount,
            },
          ];
    });

    const showStateSettingsPatterns = () => {
      modals.show({
        name: MODALS.stateSettingInheritedEntityPattern,
        config: {
          pattern: props.stateSetting.inherited_entity_pattern,
        },
      });
    };

    return {
      t,
      isInheritedMethod,
      isDependenciesMethod,
      entityStateString,
      currentCondition,
      conditionMethodSummary,
      conditionValue,
      stateMethodName,
      stateMethodSummaryText,
      counts,
      showStateSettingsPatterns,
    };
  },
};
</script>

<style lang="scss" scoped>
.state-settings-summary__count-label {
  width: 250px;
}
</style>
