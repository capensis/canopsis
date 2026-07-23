<template>
  <v-combobox
    v-model:search="search"
    ref="combobox"
    :model-value="value"
    :items="filteredItems"
    :menu-props="menuProps"
    :return-object="false"
    :custom-filter="filter"
    :label="label"
    :disabled="disabled"
    :error="error"
    :menu-icon="null"
    class="time-picker__select"
    placeholder="--:--"
    append-icon=""
    variant="underlined"
    hide-details
  />
</template>

<script>
import { normalizeTimeString } from '@/helpers/date/date-time-picker';
import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  emits: ['update:value'],
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
    allowedTime: {
      type: Function,
      required: false,
    },
  },
  data() {
    return {
      search: this.value,
    };
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
        return Array.from(
          { length: 24 },
          (item, index) => `${index < 10 ? `0${index}` : index}:00`,
        );
      }

      const { stepsInHours } = this;
      const minutesStep = Math.floor(60 / stepsInHours);

      return Array.from({ length: 24 * stepsInHours }, (item, index) => {
        const hoursIndex = Math.floor(index / stepsInHours);
        const minutesIndex = index - hoursIndex * stepsInHours;

        const hours = hoursIndex < 10 ? `0${hoursIndex}` : hoursIndex;
        const minutes = minutesIndex === 0 ? '00' : minutesIndex * minutesStep;

        return `${hours}:${minutes}`;
      });
    },

    filteredItems() {
      return this.allowedTime ? this.items.filter(this.allowedTime) : this.items;
    },
  },
  watch: {
    value(value) {
      if (value !== this.search) {
        this.search = value;
      }
    },

    search(value) {
      this.change(value);
    },
  },
  methods: {
    isValidTime(value) {
      return /^(([01][0-9])|(2[0-3])):([0-5][0-9])$/.test(value);
    },

    filter(item, queryText, itemText) {
      const text = normalizeTimeString(itemText ?? item, '');

      return text.toLocaleLowerCase().startsWith(String(queryText ?? '').toLocaleLowerCase());
    },

    change(value) {
      const nextValue = normalizeTimeString(value, '');

      if (!this.isValidTime(nextValue)) {
        return;
      }

      let preparedValue = nextValue;

      if (this.roundHours) {
        const [hours] = nextValue.split(':');

        preparedValue = `${hours}:00`;
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
  height: 30px;
  min-height: 30px;
  max-height: 30px;
  overflow: hidden;

  .v-input__control,
  .v-field,
  .v-field__field,
  .v-field__input {
    min-height: 30px;
    height: 30px;
    max-height: 30px;
  }

  .v-field__input {
    align-items: center;
    flex-wrap: nowrap;
    padding-top: 0;
    padding-bottom: 0;
    line-height: 30px;
    overflow: hidden;
  }

  .v-combobox__selection-text {
    white-space: normal;
    text-overflow: unset;
    overflow: visible;
  }

  .v-field {
    padding: 0 5px;
    line-height: inherit;
    font-size: inherit !important;
    font-weight: inherit;
  }

  .v-field__input,
  .v-field__input input,
  .v-combobox__selection,
  .v-combobox__selection-text {
    line-height: inherit !important;
    font-size: inherit !important;
    font-weight: inherit !important;
  }
}
</style>
