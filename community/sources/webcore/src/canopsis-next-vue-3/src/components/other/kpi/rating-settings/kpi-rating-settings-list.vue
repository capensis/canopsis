<template>
  <c-advanced-data-table
    :options="options"
    :items="ratingSettings"
    :loading="pending"
    :headers="headers"
    :total-items="totalItems"
    item-key="id"
    search
    advanced-pagination
    hide-actions
    @update:options="$emit('update:options', $event)"
  >
    <template #toolbar>
      <v-col cols="12">
        <v-expand-transition>
          <c-layout v-if="changedIds.length" class="mt-3">
            <v-btn color="primary" variant="outlined" @click="resetEnabledRatingSettings">
              {{ $t('common.cancel') }}
            </v-btn>
            <v-btn class="ml-2" color="primary" @click="submit">
              {{ $t('common.submit') }}
            </v-btn>
          </c-layout>
        </v-expand-transition>
      </v-col>
    </template>

    <template #enabled="{ item }">
      <c-layout align-center>
        <v-checkbox-btn
          :model-value="isEnabledRatingSetting(item)"
          :disabled="!updatable"
          hide-details
          @update:model-value="enableRatingSetting(item)"
        />
      </c-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    ratingSettings: {
      type: Array,
      required: true,
    },
    options: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      default: 0,
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

  emits: ['update:options', 'change-selected'],

  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      { value: 'enabled', width: 100, sortable: false },
      { text: t('common.criteria'), value: 'label', sortable: false },
    ]);

    return {
      headers,
    };
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
  },

  watch: {
    ratingSettings: 'resetEnabledRatingSettings',
  },

  methods: {
    submit() {
      this.$emit(
        'change-selected',
        this.changed.map(ratingSetting => ({
          ...ratingSetting,
          enabled: !ratingSetting.enabled,
        })),
      );
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
