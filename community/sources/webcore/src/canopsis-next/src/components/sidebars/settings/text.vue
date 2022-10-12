<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="form.title")
      v-divider
      field-text-editor(
        v-model="form.parameters.template",
        :title="$t('settings.templateEditor')"
      )
      v-divider
    v-btn.primary(
      :loading="submitting",
      :disabled="submitting",
      @click="submit"
    ) {{ $t('common.save') }}
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { widgetSettingsMixin } from '@/mixins/widget/settings';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldTextEditor from '@/components/sidebars/settings/fields/common/text-editor.vue';

export default {
  name: SIDE_BARS.textSettings,
  components: {
    FieldTitle,
    FieldTextEditor,
  },
  mixins: [widgetSettingsMixin],
};
</script>
