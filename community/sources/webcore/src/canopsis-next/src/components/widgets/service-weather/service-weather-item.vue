<template lang="pug">
  card-with-see-alarms-btn.cursor-pointer.service-weather-item(
    :class="itemClasses",
    :show-button="showAlarmsButton",
    :style="itemStyle",
    @click.native="handleCardClick",
    @show:alarms="$emit('show:alarms')"
  )
    v-layout.fill-height.service-weather-item__content(row, justify-space-between)
      v-flex.position-relative.fill-height
        v-layout(:class="{ 'blinking': isBlinking }", justify-start)
          c-compiled-template.service-weather-item__template.pa-3(
            :template="template",
            :context="templateContext"
          )
        v-layout.service-weather-item__toolbar.pt-1.pr-1(row, align-center)
          c-no-events-icon(:value="service.idle_since", :color="color", top)
          impact-state-indicator.mr-1(v-if="priorityEnabled", :value="service.impact_state")
          v-btn.ma-0(
            v-if="showVariablesHelpButton",
            icon,
            small,
            @click.stop="showVariablesHelpModal(service)"
          )
            v-icon(color="white", small) help
        v-icon.service-weather-item__background(size="5em", :color="color") {{ backgroundIcon }}
        v-icon.service-weather-item__secondary-icon.mb-1.mr-1(
          v-if="service.secondary_icon",
          :color="color"
        ) {{ service.secondary_icon }}
      alarm-pbehavior-counters(
        v-if="isPbehaviorCountersEnabled && hasPbehaviorCounters",
        :counters="pbehaviorCounters",
        :types="pbehaviorCountersTypes"
      )
      alarm-state-counters(
        v-if="isStateCountersEnabled",
        :counters="counters",
        :types="stateCountersTypes"
      )
</template>

<script>
import { MODALS, SERVICE_WEATHER_DEFAULT_EM_HEIGHT } from '@/constants';

import { getEntityColor } from '@/helpers/entities/entity/color';
import { getMostReadableTextColor } from '@/helpers/color';
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
    showAlarmsButton: {
      type: Boolean,
      default: false,
    },
    showVariablesHelpButton: {
      type: Boolean,
      default: false,
    },
    priorityEnabled: {
      type: Boolean,
      default: true,
    },
    countersSettings: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    templateContext() {
      return { entity: this.service };
    },

    isActionRequired() {
      return this.service.is_action_required;
    },

    isBlinking() {
      return this.isActionRequired && this.actionRequiredBlinking;
    },

    backgroundIcon() {
      if (this.isActionRequired && this.actionRequiredIcon) {
        return this.actionRequiredIcon;
      }

      return this.service.icon;
    },

    backgroundColor() {
      if (this.isActionRequired && this.actionRequiredColor) {
        return this.actionRequiredColor;
      }

      return getEntityColor(this.service, this.colorIndicator);
    },

    color() {
      if (this.isActionRequired && this.actionRequiredColor) {
        return getMostReadableTextColor(this.backgroundColor, { level: 'AA', size: 'large' });
      }

      return 'white';
    },

    itemClasses() {
      return [
        `mt-${this.margin.top}`,
        `mr-${this.margin.right}`,
        `mb-${this.margin.bottom}`,
        `ml-${this.margin.left}`,
      ];
    },

    itemHeight() {
      return SERVICE_WEATHER_DEFAULT_EM_HEIGHT + this.heightFactor;
    },

    itemStyle() {
      return {
        height: `${this.itemHeight}em`,
        backgroundColor: this.backgroundColor,
        color: this.color,
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
  overflow: hidden;

  &__content > * {
    margin-right: 2px;

    &:first-child, &:last-child {
      margin: 0;
    }
  }

  &__template {
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.2em;
  }

  &__background {
    position: absolute;
    top: 0;
    bottom: 0;
    right: 10px;
    pointer-events: none;
    opacity: 0.5;
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

    &:hover, &:focus {
      position: absolute;
    }
  }
}
</style>
