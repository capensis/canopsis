<template lang="pug">
  v-menu(
    v-if="popupData",
    v-model="isInfoPopupOpen",
    :close-on-content-click="false",
    :open-on-click="false",
    offset-x,
    lazy-with-unmount,
    lazy
  )
    div(slot="activator")
      v-layout(align-center)
        div(v-if="column.isHtml", v-html="sanitizedValue")
        div(v-else, v-bind="component.bind", v-on="component.on")
        v-btn.ma-0(data-test="alarmInfoPopupOpenButton", icon, small, @click.stop="showInfoPopup")
          v-icon(small) info
    alarm-column-cell-popup-body(
      :alarm="alarm",
      :template="popupData.template",
      @close="hideInfoPopup"
    )
  div(v-else-if="column.isHtml", v-html="sanitizedValue")
  div(v-else, v-bind="component.bind", v-on="component.on")
</template>

<script>
import { get } from 'lodash';
import sanitizeHTML from 'sanitize-html';

import { ALARM_ENTITY_FIELDS, COLOR_INDICATOR_TYPES } from '@/constants';

import { formToColumnValue } from '@/helpers/forms/widgets/alarm';

import { widgetColumnsFiltersMixin } from '@/mixins/widget/columns-filters';

import ColorIndicatorWrapper from '@/components/common/table/color-indicator-wrapper.vue';

import AlarmColumnCellPopupBody from './alarm-column-cell-popup-body.vue';
import AlarmColumnValueState from './alarm-column-value-state.vue';
import AlarmColumnValueStatus from './alarm-column-value-status.vue';
import AlarmColumnValueCategories from './alarm-column-value-categories.vue';
import AlarmColumnValueExtraDetails from './alarm-column-value-extra-details.vue';
import AlarmColumnValueLinks from './alarm-column-value-links.vue';

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
    AlarmColumnCellPopupBody,
    AlarmColumnValueState,
    AlarmColumnValueStatus,
    AlarmColumnValueCategories,
    AlarmColumnValueExtraDetails,
    AlarmColumnValueLinks,
    ColorIndicatorWrapper,
  },
  mixins: [widgetColumnsFiltersMixin],
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
    columnsFilters: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      isInfoPopupOpen: false,
    };
  },
  computed: {
    value() {
      return this.$options.filters.get(this.alarm, this.column.value, this.columnFilter, '');
    },

    sanitizedValue() {
      try {
        return sanitizeHTML(this.value, {
          allowedTags: ['h3', 'h4', 'h5', 'h6', 'blockquote', 'p', 'a', 'ul', 'ol',
            'nl', 'li', 'b', 'i', 'strong', 'em', 'strike', 'code', 'hr', 'br', 'div',
            'table', 'thead', 'caption', 'tbody', 'tr', 'th', 'td', 'pre', 'iframe', 'span', 'font', 'u'],
          allowedAttributes: {
            '*': ['style'],
            a: ['href', 'name', 'target'],
            img: ['src', 'alt'],
            font: ['color', 'size', 'face'],
          },
        });
      } catch (err) {
        console.warn(err);

        return '';
      }
    },

    popupData() {
      const popups = get(this.widget.parameters, 'infoPopups', []);

      return popups.find(popup => formToColumnValue(popup.column) === this.column.value);
    },

    columnFilter() {
      const PROPERTIES_FILTERS_MAP = {
        'v.last_update_date': value => this.$options.filters.date(value, 'long'),
        'v.creation_date': value => this.$options.filters.date(value, 'long'),
        'v.last_event_date': value => this.$options.filters.date(value, 'long'),
        'v.activation_date': value => this.$options.filters.date(value, 'long'),
        'v.state.t': value => this.$options.filters.date(value, 'long'),
        'v.status.t': value => this.$options.filters.date(value, 'long'),
        'v.resolved': value => this.$options.filters.date(value, 'long'),
        'v.duration': value => this.$options.filters.duration(value),
        'v.active_duration': value => this.$options.filters.duration(value),
        'v.current_state_duration': value => this.$options.filters.duration(value),
        'v.snooze_duration': value => this.$options.filters.duration(value),
        'v.pbh_inactive_duration': value => this.$options.filters.duration(value),
        t: value => this.$options.filters.date(value, 'long'),

        ...this.columnsFiltersMap,
      };

      return PROPERTIES_FILTERS_MAP[this.column.value];
    },

    component() {
      const PROPERTIES_COMPONENTS_MAP = {
        [ALARM_ENTITY_FIELDS.state]: {
          bind: {
            is: 'alarm-column-value-state',
            alarm: this.alarm,
          },
        },
        [ALARM_ENTITY_FIELDS.status]: {
          bind: {
            is: 'alarm-column-value-status',
            alarm: this.alarm,
          },
        },
        [ALARM_ENTITY_FIELDS.priority]: {
          bind: {
            is: 'color-indicator-wrapper',
            type: COLOR_INDICATOR_TYPES.impactState,
            entity: this.alarm.entity,
            alarm: this.alarm,
          },
        },
        [ALARM_ENTITY_FIELDS.impactState]: {
          bind: {
            is: 'color-indicator-wrapper',
            type: COLOR_INDICATOR_TYPES.impactState,
            entity: this.alarm.entity,
            alarm: this.alarm,
          },
        },
        links: {
          bind: {
            is: 'alarm-column-value-categories',
            asList: get(this.widget.parameters, 'linksCategoriesAsList.enabled', false),
            limit: get(this.widget.parameters, 'linksCategoriesAsList.limit'),
            links: this.alarm.links,
          },
        },
        [ALARM_ENTITY_FIELDS.extraDetails]: {
          bind: {
            is: 'alarm-column-value-extra-details',
            alarm: this.alarm,
          },
        },
      };

      const cell = PROPERTIES_COMPONENTS_MAP[this.column.value];

      if (cell) {
        return cell;
      }

      if (this.column.value.startsWith('links.')) {
        const links = get(this.alarm, this.column.value, []);

        return {
          bind: {
            links,

            is: 'alarm-column-value-links',
          },
        };
      }

      return {
        bind: {
          is: 'c-ellipsis',
          text: String(this.$options.filters.get(this.alarm, this.column.value, this.columnFilter, '')),
        },
      };
    },
  },
  methods: {
    showInfoPopup() {
      if (this.popupData) {
        this.isInfoPopupOpen = true;
      }
    },
    hideInfoPopup() {
      this.isInfoPopupOpen = false;
    },
  },
};
</script>
