<template>
  <v-layout
    class="alarm-counters"
    column
  >
    <v-tooltip
      v-for="counter in preparedCounters.selected"
      :key="counter.key"
      top
    >
      <template #activator="{ on }">
        <alarm-counter
          :count="counter.count"
          :icon="counter.icon"
          v-on="on"
        />
      </template>
      <span>{{ counter.name }}</span>
    </v-tooltip>
    <v-tooltip
      v-if="otherCountersValue"
      top
    >
      <template #activator="{ on }">
        <alarm-counter
          :count="otherCountersValue"
          icon="more_horiz"
          v-on="on"
        />
      </template>
      <div
        v-for="otherCounter in preparedCounters.other"
        :key="otherCounter.key"
      >
        <strong>{{ otherCounter.name }}</strong>
        <span>: {{ otherCounter.count }}</span>
      </div>
    </v-tooltip>
  </v-layout>
</template>

<script>
import AlarmCounter from './alarm-counter.vue';

export default {
  components: {
    AlarmCounter,
  },
  props: {
    counters: {
      type: Array,
      default: () => [],
    },
    types: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    preparedCounters() {
      return this.counters.reduce((acc, { count, type }) => {
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
