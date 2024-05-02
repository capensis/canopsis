<template>
  <v-select
    v-field="value"
    :label="label"
    :items="availableIcons"
    :loading="pending"
    :name="name"
    :error-messages="errors.collect(name)"
    :disabled="disabled"
    item-value="icon"
  >
    <template #selection="{ item }">
      <v-icon :class="disabled ? 'text--disabled' : ''">
        {{ item.icon }}
      </v-icon>
      <span class="ml-2">{{ item.text }}</span>
    </template>
    <template #item="{ item }">
      <v-icon>{{ item.icon }}</v-icon>
      <span class="ml-2">{{ item.text }}</span>
    </template>
  </v-select>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT, PBEHAVIOR_TYPE_TYPES, SERVICE_STATES, WEATHER_ICONS } from '@/constants';

const { mapActions: mapPbehaviorTypesActions } = createNamespacedHelpers('pbehaviorTypes');

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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pbehaviorTypes: [],
      pending: false,
    };
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    pbehaviorTypeIcons() {
      return this.pbehaviorTypes.map(pbehaviorType => ({
        icon: pbehaviorType.icon_name,
        text: pbehaviorType.name,
      }));
    },

    serviceWeatherIcons() {
      return [
        {
          icon: WEATHER_ICONS[SERVICE_STATES.ok],
          text: this.$t('serviceWeather.iconTypes.ok'),
        },
        {
          icon: WEATHER_ICONS[SERVICE_STATES.minor],
          text: this.$t('serviceWeather.iconTypes.minorOrMajor'),
        },
        {
          icon: WEATHER_ICONS[SERVICE_STATES.critical],
          text: this.$t('serviceWeather.iconTypes.critical'),
        },
      ];
    },

    availableIcons() {
      return [...this.serviceWeatherIcons, ...this.pbehaviorTypeIcons];
    },
  },
  mounted() {
    this.fetchTypesList();
  },
  methods: {
    ...mapPbehaviorTypesActions({
      fetchPbehaviorTypesListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchTypesList() {
      this.pending = true;

      const { data: types } = await this.fetchPbehaviorTypesListWithoutStore({
        params: {
          types: [PBEHAVIOR_TYPE_TYPES.inactive, PBEHAVIOR_TYPE_TYPES.maintenance, PBEHAVIOR_TYPE_TYPES.pause],
          limit: MAX_LIMIT,
        },
      });

      this.pbehaviorTypes = types;
      this.pending = false;
    },
  },
};
</script>
