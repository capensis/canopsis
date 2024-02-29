<template>
  <widget-settings-item :title="$t('settings.rootCauseSettings')">
    <c-enabled-field v-field="parameters.showRootCauseByStateClick" :label="$t('settings.showRootCauseByStateClick')" />
    <v-expand-transition>
      <c-color-indicator-field
        v-if="parameters.showRootCauseByStateClick"
        v-field="parameters.rootCauseColorIndicator"
        :label="$t('settings.entitiesColorIndicator')"
        class="mt-0"
        hide-details
      />
    </v-expand-transition>
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  model: {
    prop: 'parameters',
    event: 'input',
  },
  props: {
    parameters: {
      type: Object,
      required: true,
    },
  },
};
</script>
