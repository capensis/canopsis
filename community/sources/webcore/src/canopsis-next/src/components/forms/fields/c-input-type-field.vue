<template>
  <v-select
    class="mt-4 c-input-type-field"
    v-field="value"
    v-validate="rules"
    :items="preparedTypes"
    :label="label"
    :disabled="disabled"
    :flat="flat"
    :error="hasError"
    :name="name"
    hide-details
    dense
  >
    <template #selection="{ item }">
      <v-icon
        class="c-input-type-field__icon"
        small
      >
        {{ getInputTypeIcon(item.value) }}
      </v-icon>
    </template>
    <template #item="{ item }">
      <v-list-item-avatar class="c-input-type-field__avatar">
        <v-icon
          class="c-input-type-field__icon"
          small
        >
          {{ getInputTypeIcon(item.value) }}
        </v-icon>
      </v-list-item-avatar>
      <v-list-item-content>
        <v-list-item-title>{{ item.text }}</v-list-item-title>
      </v-list-item-content>
    </template>
  </v-select>
</template>

<script>
import { PATTERN_FIELD_TYPES } from '@/constants';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: null,
    },
    name: {
      type: String,
      default: 'type',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    flat: {
      type: Boolean,
      default: false,
    },
    types: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    hasError() {
      return this.errors.has(this.name);
    },

    preparedTypes() {
      return this.types.map(
        type => (type.text ? type : ({ ...type, text: this.$t(`common.mixedField.types.${type.value}`) })),
      );
    },
  },
  methods: {
    getInputTypeIcon(type) {
      return {
        [PATTERN_FIELD_TYPES.string]: 'title',
        [PATTERN_FIELD_TYPES.number]: 'exposure_plus_1',
        [PATTERN_FIELD_TYPES.boolean]: 'toggle_on',
        [PATTERN_FIELD_TYPES.null]: 'space_bar',
        [PATTERN_FIELD_TYPES.stringArray]: 'view_array',
      }[type];
    },
  },
};
</script>

<style lang="scss">
.c-input-type-field {
  &__icon {
    color: inherit !important;
    opacity: .6;
  }

  &__avatar {
    min-width: 30px;

    & .v-avatar {
      width: 20px !important;
      height: 20px !important;
    }
  }
}
</style>
