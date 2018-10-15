<template lang="pug">
  div
    div(v-if="component", :is="component", :alarm="alarm") {{ component.value }}
    ellipsis(
    v-else,
    :text="alarm | get(column.value, columnFilter, '')",
    :column="column.value",
    @textClicked="showPopup"
    ) {{ popupData }}
</template>

<script>
import get from 'lodash/get';
import Handlebars from 'handlebars';
import State from '@/components/other/alarm/columns-formatting/alarm-column-value-state.vue';
import ExtraDetails from '@/components/other/alarm/columns-formatting/alarm-column-value-extra-details.vue';
import Ellipsis from '@/components/tables/ellipsis.vue';
import popupMixin from '@/mixins/popup';

/**
 * Component to format alarms list columns
 *
 * @module alarm
 *
 * @prop {Object} alarm - Object representing the alarm
 * @prop {Object} widget - Object representing the widget
 * @prop {Object} column - Property concerned on the column
 */
export default {
  components: {
    State,
    ExtraDetails,
    Ellipsis,
  },
  mixins: [
    popupMixin,
  ],
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    column: {
      type: Object,
      required: true,
    },
  },
  computed: {
    popupData() {
      const popups = get(this.widget.parameters, 'infoPopups', []);

      return popups.find(popup => popup.column === this.column.value);
    },
    popupTextContent() {
      const template = Handlebars.compile(this.popupData.template);
      const context = { alarm: this.alarm.v };

      return template(context);
    },
    columnName() {
      return this.column.value.split('.')[1];
    },
    columnFilter() {
      const PROPERTIES_FILTERS_MAP = {
        'v.status.val': value => this.$t(`tables.alarmStatus.${value}`),
        'v.last_update_date': value => this.$options.filters.date(value, 'long'),
        'v.creation_date': value => this.$options.filters.date(value, 'long'),
        'v.last_event_date': value => this.$options.filters.date(value, 'long'),
        'v.state.t': value => this.$options.filters.date(value, 'long'),
        'v.status.t': value => this.$options.filters.date(value, 'long'),
        t: value => this.$options.filters.date(value, 'long'),
      };

      return PROPERTIES_FILTERS_MAP[this.column.value];
    },
    component() {
      const PROPERTIES_COMPONENTS_MAP = {
        'v.state.val': 'state',
        extra_details: 'extra-details',
      };

      return PROPERTIES_COMPONENTS_MAP[this.column.value];
    },
  },
  methods: {
    showPopup() {
      if (this.popupData) {
        this.addInfoPopup({
          text: this.popupTextContent,
        });
      }
    },
  },
};
</script>
