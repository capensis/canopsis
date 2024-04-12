<template>
  <v-combobox
    v-if="combobox"
    v-field="value"
    :label="$t('common.search')"
    :items="items"
    :menu-props="comboboxMenuProps"
    :return-object="false"
    :hide-details="false"
    append-icon=""
    item-text="search"
    item-value="search"
    hide-no-data
    single-line
    @input="submit"
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
    :hide-details="hideDetails"
    single-line
    @keydown.enter.prevent="submit(value)"
  />
</template>

<script>
import { computed } from 'vue';

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
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const comboboxMenuProps = computed(() => ({ contentClass: 'c-search-field__menu' }));

    const submit = () => emit('submit', props.value);
    const { removeItem, togglePinItem } = useSearchSavedItems(emit);

    return {
      comboboxMenuProps,

      submit,
      removeItem,
      togglePinItem,
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
