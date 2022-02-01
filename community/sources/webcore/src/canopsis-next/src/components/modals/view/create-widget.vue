<template lang="pug">
  modal-wrapper(data-test="createWidgetModal", close)
    template(slot="title")
      span {{ $t('modals.createWidget.title') }}
    template(slot="text")
      v-layout(row, wrap)
        v-flex.my-1(
          v-for="widgetType in availableWidgetTypes",
          :key="widgetType",
          xs12,
          @click="selectWidgetType(widgetType)"
        )
          v-card.widgetType(:data-test="`widget-${widgetType}`")
            v-card-title(primary-title)
              v-layout(wrap, justify-between)
                v-flex(xs11)
                  div.subheading {{ $t(`modals.createWidget.types.${widgetType}.title`) }}
                v-flex
                  v-icon {{ getIconByWidgetType(widgetType) }}
</template>

<script>
import {
  MODALS,
  WIDGET_TYPES,
  SIDE_BARS_BY_WIDGET_TYPES,
  WIDGET_TYPES_RULES,
  WIDGET_ICONS,
} from '@/constants';

import { generateWidgetByType } from '@/helpers/entities';

import { modalInnerMixin } from '@/mixins/modal/inner';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWidget,
  components: { ModalWrapper },
  mixins: [modalInnerMixin, sideBarMixin, entitiesInfoMixin],
  computed: {
    /**
     * Some widgets are only available with 'cat' edition.
     * Filter widgetTypes list to keep only available widgets, thanks to the edition
     *
     * @return {Array}
     */
    availableWidgetTypes() {
      return [
        WIDGET_TYPES.alarmList,
        WIDGET_TYPES.context,
        WIDGET_TYPES.serviceWeather,
        WIDGET_TYPES.statsCalendar,
        WIDGET_TYPES.text,
        WIDGET_TYPES.counter,
        WIDGET_TYPES.testingWeather,
      ].filter((widgetType) => {
        const rules = WIDGET_TYPES_RULES[widgetType];

        if (!rules) {
          return true;
        }

        return rules.edition && rules.edition === this.edition;
      });
    },
  },
  methods: {
    getIconByWidgetType(widgetType) {
      return WIDGET_ICONS[widgetType];
    },

    selectWidgetType(widgetType) {
      const widget = generateWidgetByType(widgetType);

      this.showSideBar({
        name: SIDE_BARS_BY_WIDGET_TYPES[widgetType],
        config: {
          widget,
          tabId: this.config.tabId,
          isNew: true,
        },
      });
      this.$modals.hide();
    },
  },
};
</script>

<style lang="scss" scoped>
  .widgetType {
    cursor: pointer;
  }
</style>
