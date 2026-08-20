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
import { computed, watch } from 'vue';

import { useModelField } from '@/hooks/form';
import { useStoreModuleHooks } from '@/hooks/store';

import materialIconsNames from '@/assets/material-icons/MaterialIcons.json';

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
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);
    const { useGetters } = useStoreModuleHooks('icon');
    const { registeredIconsById } = useGetters(['registeredIconsById']);

    const registeredIconsItems = computed(() => (
      Object.values(registeredIconsById.value)
        .map(({ title }) => ({ text: title, value: `$vuetify.icon.${title}` }))
    ));

    const materialIconsItems = computed(() => (
      materialIconsNames.map(name => ({ text: name, value: name }))
    ));

    const allIcons = computed(() => {
      if (!registeredIconsItems.value.length) {
        return materialIconsItems.value;
      }

      return [
        ...registeredIconsItems.value,
        { divider: true },
        ...materialIconsItems.value,
      ];
    });

    const rules = computed(() => ({
      required: props.required,
    }));

    watch(registeredIconsItems, (items) => {
      if (!items.some(({ value }) => value === props.value)) {
        updateModel('');
      }
    });

    return {
      registeredIconsItems,
      materialIconsItems,
      allIcons,
      rules,
    };
  },
};
</script>
