<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsSelect') }}
    v-select(
    hide-details,
    :items="items",
    multiple,
    chips,
    v-model="value"
    )
      template(
      slot="selection",
      slot-scope="{item, index}"
      )
        v-chip(v-if="index < 3")
          span {{ item.text }}
        span.ml-2.grey--text.caption(v-if="index === 3") + {{ value.length - 3 }} others
</template>

<script>
import { STATS_TYPES } from '@/constants';

export default {
  data() {
    return {
      value: [],
    };
  },
  computed: {
    items() {
      return Object.values(STATS_TYPES).map(item => ({ value: item, text: this.$t(`stats.types.${item}`) }));
    },
  },
};
</script>

