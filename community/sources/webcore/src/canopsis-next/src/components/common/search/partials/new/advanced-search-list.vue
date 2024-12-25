<template>
  <v-list>
    <v-list-item-group v-field="value">
      <template v-for="item in items">
        <v-subheader v-if="item.header" :key="item.header">
          {{ item.header }}
        </v-subheader>
        <v-list-item
          v-else
          :key="item[itemValue]"
          :item-value="item[itemValue]"
          @click="selectItem"
        >
          <v-list-item-title>
            {{ item[itemText] }}
          </v-list-item-title>
        </v-list-item>
      </template>
      <v-list-item v-if="!items.length">
        <v-list-item-title class="grey--text">
          {{ $t('advancedSearch.noDataList') }}
        </v-list-item-title>
      </v-list-item>
    </v-list-item-group>
  </v-list>
</template>

<script>
export default {
  props: {
    value: {
      type: [String, Array],
      default: '',
    },
    items: {
      type: Array,
      default: () => [],
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    itemText: {
      type: String,
      default: 'text',
    },
  },
  setup() {
    const selectItem = () => {};

    return {
      selectItem,
    };
  },
};
</script>
