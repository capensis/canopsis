<template>
  <card-with-see-alarms-btn
    :class="itemClasses"
    :show-button="showAlarmsButton"
    :style="itemStyle"
    class="cursor-pointer service-weather-item"
    @click="handleCardClick"
    @show:alarms="$emit('show:alarms')"
  >
    <c-layout
      class="service-weather-item__content"
      justify="space-between"
      fill
      row
    >
      <div class="service-weather-item__main position-relative">
        <c-layout
          :class="{ blinking: isBlinking }"
          class="service-weather-item__template-wrapper"
          justify="start"
          fill
        >
          <c-compiled-template
            :template="template"
            :context="templateContext"
            class="service-weather-item__template pa-3"
          />
        </c-layout>

        <c-layout align="center" class="service-weather-item__toolbar pt-1 pr-1">
          <c-no-events-icon :value="service.idle_since" color="white" location="top" />

          <impact-state-indicator
            v-if="priorityEnabled"
            :value="service.impact_state"
            class="mr-1"
          />

          <v-btn
            v-if="showRootCauseByStateClick"
            size="small"
            variant="text"
            icon
            @click.stop="$emit('show:root-cause')"
          >
            <v-icon color="white" icon="$root_cause" size="small" />
          </v-btn>

          <v-btn
            v-if="showVariablesHelpButton"
            size="small"
            variant="text"
            icon
            @click.stop="showVariablesHelpModal"
          >
            <v-icon color="white" icon="help" size="small" />
          </v-btn>
        </c-layout>

        <v-icon
          class="service-weather-item__background"
          color="white"
        >
          {{ backgroundIcon }}
        </v-icon>

        <v-icon
          v-if="secondaryIconEnabled && service.secondary_icon"
          class="service-weather-item__secondary-icon text-white mb-1 mr-1"
        >
          {{ service.secondary_icon }}
        </v-icon>
      </div>

      <alarm-pbehavior-counters
        v-if="isPbehaviorCountersEnabled && hasPbehaviorCounters"
        :counters="pbehaviorCounters"
        :types="pbehaviorCountersTypes"
      />

      <alarm-state-counters
        v-if="isStateCountersEnabled"
        :counters="counters"
        :types="stateCountersTypes"
      />
    </c-layout>
  </card-with-see-alarms-btn>
</template>

<script>
import { ALARM_STATES, MODALS, SERVICE_WEATHER_DEFAULT_EM_HEIGHT } from '@/constants';

import { getMostReadableTextColor } from '@/helpers/color';
import { getEntityColorClass } from '@/helpers/entities/entity/color';
import { registerHelper, unregisterHelper } from '@/helpers/handlebars';
import { alarmOrEntityTagsHelper } from '@/helpers/handlebars/tags-helper';
import { convertObjectToTreeview } from '@/helpers/treeview';

import CardWithSeeAlarmsBtn from '@/components/common/card/card-with-see-alarms-btn.vue';

import AlarmPbehaviorCounters from './alarm-pbehavior-counters.vue';
import AlarmStateCounters from './alarm-state-counters.vue';
import ImpactStateIndicator from './impact-state-indicator.vue';

