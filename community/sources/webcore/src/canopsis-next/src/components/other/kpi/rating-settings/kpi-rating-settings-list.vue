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
    <template #toolbar="">
      <v-flex xs12>
        <v-expand-transition>
          <v-layout
            v-if="changedIds.length"
            class="ml-3 mt-3"
          >
            <v-btn
              color="primary"
              outlined
              @click="resetEnabledRatingSettings"
            >
              {{ $t('common.cancel') }}
            </v-btn>
            <v-btn
              class="ml-2"
              color="primary"
              @click="submit"
            >
              {{ $t('common.submit') }}
            </v-btn>
          </v-layout>
        </v-expand-transition>
      </v-flex>
    </template>
    <template #enabled="{ item }">
      <v-layout align-center>
        <v-simple-checkbox
          :value="isEnabledRatingSetting(item)"
          :disabled="!updatable"
          hide-details
          @input="enableRatingSetting(item, $event)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
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
