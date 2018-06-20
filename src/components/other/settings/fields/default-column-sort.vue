<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.defaultSortColumn')}}
    v-container
      v-text-field(
      :value="property",
      @input="$emit('update:property', $event)",
      :placeholder="$t('settings.columnName')"
      )
      v-select(
      :value="direction",
      @input="$emit('update:direction', $event)",
      :items="directions"
      )
</template>

<script>
export default {
  props: {
    property: String,
    direction: {
      type: String,
      default: 'ASC',
    },
  },
  data() {
    return {
      directions: ['ASC', 'DESC'],
    };
  },
};
</script>
