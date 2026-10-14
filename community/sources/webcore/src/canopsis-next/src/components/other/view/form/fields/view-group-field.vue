<template>
  <c-combobox
    v-model:search="search"
    v-model:menu="menu"
    ref="combobox"
    :model-value="value"
    :items="groups"
    :label="$t('view.groupIds')"
    :hide-no-data="false"
    class="view-group-field"
    item-title="title"
    item-value="_id"
    name="group"
    return-object
    required
    @keydown.enter.capture="handleEnter"
    @update:model-value="handleModelValueUpdate"
  >
    <template #item="{ props, item }">
      <v-list-item v-bind="getItemProps(props)">
        <template v-if="item?.is_private" #prepend>
          <v-avatar class="view-group-field__avatar" size="20">
            <v-icon size="small">
              lock
            </v-icon>
          </v-avatar>
        </template>

        <v-list-item-title>
          {{ item?.title }}
        </v-list-item-title>
      </v-list-item>
    </template>

    <template #no-data>
      <v-list-item>
        <v-list-item-title v-html="$t('view.noGroupsFound')" />
      </v-list-item>
    </template>

    <template v-if="showCreateHintInAppend" #append-item>
      <v-list-item>
        <v-list-item-title v-html="$t('view.noGroupsFound')" />
      </v-list-item>
    </template>
  </c-combobox>
</template>

<script>
export default {
  props: {
    value: {
      type: [Object, String],
      required: false,
    },
    groups: {
      type: Array,
      default: () => [],
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
    normalizedSearch() {
      const search = this.search?.trim().toLowerCase();

      return search || '';
    },

    matchedGroup() {
      if (!this.normalizedSearch) {
        return null;
      }

      return this.groups.find(({ title }) => title?.trim().toLowerCase() === this.normalizedSearch);
    },

    hasExactGroup() {
      return !!this.matchedGroup;
    },

    hasFilteredGroups() {
      if (!this.normalizedSearch) {
        return false;
      }

      return this.groups.some(({ title }) => title?.trim().toLowerCase().includes(this.normalizedSearch),
      );
    },

    showCreateHintInAppend() {
      const search = this.normalizedSearch;

      if (!search) {
        return false;
      }

      return this.hasFilteredGroups && !this.hasExactGroup;
    },
  },

  methods: {
    getItemProps(props) {
      const { title, ...listItemProps } = props;

      return listItemProps;
    },

    handleModelValueUpdate(value) {
      this.$emit('update:value', value);
    },

    handleEnter(event) {
      if (!this.normalizedSearch || this.hasExactGroup) {
        return;
      }

      event.preventDefault();
      event.stopPropagation();

      this.$emit('update:value', this.search.trim());
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

<style lang="scss">
.view-group-field {
  &__avatar {
    min-width: 30px;
  }
}
</style>
