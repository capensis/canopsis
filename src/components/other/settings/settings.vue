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
      template(v-for="(field, index) in fields")
        div(:is="`field-${field}`", :key="`settings-field-${index}`")
        v-divider
    v-btn(color="green darken-4 white--text", depressed, fixed, right) {{$t('common.save')}}
</template>

<script>
import FieldTitle from '@/components/other/settings/fields/title.vue';
import FieldDefaultColumnSort from '@/components/other/settings/fields/default-column-sort.vue';
import FieldColumns from '@/components/other/settings/fields/columns.vue';
import FieldPeriodicRefresh from '@/components/other/settings/fields/periodic-refresh.vue';
import FieldDefaultElementsPerPage from '@/components/other/settings/fields/default-elements-per-page.vue';
import FieldOpenedResolvedFilter from '@/components/other/settings/fields/opened-resolved-filter.vue';
import FieldFilters from '@/components/other/settings/fields/filters.vue';
import FieldInfoPopup from '@/components/other/settings/fields/info-popup.vue';
import FieldMoreInfo from '@/components/other/settings/fields/more-info.vue';
import FieldContextEntitiesTypesFilter from '@/components/other/settings/fields/context-entities-types-filter.vue';

export default {
  components: {
    FieldTitle,
    FieldDefaultColumnSort,
    FieldColumns,
    FieldPeriodicRefresh,
    FieldDefaultElementsPerPage,
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldInfoPopup,
    FieldMoreInfo,
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
  methods: {
    close() {
      this.$emit('input', false);
    },
  },
};
</script>

<style scoped>
  .closeIcon:hover {
    cursor: pointer;
  }
</style>
