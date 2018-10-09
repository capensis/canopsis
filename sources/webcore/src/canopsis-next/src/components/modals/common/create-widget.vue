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
import { MODALS, WIDGET_TYPES, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';
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
        { title: this.$t('modals.widgetCreation.types.alarmList.title'), value: WIDGET_TYPES.alarmList },
        { title: this.$t('modals.widgetCreation.types.context.title'), value: WIDGET_TYPES.context },
        { title: this.$t('modals.widgetCreation.types.weather.title'), value: WIDGET_TYPES.weather },
        { title: this.$t('modals.widgetCreation.types.statsTable.title'), value: WIDGET_TYPES.statsTable },
      ],
    };
  },
  methods: {
    selectWidgetType(type) {
      const widget = generateWidgetByType(type);

      this.showSideBar({
        name: SIDE_BARS_BY_WIDGET_TYPES[type],
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

