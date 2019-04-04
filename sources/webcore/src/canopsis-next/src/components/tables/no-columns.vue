<template lang="pug" functional>
  .v-table__overflow
    table.v-datatable.v-table.theme--light
      tbody
        tr
          td.text-xs-center {{ $t('tables.noColumns.message') }}
</template>
