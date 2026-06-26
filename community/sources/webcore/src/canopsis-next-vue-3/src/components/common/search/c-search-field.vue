<template>
  <v-combobox
    v-if="combobox"
    ref="comboboxElement"
    :model-value="modelValue"
    :label="label || $t('common.search')"
    :items="items"
    :menu-props="comboboxMenuProps"
    :return-object="false"
    item-title="search"
    item-value="search"
    hide-no-data
    single-line
    hide-details
    @keydown.enter="keydownEnterCombobox"
    @update:model-value="submit"
    @update:search="updateModel"
  >
    <template #item="{ item }">
      <v-list-item :title="item.raw.search" />
    </template>
  </v-combobox>
  <v-text-field
    v-else
    :model-value="modelValue"
    :label="label || $t('common.search')"
    hide-details
    single-line
    @update:model-value="updateModel"
    @keydown.enter.prevent="submit(modelValue)"
  />
</template>

<script>
import { computed, ref } from 'vue';

export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    combobox: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      required: false,
    },
  },

  emits: ['update:modelValue', 'submit', 'remove:item', 'toggle-pin:item'],

  setup(props, { emit }) {
    const comboboxElement = ref();

    const comboboxMenuProps = computed(() => ({ contentClass: 'c-search-field__menu' }));

    const updateModel = value => emit('update:modelValue', value ?? '');
    const submit = () => emit('submit', props.modelValue);
    const keydownEnterCombobox = () => {
      submit();
      comboboxElement.value?.blur();
    };

    return {
      comboboxElement,

      comboboxMenuProps,

      updateModel,
      submit,
      keydownEnterCombobox,
    };
  },
};
</script>

<style lang="scss">
.c-search-field {
  padding: 0 24px;

  .v-btn--icon {
    margin: 0 6px !important;
  }

  & > :last-child .v-btn--icon {
    margin-right: -6px !important;
  }

  &__menu {
    .v-list {
      padding: 0;

      .v-list-item {
        min-height: 32px;
      }
    }
  }
}
</style>
