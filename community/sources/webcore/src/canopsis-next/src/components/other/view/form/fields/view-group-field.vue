<template lang="pug">
  v-combobox(
    v-field="value",
    v-validate="'required'",
    :items="availableGroups",
    :label="$t('view.groupIds')",
    :error-messages="errors.collect('group')",
    item-text="title",
    item-value="_id",
    name="group",
    return-object,
    blur-on-create
  )
    template(#no-data="")
      v-list-tile
        v-list-tile-content
          v-list-tile-title(v-html="$t('view.noGroupsFound')")
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Object, String],
      required: false,
    },
    groups: {
      type: Array,
      default: () => [],
    },
    private: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    availableGroups() {
      return this.groups.filter(group => group.is_private === this.private);
    },
  },
};
</script>
