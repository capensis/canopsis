<template>
  <v-menu
    v-if="column.popupTemplate"
    v-model="opened"
    :close-on-content-click="false"
    :open-on-click="false"
    class="alarms-column-cell"
    offset-x
    disable-resize
  >
    <template #activator="{ on }">
      <v-layout
        class="alarms-column-cell__layout"
        d-inline-flex
        align-center
        v-on="on"
      >
        <div
          v-if="column.isHtml"
          v-html="sanitizedValue"
          :class="textClass"
          v-on="componentOn"
        />
        <div
          v-else
          v-bind="component.bind"
          :is="component.is || component.bind?.is"
          :class="textClass"
          v-on="componentOn"
        />
        <v-btn
          :class="{ 'alarms-column-cell__show-info-btn--small': small }"
          class="ma-0 alarms-column-cell__show-info-btn"
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
      :selected-tags="selectedTags"
      :template="column.popupTemplate"
      :template-id="popupTemplateId"
      @close="hideInfoPopup"
      @select:tag="$emit('select:tag', $event)"
      @remove:tag="$emit('remove:tag', $event)"
    />
  </v-menu>
  <div
    v-else-if="column.isHtml"
    v-html="sanitizedValue"
    :class="textClass"
    v-on="componentOn"
  />
  <div
    v-else
    v-bind="component.bind"
    :is="component.is || component.bind?.is"
    :class="textClass"
    v-on="componentOn"
  />
</template>

<script>
import { get, flow, isObject } from 'lodash';
import { ref, computed, inject } from 'vue';

import { sanitizeHtml, linkifyHtml } from '@/helpers/html';
import { getAlarmWidgetColumnPopupTemplateId } from '@/helpers/entities/alarm/list';

import ColorIndicatorWrapper from '@/components/common/table/color-indicator-wrapper.vue';
import AlarmsExpandPanel from '@/components/widgets/alarm/expand-panel/alarms-expand-panel.vue';

import AlarmColumnCellPopupBody from './alarm-column-cell-popup-body.vue';
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
    AlarmsExpandPanel,
    AlarmColumnCellPopupBody,
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
    selectedTags: {
      type: Array,
      default: () => [],
    },
    small: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit, listeners }) {
    const selectAdvancedSearchField = inject('$selectAdvancedSearchField', () => {});

    const opened = ref(false);

    const value = computed(() => {
      const result = get(props.alarm, props.column.value, '');

      return props.column.filter ? props.column.filter(result) : result;
    });

    const sanitizedValue = computed(() => sanitizeHtml(linkifyHtml(String(value.value ?? ''))));

    /**
     * Applies advanced search filter for the column value when clicked
     * Selects the appropriate field and value based on column configuration
     *
     * @param {Event} event - The click event that triggered the filter
     * @returns {Event} The original event for potential chaining
     */
    const applyFilter = (event) => {
      selectAdvancedSearchField(props.column.value, props.column.isHtml ? sanitizedValue.value : value.value);

      return event;
    };

    const component = computed(() => (
      props.column.getComponent({ ...props, value: value.value, $emit: emit, $listeners: listeners })
    ));

    const componentOn = computed(() => {
      const result = { ...component.value?.on };

      if (props.column.isFilter) {
        result.click = component.value?.on?.click ? flow(applyFilter, component.value.on?.click) : applyFilter;
      }

      return result;
    });

    const popupTemplateId = computed(() => getAlarmWidgetColumnPopupTemplateId(props.widget._id, props.column.value));
    const textClass = computed(() => {
      let result = {
        'alarms-column-cell__filter': props.column.isFilter,
      };

      if (isObject(component.value?.bind?.class)) {
        result = { ...result, ...component.value?.bind?.class };
      } else {
        result[component.value?.bind?.class] = !!component.value?.bind?.class;
      }

      return result;
    });

    const showInfoPopup = () => opened.value = true;
    const hideInfoPopup = () => opened.value = false;

    return {
      opened,
      component,
      sanitizedValue,
      popupTemplateId,
      textClass,
      showInfoPopup,
      hideInfoPopup,
      componentOn,
    };
  },
};
</script>

<style lang="scss">
.alarms-column-cell {
  &__show-info-btn {
    flex-shrink: 0 !important;

    &--small {
      width: 22px;
      height: 22px;
      max-width: 22px;
      max-height: 22px;
    }
  }

  &__layout {
    max-width: 100%;
  }

  &__filter {
    cursor: pointer !important;

    &:not(.c-alarm-state-chip):hover, &.c-alarm-state-chip:hover .chip {
      text-decoration: underline !important;
    }
  }
}
</style>
