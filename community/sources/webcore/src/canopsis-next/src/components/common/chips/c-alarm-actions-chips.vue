<template>
  <v-layout
    :key="wrapperKey"
    :class="{ 'c-alarm-actions-chips--small': small }"
    class="c-alarm-actions-chips"
    wrap
    align-center
  >
    <c-chip
      v-for="item in inlineItems"
      :key="item[itemValue]"
      :class="itemClass"
      :color="item.color"
      :small="small"
      :closable="isClosableItem(item)"
      :outlined="outlined"
      :text-color="textColor"
      @click="selectItem(item)"
      @close="closeItem(item)"
    >
      <slot
        :item="item"
        name="item"
      >
        <span>{{ item[itemText] }}</span>
      </slot>
    </c-chip>
    <v-menu
      v-if="dropDownItems.length"
      key="more"
      max-height="400px"
      bottom
      left
      @input="$emit('activate')"
    >
      <template #activator="{ on }">
        <v-btn
          class="c-alarm-actions-chips__more-btn ma-0"
          color="grey"
          icon
          v-on="on"
        >
          <v-icon
            color="white"
            size="20"
          >
            more_horiz
          </v-icon>
        </v-btn>
      </template>
      <v-card>
        <v-card-text>
          <v-layout
            :class="{ 'c-alarm-actions-chips--small': small }"
            class="c-alarm-actions-chips__more"
            wrap
          >
            <c-chip
              v-for="item in dropDownItems"
              :key="item[itemValue]"
              :class="itemClass"
              :color="item.color"
              :closable="isClosableItem(item)"
              :text-color="textColor"
              :outlined="outlined"
              class="mx-0"
              @click="selectItem(item)"
              @close="closeItem(item)"
            >
              <slot
                :item="item"
                name="item"
              >
                <span>{{ item[itemText] }}</span>
              </slot>
            </c-chip>
          </v-layout>
        </v-card-text>
      </v-card>
    </v-menu>
  </v-layout>
</template>

<script>
import { computed, ref, watch } from 'vue';

import { uid } from '@/helpers/uid';

export default {
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    activeItem: {
      type: String,
      required: false,
    },
    activeItems: {
      type: Array,
      default: () => [],
    },
    inlineCount: {
      type: [Number, String],
      default: 2,
    },
    closable: {
      type: Boolean,
      default: false,
    },
    closableActive: {
      type: Boolean,
      default: false,
    },
    small: {
      type: Boolean,
      default: false,
    },
    outlined: {
      type: Boolean,
      default: false,
    },
    itemValue: {
      type: String,
      default: 'text',
    },
    itemText: {
      type: String,
      default: 'text',
    },
    itemClass: {
      type: [String, Object],
      required: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    textColor: {
      type: String,
      default: 'white',
    },
  },
  setup(props, { emit }) {
    const wrapperKey = ref(uid());

    const sortedItems = computed(() => [...props.items].sort((first, second) => {
      if (first[props.itemValue] === props.activeItem || props.activeItems.includes(first[props.itemValue])) {
        return -1;
      }

      if (second[props.itemValue] === props.activeItem || props.activeItems.includes(second[props.itemValue])) {
        return 0;
      }

      if (first[props.itemText] < second[props.itemText]) {
        return -1;
      }

      if (first[props.itemText] > second[props.itemText]) {
        return 1;
      }

      return 0;
    }));

    const inlineItems = computed(() => sortedItems.value.slice(0, props.inlineCount));

    const dropDownItems = computed(() => sortedItems.value.slice(props.inlineCount));

    watch(inlineItems, () => wrapperKey.value = uid());

    /**
     * Emits `select` with the full item or its value (per `returnObject`), unless the item is already in `activeItems`.

     * @param {Object} item - Chip item; value key is `itemValue`.
     */
    const selectItem = (item) => {
      if (props.activeItems.includes(item[props.itemValue])) {
        return;
      }

      emit('select', props.returnObject ? item : item[props.itemValue]);
    };

    /**
     * Emits `close` with the full item or its value (per `returnObject`).

     * @param {Object} item - Chip item; value key is `itemValue`.
     */
    const closeItem = item => emit('close', props.returnObject ? item : item[props.itemValue]);

    /**
     * Whether the chip should show a close control: when `closable`, or when `closableActive` and the item is active.

     * @param {Object} item - Chip item; compared using `itemValue` and `activeItem` / `activeItems`.
     * @returns {boolean}
     */
    const isClosableItem = (item) => {
      const isActive = props.activeItem === item[props.itemValue]
        || props.activeItems.includes(item[props.itemValue]);

      return props.closable || (props.closableActive && isActive);
    };

    return {
      wrapperKey,
      inlineItems,
      dropDownItems,
      selectItem,
      closeItem,
      isClosableItem,
    };
  },
};
</script>

<style lang="scss">
.c-alarm-actions-chips {
  &, &__more {
    column-gap: 8px;
    row-gap: 4px;
  }

  &--small {
    column-gap: 4px;
  }

  & &__more-btn {
    width: 24px;
    height: 24px;

    .theme--light & {
      background-color: var(--v-application-background-darken2);
    }

    .theme--dark & {
      background-color: var(--v-application-background-lighten4);
    }
  }
}
</style>
