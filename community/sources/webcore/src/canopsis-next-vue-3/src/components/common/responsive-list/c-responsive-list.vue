<template>
  <c-layout class="c-responsive-list" align="center" gap="0">
    <div ref="listContainer" :class="layoutClass" class="c-responsive-list__container">
      <div
        v-for="item in visibleItems"
        :key="item.key"
        :ref="(element) => setItemRef(item.key, element)"
      >
        <slot :item="item">
          {{ getValue(item) }}
        </slot>
      </div>
    </div>

    <v-menu v-if="shownMenu" location="bottom end" origin="top end">
      <template #activator="{ props: activatorProps }">
        <v-btn
          class="c-responsive-list__menu-button text-white"
          density="compact"
          variant="text"
          icon
          v-bind="activatorProps"
        >
          <v-icon size="small">more_vert</v-icon>
        </v-btn>
      </template>

      <v-list class="py-0">
        <div v-for="item in hiddenItems" :key="item.key" class="bg-white">
          <slot :item="item">
            {{ getValue(item) }}
          </slot>
        </div>
      </v-list>
    </v-menu>
  </c-layout>
</template>

<script>
export default {
  name: 'CResponsiveList',

  props: {
    items: {
      type: Array,
      default: () => [],
    },
    itemKey: {
      type: String,
      default: 'id',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
  },

  data() {
    return {
      visible: false,
      containerWidth: 0,
      itemsWidthByKey: {},
      itemRefsByKey: {},
      resizeObserver: null,
    };
  },

  computed: {
    itemsWithKey() {
      return this.items.map((item) => ({
        ...item,
        key: this.getKey(item),
      }));
    },

    shownMenu() {
      return this.itemsWithKey.length !== this.visibleItems.length;
    },

    visibleItems() {
      return this.visible ? this.getVisibleItems(this.itemsWithKey) : this.itemsWithKey;
    },

    hiddenItems() {
      return this.itemsWithKey.slice(this.visibleItems.length);
    },

    layoutClass() {
      return { 'c-responsive-list__container--hidden': !this.visible };
    },
  },

  watch: {
    items: {
      handler() {
        this.setItemsWidth();
      },
      deep: true,
    },

    $route() {
      this.setItemsWidth();
    },
  },

  mounted() {
    this.resizeObserver = new ResizeObserver(this.setItemsWidth);
    this.resizeObserver.observe(this.$refs.listContainer);

    this.setItemsWidth();
  },

  beforeUnmount() {
    this.resizeObserver?.disconnect();
  },

  methods: {
    getKey(item) {
      return item[this.itemKey];
    },

    getValue(item) {
      return item[this.itemValue];
    },

    setItemRef(key, element) {
      if (element) {
        this.itemRefsByKey[key] = element;
      } else {
        delete this.itemRefsByKey[key];
      }
    },

    getVisibleItems(items = []) {
      let availableWidth = this.containerWidth;
      const results = [];

      for (const item of items) {
        const key = this.getKey(item);

        availableWidth -= this.itemsWidthByKey[key] || 0;

        if (availableWidth <= 0) {
          break;
        }

        results.push(item);
      }

      return results;
    },

    getItemsWidthByKey(items = []) {
      return items.reduce((acc, item) => {
        const key = this.getKey(item);
        const element = this.itemRefsByKey[key];

        if (element) {
          const { width } = element.getBoundingClientRect();

          acc[key] = width;
        }

        return acc;
      }, {});
    },

    setItemsWidth() {
      this.setContainerWidth();
      this.visible = false;

      this.$nextTick(() => {
        this.itemsWidthByKey = this.getItemsWidthByKey(this.itemsWithKey);
        this.visible = true;
      });
    },

    setContainerWidth() {
      const { width: availableWidth } = this.$refs.listContainer.getBoundingClientRect();

      this.containerWidth = availableWidth;
    },
  },
};
</script>

<style lang="scss" scoped>
.c-responsive-list {
  min-width: 0;
  overflow: hidden;
  width: 100%;

  &__container {
    align-items: center;
    display: flex;
    min-width: 0;
    width: 100%;

    &--hidden {
      visibility: hidden;
    }
  }

  &__menu-button {
    height: 32px;
    width: 32px;
  }
}
</style>
