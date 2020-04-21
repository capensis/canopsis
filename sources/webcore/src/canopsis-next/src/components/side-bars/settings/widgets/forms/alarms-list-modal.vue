<template lang="pug">
  v-list-group(data-test="widgetAlarmsList")
    v-list-tile(slot="activator") {{ $t('settings.titles.alarmListSettings') }}
    v-list.grey.lighten-4.px-2.py-0(expand)
      field-columns(v-field="form.widgetColumns", withHtml)
      v-divider
      field-default-elements-per-page(v-field="form.itemsPerPage")
      v-divider
      field-info-popup(
        v-field="form.infoPopups",
        :columns="form.widgetColumns",
        data-test="widgetInfoPopup"
      )
      v-divider
      field-text-editor(
        v-model="form.moreInfoTemplate",
        :title="$t('settings.moreInfosModal')",
        data-test="widgetMoreInfoTemplate"
      )
</template>

<script>
import FieldColumns from '../fields/common/columns.vue';
import FieldInfoPopup from '../fields/alarm/info-popup.vue';
import FieldTextEditor from '../fields/common/text-editor.vue';
import FieldDefaultElementsPerPage from '../fields/common/default-elements-per-page.vue';

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
