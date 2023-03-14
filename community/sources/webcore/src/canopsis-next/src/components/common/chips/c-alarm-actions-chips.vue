<template lang="pug">
  v-layout(row, wrap, align-center)
    c-alarm-action-chip(
      v-for="item in inlineItems",
      :key="item[itemValue]",
      :class="itemClass",
      :color="item.color",
      :small="small",
      :closable="closable",
      @click="selectItem(item)",
      @close="closeItem(item)"
    )
      slot(name="item", :item="item")
        span {{ item[itemText] }}
    v-menu(
      v-if="dropDownItems.length",
      bottom,
      left,
      @input="$emit('activate')"
    )
      template(#activator="{ on }")
        v-btn(v-on="on", color="grey", icon, small)
          v-icon(color="white", small) more_horiz
      v-card
        v-card-text
          c-alarm-action-chip(
            v-for="item in dropDownItems",
            :key="item[itemValue]",
            :class="itemClass",
            :color="item.color",
            :closable="closable",
            @click="selectItem(item)",
            @close="closeItem(item)"
          )
            slot(name="item", :item="item")
              span {{ item[itemText] }}
</template>

<script>
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
    inlineCount: {
      type: [Number, String],
      default: 2,
    },
    closable: {
      type: Boolean,
      default: false,
    },
    small: {
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
      type: String,
      required: false,
    },
  },
  computed: {
    sortedItems() {
      return [...this.items].sort((first, second) => {
        if (first[this.itemValue] === this.activeItem) {
          return -1;
        }

        if (second[this.itemValue] === this.activeItem) {
          return 0;
        }

        if (first[this.itemText] < second[this.itemText]) {
          return -1;
        }

        if (first[this.itemText] > second[this.itemText]) {
          return 1;
        }

        return 0;
      });
    },

    inlineItems() {
      return this.sortedItems.slice(0, this.inlineCount);
    },

    dropDownItems() {
      return this.sortedItems.slice(this.inlineCount);
    },
  },
  methods: {
    selectItem(item) {
      this.$emit('select', item[this.itemValue]);
    },

    closeItem(item) {
      this.$emit('close', item[this.itemValue]);
    },
  },
};
</script>
