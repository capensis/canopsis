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
import uuid from '@/helpers/uuid';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS, WIDGET_TYPES } from '@/constants';

const generateWidgetByType = (type) => {
  const id = uuid(`widget_${type}`);
  const widget = {
    id,
    widgetId: id,
    title: null,
    preference_id: uuid(),
    xtype: type,
    tagName: null,
    mixins: [],
    default_sort_column: {
      direction: 'ASC',
    },
    columns: [],
    popup: [],
  };

  if (type === WIDGET_TYPES.alarmList) { // TODO: move into constants
    widget.alarms_state_filter = null;
    widget.hide_resources = false;
    widget.widget_columns = [];
    widget.columns = [
      'connector_name',
      'component',
      'resource',
      'state',
      'status',
      'last_update_date',
      'extra_details',
    ];
  }

  return widget;
};

/**
 * Modal to add a time filter on alarm-list
 */
export default {
  name: MODALS.createWidget,
  mixins: [modalInnerMixin],
  data() {
    return {
      // TODO: add correct value
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
