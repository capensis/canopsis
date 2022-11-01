<template lang="pug">
  v-select(
    v-field="value",
    :label="label",
    :items="availableIcons",
    :loading="pending",
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
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT, PBEHAVIOR_TYPE_TYPES, WEATHER_ICONS } from '@/constants';

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

    availableIcons() {
      return Object.entries(WEATHER_ICONS)
        .map(([value, icon]) => ({
          icon,
          text: this.$te(`common.stateTypes.${value}`)
            ? this.$t(`common.stateTypes.${value}`)
            : this.$t(`serviceWeather.iconTypes.${value}`),
        }))
        .concat(this.pbehaviorTypeIcons);
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
