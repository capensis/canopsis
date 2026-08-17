<template>
  <c-layout class="text-subtitle-1" column>
    <span>
      {{ $t('stateSetting.computeMethod') }}:
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
          {{ $t('stateSetting.stateIsInheritFrom', { name: entity.name }) }}
        </p>
        <c-layout v-else-if="isDependenciesMethod" column>
          <span class="mb-2">
            {{
              $t('stateSetting.entityThresholdSummary', {
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

import { ALARM_STATES, STATE_SETTING_METHODS, STATE_SETTING_THRESHOLDS_METHODS } from '@/constants';

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

  computed: {
    isInheritedMethod() {
      return this.stateSetting?.method === STATE_SETTING_METHODS.inherited;
    },

    isDependenciesMethod() {
      return this.stateSetting?.method === STATE_SETTING_METHODS.dependencies;
    },

    entityState() {
      return this.entity.state;
    },

    entityStateString() {
      return this.$t(`common.stateTypes.${this.entityState}`);
    },

    currentCondition() {
      const { ok, minor, major, critical } = this.stateSetting?.state_thresholds ?? {};

      return {
        [ALARM_STATES.ok]: ok,
        [ALARM_STATES.minor]: minor,
        [ALARM_STATES.major]: major,
        [ALARM_STATES.critical]: critical,
      }[this.entityState];
    },

    isConditionShareMethod() {
      return this.currentCondition?.method === STATE_SETTING_THRESHOLDS_METHODS.share;
    },

    conditionMethodSummary() {
      return this.$t(`stateSetting.thresholdConditions.${this.currentCondition.cond}`).toLowerCase();
    },

    conditionValue() {
      return `${this.currentCondition.value}${this.isConditionShareMethod ? '%' : ''}`;
    },

    stateMethodName() {
      return this.stateSetting?.title || this.$t('common.event', 2);
    },

    stateMethodSummaryText() {
      if (!this.stateSetting || !this.currentCondition) {
        return '';
      }

      if (this.isInheritedMethod) {
        return this.$t('stateSetting.stateIsInheritFrom');
      }

      return this.$t('stateSetting.entityThresholdSummary', {
        state: this.entityStateString,
        method: this.currentCondition.method,
        condition: this.conditionMethodSummary,
        dependenciesEntitiesState: this.currentCondition.state,
        value: this.conditionValue,
      });
    },

    counts() {
      const {
        depends_count: dependsCount,
        threshold_state_depends_count: thresholdStateDependsCount,
        threshold_state: thresholdState,
      } = this.stateSetting ?? {};

      return isUndefined(dependsCount)
        ? []
        : [
            {
              message: this.$t('stateSetting.dependsCount'),
              count: dependsCount,
            },
            {
              message: this.$t('stateSetting.stateDependsCount', { state: thresholdState }),
              count: thresholdStateDependsCount,
            },
          ];
    },
  },
};
</script>

<style lang="scss" scoped>
.state-settings-summary__count-label {
  width: 250px;
}
</style>
