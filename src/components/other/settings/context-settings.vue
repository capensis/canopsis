<template lang="pug">
  v-navigation-drawer(
  :temporary="$mq === 'mobile' || $mq === 'tablet'",
  :value="value",
  stateless,
  clipped,
  right,
  app
  )
    v-toolbar(color="blue darken-4")
      v-list
        v-list-tile
          v-list-tile-title.white--text.text-xs-center {{ title }}
      v-icon.closeIcon(@click.stop="close", color="white") close
    v-divider
    v-list.pt-0(expand)
      div
        field-title
        v-divider
      div
        field-default-column-sort(
        :direction.sync="settings.default_sort_column.direction",
        :property.sync="settings.default_sort_column.property"
        )
        v-divider
      div
        field-context-entities-types-filter
        v-divider
    v-btn(@click="submit", color="green darken-4 white--text", depressed, fixed, right) {{$t('common.save')}}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldDefaultColumnSort from '@/components/other/settings/fields/default-column-sort.vue';
import FieldContextEntitiesTypesFilter from '@/components/other/settings/fields/context-entities-types-filter.vue';

const storeExample = {
  default_sort_column: {
    direction: 'DESC',
    property: 'creation_date',
  },
};

export default {
  components: {
    FieldTitle,
    FieldDefaultColumnSort,
    FieldContextEntitiesTypesFilter,
  },
  props: {
    value: {
      type: Boolean,
      default: false,
    },
    title: {
      type: String,
      default: '',
    },
    fields: {
      type: Array,
      default() {
        return [];
      },
    },
  },
  data() {
    return {
      settings: cloneDeep(storeExample),
    };
  },
  methods: {
    close() {
      this.$emit('input', false);
    },
    submit() {
      console.warn(this.settings);
    },
  },
};
</script>

<style scoped>
  .closeIcon:hover {
    cursor: pointer;
  }
</style>
