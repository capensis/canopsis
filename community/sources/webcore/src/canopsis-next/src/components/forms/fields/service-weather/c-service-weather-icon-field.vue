<template lang="pug">
  v-select(
    v-field="value",
    :label="label",
    :items="availableIcons",
    :name="name",
    :error-messages="errors.collect(name)",
    item-value="icon"
  )
    template(#selection="{ item }")
      v-icon {{ item.icon }}
      span.ml-2 {{ item.text }}
    template(#item="{ item }")
      v-icon {{ item.icon }}
      span.ml-2 {{ item.text }}
</template>

<script>
import { WEATHER_ICONS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'icon',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    availableIcons() {
      return Object.entries(WEATHER_ICONS).map(([value, icon]) => ({
        icon,
        text: this.$te(`common.stateTypes.${value}`)
          ? this.$t(`common.stateTypes.${value}`)
          : this.$t(`serviceWeather.iconTypes.${value}`),
      }));
    },
  },
};
</script>
