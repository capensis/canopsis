<template lang="pug">
  div
    v-subheader {{ $t('tables.alarmGeneral.title') }}
    v-data-table(
    :headers="headers",
    :items="[item]",
    hide-actions
    )
      template(slot="headerCell", slot-scope="props")
        span {{ $t(props.header.text) }}
      template(slot="items", slot-scope="props")
        td.text-xs-left {{ props.item.v.state.a }}
        td.text-xs-left {{ props.item.v.connector }}
        td.text-xs-left {{ props.item.v.component }}
        td.text-xs-left {{ props.item.v.resource }}
</template>

<script>
export default {
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
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
    };
  },
};
</script>
