<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title
      v-divider
    v-btn(@click="", color="green darken-4 white--text", depressed) {{ $t('common.save') }}
</template>

<script>
import { SIDE_BARS } from '@/constants';
import widgetSettingsMixin from '@/mixins/widget/settings';

import FieldTitle from '../partial/fields/title.vue';

export default {
  name: SIDE_BARS.statsTableSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldTitle,
  },
  mixins: [widgetSettingsMixin],
};

</script>
