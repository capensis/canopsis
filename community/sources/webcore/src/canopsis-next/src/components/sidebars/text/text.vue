<template>
  <widget-settings
    :submitting="submitting"
    @submit="submit"
  >
    <field-title v-model="form.title" />
    <v-divider />
    <field-text-editor
      v-model="form.parameters.template"
      :title="$t('common.template')"
    />
    <v-divider />
  </widget-settings>
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import FieldTitle from '../form/fields/title.vue';
import FieldTextEditor from '../form/fields/text-editor.vue';
import WidgetSettings from '../partials/widget-settings.vue';

export default {
  name: SIDE_BARS.textSettings,
  components: {
    WidgetSettings,
    FieldTitle,
    FieldTextEditor,
  },
  mixins: [widgetSettingsMixin],
};
</script>
