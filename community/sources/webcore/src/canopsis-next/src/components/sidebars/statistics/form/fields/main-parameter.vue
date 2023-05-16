<template lang="pug">
  widget-settings-item(:title="$t('settings.mainParameter.title')")
    c-progress-overlay(:pending="pending", transition)
    v-radio-group(
      v-field="value",
      v-validate="'required'",
      :error-messages="errors.collect('mainParameter')",
      name="mainParameter"
    )
      v-layout(
        v-for="{ id, label, tooltip } of mainRatingSettingsWithCustom",
        :key="id",
        row
      )
        v-radio.my-0(:value="id", :label="label", color="primary")
        c-help-icon(v-if="tooltip", :text="tooltip", icon="help", top)
    template(v-if="isCustomParameter")
      v-text-field(
        v-validate="'required'",
        :label="$t('settings.columnName')",
        :error-messages="errors.collect('columnName')",
        name="columnName"
      )
      field-filters-list(
        :filters="patterns",
        with-entity,
        addable,
        editable,
        @input="$emit('update:patterns', $event)"
      )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { KPI_ENTITY_RATING_SETTINGS_CUSTOM_PARAMETER, KPI_RATING_SETTINGS_TYPES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import FieldFiltersList from '@/components/sidebars/form/fields/filters-list.vue';

const { mapActions } = createNamespacedHelpers('ratingSettings');

export default {
  inject: ['$validator'],
  components: {
    WidgetSettingsItem,
    FieldFiltersList,
  },
  props: {
    value: {
      type: [String, Number],
      required: false,
    },
    patterns: {
      type: Array,
      default: () => [],
    },
    type: {
      type: Number,
      default: KPI_RATING_SETTINGS_TYPES.entity,
    },
  },
  data() {
    return {
      pending: false,
      mainRatingSettings: [],
    };
  },
  computed: {
    isEntityType() {
      return this.type === KPI_RATING_SETTINGS_TYPES.entity;
    },

    isCustomParameter() {
      return this.value === KPI_ENTITY_RATING_SETTINGS_CUSTOM_PARAMETER;
    },

    mainRatingSettingsWithCustom() {
      return this.isEntityType && !this.pending
        ? [
          ...this.mainRatingSettings,
          {
            id: KPI_ENTITY_RATING_SETTINGS_CUSTOM_PARAMETER,
            label: this.$t('settings.mainParameter.custom.label'), // TODO: translation
            tooltip: this.$t('settings.mainParameter.custom.tooltip'),
          },
        ]
        : this.mainRatingSettings;
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions({ fetchRatingSettingsWithoutStore: 'fetchListWithoutStore' }),

    async fetchList() {
      this.pending = true;

      const { data: mainRatingSettings = [] } = await this.fetchRatingSettingsWithoutStore({
        params: {
          type: this.type,
          main: true,
          paginate: false,
        },
      });

      this.mainRatingSettings = mainRatingSettings;
      this.pending = false;
    },
  },
};
</script>
