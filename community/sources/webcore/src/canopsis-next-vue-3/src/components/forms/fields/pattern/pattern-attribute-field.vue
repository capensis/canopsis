<template>
  <c-autocomplete
    :model-value="itemValue"
    :items="items"
    :disabled="disabled"
    :label="label || $t('common.searchBy')"
    :name="name"
    :required="required"
    :item-props="getItemProps"
    item-title="text"
    item-value="value"
    @update:model-value="updateValue"
  >
    <template #item="{ props: itemProps, item }">
      <v-list-item v-bind="itemProps">
        <template v-if="item.raw?.alias" #append>
          <c-simple-tooltip :content="`infos.${item.raw.originalValue}.value`" top>
            <template #activator="{ props: activatorProps }">
              <v-icon
class="ml-1"
color="primary"
size="small"
v-bind="activatorProps">
                alternate_email
              </v-icon>
            </template>
          </c-simple-tooltip>
        </template>
      </v-list-item>
    </template>
    <template #selection="{ item }">
      <v-icon v-if="selectionItem(item)?.alias" class="mr-1" size="small">
alternate_email
</v-icon>
      <span class="text-truncate">{{ selectionTitle(item) }}</span>
    </template>
  </c-autocomplete>
</template>

<script>
export default {
  emits: ['update:value'],
  props: {
    value: {
      type: [String, Object],
      required: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'attribute',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    itemsByValue() {
      return Object.fromEntries(this.items.map(item => [item.value, item]));
    },

    itemValue() {
      const value = this.value?.value ?? this.value;

      if (value == null || value === '') {
        return null;
      }

      return this.itemsByValue[value]?.value ?? value;
    },

    selectedItem() {
      if (this.itemValue == null || this.itemValue === '') {
        return null;
      }

      return this.itemsByValue[this.itemValue] ?? null;
    },
  },
  methods: {
    getItemProps(item) {
      return {
        disabled: item.raw?.options?.disabled,
      };
    },

    selectionItem(item) {
      return item?.raw ?? this.selectedItem;
    },

    selectionTitle(item) {
      return item?.title ?? this.selectedItem?.text ?? '';
    },

    updateValue(value) {
      const selectedItem = this.itemsByValue[value?.value ?? value] ?? value;

      if (!this.returnObject) {
        this.$emit('update:value', selectedItem?.value ?? selectedItem);

        return;
      }

      this.$emit('update:value', selectedItem);
    },
  },
};
</script>
