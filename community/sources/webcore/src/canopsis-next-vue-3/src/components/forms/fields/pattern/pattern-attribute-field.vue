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
        <v-list-item-title>
          {{ item?.text }}
        </v-list-item-title>

        <template v-if="item?.alias" #append>
          <c-simple-tooltip :content="`infos.${item.originalValue}.value`" top>
            <template #activator="{ props: activatorProps }">
              <v-icon
                class="ml-1"
                color="primary"
                size="small"
                v-bind="activatorProps"
              >
                alternate_email
              </v-icon>
            </template>
          </c-simple-tooltip>
        </template>
      </v-list-item>
    </template>
    <template v-if="selectedItem?.alias" #prepend-inner>
      <v-icon class="mt-1" color="primary" size="small">
        alternate_email
      </v-icon>
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
        disabled: item?.options?.disabled,
      };
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
