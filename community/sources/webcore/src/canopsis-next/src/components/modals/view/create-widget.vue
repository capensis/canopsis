<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('modals.createWidget.title') }}
    template(#text="")
      v-layout(row, wrap)
        v-flex.my-1(
          v-for="type in availableTypes",
          :key="type",
          xs12,
          @click="selectType(type)"
        )
          v-card.cursor-pointer
            v-card-title(primary-title)
              v-layout(wrap, justify-between)
                v-flex(xs11)
                  div.subheading {{ $t(`modals.createWidget.types.${type}.title`) }}
                v-flex
                  v-icon {{ getIconByType(type) }}
</template>

<script>
import {
  MODALS,
  WIDGET_TYPES,
  SIDE_BARS_BY_WIDGET_TYPES,
  WIDGET_TYPES_RULES,
  WIDGET_ICONS,
} from '@/constants';

import { getNewWidgetGridParametersY } from '@/helpers/grid-layout';
import { getEmptyWidgetByType } from '@/helpers/forms/widgets/common';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWidget,
  components: { ModalWrapper },
  mixins: [modalInnerMixin, entitiesInfoMixin],
  computed: {
    /**
     * Some widgets are only available with 'pro' edition.
     * Filter widgetTypes list to keep only available widgets, thanks to the edition
     *
     * @return {Array}
     */
    availableTypes() {
      return [
        WIDGET_TYPES.alarmList,
        WIDGET_TYPES.context,
        WIDGET_TYPES.serviceWeather,
        WIDGET_TYPES.statsCalendar,
        WIDGET_TYPES.text,
        WIDGET_TYPES.counter,
        WIDGET_TYPES.testingWeather,
        WIDGET_TYPES.map,
      ].filter((widgetType) => {
        const rules = WIDGET_TYPES_RULES[widgetType];

        if (!rules) {
          return true;
        }

        return rules.edition && rules.edition === this.edition;
      });
    },

    tabWidgets() {
      return this.config.tab?.widgets ?? [];
    },
  },
  methods: {
    getIconByType(type) {
      return WIDGET_ICONS[type];
    },

    getWidgetWithUpdatedGridParametersByType(type) {
      const { tab } = this.config;
      const widget = getEmptyWidgetByType(type);
      const { mobile, tablet, desktop } = getNewWidgetGridParametersY(tab.widgets);

      widget.grid_parameters.mobile.y = mobile;
      widget.grid_parameters.tablet.y = tablet;
      widget.grid_parameters.desktop.y = desktop;
      widget.tab = tab._id;

      return widget;
    },

    selectType(type) {
      this.$sidebar.show({
        name: SIDE_BARS_BY_WIDGET_TYPES[type],
        config: {
          widget: this.getWidgetWithUpdatedGridParametersByType(type),
        },
      });

      this.$modals.hide();
    },
  },
};
</script>
