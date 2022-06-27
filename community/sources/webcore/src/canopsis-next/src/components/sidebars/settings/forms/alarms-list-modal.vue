<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t(`settings.titles.${$constants.SIDE_BARS.alarmSettings}`) }}
    v-list.grey.lighten-4.px-2.py-0(expand)
      field-columns(
        v-field="form.widgetColumns",
        :label="$t('settings.columnNames')",
        with-html,
        with-state
      )
      v-divider
      field-default-elements-per-page(v-field="form.itemsPerPage")
      v-divider
      field-info-popup(
        v-field="form.infoPopups",
        :columns="form.widgetColumns"
      )
      v-divider
      field-text-editor(
        v-model="form.moreInfoTemplate",
        :title="$t('settings.moreInfosModal')"
      )
</template>

<script>
import FieldColumns from '@/components/sidebars/settings/fields/common/columns.vue';
import FieldInfoPopup from '@/components/sidebars/settings/fields/alarm/info-popup.vue';
import FieldTextEditor from '@/components/sidebars/settings/fields/common/text-editor.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/settings/fields/common/default-elements-per-page.vue';

export default {
  components: {
    FieldColumns,
    FieldInfoPopup,
    FieldTextEditor,
    FieldDefaultElementsPerPage,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
