<template lang='pug'>
  v-card
    v-card-title.green.darken-4.white--text
      v-layout(justify-space-between, align-center)
        h2 {{ $t('modals.widgetCreation.title') }}
        v-btn(@click='hideModal', icon, small)
          v-icon.white--text close
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
              v-layout(wrap)
                v-flex(xs12)
                  div.subheading {{ widgetType.title }}
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
        { title: this.$t('modals.widgetCreation.types.alarmList.title'), value: this.$constants.WIDGET_TYPES.alarmList },
        { title: this.$t('modals.widgetCreation.types.context.title'), value: this.$constants.WIDGET_TYPES.context },
        { title: this.$t('modals.widgetCreation.types.weather.title'), value: this.$constants.WIDGET_TYPES.weather },
        { title: this.$t('modals.widgetCreation.types.statsHistogram.title'), value: this.$constants.WIDGET_TYPES.statsHistogram },
        { title: this.$t('modals.widgetCreation.types.statsTable.title'), value: this.$constants.WIDGET_TYPES.statsTable },
        { title: this.$t('modals.widgetCreation.types.statsCalendar.title'), value: this.$constants.WIDGET_TYPES.statsCalendar },
        { title: this.$t('modals.widgetCreation.types.statsCurves.title'), value: this.$constants.WIDGET_TYPES.statsCurves },
        { title: this.$t('modals.widgetCreation.types.statsNumber.title'), value: this.$constants.WIDGET_TYPES.statsNumber },
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

