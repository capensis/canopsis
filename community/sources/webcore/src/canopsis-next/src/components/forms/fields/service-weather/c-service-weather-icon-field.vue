<template lang="pug">
  v-select(
    v-field="value",
    :label="label",
    :items="availableIcons",
    :name="name",
    :error-messages="errors.collect(name)"
  )
    template(#selection="{ item }")
      v-icon {{ item.icon }}
      span.ml-2 {{ item.text }}
    template(#item="{ item }")
      v-icon {{ item.icon }}
      span.ml-2 {{ item.text }}
</template>

<script>
import { PBEHAVIOR_TYPE_TYPES, WEATHER_ICONS } from '@/constants';

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
    types: {
      type: Array,
      default: () => [
        PBEHAVIOR_TYPE_TYPES.maintenance,
        PBEHAVIOR_TYPE_TYPES.inactive,
        PBEHAVIOR_TYPE_TYPES.pause,
      ],
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    availableIcons() {
      return this.types.map(type => ({
        value: type,
        icon: WEATHER_ICONS[type],
        text: this.$t(`serviceWeather.iconTypes.${type}`),
      }));
    },
  },
};
</script>
