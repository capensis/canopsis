<template lang="pug">
  v-layout.pl-3(justify-end, align-center)
    div {{ $t('common.widgetId') }}: {{ widgetId }}
    v-tooltip(left)
      v-btn(
        v-clipboard:copy="widgetId",
        v-clipboard:success="addWidgetIdCopiedSuccessPopup",
        v-clipboard:error="addWidgetIdCopiedErrorPopup",
        slot="activator",
        small,
        fab,
        icon,
        depressed
      )
        v-icon file_copy
      span {{ $t('modals.variablesHelp.copyToClipboard') }}
</template>

<script>

export default {
  props: {
    widgetId: {
      type: String,
      required: true,
    },
  },
  methods: {
    addWidgetIdCopiedSuccessPopup() {
      this.$popups.success({ text: this.$t('success.widgetIdCopied') });
    },

    addWidgetIdCopiedErrorPopup() {
      this.$popups.error({ text: this.$t('errors.default') });
    },
  },
};
</script>
