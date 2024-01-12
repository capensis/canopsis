<template>
  <widget-settings-item :title="$t('settings.mainParameter.title')">
    <c-progress-overlay
      :pending="pending"
      transition
    />
    <v-radio-group
      v-field="form.criteria"
      v-validate="'required'"
      :error-messages="errors.collect('mainParameter')"
      name="mainParameter"
    >
      <v-layout
        v-for="{ id, label, tooltip } of mainRatingSettingsWithCustom"
        :key="id"
      >
        <v-radio
          class="my-0"
          :value="id"
          :label="label"
          color="primary"
        />
        <c-help-icon
          v-if="tooltip"
          :text="tooltip"
          icon="help"
          top
        />
      </v-layout>
    </v-radio-group>
    <template v-if="isCustomCriteria">
      <c-name-field
        v-field="form.columnName"
        :label="$t('settings.columnName')"
        name="columnName"
        required
      />
      <field-filters-list
        v-field="form.patterns"
        name="patterns"
        with-entity
        addable
        editable
        required
      />
    </template>
  </widget-settings-item>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { KPI_RATING_SETTINGS_TYPES, KPI_ENTITY_RATING_SETTINGS_CUSTOM_CRITERIA } from '@/constants';

import { isCustomCriteria } from '@/helpers/entities/metric/form';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import FieldFiltersList from '@/components/sidebars/form/fields/filters-list.vue';

const { mapActions } = createNamespacedHelpers('ratingSettings');

export default {
  inject: ['$validator'],
  components: {
    WidgetSettingsItem,
    FieldFiltersList,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
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

    isCustomCriteria() {
      return isCustomCriteria(this.form.criteria);
    },

    mainRatingSettingsWithCustom() {
      return this.isEntityType && !this.pending
        ? [
          ...this.mainRatingSettings,
          {
            id: KPI_ENTITY_RATING_SETTINGS_CUSTOM_CRITERIA,
            label: this.$t('settings.mainParameter.custom.label'),
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
