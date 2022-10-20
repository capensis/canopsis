<template lang="pug">
  v-layout.alarm-counters(column)
    v-tooltip(
      v-for="counter in preparedCounters.selected",
      :key="counter.key",
      top
    )
      template(#activator="{ on }")
        alarm-counter(
          v-on="on",
          :count="counter.count",
          :icon="counter.icon"
        )
      span {{ counter.name }}
    v-tooltip(:disabled="!otherCountersValue", top)
      template(#activator="{ on }")
        alarm-counter(
          v-on="on",
          :count="otherCountersValue",
          icon="more_horiz"
        )
      div(v-for="otherCounter in preparedCounters.other", :key="otherCounter.key")
        strong {{ otherCounter.name }}
        span : {{ otherCounter.count }}
</template>

<script>
import AlarmCounter from './alarm-counter.vue';

export default {
  components: {
    AlarmCounter,
  },
  props: {
    counters: {
      type: Object,
      default: () => ({}),
    },
    types: {
      type: Array,
      required: true,
    },
  },
  computed: {
    preparedCounters() {
      return (this.counters?.pbh_types ?? []).reduce((acc, { count, type }) => {
        if (this.types.includes(type._id)) {
          acc.selected.push({
            key: type._id,
            name: type.name,
            icon: type.icon_name,
            count,
          });
        } else {
          acc.other.push({
            key: type._id,
            name: type.name,
            count,
          });
        }

        return acc;
      }, {
        selected: [],
        other: [],
      });
    },

    otherCountersValue() {
      return this.preparedCounters.other.reduce((acc, { count }) => acc + count, 0);
    },
  },
};
</script>
