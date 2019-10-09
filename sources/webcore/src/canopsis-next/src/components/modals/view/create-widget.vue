<template lang='pug'>
  v-card(data-test="createWidgetModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.widgetCreation.title') }}
    v-card-text
      v-layout(row, wrap)
        v-flex.my-1(
          xs12,
          v-for="widget in availableWidgets",
          :key="widget.value",
          @click="selectWidgetType(widget.value)"
        )
          v-card.widgetType(:data-test="`widget-${widget.value}`")
            v-card-title(primary-title)
              v-layout(wrap, justify-between)
                v-flex(xs11)
                  div.subheading {{ widget.title }}
                v-flex
                  v-icon {{ widget.icon }}
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
</template>

<script>
import { MODALS, WIDGET_TYPES, SIDE_BARS_BY_WIDGET_TYPES, CANOPSIS_EDITION } from '@/constants';

import { generateWidgetByType } from '@/helpers/entities';

import modalInnerMixin from '@/mixins/modal/inner';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import entitiesInfoMixin from '@/mixins/entities/info';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWidget,
  mixins: [modalInnerMixin, sideBarMixin, entitiesInfoMixin],
  data() {
    return {
      widgetsTypes: [
        {
          title: this.$t('modals.widgetCreation.types.alarmList.title'),
          value: WIDGET_TYPES.alarmList,
          icon: 'view_list',
          group: 'base',
        },
        {
          title: this.$t('modals.widgetCreation.types.context.title'),
          value: WIDGET_TYPES.context,
          icon: 'view_list',
          group: 'base',
        },
        {
          title: this.$t('modals.widgetCreation.types.weather.title'),
          value: WIDGET_TYPES.weather,
          icon: 'view_module',
          group: 'base',
        },
        {
          title: this.$t('modals.widgetCreation.types.statsHistogram.title'),
          value: WIDGET_TYPES.statsHistogram,
          icon: 'bar_chart',
          group: 'stat',
        },
        {
          title: this.$t('modals.widgetCreation.types.statsTable.title'),
          value: WIDGET_TYPES.statsTable,
          icon: 'table_chart',
          group: 'stat',
        },
        {
          title: this.$t('modals.widgetCreation.types.statsCalendar.title'),
          value: WIDGET_TYPES.statsCalendar,
          icon: 'calendar_today',
          group: 'stat',
        },
        {
          title: this.$t('modals.widgetCreation.types.statsCurves.title'),
          value: WIDGET_TYPES.statsCurves,
          icon: 'show_chart',
          group: 'stat',
        },
        {
          title: this.$t('modals.widgetCreation.types.statsNumber.title'),
          value: WIDGET_TYPES.statsNumber,
          icon: 'table_chart',
          group: 'stat',
        },
        {
          title: this.$t('modals.widgetCreation.types.statsPareto.title'),
          value: WIDGET_TYPES.statsPareto,
          icon: 'multiline_chart',
          group: 'stat',
        },
        {
          title: this.$t('modals.widgetCreation.types.text.title'),
          value: WIDGET_TYPES.text,
          icon: 'view_headline',
          group: 'base',
        },
      ],
    };
  },
  computed: {
    availableWidgets() {
      if (this.edition === CANOPSIS_EDITION.core) {
        return this.widgetsTypes.filter(widget => widget.group !== 'stat');
      }

      return this.widgetsTypes;
    },
  },
  methods: {
    selectWidgetType(type) {
      const widget = generateWidgetByType(type);

      this.showSideBar({
        name: SIDE_BARS_BY_WIDGET_TYPES[type],
        config: {
          widget,
          tabId: this.config.tabId,
          isNew: true,
        },
      });
      this.hideModal();
    },
  },
};
</script>

<style lang="scss" scoped>
  .widgetType {
    cursor: pointer,
  }
</style>

