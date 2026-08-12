<template>
  <c-radio-group
    v-if="column"
    :model-value="value"
    :name="name"
    class="density__radio-group"
    @update:model-value="updateModel"
  >
    <c-layout v-for="type in types" :key="type.value" class="mb-3">
      <div class="w-50">
        <v-radio :value="type.value" :label="type.text" color="primary" />
      </div>
      <div class="w-50">
        <v-icon class="density__icon">
          {{ type.icon }}
        </v-icon>
      </div>
    </c-layout>
  </c-radio-group>
  <v-btn-toggle
    v-else
    :model-value="value"
    :name="name"
    :density="density"
    class="density__btn-toggle"
    mandatory
    @update:model-value="updateModel"
  >
    <v-tooltip v-for="type in types" :key="type.value" location="top">
      <template #activator="{ props: tooltipProps }">
        <v-btn
          v-bind="tooltipProps"
          :value="type.value"
          :density="density"
          :size="buttonSize"
          class="ma-0"
          variant="text"
        >
          <v-icon :size="iconSize">
            {{ type.icon }}
          </v-icon>
        </v-btn>
      </template>
      <span>{{ type.text }}</span>
    </v-tooltip>
  </v-btn-toggle>
</template>

<script>
import { DENSE_TYPES } from '@/constants';

export default {
  props: {
    value: {
      type: Number,
      default: DENSE_TYPES.large,
    },
    name: {
      type: String,
      default: 'dense',
    },
    column: {
      type: Boolean,
      default: false,
    },
    buttonSize: {
      type: String,
      default: 'small',
    },
    iconSize: {
      type: [String, Number],
      default: 'small',
    },
    density: {
      type: String,
      default: 'default',
    },
  },

  emits: ['change'],

  computed: {
    types() {
      return [
        {
          value: DENSE_TYPES.large,
          icon: 'density_large',
          text: this.$t('settings.density.comfort'),
        },
        {
          value: DENSE_TYPES.medium,
          icon: 'density_medium',
          text: this.$t('settings.density.compact'),
        },
        {
          value: DENSE_TYPES.small,
          icon: 'density_small',
          text: this.$t('settings.density.ultraCompact'),
        },
      ];
    },
  },

  methods: {
    updateModel(value) {
      this.$emit('change', value);
    },
  },
};
</script>

<style lang="scss">
.density {
  &__btn-toggle {
    box-shadow: none;
  }

  &__radio-group .v-input__control {
    width: 100%;
  }

  &__icon {
    padding: 2px;
    border-radius: 5px;
    border: 1px solid #707070;
  }
}
</style>
