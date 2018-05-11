<template lang="pug">
  div
    v-subheader {{ $t('tables.alarmGeneral.title') }}
    v-data-table(
    :headers="headers",
    :items="items",
    hide-actions
    )
      template(slot="headerCell", slot-scope="props")
        span {{ $t(props.header.text) }}
      template(slot="items", slot-scope="props")
        td.text-xs-left {{ props.item.name }}
        td.text-xs-left {{ props.item.name }}
        td.text-xs-left {{ props.item.name }}
        td.text-xs-left {{ props.item.name }}
</template>

<script>
export default {
  data() {
    return {
      headers: [
        {
          text: 'tables.alarmGeneral.author',
          sortable: false,
        },
        {
          text: 'tables.alarmGeneral.connector',
          sortable: false,
        },
        {
          text: 'tables.alarmGeneral.component',
          sortable: false,
        },
        {
          text: 'tables.alarmGeneral.resource',
          sortable: false,
        },
      ],
      items: [
        {
          name: 'something',
        },
      ],
    };
  },
};
</script>
