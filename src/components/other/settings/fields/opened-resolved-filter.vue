<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.filterOnOpenResolved')}}
    v-container
      v-layout
        v-radio-group(v-model="value", @change="$emit('input', $event)")
          v-radio(:value="ALARM_FILTER_STATES.opened", :label="$t('settings.open')", hide-details)
          v-radio(:value="ALARM_FILTER_STATES.resolved", :label="$t('settings.resolved')", hide-details)
</template>

<script>
import { ALARM_FILTER_STATES } from '@/constants';

export default {
  props: {
    value: {
      type: String,
      default: ALARM_FILTER_STATES.opened,
    },
  },
  data() {
    return {
      ALARM_FILTER_STATES,
    };
  },
};
</script>
