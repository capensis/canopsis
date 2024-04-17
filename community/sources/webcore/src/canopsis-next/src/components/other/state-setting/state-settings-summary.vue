<template>
  <v-layout class="text-subtitle-1" column>
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
      <b v-else>{{ stateMethodName }}</b>
    </span>
    <v-expand-transition>
      <div v-if="!pending && stateMethodSummaryText">
        <p v-if="isInheritedMethod">
          <i18n path="stateSetting.stateIsInheritFrom" tag="span">
            <b place="name">{{ entity.name }}</b>
          </i18n>
          <v-btn
            class="ml-2"
            color="primary"
            outlined
            small
            @click="showStateSettingsPatterns"
          >
            {{ $t('stateSetting.seeFilterPattern') }}
          </v-btn>
        </p>
        <v-layout v-else-if="isDependenciesMethod" column>
          <i18n class="mb-2" path="stateSetting.entityThresholdSummary">
            <b place="name">{{ entity.name }}</b>
            <b place="state">{{ entityStateString }}</b>
            <span place="method">{{ currentCondition.method }}</span>
            <span place="condition">{{ conditionMethodSummary }}</span>
            <b place="dependenciesEntitiesState">{{ currentCondition.state }}</b>
            <b place="value">{{ conditionValue }}</b>
          </i18n>
          <v-layout
            v-for="{ message, count } in counts"
            :key="message"
            class="text-body-2 font-weight-regular"
            align-center
          >
            <v-flex
              lg3
              md4
              sm5
              xs6
            >
              {{ message }}:
            </v-flex>
            <v-flex><b>{{ count }}</b></v-flex>
          </v-layout>
        </v-layout>
      </div>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import { isUndefined } from 'lodash';

import { ALARM_STATES, MODALS, STATE_SETTING_METHODS, STATE_SETTING_THRESHOLDS_METHODS } from '@/constants';

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
      const {
        ok,
        minor,
        major,
        critical,
      } = this.stateSetting.state_thresholds ?? {};

      return {
        [ALARM_STATES.ok]: ok,
        [ALARM_STATES.minor]: minor,
        [ALARM_STATES.major]: major,
        [ALARM_STATES.critical]: critical,
      }[this.entityState];
    },

    isConditionShareMethod() {
      return this.currentCondition.method === STATE_SETTING_THRESHOLDS_METHODS.share;
    },

    conditionMethodSummary() {
      return this.$t(`stateSetting.thresholdConditions.${this.currentCondition.cond}`).toLowerCase();
    },

    conditionValue() {
      return `${this.currentCondition.value}${this.isConditionShareMethod ? '%' : ''}`;
    },

    stateMethodName() {
      return this.stateSetting?.title || this.$tc('common.event', 2);
    },

    stateMethodSummaryText() {
      if (!this.stateSetting) {
        return '';
      }

      if (this.isInheritedMethod) {
        return this.$t('stateSetting.stateIsInheritFrom');
      }

      const {
        ok,
        minor,
        major,
        critical,
      } = this.stateSetting.state_thresholds ?? {};

      const currentCondition = {
        [ALARM_STATES.ok]: ok,
        [ALARM_STATES.minor]: minor,
        [ALARM_STATES.major]: major,
        [ALARM_STATES.critical]: critical,
      }[this.entityState];

      if (!currentCondition) {
        return '';
      }

      return this.$t('stateSetting.entityThresholdSummary', {
        state: this.entityStateString,
        method: currentCondition.method,
        condition: this.$t(`stateSetting.thresholdConditions.${currentCondition.cond}`).toLowerCase(),
        dependenciesEntitiesState: currentCondition.state,
        value: `${currentCondition.value}${this.isConditionShareMethod ? '%' : ''}`,
      });
    },

    counts() {
      const {
        depends_count: dependsCount,
        threshold_state_depends_count: thresholdStateDependsCount,
        threshold_state: thresholdState,
      } = this.stateSetting;

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
  methods: {
    showStateSettingsPatterns() {
      this.$modals.show({
        name: MODALS.stateSettingInheritedEntityPattern,
        config: {
          pattern: this.stateSetting.inherited_entity_pattern,
        },
      });
    },
  },
};
</script>
