<template>
  <c-enabled-field
    v-field="value"
    hide-details
  >
    <template #label="">
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
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
