<template lang="pug">
  div
    v-menu(
      v-if="popupData",
      v-model="isInfoPopupOpen",
      :close-on-content-click="false",
      :open-on-click="false",
      offset-x
    )
      div(slot="activator")
        v-layout(align-center)
          div(v-if="column.isHtml", v-html="sanitizedValue")
          div(v-else, v-bind="component.bind", v-on="component.on")
          v-btn.ma-0(icon, small, @click.stop="showInfoPopup")
            v-icon(small) info
      v-card(dark)
        v-card-title.primary.pa-2.white--text
          v-layout(justify-space-between, align-center)
            h4 {{ $t('alarmList.infoPopup') }}
            v-btn.ma-0.ml-3(icon, small, @click="hideInfoPopup", color="white")
              v-icon(small, color="error") close
        v-card-text.pa-2(v-html="popupTextContent")
    div(v-else-if="column.isHtml", v-html="sanitizedValue")
    div(v-else, v-bind="component.bind", v-on="component.on")
</template>

<script>
import { get } from 'lodash';

import { compile } from '@/helpers/handlebars';
import popupMixin from '@/mixins/popup';

import Ellipsis from '@/components/tables/ellipsis.vue';

import AlarmColumnValueState from './alarm-column-value-state.vue';
import AlarmColumnValueLinks from './alarm-column-value-links.vue';
import AlarmColumnValueLink from './alarm-column-value-link.vue';
import AlarmColumnValueExtraDetails from './alarm-column-value-extra-details.vue';

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
    Ellipsis,
    AlarmColumnValueState,
    AlarmColumnValueLinks,
    AlarmColumnValueLink,
    AlarmColumnValueExtraDetails,
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
    value() {
      return this.$options.filters.get(this.alarm, this.column.value, this.columnFilter, '');
    },

    sanitizedValue() {
      try {
        return this.$sanitize(this.value, {
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
        'v.duration': value => this.$options.filters.duration({ value }),
        'v.current_state_duration': value => this.$options.filters.duration({ value }),
        t: value => this.$options.filters.date(value, 'long'),
      };

      return PROPERTIES_FILTERS_MAP[this.column.value];
    },

    component() {
      const PROPERTIES_COMPONENTS_MAP = {
        'v.state.val': {
          bind: {
            is: 'alarm-column-value-state',
            alarm: this.alarm,
          },
        },
        links: {
          bind: {
            is: 'alarm-column-value-links',
            links: this.alarm.links,
          },
        },
        extra_details: {
          bind: {
            is: 'alarm-column-value-extra-details',
            alarm: this.alarm,
          },
        },
      };

      if (PROPERTIES_COMPONENTS_MAP[this.column.value]) {
        return PROPERTIES_COMPONENTS_MAP[this.column.value];
      }

      if (this.column.value.startsWith('links.')) {
        const category = this.column.value.slice(6);
        const links = {
          [category]: this.$options.filters.get(this.alarm, this.column.value, null, []),
        };

        return {
          bind: {
            links,

            is: 'alarm-column-value-links',
          },
        };
      }

      return {
        bind: {
          is: 'ellipsis',
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
