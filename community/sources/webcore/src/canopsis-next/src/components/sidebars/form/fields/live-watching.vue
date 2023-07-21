<template lang="pug">
  c-enabled-field(
    v-field="value",
    :label="$t('settings.fields.liveWatching')",
    hide-details
  )
    template(#append="")
      c-help-icon(
        :text="$t('settings.fields.liveWatchingTooltip')",
        color="info",
        top
      )
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
