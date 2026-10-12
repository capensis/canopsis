<template>
  <widget-settings-item
    :title="$t('settings.filenameRecognition')"
    optional
  >
    <c-file-name-mask-field
      :model-value="modelValue.screenshot_filemask"
      :title="$t('settings.screenshotMask.title')"
      :help-text="$t('settings.screenshotMask.helpText')"
      :placeholder="$constants.TEST_CASE_FILE_MASK"
      name="screenshot_filemask"
      required
      @update:model-value="updateField('screenshot_filemask', $event)"
    />
    <c-file-name-mask-field
      :model-value="modelValue.video_filemask"
      :title="$t('settings.videoMask.title')"
      :help-text="$t('settings.videoMask.helpText')"
      :placeholder="$constants.TEST_CASE_FILE_MASK"
      name="video_filemask"
      required
      @update:model-value="updateField('video_filemask', $event)"
    />
    <c-file-name-mask-field
      :model-value="modelValue.report_fileregexp"
      :title="$t('settings.reportFileRegexp.title')"
      :help-text="$t('settings.reportFileRegexp.helpText')"
      :placeholder="$constants.TEST_CASE_FILE_NAME"
      name="report_fileregexp"
      @update:model-value="updateField('report_fileregexp', $event)"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:modelValue'],

  methods: {
    updateField(field, value) {
      this.$emit('update:modelValue', {
        ...this.modelValue,
        [field]: value,
      });
    },
  },
};
</script>
