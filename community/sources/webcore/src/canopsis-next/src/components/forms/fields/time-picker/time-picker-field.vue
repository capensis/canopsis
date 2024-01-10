<template>
  <v-combobox
    class="time-picker__select"
    ref="combobox"
    :value="value"
    :items="items"
    :menu-props="menuProps"
    :return-object="false"
    :filter="filter"
    :label="label"
    :disabled="disabled"
    :error="error"
    placeholder="−−:−−"
    append-icon=""
    hide-details
    @change="change"
  />
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

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
    disabled: {
      type: Boolean,
      default: false,
    },
    error: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    menuProps() {
      return {
        auto: true,
        minWidth: 90,
        maxHeight: 200,
        scrollCalculator: this.scrollCalculator,
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
    scrollCalculator(el) {
      if (!this.value) {
        return el.scrollTop;
      }

      const maxScrollTop = el.scrollHeight - el.offsetHeight;
      const index = this.items.findIndex(item => item >= this.value);
      const elements = el.querySelectorAll('.v-list-item');

      const activeTile = elements[index === -1 ? elements.length - 1 : index];

      const newScrollTop = (activeTile.offsetTop - (el.offsetHeight / 2)) + (activeTile.offsetHeight / 2);

      return Math.min(maxScrollTop, Math.max(0, newScrollTop));
    },

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

<style lang="scss">
.time-picker__select {
  display: inline-block;
  width: 56px;
  min-width: 56px;

  &.v-select--is-menu-active .v-input__slot {
    background: #686868;
  }

  .v-input__slot {
    padding: 0 5px;
  }
}
</style>
