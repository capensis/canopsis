<template>
  <v-menu
    :min-width="listMinWidth || undefined"
    :nudge-bottom="1"
    :transition="false"
    content-class="c-alarm-advanced-search__history-menu"
    bottom
    offset-y
    @input="handleMenuToggle"
  >
    <template #activator="{ on }">
      <c-action-btn
        :tooltip="$t('common.search')"
        :disabled="!searches.length"
        icon="history"
        v-on="on"
      />
    </template>
    <div
      ref="scrollContainerElement"
      class="c-alarm-advanced-search__history-list"
    >
      <v-list
        class="pa-0"
        dense
      >
        <div
          v-for="search in items"
          :key="search._id"
          :ref="setItemRef(search._id)"
          class="c-alarm-advanced-search__history-item"
        >
          <v-list-item
            v-if="isItemVisible(search._id)"
            @click="select(search)"
          >
            <v-list-item-content class="pa-0">
              <alarm-advanced-search-rules
                :rules="search.rules"
                :attributes="attributes"
                disabled
              />
            </v-list-item-content>
            <v-list-item-action>
              <advanced-search-history-item-btns
                :id="search._id"
                :pinned="search.pinned"
                @remove="remove"
                @toggle-pin="togglePin"
              />
            </v-list-item-action>
          </v-list-item>
          <div
            v-else
            :style="getItemPlaceholderStyle(search._id)"
            class="c-alarm-advanced-search__history-item__placeholder"
          />
        </div>
        <div
          ref="appendItemElement"
          class="c-alarm-advanced-search__history-list__append-item"
        />
      </v-list>
    </div>
  </v-menu>
</template>
<script>
import { computed, nextTick } from 'vue';

import { useAlarmAdvancedSearchHistory } from '../hooks/alarm-advanced-search-history';

import AdvancedSearchHistoryItemBtns from './advanced-search-history-item-btns.vue';
import AlarmAdvancedSearchRules from './alarm-advanced-search-rules.vue';

export default {
  components: { AdvancedSearchHistoryItemBtns, AlarmAdvancedSearchRules },
  props: {
    searches: {
      type: Array,
      default: () => [],
    },
    attributes: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const {
      items,
      scrollContainerElement,
      appendItemElement,
      setItemRef,
      isItemVisible,
      getItemPlaceholderStyle,
      listMinWidth,
      updateListWidth,
    } = useAlarmAdvancedSearchHistory(computed(() => props.searches));

    /**
     * Recalculates list width after the menu content is mounted.
     *
     * @param {boolean} isOpen
     */
    const handleMenuToggle = (isOpen) => {
      if (!isOpen) {
        return;
      }

      nextTick(() => {
        nextTick(() => updateListWidth());
      });
    };

    /**
     * Emits a 'select' event with the specified search configuration.
     *
     * @param {Object} search - The search configuration object to be selected.
     */
    const select = search => emit('select', search);

    /**
     * Emits a 'remove' event with the specified identifier.
     *
     * @param {string} id - The unique identifier of the item to be removed.
     */
    const remove = id => emit('remove', id);

    /**
     * Emits a 'toggle-pin' event with the specified identifier.
     *
     * @param {string} id - The unique identifier of the item whose pin status is to be toggled.
     */
    const togglePin = id => emit('toggle-pin', id);

    return {
      items,
      scrollContainerElement,
      appendItemElement,
      setItemRef,
      isItemVisible,
      getItemPlaceholderStyle,
      listMinWidth,
      handleMenuToggle,

      select,
      remove,
      togglePin,
    };
  },
};
</script>

<style lang="scss">
.c-alarm-advanced-search__history-menu {
  overflow-y: hidden !important;
  width: max-content;
  min-width: max-content;

  @media (min-width: 1264px) {
    max-width: 95vw;
  }
}

.c-alarm-advanced-search__history-list {
  max-height: 95vh;
  overflow-y: auto;
  width: max-content;
  min-width: 100%;
}

.c-alarm-advanced-search__history-item {
  width: max-content;
  min-width: 100%;
}

.c-alarm-advanced-search__history-item__placeholder {
  width: 100%;
}

.c-alarm-advanced-search__history-list__append-item {
  height: 1px;
}
</style>
