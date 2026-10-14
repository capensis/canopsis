<template>
  <c-combobox
    v-model:search="search"
    v-model:menu="menu"
    ref="combobox"
    :model-value="value"
    :items="availableItems"
    :label="label"
    :hide-no-data="false"
    :name="name"
    item-title="text"
    item-value="value"
    return-object
    required
    @keydown.enter.capture="handleEnter"
    @update:model-value="handleModelValueUpdate"
  >
    <template #item="{ props, item }">
      <v-list-item v-bind="getItemProps(props)">
        <v-list-item-title>{{ getItem(item).text }}</v-list-item-title>
      </v-list-item>
    </template>

    <template #selection="{ item }">
      {{ getSelectedItemText(item) }}
    </template>

    <template #no-data>
      <v-list-item>
        <v-list-item-title v-html="$t('kpi.noData')" />
      </v-list-item>
    </template>

    <template v-if="showCreateHintInAppend" #append-item>
      <v-list-item>
        <v-list-item-title v-html="$t('kpi.noData')" />
      </v-list-item>
    </template>
  </c-combobox>
</template>

<script>
export default {
  props: {
    value: {
      type: [Object, String],
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    parameters: {
      type: Array,
      default: () => [],
    },
    disabledParameters: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'metric',
    },
  },

  emits: ['update:value'],

  data() {
    return {
      menu: false,
      search: '',
    };
  },

  computed: {
    currentValue() {
      return this.value?.value ?? this.value ?? '';
    },

    availableItems() {
      return this.parameters.map(value => ({
        value,
        text: this.$t(`alarm.metrics.${value}`),
        disabled: this.disabledParameters.includes(value) && value !== this.currentValue,
      }));
    },

    normalizedSearch() {
      return this.search?.trim().toLowerCase() || '';
    },

    matchedMetric() {
      if (!this.normalizedSearch) {
        return null;
      }

      return this.availableItems.find(
        ({ text }) => text?.trim().toLowerCase() === this.normalizedSearch,
      );
    },

    hasFilteredItems() {
      if (!this.normalizedSearch) {
        return false;
      }

      return this.availableItems.some(({ text }) => text?.trim().toLowerCase().includes(this.normalizedSearch),
      );
    },

    showCreateHintInAppend() {
      return this.hasFilteredItems && !this.matchedMetric;
    },
  },

  methods: {
    getItem(item) {
      return item?.raw ?? item ?? {};
    },

    getSelectedItemText(item) {
      const selectedItem = this.getItem(item);

      return selectedItem.text
        ?? this.availableItems.find(({ value }) => value === this.currentValue)?.text
        ?? this.currentValue;
    },

    getItemProps(props) {
      const { title, ...listItemProps } = props;

      return listItemProps;
    },

    handleModelValueUpdate(value) {
      this.$emit('update:value', value?.value ?? value ?? '');
    },

    handleEnter(event) {
      if (!this.normalizedSearch || this.matchedMetric) {
        return;
      }

      event.preventDefault();
      event.stopPropagation();
      this.$emit('update:value', this.matchedMetric?.value ?? this.search.trim());
      this.closeMenuAndBlur();
    },

    closeMenuAndBlur() {
      this.menu = false;
      this.blurInput();

      this.$nextTick(() => {
        this.menu = false;
        this.blurInput();
      });
    },

    blurInput() {
      this.$refs.combobox?.$el?.querySelector('input')?.blur();
    },
  },
};
</script>
