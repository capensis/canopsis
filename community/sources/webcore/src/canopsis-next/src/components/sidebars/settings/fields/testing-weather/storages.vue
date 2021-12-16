<template lang="pug">
  v-list-group(:disabled="disabled", hide-on-disabled)
    v-list-tile(slot="activator") {{ title }}
      div.font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container.pr-2
      c-storages-field(
        v-field="storages",
        v-on="$listeners",
        :label="$t('settings.serverStorage')",
        :disabled="disabled",
        :help-text="helpText"
      )
</template>

<script>
export default {
  model: {
    prop: 'storages',
    event: 'input',
  },
  props: {
    storages: {
      type: Array,
      default: () => [],
    },
    helpText: {
      type: String,
      required: false,
    },
    title: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
