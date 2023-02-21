<template lang="pug">
  v-card.white--text.cursor-pointer.weather-item(
    :class="itemClasses",
    :style="itemStyle",
    tile,
    @click.native="showAdditionalInfoModal"
  )
    v-layout.fill-height.weather-item__content(row, justify-space-between)
      v-flex.position-relative.fill-height
        v-layout(:class="{ 'blinking': isBlinking }", justify-start)
          c-runtime-template.weather-item__service-name.pa-3(:template="compiledTemplate")
        v-layout.weather-item__toolbar.pt-1.pr-1(row, align-center)
          c-no-events-icon(:value="service.idle_since", color="white", top)
          impact-state-indicator.mr-1(v-if="isPriorityEnabled", :value="service.impact_state")
        v-icon.weather-item__background.white--text(size="5em") {{ service.icon }}
        v-icon.weather-item__secondary-icon.mb-1.mr-1(
          v-if="service.secondary_icon",
          color="white"
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
    v-btn.see-alarms-btn(
      v-if="isBothModalType && hasAlarmsListAccess",
      flat,
      @click.stop="showAlarmListModal"
    ) {{ $t('serviceWeather.seeAlarms') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import {
  MODALS,
  USERS_PERMISSIONS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
} from '@/constants';

import { compile } from '@/helpers/handlebars';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities';
import { getEntityColor } from '@/helpers/color';

import { authMixin } from '@/mixins/auth';
import { entitiesServiceEntityMixin } from '@/mixins/entities/service-entity';

import AlarmPbehaviorCounters from './alarm-pbehavior-counters.vue';
import AlarmStateCounters from './alarm-state-counters.vue';
import ImpactStateIndicator from './impact-state-indicator.vue';

const { mapActions } = createNamespacedHelpers('service');

export default {
  components: {
    AlarmPbehaviorCounters,
    AlarmStateCounters,
    ImpactStateIndicator,
  },
  mixins: [authMixin, entitiesServiceEntityMixin],
  props: {
    service: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.widget.parameters.blockTemplate ?? '', { entity: this.service });

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
  computed: {
    hasMoreInfosAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.moreInfos);
    },

    hasAlarmsListAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.alarmsList);
    },

    color() {
      return getEntityColor(this.service, this.widget.parameters.colorIndicator);
    },

    itemClasses() {
      const classes = [
        `mt-${this.widget.parameters.margin.top}`,
        `mr-${this.widget.parameters.margin.right}`,
        `mb-${this.widget.parameters.margin.bottom}`,
        `ml-${this.widget.parameters.margin.left}`,
      ];

      if (this.isBothModalType && this.hasAlarmsListAccess) {
        classes.push('v-card__with-see-alarms-btn');
      }

      return classes;
    },

    itemHeight() {
      return 4 + this.widget.parameters.heightFactor;
    },

    itemStyle() {
      return {
        height: `${this.itemHeight}em`,
        backgroundColor: this.color,
      };
    },

    isBlinking() {
      return this.service.is_action_required;
    },

    isBothModalType() {
      return this.widget.parameters.modalType === SERVICE_WEATHER_WIDGET_MODAL_TYPES.both;
    },

    isAlarmListModalType() {
      return this.widget.parameters.modalType === SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList;
    },

    isPriorityEnabled() {
      return this.widget.parameters.isPriorityEnabled ?? true;
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
      return this.widget.parameters.counters?.pbehavior_enabled;
    },

    pbehaviorCountersTypes() {
      return this.widget.parameters.counters?.pbehavior_types ?? [];
    },

    isStateCountersEnabled() {
      return this.widget.parameters.counters?.state_enabled;
    },

    stateCountersTypes() {
      return this.widget.parameters.counters?.state_types ?? [];
    },
  },
  methods: {
    ...mapActions({
      fetchServiceAlarmsWithoutStore: 'fetchAlarmsWithoutStore',
    }),

    showAdditionalInfoModal(e) {
      if (e.target.tagName !== 'A' || !e.target.href) {
        if (this.isAlarmListModalType && this.hasAlarmsListAccess) {
          this.showAlarmListModal();
        } else if (!this.isAlarmListModalType && this.hasMoreInfosAccess) {
          this.showMainInfoModal();
        }
      }
    },

    showMainInfoModal() {
      this.$modals.show({
        name: MODALS.serviceEntities,
        config: {
          color: this.color,
          service: this.service,
          widgetParameters: this.widget.parameters,
        },
      });
    },

    async showAlarmListModal() {
      try {
        const widget = generatePreparedDefaultAlarmListWidget();

        widget.parameters = {
          ...widget.parameters,
          ...this.widget.parameters.alarmsList,

          serviceDependenciesColumns: this.widget.parameters.serviceDependenciesColumns,
        };

        this.$modals.show({
          name: MODALS.alarmsList,
          config: {
            widget,
            title: this.$t('modals.alarmsList.prefixTitle', { prefix: this.service.name }),
            fetchList: params => this.fetchServiceAlarmsWithoutStore({ id: this.service._id, params }),
          },
        });
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.weather-item {
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

  &__content > * {
    margin-right: 2px;

    &:first-child, &:last-child {
      margin: 0;
    }
  }
}
</style>
