<template lang="pug">
  widget-settings(:submitting="submitting", @submit="submit")
    field-title(v-model="form.title")
    v-divider
    field-preset(v-model="form.parameters", :type="form.type")
    v-divider
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldPreset from '@/components/sidebars/settings/fields/chart/preset.vue';

import WidgetSettings from './partials/widget-settings.vue';

export default {
  name: SIDE_BARS.textSettings,
  components: {
    WidgetSettings,
    FieldTitle,
    FieldPreset,
  },
  mixins: [widgetSettingsMixin],
};
</script>
