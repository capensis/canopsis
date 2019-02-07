<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ title }}
      .font-italic.caption.ml-1 ({{ $t('common.default') }}: 1 {{ $tc('common.times.day') }})
    v-container
      v-layout
        v-flex
          v-text-field.pt-0(
          type="number",
          :value="durationValue",
          @input="updateDurationValue",
          v-validate="'required|integer|min:1'",
          data-vv-name="durationValue",
          :error-messages="errors.collect('durationValue')",
          )
        v-flex
          v-select.pt-0(
            :items="units",
            :value="durationUnit",
            @input="updateDurationUnit"
            v-validate="'required'",
            data-vv-name="durationUnit",
            :error-messages="errors.collect('durationUnit')"
          )
</template>

<script>
import { find } from 'lodash';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      units: [
        {
          text: this.$tc('common.times.hour'),
          value: 'h',
        },
        {
          text: this.$tc('common.times.day'),
          value: 'd',
        },
        {
          text: this.$tc('common.times.week'),
          value: 'w',
        },
        {
          text: this.$tc('common.times.month'),
          value: 'm',
        },
      ],
    };
  },
  computed: {
    durationValue() {
      return parseInt(this.value.slice(0, this.value.length - 1), 10);
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
