<template lang="pug">
  widget-settings(:submitting="submitting", @submit="submit")
    field-title(v-model="form.title")
    v-divider
    field-text-editor(
      v-model="form.parameters.template",
      :title="$t('common.template')"
    )
    v-divider
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldTextEditor from '@/components/sidebars/settings/fields/common/text-editor.vue';

import WidgetSettings from './partials/widget-settings.vue';

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
