<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.filenameRecognition') }}
      span.font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container
      c-file-name-mask-field(
        v-field="value.screenshot_filemask",
        :title="$t('settings.screenshotMask.title')",
        :help-text="$t('settings.screenshotMask.helpText')",
        :placeholder="$constants.TEST_CASE_FILE_MASK",
        name="screenshot_filemask"
      )
      c-file-name-mask-field(
        v-field="value.video_filemask",
        :title="$t('settings.videoMask.title')",
        :help-text="$t('settings.videoMask.helpText')",
        :placeholder="$constants.TEST_CASE_FILE_MASK",
        name="video_filemask"
      )
</template>

<script>
export default {
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
};
</script>
