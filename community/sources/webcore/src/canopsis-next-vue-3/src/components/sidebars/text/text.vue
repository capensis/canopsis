<template>
  <widget-settings :submitting="submitting" :dirty="hasChanges" @submit="submit">
    <field-title v-model="form.title" />

    <v-divider />

    <field-text-editor
      v-model="form.parameters.template"
      :title="t('common.template', 1)"
      :sanitize-options="sanitizeOptions"
    />

    <v-divider />
  </widget-settings>
</template>

<script setup>
import { toRef } from 'vue';

import { SIDE_BARS, TEXT_WIDGET_SANITIZE_OPTIONS } from '@/constants';
import { useI18n } from '@/hooks/i18n';
import { useWidgetSettings } from '@/hooks/widget/settings';

import FieldTextEditor from '../form/fields/text-editor.vue';
import FieldTitle from '../form/fields/title.vue';
import WidgetSettings from '../partials/widget-settings.vue';

defineOptions({
  name: SIDE_BARS.textSettings,
});

const props = defineProps({
  sidebar: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();

const sanitizeOptions = TEXT_WIDGET_SANITIZE_OPTIONS;
const sidebar = toRef(props, 'sidebar');

const { form, hasChanges, submitting, close, submit } = useWidgetSettings({
  sidebar,
});

defineExpose({
  close,
});
</script>
