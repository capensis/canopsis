<template lang="pug">
  c-advanced-data-table.white(
    :pagination="pagination",
    :items="ratingSettings",
    :loading="pending",
    :headers="headers",
    :total-items="totalItems",
    item-key="id",
    search,
    advanced-pagination,
    hide-actions,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#toolbar="")
      v-flex(xs12)
        v-expand-transition
          v-layout.ml-3(v-if="isEnabledRatingSettingsChanged")
            v-btn(
              outline,
              color="primary",
              @click="resetEnabledRatingSettings"
            ) {{ $t('common.cancel') }}
            v-btn(
              color="primary",
              @click="$emit('enable-selected', enabled)"
            ) {{ $t('common.submit') }}

    template(#enabled="props")
      v-layout(row, align-center)
        v-checkbox-functional(
          :input-value="isEnabledRatingSetting(props.item)",
          :disabled="!updatable",
          hide-details,
          @change="enableRatingSetting(props.item, $event)"
        )
</template>

<script>
export default {
  props: {
    ratingSettings: {
      type: Array,
      required: true,
    },
    pagination: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      enabled: this.getEnabledRatingSettingsIds(),
    };
  },
  computed: {
    enabledIds() {
      return this.enabled.map(({ id }) => id);
    },

    isEnabledRatingSettingsChanged() {
      const submittedRatingSettingsIds = this.getEnabledRatingSettingsIds()
        .map(({ id }) => id);

      return this.enabledIds.length !== submittedRatingSettingsIds.length
        || !this.enabledIds.every(id => submittedRatingSettingsIds.includes(id));
    },

    headers() {
      return [
        { value: 'enabled', width: 100, sortable: false },
        { text: this.$t('common.criteria'), value: 'label', sortable: false },
      ];
    },
  },
  watch: {
    ratingSettings: 'resetEnabledRatingSettings',
  },
  methods: {
    resetEnabledRatingSettings() {
      this.enabled = this.getEnabledRatingSettingsIds();
    },

    enableRatingSetting(ratingSetting, checked) {
      if (checked) {
        this.enabled.push(ratingSetting);
      } else {
        this.enabled = this.enabled.filter(item => item.id !== ratingSetting.id);
      }
    },

    getEnabledRatingSettingsIds() {
      return this.ratingSettings.filter(item => item.enabled);
    },

    isEnabledRatingSetting(ratingSetting) {
      return this.enabledIds.includes(ratingSetting.id);
    },
  },
};
</script>
