<template>
  <widget-settings-item :title="$t('settings.rootCauseSettings')">
    <c-enabled-field
      :model-value="modelValue.showRootCauseByStateClick"
      :label="$t('settings.showRootCauseByStateClick')"
      @update:model-value="updateField('showRootCauseByStateClick', $event)"
    />
    <v-expand-transition>
      <c-color-indicator-field
        v-if="modelValue.showRootCauseByStateClick"
        :model-value="modelValue.rootCauseColorIndicator"
        :label="$t('settings.entitiesColorIndicator')"
        class="mt-0"
        hide-details
        @update:model-value="updateField('rootCauseColorIndicator', $event)"
      />
    </v-expand-transition>
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import { useModelField } from '@/hooks/form/model-field';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    return {
      updateField,
    };
  },
};
</script>
