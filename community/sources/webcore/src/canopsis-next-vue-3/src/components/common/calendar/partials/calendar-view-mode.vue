<template>
  <v-menu>
    <template #activator="{ props: menuProps }">
      <v-btn v-bind="menuProps" class="mx-2 my-1" variant="text">
        {{ currentType }}
        <v-icon>arrow_drop_down</v-icon>
      </v-btn>
    </template>

    <v-list>
      <v-list-item
        v-for="type in types"
        :key="type.value"
        @click="$emit('update:modelValue', type.value)"
      >
        <v-list-item-title>{{ type.text }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      required: true,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    return { t };
  },

  computed: {
    currentType() {
      return this.t(`calendar.${this.modelValue}`);
    },

    types() {
      return ['day', 'week', 'month'].map((type) => ({
        text: this.t(`calendar.${type}`),
        value: type,
      }));
    },
  },
};
</script>
