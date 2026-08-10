<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ $t('modals.createWidget.title') }}</span>
    </template>

    <template #text>
      <c-layout column>
        <v-card
          v-for="{ type, text, icon, children } in availableTypes"
          :key="type"
          class="my-1 cursor-pointer"
        >
          <v-menu v-if="children" location="top">
            <template #activator="{ props: activatorProps }">
              <v-card-title v-bind="activatorProps" class="py-3">
                <v-row align="center" justify="space-between" no-gutters>
                  <v-col cols="11">
                    <div class="text-subtitle-1">
                      {{ text }}
                    </div>
                  </v-col>
                  <v-col class="text-right">
                    <v-icon :icon="icon" />
                  </v-col>
                </v-row>
              </v-card-title>
            </template>

            <v-list>
              <v-list-item
                v-for="child in children"
                :key="child.type"
                @click="selectType(child.type)"
              >
                <template #prepend>
                  <v-icon :icon="child.icon" />
                </template>

                <v-list-item-title>{{ child.text }}</v-list-item-title>
              </v-list-item>
            </v-list>
          </v-menu>

          <v-card-title v-else class="py-3" @click="selectType(type)">
            <v-row align="center" justify="space-between" no-gutters>
              <v-col cols="11">
                <div class="text-subtitle-1">
                  {{ text }}
                </div>
              </v-col>
              <v-col class="text-right">
                <v-icon :icon="icon" />
              </v-col>
            </v-row>
          </v-card-title>
        </v-card>
      </c-layout>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { mapState } from 'pinia';

import {
  MODALS,
  SIDE_BARS_BY_WIDGET_TYPES,
  TOP_LEVEL_WIDGET_TYPES,
  WIDGET_ICONS,
  WIDGET_TYPES,
  WIDGET_TYPES_RULES,
} from '@/constants';
import { getEmptyWidgetByType } from '@/helpers/entities/widget/form';
import { calculateNewWidgetGridParametersY } from '@/helpers/entities/widget/grid';
import { useAppInfoStore } from '@/stores/modules/app-info';

export default {
  name: MODALS.createWidget,
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  computed: {
    ...mapState(useAppInfoStore, ['edition']),

    config() {
      return this.modal.config ?? {};
    },

    availableTypes() {
      const widgetTypePreparer = type => ({
        type,
        text: this.$t(`modals.createWidget.types.${type}.title`),
        icon: WIDGET_ICONS[type],
      });

      const types = TOP_LEVEL_WIDGET_TYPES.map(widgetTypePreparer);

      types.push(
        {
          type: 'chart',
          text: this.$t('modals.createWidget.types.chart.title'),
          icon: 'stacked_bar_chart',
          children: [
            WIDGET_TYPES.barChart,
            WIDGET_TYPES.lineChart,
            WIDGET_TYPES.pieChart,
            WIDGET_TYPES.numbers,
          ].map(widgetTypePreparer),
        },
        {
          type: 'report',
          text: this.$t('modals.createWidget.types.report.title'),
          icon: 'table_chart',
          children: [WIDGET_TYPES.userStatistics, WIDGET_TYPES.alarmStatistics].map(
            widgetTypePreparer,
          ),
        },
      );

      return types.reduce((acc, type) => {
        if (!type.children) {
          const rule = WIDGET_TYPES_RULES[type.type];

          if (!rule?.edition || rule.edition === this.edition) {
            acc.push(type);
          }

          return acc;
        }

        const children = type.children.filter((child) => {
          const rule = WIDGET_TYPES_RULES[child.type];

          return !rule?.edition || rule.edition === this.edition;
        });

        if (children.length) {
          acc.push({
            ...type,
            children,
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
      const widget = this.getWidgetWithUpdatedGridParametersByType(type);

      this.$modals.hide({ id: this.modal.id });

      window.requestAnimationFrame(() => {
        this.$sidebar.show({
          name: SIDE_BARS_BY_WIDGET_TYPES[type],
          config: {
            widget,
          },
        });
      });
    },
  },
};
</script>
