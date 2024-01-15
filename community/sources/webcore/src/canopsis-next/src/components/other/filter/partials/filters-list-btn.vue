<template>
  <c-action-btn
    :tooltip="$t('filter.selector.buttons.list')"
    icon="filter_list"
    @click="showFiltersListModal"
  />
</template>

<script>
import { MODALS } from '@/constants';

export default {
  props: {
    widgetId: {
      type: String,
      required: true,
    },
    private: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    editable: {
      type: Boolean,
      default: false,
    },
    withAlarm: {
      type: Boolean,
      default: false,
    },
    withEntity: {
      type: Boolean,
      default: false,
    },
    withPbehavior: {
      type: Boolean,
      default: false,
    },
    withServiceWeather: {
      type: Boolean,
      default: false,
    },
    entityTypes: {
      type: Array,
      required: false,
    },
    entityCountersType: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    showFiltersListModal() {
      this.$modals.show({
        name: MODALS.filtersList,
        config: {
          widgetId: this.widgetId,
          private: this.private,
          addable: this.addable,
          editable: this.editable,
          withAlarm: this.withAlarm,
          withEntity: this.withEntity,
          withPbehavior: this.withPbehavior,
          withServiceWeather: this.withServiceWeather,
          entityTypes: this.entityTypes,
          entityCountersType: this.entityCountersType,
        },
      });
    },
  },
};
</script>
