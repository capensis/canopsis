<template>
  <v-select
    v-field="value"
    v-validate="rules"
    :items="themes"
    :label="label || $tc('common.theme')"
    :loading="pending"
    :disabled="disabled"
    :name="name"
    :error-messages="errors.collect(name)"
    :hide-details="hideDetails"
    item-text="name"
    item-value="_id"
  />
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

const { mapActions: mapThemeActions } = createNamespacedHelpers('theme');

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: [Object, String],
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'map',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pending: false,
      themes: [],
    };
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapThemeActions({
      fetchThemesListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchList() {
      this.pending = true;

      try {
        const { data: themes } = await this.fetchThemesListWithoutStore({ params: { limit: MAX_LIMIT } });

        this.themes = themes;
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
