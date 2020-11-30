<template lang="pug">
  v-combobox.time-picker__select(
    ref="combobox",
    :value="value",
    :items="items",
    :menu-props="menuProps",
    :return-object="false",
    :filter="filter",
    :label="label",
    append-icon="",
    hide-details,
    @change="change"
  )
</template>

<script>
import formBaseMixin from '@/mixins/form/base';

/**
 * TODO: Move into another place
 *
 * @type {number}
 */
const HOURS_IN_DAY = 24;
const MINUTES_STEP = 15;
const MINUTES_STEPS_IN_HOUR = Math.floor(60 / MINUTES_STEP);

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '12:00',
    },
    label: {
      type: String,
      default: null,
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    menuProps() {
      return {
        minWidth: 90,
        maxHeight: 200,
      };
    },

    items() {
      if (this.roundHours) {
        return new Array(HOURS_IN_DAY).fill(0).map((item, index) => `${index < 10 ? `0${index}` : index}:00`);
      }

      return new Array(HOURS_IN_DAY * MINUTES_STEPS_IN_HOUR).fill(0).map((item, index) => {
        const hoursIndex = Math.floor(index / MINUTES_STEPS_IN_HOUR);
        const minutesIndex = index - (hoursIndex * MINUTES_STEPS_IN_HOUR);

        const hours = hoursIndex < 10 ? `0${hoursIndex}` : hoursIndex;
        const minutes = minutesIndex === 0 ? '00' : minutesIndex * MINUTES_STEP;

        return `${hours}:${minutes}`;
      });
    },
  },
  methods: {
    filter(item, queryText, itemText) {
      return itemText.toLocaleLowerCase().startsWith(queryText.toLocaleLowerCase());
    },
    change(value) {
      const result = value.match(/(\d{2}):(\d{2})/);
      let preparedValue = value;

      if (result && this.roundHours) {
        const [, hours] = result;

        preparedValue = `${hours}:00`;
      } else if (!result) {
        preparedValue = this.value;
      }

      if (value !== preparedValue) {
        this.$refs.combobox.setValue(preparedValue);
        this.$refs.combobox.setSearch('');
      }

      if (this.value !== preparedValue) {
        this.updateModel(preparedValue);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.time-picker__select {
  display: inline-block;
  width: 56px;

  &.v-select--is-menu-active .v-input__slot {
    background: #686868;
  }

  & /deep/ {
    .v-select__slot {
      padding: 0 5px;
    }
  }
}
</style>
