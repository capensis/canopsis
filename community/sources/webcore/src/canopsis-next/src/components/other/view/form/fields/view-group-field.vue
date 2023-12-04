<template lang="pug">
  v-combobox.view-group-field(
    v-field="value",
    v-validate="'required'",
    :items="groups",
    :label="$t('view.groupIds')",
    :error-messages="errors.collect('group')",
    item-text="title",
    item-value="_id",
    name="group",
    return-object,
    blur-on-create
  )
    template(#item="{ item }")
      v-list-tile-avatar.view-group-field__avatar(v-if="item.is_private", size="20")
        v-icon(small) lock
      v-list-tile-content
        v-list-tile-title {{ item.title }}
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
