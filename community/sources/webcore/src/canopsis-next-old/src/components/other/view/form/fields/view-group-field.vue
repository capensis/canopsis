<template>
  <v-combobox
    v-field="value"
    v-validate="'required'"
    :items="groups"
    :label="$t('view.groupIds')"
    :error-messages="errors.collect('group')"
    class="view-group-field"
    item-text="title"
    item-value="_id"
    name="group"
    return-object
    blur-on-create
  >
    <template #item="{ item }">
      <v-list-item-avatar
        v-if="item.is_private"
        class="view-group-field__avatar"
        size="20"
      >
        <v-icon small>
          lock
        </v-icon>
      </v-list-item-avatar>
      <v-list-item-content>
        <v-list-item-title>
          {{ item.title }}
        </v-list-item-title>
      </v-list-item-content>
    </template>
    <template #no-data="">
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title v-html="$t('view.noGroupsFound')" />
        </v-list-item-content>
      </v-list-item>
    </template>
  </v-combobox>
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
