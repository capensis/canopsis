<template>
  <c-enabled-field
    :model-value="modelValue"
    hide-details
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template #label>
      {{ $t('settings.liveWatching') }}
      <c-help-icon
        :text="$t('settings.liveWatchingTooltip')"
        icon-class="mt-1 ml-2"
        color="info"
        top
      />
    </template>
  </c-enabled-field>
</template>

<script>
export default {
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],
};
</script>
