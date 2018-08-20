<template lang='pug'>
  v-card
    v-card-title
      v-layout(justify-space-between, align-center)
        h2 Select a widget
        v-btn(@click='hideModal', icon, small)
          v-icon close
    v-card-text
      v-layout(row)
        v-flex(xs12 sm8 offset-sm2)
          v-list
            v-list-tile(
            v-for='widgetType in widgetsTypes',
            :key='widgetType.title',
            @click='selectWidgetType(widgetType.value)'
            )
              v-list-tile-action
                v-icon {{ widgetType.icon }}
              v-list-tile-content
                v-list-tile-title {{ widgetType.title }}
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
        { title: WIDGET_TYPES.alarmList, icon: 'list', value: WIDGET_TYPES.alarmList },
        { title: WIDGET_TYPES.context, icon: 'list', value: WIDGET_TYPES.context },
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
