<template lang="pug">
  v-card.white--text.cursor-pointer(
  :class="getItemClasses",
  :style="{ height: itemHeight + 'em'}",
  tile,
  @click.native="showAdditionalInfoModal"
  )
    div(:class="{ blinking: isBlinking }", )
      v-layout(justify-start)
        v-icon.px-3.py-2.white--text(size="2em") {{ format.icon }}
        div.watcherName.pt-3(v-html="compiledTemplate")
        v-btn.pauseIcon.white(v-if="watcher.active_pb_some && !watcher.active_pb_all", fab, icon, small)
          v-icon pause
</template>

<script>
import find from 'lodash/find';

import {
  MODALS,
  WIDGET_TYPES,
  WATCHER_STATES_COLORS,
  WATCHER_PBEHAVIOR_COLOR,
  PBEHAVIOR_TYPES,
  WEATHER_ICONS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
} from '@/constants';

import compile from '@/helpers/handlebars';
import { generateWidgetByType } from '@/helpers/entities';

import modalMixin from '@/mixins/modal';
import entitiesWatcherEntityMixin from '@/mixins/entities/watcher-entity';

export default {
  mixins: [modalMixin, entitiesWatcherEntityMixin],
  props: {
    watcher: {
      type: Object,
      required: true,
    },
    template: {
      type: String,
    },
    widget: {
      type: Object,
    },
  },
  computed: {
    isPaused() {
      return this.watcher.active_pb_all;
    },
    hasWatcherPbehavior() {
      return this.watcher.active_pb_watcher;
    },
    format() {
      if (!this.isPaused && !this.hasWatcherPbehavior) {
        const state = this.watcher.state.val;

        return {
          icon: WEATHER_ICONS[state],
          color: WATCHER_STATES_COLORS[state],
        };
      }

      const pbehaviors = this.hasWatcherPbehavior ? this.watcher.watcher_pbehavior : this.watcher.pbehavior;

      const maintenancePbehavior = find(pbehaviors, { type_: PBEHAVIOR_TYPES.maintenance });
      const outOfSurveillancePbehavior = find(pbehaviors, { type_: PBEHAVIOR_TYPES.outOfSurveillance });

      let icon = WEATHER_ICONS.pause;

      if (maintenancePbehavior) {
        icon = WEATHER_ICONS.maintenance;
      } else if (outOfSurveillancePbehavior) {
        icon = WEATHER_ICONS.outOfSurveillance;
      }

      if (this.isPaused && !this.hasWatcherPbehavior) {
        icon = WEATHER_ICONS.pause;
      }

      return {
        color: WATCHER_PBEHAVIOR_COLOR,
        icon,
      };
    },
    compiledTemplate() {
      return compile(this.template, { watcher: this.watcher });
    },
    getItemClasses() {
      return [
        this.format.color,
        `mt-${this.widget.parameters.margin.top}`,
        `mr-${this.widget.parameters.margin.right}`,
        `mb-${this.widget.parameters.margin.bottom}`,
        `ml-${this.widget.parameters.margin.left}`,
      ];
    },
    itemHeight() {
      return 4 + this.widget.parameters.heightFactor;
    },
    isBlinking() {
      return this.watcher.alerts_not_ack;
    },
  },
  methods: {
    showAdditionalInfoModal() {
      if (this.widget.parameters.modalType === SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList) {
        this.showAlarmListModal();
      } else {
        this.showMainInfoModal();
      }
    },

    showMainInfoModal() {
      this.showModal({
        name: MODALS.watcher,
        config: {
          watcherId: this.watcher.entity_id,
          entityTemplate: this.widget.parameters.entityTemplate,
          modalTemplate: this.widget.parameters.modalTemplate,
        },
      });
    },

    async showAlarmListModal() {
      const entities = await this.fetchWatcherEntitiesListWithoutStore({ watcherId: this.watcher.entity_id });
      const widget = generateWidgetByType(WIDGET_TYPES.alarmList);
      const watcherFilter = {
        title: this.watcher.display_name,
        filter: {
          $and: [
            {
              entity: { $in: entities.map(entity => entity.entity_id) },
            },
          ],
        },
      };

      const widgetParameters = {
        widgetColumns: widget.parameters.widgetColumns.map(column => ({
          label: column.label,
          value: column.value.replace('alarm.', 'v.'),
        })),
        mainFilter: watcherFilter,
        viewFilters: [watcherFilter],
      };

      this.showModal({
        name: MODALS.alarmsList,
        config: {
          query: {
            filter: watcherFilter.filter,
          },
          widget: {
            ...widget,

            parameters: {
              ...widget.parameters,
              ...widgetParameters,
            },
          },
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .pauseIcon {
    position: absolute;
    right: 0;
    bottom: 0;
    cursor: inherit;
  }

  .watcherName {
    color: white;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1em;
  }

  @keyframes blink {
    0% { opacity: 1 }
    50% { opacity: 0.3 }
  }

  .blinking {
    animation: blink 2s linear infinite;
  }

  .cursor-pointer {
    cursor: pointer;
  }
</style>
