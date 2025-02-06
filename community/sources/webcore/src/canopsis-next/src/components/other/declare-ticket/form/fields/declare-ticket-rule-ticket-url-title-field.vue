<template lang="pug">
  v-text-field(
    v-field="value",
    :label="$t('declareTicket.ticketURLTitle')"
  )
    template(#append="")
      c-help-icon(
        :text="$t('declareTicket.ticketURLTitleHelp')",
        icon="help",
        color="grey darken-1",
        left
      )
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
  },
};
</script>
