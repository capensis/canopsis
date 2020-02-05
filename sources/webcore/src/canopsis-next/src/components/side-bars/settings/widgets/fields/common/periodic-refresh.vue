<template lang="pug">
  v-list-group(data-test="periodicRefresh")
    v-list-tile(slot="activator") {{$t('settings.periodicRefresh')}}
      .font-italic.caption.ml-1 ({{ $t('common.optional') }})
    v-container
      periodic-refresh-field(v-field="value")
</template>

<script>
import { DEFAULT_PERIODIC_REFRESH } from '@/constants';

import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';

export default {
  components: { PeriodicRefreshField },
  props: {
    value: {
      type: Object,
      default: DEFAULT_PERIODIC_REFRESH,
    },
  },
};
</script>
