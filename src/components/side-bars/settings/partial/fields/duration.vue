<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.duration') }}
    v-container
      v-layout
        v-flex
          v-text-field.pt-0(
          type="number",
          :value="durationValue",
          @input="updateDurationValue",
          hide-details
          )
        v-flex
          v-select.pt-0(
            :items="units",
            :value="durationUnit",
            hide-details,
            @input="updateDurationUnit"
          )
</template>

<script>
import find from 'lodash/find';

export default {
  props: {
    value: {
      type: String,
    },
  },
  data() {
    return {
      units: [
        {
          text: 'Hour',
          value: 'h',
        },
        {
          text: 'Day',
          value: 'd',
        },
        {
          text: 'Week',
          value: 'w',
        },
        {
          text: 'Month',
          value: 'm',
        },
      ],
    };
  },
  computed: {
    durationValue() {
      return this.value.slice(0, this.value.length - 1);
    },
    durationUnit() {
      return find(this.units, item => item.value === this.value.slice(-1));
    },
  },
  methods: {
    updateDurationValue(event) {
      this.$emit('input', `${event}${this.durationUnit.value}`);
    },
    updateDurationUnit(event) {
      this.$emit('input', `${this.durationValue}${event}`);
    },
  },
};
</script>
