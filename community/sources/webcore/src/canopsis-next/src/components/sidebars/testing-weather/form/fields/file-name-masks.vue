<template>
  <widget-settings-item
    :title="$t('settings.filenameRecognition')"
    optional
  >
    <c-file-name-mask-field
      v-field="value.screenshot_filemask"
      :title="$t('settings.screenshotMask.title')"
      :help-text="$t('settings.screenshotMask.helpText')"
      :placeholder="$constants.TEST_CASE_FILE_MASK"
      name="screenshot_filemask"
      required
    />
    <c-file-name-mask-field
      v-field="value.video_filemask"
      :title="$t('settings.videoMask.title')"
      :help-text="$t('settings.videoMask.helpText')"
      :placeholder="$constants.TEST_CASE_FILE_MASK"
      name="video_filemask"
      required
    />
    <c-file-name-mask-field
      v-field="value.report_fileregexp"
      :title="$t('settings.reportFileRegexp.title')"
      :help-text="$t('settings.reportFileRegexp.helpText')"
      :placeholder="$constants.TEST_CASE_FILE_NAME"
      name="report_fileregexp"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
};
</script>
