<template>
  <v-layout column>
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
          <i-18n path="stateSetting.entityThresholdSummary">
            <b place="name">{{ entity.name }}</b>
            <b place="state">{{ entityStateString }}</b>
            <span place="method">{{ currentCondition.method }}</span>
            <span place="condition">{{ conditionMethodSummary }}</span>
            <b place="impactingEntitiesState">{{ currentCondition.state }}</b>
            <b place="value">{{ conditionValue }}</b>
          </i-18n>
        </v-layout>
      </div>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import { ENTITIES_STATES, JUNIT_STATE_SETTING_METHODS, MODALS, STATE_SETTING_METHODS } from '@/constants';

import { isEntityEventsStateSettings } from '@/helpers/entities/entity/entity';

export default {
  props: {
    entity: {
      type: Object,
      required: true,
    },
    stateSetting: {
      type: Object,
      required: false,
    },
    pending: {
      type: Boolean,
      required: false,
    },
  },
  computed: {
    isEventsStateSettings() {
      return isEntityEventsStateSettings(this.entity);
    },

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
      } = this.stateSetting.state_thresholds;

      return {
        [ENTITIES_STATES.ok]: ok,
        [ENTITIES_STATES.minor]: minor,
        [ENTITIES_STATES.major]: major,
        [ENTITIES_STATES.critical]: critical,
      }[this.entityState];
    },

    conditionMethodSummary() {
      return this.$t(`stateSetting.thresholdConditions.${this.currentCondition.cond}`).toLowerCase();
    },

    conditionValue() {
      return `${this.currentCondition.value}${this.isShareMethod ? '%' : ''}`;
    },

    stateMethodName() {
      if (this.isEventsStateSettings) {
        return this.$tc('common.event', 2);
      }

      return this.stateSetting?.title || this.$t(`stateSetting.junit.methods.${JUNIT_STATE_SETTING_METHODS.worst}`);
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
      } = this.stateSetting.state_thresholds;

      const currentCondition = {
        [ENTITIES_STATES.ok]: ok,
        [ENTITIES_STATES.minor]: minor,
        [ENTITIES_STATES.major]: major,
        [ENTITIES_STATES.critical]: critical,
      }[this.entityState];

      return this.$t('stateSetting.entityThresholdSummary', {
        state: this.entityStateString,
        method: currentCondition.method,
        condition: this.$t(`stateSetting.thresholdConditions.${currentCondition.cond}`).toLowerCase(),
        impactingEntitiesState: currentCondition.state,
        value: `${currentCondition.value}${this.isShareMethod ? '%' : ''}`,
      });
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
