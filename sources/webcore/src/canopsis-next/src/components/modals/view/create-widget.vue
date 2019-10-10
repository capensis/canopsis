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
          :key="widget",
          @click="selectWidgetType(widget)"
        )
          v-card.widgetType(:data-test="`widget-${widget}`")
            v-card-title(primary-title)
              v-layout(wrap, justify-between)
                v-flex(xs11)
                  div.subheading {{ $t(`modals.widgetCreation.types.${widget}.title`) }}
                v-flex
                  v-icon {{ iconByWidgetType(widget) }}
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
</template>

<script>
import { MODALS, WIDGET_TYPES, SIDE_BARS_BY_WIDGET_TYPES, WIDGET_TYPES_RULES, WIDGET_ICONS } from '@/constants';

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
  computed: {
    /**
     * Some widgets are only available with 'cat' edition.
     * Filter widgetTypes list to keep only available widgets, thanks to the edition
     *
     * @return {Array}
     */
    availableWidgets() {
      return Object.keys(WIDGET_TYPES).filter((widget) => {
        const rules = WIDGET_TYPES_RULES[WIDGET_TYPES[widget]];

        if (!rules) {
          return true;
        }

        return rules.edition && rules.edition === this.edition;
      });
    },

    iconByWidgetType() {
      return type => WIDGET_ICONS[WIDGET_TYPES[type]];
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

