<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ $t('modals.createWidget.title') }}
    template(#text="")
      v-layout(column)
        v-card.my-1.cursor-pointer(
          v-for="{ type, text, icon } in availableTypes",
          :key="type",
          @click="selectType(type)"
        )
          v-card-title(primary-title)
            v-layout(wrap, justify-between)
              v-flex(xs11)
                div.subheading {{ text }}
              v-flex
                v-icon {{ icon }}
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
      return Object.values(WIDGET_TYPES).filter((type) => {
        const rules = WIDGET_TYPES_RULES[type];

        if (!rules) {
          return true;
        }

        return rules.edition && rules.edition === this.edition;
      }).map(type => ({
        type,
        text: this.$t(`modals.createWidget.types.${type}.title`),
        icon: WIDGET_ICONS[type],
      }));
    },
  },
  methods: {
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