export default {
  components: {
    CardWithSeeAlarmsBtn,
    AlarmPbehaviorCounters,
    AlarmStateCounters,
    ImpactStateIndicator,
  },

  props: {
    service: {
      type: Object,
      required: true,
    },
    margin: {
      type: Object,
      default: () => ({
        top: 0,
        right: 0,
        bottom: 0,
        left: 0,
      }),
    },
    template: {
      type: String,
      default: '',
    },
    heightFactor: {
      type: Number,
      default: 0,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    actionRequiredBlinking: {
      type: Boolean,
      default: true,
    },
    actionRequiredIcon: {
      type: String,
      required: false,
    },
    actionRequiredColor: {
      type: String,
      required: false,
    },
    noActionRequiredIcon: {
      type: String,
      required: false,
    },
    noActionRequiredColor: {
      type: String,
      required: false,
    },
    showAlarmsButton: {
      type: Boolean,
      default: false,
    },
    showVariablesHelpButton: {
      type: Boolean,
      default: false,
    },
    showRootCauseByStateClick: {
      type: Boolean,
      default: false,
    },
    priorityEnabled: {
      type: Boolean,
      default: true,
    },
    secondaryIconEnabled: {
      type: Boolean,
      default: true,
    },
    countersSettings: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['show:service', 'show:alarms', 'show:root-cause'],

  computed: {
    templateContext() {
      return { entity: this.service };
    },

    isActionRequired() {
      return this.service.is_action_required;
    },

    isNoActionRequired() {
      return !this.isActionRequired && this.service?.state?.val > ALARM_STATES.ok;
    },

    isBlinking() {
      return this.isActionRequired && this.actionRequiredBlinking;
    },

    backgroundIcon() {
      return (
        {
          [this.isActionRequired]: this.actionRequiredIcon,
          [this.isNoActionRequired]: this.noActionRequiredIcon,
        }.true || this.service.icon
      );
    },

    backgroundColorByActionRequired() {
      return {
        [this.isActionRequired]: this.actionRequiredColor,
        [this.isNoActionRequired]: this.noActionRequiredColor,
      }.true;
    },

    colorByActionRequired() {
      if (!this.backgroundColorByActionRequired) {
        return undefined;
      }

      return getMostReadableTextColor(this.backgroundColorByActionRequired, {
        level: 'AA',
        size: 'large',
      });
    },

    itemClasses() {
      return [
        `mt-${this.margin.top}`,
        `mr-${this.margin.right}`,
        `mb-${this.margin.bottom}`,
        `ml-${this.margin.left}`,
        getEntityColorClass(this.service, this.colorIndicator),
      ];
    },

    itemHeight() {
      return SERVICE_WEATHER_DEFAULT_EM_HEIGHT + this.heightFactor;
    },

    itemStyle() {
      return {
        height: `${this.itemHeight}em`,
        backgroundColor: this.backgroundColorByActionRequired,
        color: this.colorByActionRequired,
      };
    },

    counters() {
      return this.service.counters ?? {};
    },

    pbehaviorCounters() {
      return this.counters?.pbh_types ?? [];
    },

    hasPbehaviorCounters() {
      return this.pbehaviorCounters.length;
    },

    isPbehaviorCountersEnabled() {
      return this.countersSettings?.pbehavior_enabled;
    },

    pbehaviorCountersTypes() {
      return this.countersSettings?.pbehavior_types ?? [];
    },

    isStateCountersEnabled() {
      return this.countersSettings?.state_enabled;
    },

    stateCountersTypes() {
      return this.countersSettings?.state_types ?? [];
    },
  },

  beforeCreate() {
    registerHelper('tags', alarmOrEntityTagsHelper);
  },

  beforeUnmount() {
    unregisterHelper('tags');
  },

  methods: {
    handleCardClick(event) {
      if (event.target.tagName !== 'A' || !event.target.href) {
        this.$emit('show:service');
      }
    },

    showVariablesHelpModal() {
      const entityFields = convertObjectToTreeview(this.service, 'entity');
      const variables = [entityFields];

      this.$modals.show({
        name: MODALS.variablesHelp,
        config: {
          variables,
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.service-weather-item {
  width: 100%;
  overflow: hidden;
  background-image: linear-gradient(90deg, rgba(255, 255, 255, 0) 0%, rgba(0, 0, 0, 0.3) 100%);

  &__content {
    width: 100%;
    height: 100%;
  }

  &__main {
    flex: 1 1 auto;
    min-width: 0;
    height: 100%;
  }

  &__template-wrapper {
    width: 100%;
    height: 100%;
  }

  &__content > * {
    margin-right: 2px;

    &:first-child,
    &:last-child {
      margin: 0;
    }
  }

  &__template {
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.2em;

    :deep(hr) {
      border-color: currentColor;
      background-color: currentColor;
    }
  }

  &__background {
    position: absolute;
    top: 50%;
    bottom: 0;
    right: 10px;
    font-size: 5em !important;
    pointer-events: none;
    opacity: 0.5;
    transform: translateY(-50%);
  }

  &__toolbar {
    position: absolute;
    top: 0;
    right: 0;
    z-index: 1;
  }

  &__secondary-icon {
    position: absolute;
    right: 0;
    bottom: 0;
    cursor: inherit;

    &:hover,
    &:focus {
      position: absolute;
    }
  }
}
</style>
