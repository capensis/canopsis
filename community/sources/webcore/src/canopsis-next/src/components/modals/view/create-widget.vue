<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ $t('modals.createWidget.title') }}</span>
    </template>
    <template #text="">
      <v-layout column>
        <v-card
          v-for="{ type, text, icon, on, children } in availableTypes"
          :key="type"
          class="my-1 cursor-pointer"
          v-on="on"
        >
          <v-menu
            :disabled="!children"
            top
            offset-y
          >
            <template #activator="{ on: menuOn }">
              <v-card-title
                primary-title
                v-on="menuOn"
              >
                <v-layout
                  wrap
                  justify-between
                >
                  <v-flex xs11>
                    <div class="text-subtitle-1">
                      {{ text }}
                    </div>
                  </v-flex>
                  <v-flex>
                    <v-icon>{{ icon }}</v-icon>
                  </v-flex>
                </v-layout>
              </v-card-title>
            </template>
            <v-list>
              <v-list-item
                v-for="child in children"
                :key="child.type"
                v-on="child.on"
              >
                <v-list-item-avatar>
                  <v-icon>{{ child.icon }}</v-icon>
                </v-list-item-avatar>
                <v-list-item-title>{{ child.text }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>
        </v-card>
      </v-layout>
    </template>
  </modal-wrapper>
</template>

<script>
import {
  MODALS,
  WIDGET_TYPES,
  SIDE_BARS_BY_WIDGET_TYPES,
  WIDGET_TYPES_RULES,
  WIDGET_ICONS,
  TOP_LEVEL_WIDGET_TYPES,
} from '@/constants';

import { calculateNewWidgetGridParametersY } from '@/helpers/entities/widget/grid';
import { getEmptyWidgetByType } from '@/helpers/entities/widget/form';

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
      const widgetTypePreparer = type => ({
        type,
        text: this.$t(`modals.createWidget.types.${type}.title`),
        icon: WIDGET_ICONS[type],
        on: {
          click: () => this.selectType(type),
        },
      });

      const types = TOP_LEVEL_WIDGET_TYPES.map(widgetTypePreparer);

      types.push({
        type: 'chart',
        text: this.$t('modals.createWidget.types.chart.title'),
        icon: 'stacked_bar_chart',
        children: [
          WIDGET_TYPES.barChart,
          WIDGET_TYPES.lineChart,
          WIDGET_TYPES.pieChart,
          WIDGET_TYPES.numbers,
        ].map(widgetTypePreparer),
      }, {
        type: 'report',
        text: this.$t('modals.createWidget.types.report.title'),
        icon: 'table_chart',
        children: [
          WIDGET_TYPES.userStatistics,
          WIDGET_TYPES.alarmStatistics,
        ].map(widgetTypePreparer),
      });

      return types.reduce((acc, type) => {
        if (!type.children) {
          const rule = WIDGET_TYPES_RULES[type];

          if (!rule?.edition || rule?.edition === this.edition) {
            acc.push(type);
          }

          return acc;
        }

        const filteredChildren = type.children.filter((child) => {
          const rule = WIDGET_TYPES_RULES[child.type];

          return !rule?.edition || rule?.edition === this.edition;
        });

        if (filteredChildren.length) {
          acc.push({
            ...type,
            children: filteredChildren,
          });
        }

        return acc;
      }, []);
    },
  },
  methods: {
    getWidgetWithUpdatedGridParametersByType(type) {
      const { tab } = this.config;
      const widget = getEmptyWidgetByType(type);
      const { mobile, tablet, desktop } = calculateNewWidgetGridParametersY(tab.widgets);

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
