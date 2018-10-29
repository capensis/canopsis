<template lang='pug'>
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.widgetCreation.title') }}
    v-card-text
      v-layout(row, wrap)
        v-flex.my-1(
        xs12,
        v-for="widgetType in widgetsTypes",
        :key="widgetType.value",
        @click="selectWidgetType(widgetType.value)"
        )
          v-card.widgetType
            v-card-title(primary-title)
              v-layout(wrap, justify-between)
                v-flex(xs11)
                  div.subheading {{ widgetType.title }}
                v-flex
                  v-icon {{ widgetType.icon }}
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
</template>

<script>
import { MODALS } from '@/constants';
import { generateWidgetByType } from '@/helpers/entities';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import sideBarMixin from '@/mixins/side-bar/side-bar';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWidget,
  mixins: [modalInnerMixin, sideBarMixin],
  data() {
    return {
      widgetsTypes: [
        { title: this.$t('modals.widgetCreation.types.alarmList.title'), value: this.$constants.WIDGET_TYPES.alarmList, icon: 'view_list' },
        { title: this.$t('modals.widgetCreation.types.context.title'), value: this.$constants.WIDGET_TYPES.context, icon: 'view_list' },
        { title: this.$t('modals.widgetCreation.types.weather.title'), value: this.$constants.WIDGET_TYPES.weather, icon: 'view_module' },
        { title: this.$t('modals.widgetCreation.types.statsHistogram.title'), value: this.$constants.WIDGET_TYPES.statsHistogram, icon: 'bar_chart' },
        { title: this.$t('modals.widgetCreation.types.statsTable.title'), value: this.$constants.WIDGET_TYPES.statsTable, icon: 'table_chart' },
        { title: this.$t('modals.widgetCreation.types.statsCalendar.title'), value: this.$constants.WIDGET_TYPES.statsCalendar, icon: 'calendar_today' },
        { title: this.$t('modals.widgetCreation.types.statsCurves.title'), value: this.$constants.WIDGET_TYPES.statsCurves, icon: 'show_chart' },
        { title: this.$t('modals.widgetCreation.types.statsNumber.title'), value: this.$constants.WIDGET_TYPES.statsNumber, icon: 'table_chart' },
      ],
    };
  },
  methods: {
    selectWidgetType(type) {
      const widget = generateWidgetByType(type);

      this.showSideBar({
        name: this.$constants.SIDE_BARS_BY_WIDGET_TYPES[type],
        config: {
          widget,
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

