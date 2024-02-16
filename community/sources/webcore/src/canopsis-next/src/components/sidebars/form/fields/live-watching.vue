<template>
  <c-enabled-field
    v-field="value"
    :label="$t('settings.liveWatching')"
    hide-details
  >
    <template #append="">
      <c-help-icon
        :text="$t('settings.liveWatchingTooltip')"
        icon-class="mt-1"
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
