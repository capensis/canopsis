<template lang="pug">
  v-combobox.time-picker__select(
    ref="combobox",
    :value="value",
    :items="items",
    :menu-props="menuProps",
    :return-object="false",
    :filter="filter",
    :label="label",
    placeholder="−−:−−",
    append-icon="",
    hide-details,
    @change="change"
  )
</template>

<script>
import formBaseMixin from '@/mixins/form/base';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: null,
    },
    label: {
      type: String,
      default: null,
    },
    stepsInHours: {
      type: Number,
      default: 4,
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
        return new Array(24).fill(0).map((item, index) => `${index < 10 ? `0${index}` : index}:00`);
      }

      const { stepsInHours } = this;
      const minutesStep = Math.floor(60 / stepsInHours);

      return new Array(24 * stepsInHours).fill(0).map((item, index) => {
        const hoursIndex = Math.floor(index / stepsInHours);
        const minutesIndex = index - (hoursIndex * stepsInHours);

        const hours = hoursIndex < 10 ? `0${hoursIndex}` : hoursIndex;
        const minutes = minutesIndex === 0 ? '00' : minutesIndex * minutesStep;

        return `${hours}:${minutes}`;
      });
    },
  },
  methods: {
    filter(item, queryText, itemText) {
      return itemText.toLocaleLowerCase().startsWith(queryText.toLocaleLowerCase());
    },

    change(value) {
      const result = value.match(/(([01][0-9])|(2[0-3])):([0-5][0-9])/);
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
