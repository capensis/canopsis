<template lang="pug">
  c-advanced-data-table(
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
          v-layout.ml-3(v-if="changedIds.length")
            v-btn(
              outline,
              color="primary",
              @click="resetEnabledRatingSettings"
            ) {{ $t('common.cancel') }}
            v-btn(
              color="primary",
              @click="submit"
            ) {{ $t('common.submit') }}

    template(#enabled="{ item }")
      v-layout(row, align-center)
        v-checkbox-functional(
          :input-value="isEnabledRatingSetting(item)",
          :disabled="!updatable",
          hide-details,
          @change="enableRatingSetting(item, $event)"
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
      changed: [],
    };
  },
  computed: {
    changedIds() {
      return this.changed.map(({ id }) => id);
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
    submit() {
      this.$emit('change-selected', this.changed.map(ratingSetting => ({
        ...ratingSetting,
        enabled: !ratingSetting.enabled,
      })));
    },

    isRatingSettingChanged(ratingSetting) {
      return this.changedIds.includes(ratingSetting.id);
    },

    resetEnabledRatingSettings() {
      this.changed = [];
    },

    enableRatingSetting(ratingSetting) {
      if (!this.isRatingSettingChanged(ratingSetting)) {
        this.changed.push(ratingSetting);
      } else {
        this.changed = this.changed.filter(item => item.id !== ratingSetting.id);
      }
    },

    isEnabledRatingSetting(ratingSetting) {
      return this.isRatingSettingChanged(ratingSetting)
        ? !ratingSetting.enabled
        : ratingSetting.enabled;
    },
  },
};
</script>
