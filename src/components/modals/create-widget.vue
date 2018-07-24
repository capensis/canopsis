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
            @click='selectWidgetType(widgetType.value)',
            :key='widgetType.title'
            )
              v-list-tile-action
                v-icon {{ widgetType.icon }}
              v-list-tile-content
                v-list-tile-title {{ widgetType.title }}
</template>

<script>
import { MODALS, WIDGET_TYPES } from '@/constants';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { generateWidgetByType } from '@/helpers/entities';

/**
 * Modal to create widget
 */
export default {
  name: MODALS.createWidget,
  mixins: [modalInnerMixin],
  data() {
    return {
      widgetsTypes: [
        { title: WIDGET_TYPES.alarmList, icon: 'list', value: 'listalarm' },
        { title: WIDGET_TYPES.context, icon: 'list', value: 'crudcontext' },
      ],
    };
  },
  methods: {
    selectWidgetType(type) {
      const widgetWrapper = generateWidgetByType(type);

      if (this.config.action) {
        this.config.action(widgetWrapper);
      }

      this.hideModal();
    },
  },
};
</script>
