<template>
  <c-page>
    <template #header>
      {{ $t('parameters.tabs.importExportViews') }}
    </template>
    <views-import-export />
  </c-page>
</template>

<script>
import ViewsImportExport from '@/components/other/view/views-import-export.vue';

export default {
  components: { ViewsImportExport },
};
</script>
