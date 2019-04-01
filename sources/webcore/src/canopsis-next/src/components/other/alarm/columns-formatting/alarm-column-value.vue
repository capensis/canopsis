<template lang="pug">
  div
    v-menu(
    v-if="popupData",
    v-model="isInfoPopupOpen",
    :close-on-content-click="false",
    :open-on-click="false",
    offset-y
    )
      div(slot="activator")
        div(v-bind="component.bind", v-on="component.on")
      v-card(dark)
        v-card-title.primary.pa-2.white--text
          h4 {{ $t('alarmList.infoPopup') }}
        v-card-text.pa-2(v-html="popupTextContent")
    div(v-else, v-bind="component.bind", v-on="component.on")
</template>

<script>
import { get } from 'lodash';

import { compile } from '@/helpers/handlebars';
import popupMixin from '@/mixins/popup';

import State from '@/components/other/alarm/columns-formatting/alarm-column-value-state.vue';
import Links from '@/components/other/alarm/columns-formatting/alarm-column-value-links.vue';
import Link from '@/components/other/alarm/columns-formatting/alarm-column-value-link.vue';
import ExtraDetails from '@/components/other/alarm/columns-formatting/alarm-column-value-extra-details.vue';
import Ellipsis from '@/components/tables/ellipsis.vue';

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
    Links,
    Link,
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
  data() {
    return {
      isInfoPopupOpen: false,
    };
  },
  computed: {
    popupData() {
      const popups = get(this.widget.parameters, 'infoPopups', []);

      return popups.find(popup => popup.column === this.column.value);
    },
    popupTextContent() {
      if (this.popupData) {
        return compile(this.popupData.template, { alarm: this.alarm, entity: this.alarm.entity || {} });
      }
      return '';
    },
    columnFilter() {
      const PROPERTIES_FILTERS_MAP = {
        'v.status.val': value => this.$t(`tables.alarmStatus.${value}`),
        'v.last_update_date': value => this.$options.filters.date(value, 'long'),
        'v.creation_date': value => this.$options.filters.date(value, 'long'),
        'v.last_event_date': value => this.$options.filters.date(value, 'long'),
        'v.state.t': value => this.$options.filters.date(value, 'long'),
        'v.status.t': value => this.$options.filters.date(value, 'long'),
        'v.resolved': value => this.$options.filters.date(value, 'long'),
        t: value => this.$options.filters.date(value, 'long'),
      };

      return PROPERTIES_FILTERS_MAP[this.column.value];
    },
    component() {
      const PROPERTIES_COMPONENTS_MAP = {
        'v.state.val': {
          bind: {
            is: 'state',
            alarm: this.alarm,
          },
        },
        links: {
          bind: {
            is: 'links',
            alarm: this.alarm,
          },
        },
        extra_details: {
          bind: {
            is: 'extra-details',
            alarm: this.alarm,
          },
        },
      };

      if (PROPERTIES_COMPONENTS_MAP[this.column.value]) {
        return PROPERTIES_COMPONENTS_MAP[this.column.value];
      }

      if (this.column.value.startsWith('links.')) {
        return {
          bind: {
            is: 'link',
            link: this.$options.filters.get(this.alarm, this.column.value),
          },
        };
      }

      return {
        bind: {
          is: 'ellipsis',
          text: this.$options.filters.get(this.alarm, this.column.value, this.columnFilter, ''),
        },
        on: {
          textClicked: this.showInfoPopup,
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
  },
};
</script>
