<template>
  <widget-settings-item :title="$t('view.periodicRefresh')">
    <periodic-refresh-field
      :periodic-refresh="modelValue"
      name="parameters.periodic_refresh"
      @update:periodic-refresh="$emit('update:modelValue', $event)"
    />
    <live-watching-field
      v-if="withLiveWatching"
      :model-value="liveWatching"
      @update:model-value="$emit('update:liveWatching', $event)"
    />
  </widget-settings-item>
</template>

<script>
import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

import LiveWatchingField from './live-watching.vue';

export default {
  components: {
    PeriodicRefreshField,
    WidgetSettingsItem,
    LiveWatchingField,
  },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    liveWatching: {
      type: Boolean,
      default: false,
    },
    withLiveWatching: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue', 'update:liveWatching'],
};
</script>
