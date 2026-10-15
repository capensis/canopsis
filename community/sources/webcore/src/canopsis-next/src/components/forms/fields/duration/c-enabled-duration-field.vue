<template>
  <v-row density="comfortable">
    <v-col cols="7">
      <component
        :is="component"
        :model-value="duration.enabled"
        :label="label"
        :name="enabledFieldName"
        color="primary"
        class="pb-2"
        hide-details
        @update:model-value="updateEnabled"
      >
        <template #label>
          <v-row class="w-100" align="center" no-gutters>
            <v-col>
              {{ label }}
              <c-help-icon
                v-if="helpText"
                :text="helpText"
                icon-class="ml-2"
                color="info"
                max-width="300"
                top
              />
            </v-col>
            <v-col cols="auto">
              <v-fade-transition>
                <span v-if="suffix && !hideDuration" class="pr-2 pl-1 text-lowercase">
                  {{ suffix }}
                </span>
              </v-fade-transition>
            </v-col>
          </v-row>
        </template>
      </component>
    </v-col>

    <v-fade-transition>
      <v-col v-if="!hideDuration" cols="5">
        <c-duration-field
          :duration="duration"
          :units-label="$t('common.unit')"
          :disabled="!duration.enabled"
          :required="duration.enabled"
          :units="timeUnits"
          :name="name"
          :min="min"
          @update:duration="updateModel"
        />
      </v-col>
    </v-fade-transition>
  </v-row>
</template>

<script>
import { AVAILABLE_TIME_UNITS } from '@/constants';

import { convertUnit } from '@/helpers/date/duration';

export default {
  props: {
    duration: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      required: true,
    },
    helpText: {
      type: String,
      required: false,
    },
    suffix: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'duration',
    },
    units: {
      type: Array,
      required: false,
    },
    after: {
      type: Object,
      required: false,
    },
    switcher: {
      type: Boolean,
      default: false,
    },
    hideValueOnFalse: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:duration'],

  computed: {
    enabledFieldName() {
      return `${this.name}.enabled`;
    },

    hideDuration() {
      return this.hideValueOnFalse && !this.duration.enabled;
    },

    component() {
      return this.switcher ? 'c-enabled-field' : 'c-checkbox';
    },

    timeUnits() {
      const units = this.units || [
        AVAILABLE_TIME_UNITS.day,
        AVAILABLE_TIME_UNITS.week,
        AVAILABLE_TIME_UNITS.month,
        AVAILABLE_TIME_UNITS.year,
      ];

      return units.map(({ value, text }) => ({
        value,
        text: this.$t(text, this.duration?.value || 0),
      }));
    },

    min() {
      if (!this.duration?.enabled || !this.after) {
        return 1;
      }

      return Math.floor(convertUnit(this.after.value, this.after.unit, this.duration.unit)) + 1;
    },
  },

  methods: {
    updateModel(duration) {
      this.$emit('update:duration', duration);
    },

    updateEnabled(enabled) {
      this.updateModel({
        ...this.duration,
        enabled,
      });
    },
  },
};
</script>
