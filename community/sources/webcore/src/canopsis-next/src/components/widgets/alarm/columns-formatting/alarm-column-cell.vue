<template>
  <v-menu
    class="alarms-column-cell"
    v-if="column.popupTemplate"
    v-model="opened"
    :close-on-content-click="false"
    :open-on-click="false"
    offset-x
    lazy
  >
    <template #activator="{ on }">
      <v-layout
        v-on="on"
        d-inline-flex
        align-center
      >
        <c-compiled-template
          v-if="column.isHtml"
          :template="value"
        />
        <div
          v-else
          v-bind="component.bind"
          v-on="component.on"
          :is="component.is"
        />
        <v-btn
          class="ma-0 alarms-column-cell__show-info-btn"
          :class="{ 'alarms-column-cell__show-info-btn--small': small }"
          icon
          small
          @click.stop="showInfoPopup"
        >
          <v-icon small>
            info
          </v-icon>
        </v-btn>
      </v-layout>
    </template>
    <alarm-column-cell-popup-body
      :alarm="alarm"
      :template="column.popupTemplate"
      @close="hideInfoPopup"
    />
  </v-menu>
  <c-compiled-template
    v-else-if="column.isHtml"
    :template="value"
  />
  <div
    v-else
    v-bind="component.bind"
    v-on="component.on"
    :is="component.is"
  />
</template>

<script>
import { get } from 'lodash';

import ColorIndicatorWrapper from '@/components/common/table/color-indicator-wrapper.vue';

import AlarmColumnCellPopupBody from './alarm-column-cell-popup-body.vue';
import AlarmColumnValueState from './alarm-column-value-state.vue';
import AlarmColumnValueStatus from './alarm-column-value-status.vue';
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
    AlarmColumnCellPopupBody,
    AlarmColumnValueState,
    AlarmColumnValueStatus,
    AlarmColumnValueExtraDetails,
    ColorIndicatorWrapper,
  },
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
    selectedTag: {
      type: String,
      default: '',
    },
    small: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      opened: false,
    };
  },
  computed: {
    value() {
      const value = get(this.alarm, this.column.value, '');

      return this.column.filter ? this.column.filter(value) : value;
    },

    component() {
      return this.column.getComponent(this);
    },
  },
  methods: {
    showInfoPopup() {
      this.opened = true;
    },

    hideInfoPopup() {
      this.opened = false;
    },
  },
};
</script>

<style lang="scss">
.alarms-column-cell {
  &__show-info-btn {
    &--small {
      width: 22px;
      max-width: 22px;
      height: 22px;
    }
  }
}
</style>
