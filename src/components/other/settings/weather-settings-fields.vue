<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title")
      div DataSet
      div Weather template
      div Modal template
      div Entity template
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{ $t('common.save') }}
</template>

<script>
import FieldTitle from '@/components/other/settings/fields/title.vue';

import widgetSettingsMixin from '@/mixins/widget/settings';

export default {
  components: {
    FieldTitle,
  },
  mixins: [widgetSettingsMixin],
  data() {
    return {
      settings: {
        widget: {
          title: this.widget.title,
        },
      },
    };
  },
};
</script>

