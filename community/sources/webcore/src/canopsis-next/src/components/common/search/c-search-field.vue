<template>
  <v-combobox
    v-if="combobox"
    v-field="value"
    ref="comboboxElement"
    :label="$t('common.search')"
    :items="items"
    :menu-props="comboboxMenuProps"
    :return-object="false"
    append-icon=""
    item-text="search"
    item-value="search"
    hide-no-data
    single-line
    hide-details
    @keydown.enter="keydownEnterCombobox"
    @input="submit"
    @update:search-input="updateModel"
  >
    <template #item="{ item }">
      <v-list-item-content>
        <v-list-item-title class="pr-2">
          {{ item.search }}
        </v-list-item-title>
      </v-list-item-content>
      <v-list-item-action>
        <v-layout>
          <v-btn
            small
            icon
            @click.stop="removeItem(item.search)"
          >
            <v-icon
              color="grey"
              small
            >
              delete
            </v-icon>
          </v-btn>
          <v-btn
            :class="{ 'c-search-field__item__pinned': item.pinned }"
            small
            icon
            @click.stop="togglePinItem(item.search)"
          >
            <v-icon
              :color="item.pinned ? 'inherit' : 'grey'"
              small
            >
              $vuetify.icons.push_pin
            </v-icon>
          </v-btn>
        </v-layout>
      </v-list-item-action>
    </template>
  </v-combobox>
  <v-text-field
    v-else
    v-field="value"
    :label="$t('common.search')"
    hide-details
    single-line
    @keydown.enter.prevent="submit(value)"
  />
</template>

<script>
import { computed, ref } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

import { useSearchSavedItems } from './hooks/search';

export default {
  props: {
    value: {
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
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);

    const comboboxElement = ref();

    const comboboxMenuProps = computed(() => ({ contentClass: 'c-search-field__menu' }));

    const submit = () => emit('submit', props.value);
    const { removeItem, togglePinItem } = useSearchSavedItems(emit);

    const keydownEnterCombobox = () => {
      submit();
      comboboxElement.value?.blur();
    };

    return {
      comboboxElement,

      comboboxMenuProps,

      updateModel,
      submit,
      removeItem,
      togglePinItem,
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
        height: 32px;

        .v-btn:not(.c-search-field__item__pinned) {
          opacity: 0;
        }

        &:hover .v-btn {
          opacity: 1;
        }
      }
    }
  }
}
</style>
