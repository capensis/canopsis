<template>
  <v-autocomplete
    v-field="value"
    v-validate="rules"
    :label="label"
    :hint="hint"
    :items="allIcons"
    :name="name"
    :error-messages="errors.collect(name)"
    :disabled="disabled"
    persistent-hint
  >
    <template #selection="{ item }">
      <v-icon>{{ item.value }}</v-icon>
      <span class="ml-2">{{ item.text }}</span>
    </template>
    <template #item="{ item }">
      <v-icon>{{ item.value }}</v-icon>
      <span class="ml-2">{{ item.text }}</span>
    </template>
    <template #no-data="">
      <slot name="no-data" />
    </template>
  </v-autocomplete>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import materialIconNameByCode from '@/assets/material-icons/MaterialIcons-Regular.json';

const { mapGetters } = createNamespacedHelpers('icon');

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    hint: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'icon_name',
    },
    required: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ...mapGetters(['registeredIconsById']),

    registeredIconsItems() {
      return Object.values(this.registeredIconsById)
        .map(({ title }) => ({ text: title, value: `$vuetify.icon.${title}` }));
    },

    materialIconsItems() {
      return Object.keys(materialIconNameByCode).map(name => ({ text: name, value: name }));
    },

    allIcons() {
      if (!this.registeredIconsItems.length) {
        return this.materialIconsItems;
      }

      return [
        ...this.registeredIconsItems,
        { divider: true },
        ...this.materialIconsItems,
      ];
    },

    rules() {
      return {
        required: this.required,
      };
    },
  },
};
</script>
