<template>
  <c-layout wrap>
    <div class="c-enabled-duration-field__enabled">
      <component
        :is="component"
        :model-value="duration.enabled"
        :label="label"
        :name="enabledFieldName"
        color="primary"
        class="pb-2"
        @update:model-value="updateEnabled"
      >
        <template #label>
          <c-layout class="full-width" align="center">
            {{ label }}
            <c-help-icon
              v-if="helpText"
              :text="helpText"
              icon-class="ml-2"
              color="info"
              max-width="300"
              top
            />
            <v-fade-transition>
              <span v-if="suffix && !hideDuration" class="pr-2 pl-1 text-lowercase">
                {{ suffix }}
              </span>
            </v-fade-transition>
          </c-layout>
        </template>
      </component>
    </div>

    <v-fade-transition>
      <div v-if="!hideDuration" class="c-enabled-duration-field__duration">
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
      </div>
    </v-fade-transition>
  </c-layout>
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

<style lang="scss" scoped>
.c-enabled-duration-field__enabled {
  flex: 1 1 58%;
  min-width: 0;
}

.c-enabled-duration-field__duration {
  flex: 1 1 42%;
  min-width: 0;
}
</style>
