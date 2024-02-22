<template>
  <v-layout class="c-responsive-list">
    <div
      v-resize="setContainerWidth"
      ref="listContainer"
      :class="layoutClass"
      class="c-responsive-list__container"
    >
      <div
        v-for="item in visibleItems"
        :key="item.key"
        :ref="`item:${item.key}`"
      >
        <slot :item="item">
          {{ getValue(item) }}
        </slot>
      </div>
    </div>
    <v-menu
      v-if="shownMenu"
      bottom
      offset-y
    >
      <template #activator="{ on }">
        <v-btn
          class="ma-1"
          icon
          small
          v-on="on"
        >
          <v-icon
            color="white"
            small
          >
            more_vert
          </v-icon>
        </v-btn>
      </template>
      <div
        v-for="item in hiddenItems"
        :key="item.key"
        class="white"
      >
        <slot :item="item">
          {{ getValue(item) }}
        </slot>
      </div>
    </v-menu>
  </v-layout>
</template>

<script>
export default {
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
    };
  },
  computed: {
    itemsWithKey() {
      return this.items.map(item => ({
        ...item,
        key: this.getKey(item),
      }));
    },

    shownMenu() {
      return this.itemsWithKey.length !== this.visibleItems.length;
    },

    visibleItems() {
      return this.visible
        ? this.getVisibleItems(this.itemsWithKey)
        : this.itemsWithKey;
    },

    hiddenItems() {
      return this.itemsWithKey.slice(this.visibleItems.length);
    },

    layoutClass() {
      return { 'c-responsive-list__container--hidden': !this.visible };
    },
  },
  watch: {
    items() {
      this.setItemsWidth();
    },

    $route() {
      this.setItemsWidth();
    },
  },
  mounted() {
    this.setItemsWidth();
  },
  methods: {
    getKey(item) {
      return item[this.itemKey];
    },

    getValue(item) {
      return item[this.itemValue];
    },

    getVisibleItems(items = []) {
      let availableWidth = this.containerWidth;
      const results = [];

      // eslint-disable-next-line no-restricted-syntax
      for (const item of items) {
        const key = this.getKey(item);

        availableWidth -= this.itemsWidthByKey[key];

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
        const [element] = this.$refs[`item:${key}`] || [];

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
        this.itemsWidthByKey = this.getItemsWidthByKey(this.items);

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
  overflow: hidden;

  &__container {
    width: 100%;
    display: flex;
    align-items: center;

    &--hidden {
      visibility: hidden;
    }
  }
}
</style>
