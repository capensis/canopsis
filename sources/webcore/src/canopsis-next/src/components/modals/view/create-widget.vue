<template lang="pug">
  modal-wrapper(data-test="createWidgetModal")
    template(slot="title")
      span {{ $t('modals.widgetCreation.title') }}
    template(slot="text")
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
    template(slot="actions")
      v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
</template>

<script>
import { MODALS, WIDGET_TYPES, SIDE_BARS_BY_WIDGET_TYPES, WIDGET_TYPES_RULES, WIDGET_ICONS } from '@/constants';

import { generateWidgetByType } from '@/helpers/entities';

import modalInnerMixin from '@/mixins/modal/inner';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import entitiesInfoMixin from '@/mixins/entities/info';

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
      const widget = generateWidgetByType(WIDGET_TYPES[type]);

      this.showSideBar({
        name: SIDE_BARS_BY_WIDGET_TYPES[WIDGET_TYPES[type]],
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
    cursor: pointer,
  }
</style>

